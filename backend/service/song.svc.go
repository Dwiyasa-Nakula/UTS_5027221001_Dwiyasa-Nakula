package service

import (
	"context"
	"log"

	"github.com/Dwiyasa-Nakula/master/backend/genproto/musicplaylist"
	"github.com/Dwiyasa-Nakula/master/backend/model"
	"github.com/Dwiyasa-Nakula/master/backend/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// SongService handles gRPC requests related to songs.
type SongService struct {
	musicplaylist.UnimplementedSongApiServer // Embed the generated gRPC server interface
	repo *repository.SongRepository          // Repository to interact with the database
}

// NewSongService creates a new instance of SongService.
func NewSongService(repo *repository.SongRepository) *SongService {
	return &SongService{
		repo: repo,
	}
}

// CreateSong creates a new song.
// It takes a context and a musicplaylist.Song as input.
// It returns the created song along with any error encountered.
func (s *SongService) CreateSong(ctx context.Context, tm *musicplaylist.Song) (*musicplaylist.Song, error) {
	log.Printf("CreateSong(%v) \n", tm)

	// Convert the received gRPC song to a model song
	newSong := &model.Song{ 
		Title:    	tm.Title,
		Artist: 	tm.Artist,
		Album: 		tm.Album,
		Duration: 	tm.Duration,
		Link:	 	tm.Link,
	}

	// Save the new song in the repository
	song, err := s.repo.Save(newSong)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// Convert the model song back to a gRPC song and return
	return s.toSong(&song), nil
}

// ListSongs retrieves a list of all songs.
// It takes a context and an empty message as input.
// It returns a list of songs along with any error encountered.
func (s *SongService) ListSongs(ctx context.Context, e *empty.Empty) (*musicplaylist.SongList, error) {
	log.Printf("ListSongs() \n")

	// Retrieve all songs from the repository
	var totas []*musicplaylist.Song
	Songs, err := s.repo.FindAll()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// Convert each model song to a gRPC song
	for _, u := range Songs {
		totas = append(totas, s.toSong(&u))
	}

	// Create a gRPC song list and return
	SongList := &musicplaylist.SongList{
		List: totas,
	}

	return SongList, nil
}

// UpdateSong updates an existing song.
// It takes a context and a musicplaylist.Song as input.
// It returns the updated song along with any error encountered.
func (s *SongService) UpdateSong(ctx context.Context, tm *musicplaylist.Song) (*musicplaylist.Song, error) {
	log.Printf("UpdateSong(%v) \n", tm)

	// Check if the song ID is provided
	if tm.Id == "" {
		return nil, status.Error(codes.FailedPrecondition, "UpdateSong must provide songID")
	}

	// Convert the song ID to an ObjectID
	songID, err := primitive.ObjectIDFromHex(tm.Id)
	if err != nil {
		log.Printf("Invalid SongID(%s) \n", tm.Id)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Create a model song with updated fields
	updateSong := &model.Song{
		ID:          songID,
		Title:       tm.Title,
		Artist: 	 tm.Artist,
		Album: 		 tm.Album,
		Duration: 	 tm.Duration,
		Link: 		 tm.Link,
	}

	// Update the song in the repository
	song, err := s.repo.Update(updateSong)
	if err != nil {
		log.Printf("Fail UpdateSong %v \n", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Convert the updated model song back to a gRPC song and return
	return s.toSong(&song), nil
}

// DeleteSong deletes an existing song.
// It takes a context and a string value (song ID) as input.
// It returns a boolean indicating the deletion success along with any error encountered.
func (s *SongService) DeleteSong(ctx context.Context, id *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("DeleteSong(%s) \n", id.GetValue())

	// Delete the song from the repository
	deleted, err := s.repo.Delete(id.GetValue())
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// Return a boolean indicating the deletion success
	return &wrapperspb.BoolValue{Value: deleted}, nil
}

// toSong converts a model.Song to a musicplaylist.Song.
// It takes a model song as input and returns the equivalent gRPC song.
func (s *SongService) toSong(u *model.Song) *musicplaylist.Song {
	tota := &musicplaylist.Song{
		Id:          u.ID.Hex(),
		Title:       u.Title,
		Artist: 	 u.Artist,
		Album: 		 u.Album,
		Duration: 	 u.Duration,
		Link:	 	 u.Link,
	}
	return tota
}