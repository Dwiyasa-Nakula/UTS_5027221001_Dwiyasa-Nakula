package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Dwiyasa-Nakula/master/backend/genproto/musicplaylist"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// httpServer represents an HTTP server.
type httpServer struct {
	addr string
}

// NewHttpServer creates a new instance of httpServer.
func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

// init initializes the configuration.
func init() {
	args := os.Args[1:]
	var configname string = "default-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}
	log.Printf("loading config file %s.yml \n", configname)

	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Fatal error config file: " + err.Error())
	}
}

// Run starts the HTTP server.
func (s *httpServer) Run() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/create", s.handleCreate)
	http.HandleFunc("/update", s.handleUpdate)
	http.HandleFunc("/delete", s.handleDelete)
	http.HandleFunc("/playlist", s.handleList)
	log.Printf("Starting HTTP server on localhost%s/playlist\n", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

// handleIndex handles requests to the index page.
func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(songsTemplate))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleCreate handles requests to create a new song.
func (s *httpServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from HTML form.
	title := r.FormValue("title")
	artist := r.FormValue("artist")
	album := r.FormValue("album")
	duration := r.FormValue("duration")
	link := r.FormValue("link")

	// Initialize gRPC connection.
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Create song client.
	songClient := musicplaylist.NewSongApiClient(client)

	// Create new song.
	_, err = songClient.CreateSong(context.Background(), &musicplaylist.Song{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
		Link:     link,
	})
	if err != nil {
		http.Error(w, "Failed to create song", http.StatusInternalServerError)
		return
	}

	// Redirect to list page.
	http.Redirect(w, r, "/playlist", http.StatusSeeOther)
}

// handleUpdate handles requests to update an existing song.
func (s *httpServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// Retrieve data from HTML form.
	id := r.FormValue("id")
	title := r.FormValue("title")
	artist := r.FormValue("artist")
	album := r.FormValue("album")
	duration := r.FormValue("duration")
	link := r.FormValue("link")

	// Initialize gRPC connection.
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Create song client.
	songClient := musicplaylist.NewSongApiClient(client)

	// Update song.
	_, err = songClient.UpdateSong(context.Background(), &musicplaylist.Song{
		Id:       id,
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
		Link:     link,
	})
	if err != nil {
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	// Redirect to list page.
	http.Redirect(w, r, "/playlist", http.StatusSeeOther)
}

// handleDelete handles requests to delete an existing song.
func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Get song ID from URL parameter.
	id := r.URL.Query().Get("id")

	// Initialize gRPC connection.
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Create song client.
	songClient := musicplaylist.NewSongApiClient(client)

	// Delete song.
	_, err = songClient.DeleteSong(context.Background(), &wrapperspb.StringValue{Value: id})
	if err != nil {
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	// Redirect to list page.
	http.Redirect(w, r, "/playlist", http.StatusSeeOther)
}

// handleList handles requests to list all songs.
func (s *httpServer) handleList(w http.ResponseWriter, r *http.Request) {
	// Initialize gRPC connection.
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Create song client.
	songClient := musicplaylist.NewSongApiClient(client)

	// Fetch list of songs from server.
	songs, err := songClient.ListSongs(context.Background(), &emptypb.Empty{})
	if err != nil {
		http.Error(w, "Failed to fetch songs: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to fetch songs: %v\n", err)
		return
	}

	// Prepare song data for display in HTML page.
	type ViewData struct {
		Songs []*musicplaylist.Song
	}
	data := ViewData{
		Songs: songs.List,
	}

	// Create HTML template.
	tmpl := template.Must(template.New("index").Parse(songsTemplate))

	// Display HTML template with prepared song data.
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//run the local server
func main() {
	httpServer := NewHttpServer(":9999")
	httpServer.Run()
}

// songsTemplate defines the HTML template for displaying the song list.
var songsTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Music Playlist</title>
    <style>
	body {
		font-family: Arial, sans-serif;
		margin: 0;
		padding: 0;
		background-color: #000000;
		background-image: 
			radial-gradient(at 47% 33%, hsl(162.00, 77%, 40%) 0, transparent 59%), 
			radial-gradient(at 82% 65%, hsl(218.00, 39%, 11%) 0, transparent 55%);
	}
	.container {
		max-width: 800px;
		margin: 20px auto;
		padding: 20px;
		backdrop-filter: blur(6px) saturate(200%);
		-webkit-backdrop-filter: blur(6px) saturate(200%);
		background-color: rgba(17, 25, 40, 0.8);
		border-radius: 12px;
		border: 1px solid rgba(255, 255, 255, 0.125);
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
	}
	h1 {
		color: #fff;
	}
	h2 {
		color: #fff;
	}
	form {
		margin-bottom: 20px;
	}
	label {
		display: block;
		margin-bottom: 5px;
		color: #fff;
	}
	input[type="text"],
	textarea {
		width: 100%;
		padding: 10px;
		margin-bottom: 10px;
		border: 1px solid #ccc;
		border-radius: 4px;
		box-sizing: border-box;
	}
	input[type="submit"] {
		background-color: #4caf50;
		color: white;
		border: none;
		border-radius: 4px;
		padding: 10px 20px;
		cursor: pointer;
	}
	input[type="submit"]:hover {
		background-color: #45a049;
	}
	ul {
		list-style-type: none;
		padding: 0;
	}
	li {
		padding: 10px;
		margin-bottom: 5px;
		background-color: #f9f9f9;
		border-radius: 4px;
		box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
	}
	a {
		text-decoration: none;
		color: #4caf50;
	}
	a:hover {
		text-decoration: underline;
	}
	.refresh-btn {
		background-color: #008CBA;
		color: white;
		border: none;
		border-radius: 4px;
		padding: 10px 20px;
		cursor: pointer;
		text-decoration: none;
	}
	.refresh-btn:hover {
		background-color: #005f6b;
	}
	.update-form {
		display: none;
		margin-bottom: 10px;
	}
	.update-form input[type="submit"] {
		background-color: #4caf50;
		color: white;
		border: none;
		border-radius: 4px;
		padding: 10px 20px;
		cursor: pointer;
	}
	.back-btn {
		background-color: #f44336;
		color: white;
		border: none;
		border-radius: 4px;
		padding: 10px 20px;
		cursor: pointer;
		text-decoration: none;
		margin-right: 10px; /* Margin kanan untuk memberi jarak dari tombol "Update" */
	}
	.back-btn:hover {
		background-color: #d32f2f;
	}
	.grid-form {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 10px;
	}

	.form-group {
		grid-column: span 1;
	}

	.submit-group {
		grid-column: span 2;
	}
	.action-buttons {
	float: right;
	}
    </style>
</head>
<body>
<div class="container">
    <h1>Music Playlist</h1>
    <form action="/create" method="post" class="grid-form">
        <div class="form-group">
            <label for="title">Track Title:</label>
            <input type="text" id="title" name="title" required>
        </div>
        <div class="form-group">
            <label for="artist">Artist:</label>
            <input type="text" id="artist" name="artist" required>
        </div>
        <div class="form-group">
            <label for="album">Album:</label>
            <input type="text" id="album" name="album" required>
        </div>
        <div class="form-group">
            <label for="duration">Duration:</label>
            <input type="text" id="duration" name="duration" required>
        </div>
		<div class="form-group">
            <label for="link">Soundcloud Track Number: (embaded)</label>
            <input type="text" id="link" name="link" required>
        </div>
        <div class="form-group submit-group">
            <input type="submit" value="Add Track">
        </div>
    </form>    
    <hr>
    <h2>Playlist</h2>
    <a href="/playlist" class="refresh-btn">Refresh Playlist</a>
    <ul>
        {{if not (eq (len .Songs) 0)}}
            {{range .Songs}}
            <li>
				<span>{{.Title}} - {{.Artist}} - {{.Album}} - {{.Duration}}</span>
				<div class="action-buttons">
					<a href="#" onclick="showUpdateForm('{{.Id}}')">Update</a> 
					<a style="color: #d32f2f;" href="/delete?id={{.Id}}">Delete</a>
				</div>
				<iframe width="100%" height="166" scrolling="no" frameborder="no" allow="autoplay"
                    src="https://w.soundcloud.com/player/?url=https%3A//api.soundcloud.com/tracks/{{.Link}}&amp;color=%23ff5500&amp;auto_play=false&amp;hide_related=true&amp;show_comments=false&amp;show_user=false&amp;show_reposts=false&amp;show_teaser=true&amp;visual=true">
                </iframe>
                <form id="updateForm{{.Id}}" class="update-form" action="/update" method="post">
					<input type="hidden" name="id" value="{{.Id}}">
					<label for="title{{.Id}}">New Title:</label><br>
					<input type="text" id="title{{.Id}}" name="title" value="{{.Title}}"><br>
					<label for="artist{{.Id}}">New Artist:</label><br>
					<input id="artist{{.Id}}" name="artist">{{.Artist}}</input><br><br>
                    <label for="album{{.Id}}">New Album:</label><br>
					<input type="text" id="album{{.Id}}" name="album" value="{{.Album}}"><br>
					<label for="duration{{.Id}}">New Duration:</label><br>
					<input id="duration{{.Id}}" name="duration">{{.Duration}}</input><br><br>
					<label for="link{{.Id}}">New Link:</label><br>
					<input id="link{{.Id}}" name="link">{{.Link}}</input><br><br>
					<input type="submit" value="Update Song">
					<a href="/playlist" class="back-btn">Back</a>
				</form>
            </li>
            {{end}}
            {{else}}
                <li>No songs available</li>
            {{end}}
    </ul>
</div>

<script>
    // Function to show the update track form
    function showUpdateForm(trackId) {
        var formId = 'updateForm' + trackId;
        var form = document.getElementById(formId);
        if (form.style.display === 'none') {
            form.style.display = 'block';
        } else {
            form.style.display = 'none';
        }
    }
</script>
</body>
</html>`