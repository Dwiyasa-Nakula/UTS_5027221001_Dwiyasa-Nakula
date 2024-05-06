package repository

import (
	"context"
	"log"
	"time"

	"github.com/Dwiyasa-Nakula/master/backend/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SongRepository handles operations related to songs in the database.
type SongRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// NewSongRepo creates a new instance of SongRepository.
func NewSongRepo(db *mongo.Database) *SongRepository {
	return &SongRepository{
		db:  db,
		col: db.Collection(model.SongCollection),
	}
}

// Save inserts a new song into the database.
// It takes a pointer to a model.Song as input and returns the saved song along with any error encountered.
func (r *SongRepository) Save(u *model.Song) (model.Song, error) {
	log.Printf("Save(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	var song model.Song
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return song, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&song)
	if err != nil {
		log.Println(err)
		return song, err
	}

	return song, nil
}

// FindAll retrieves all songs from the database.
// It returns a slice of songs along with any error encountered.
func (r *SongRepository) FindAll() ([]model.Song, error) {
	log.Println("FindAll()")
	ctx, cancel := timeoutContext()
	defer cancel()

	var songs []model.Song
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return songs, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var song model.Song
		err := cur.Decode(&song)
		if err != nil {
			log.Println(err)
		}
		songs = append(songs, song)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return songs, nil
}

// Update updates an existing song in the database.
// It takes a pointer to a model.Song as input and returns the updated song along with any error encountered.
func (r *SongRepository) Update(u *model.Song) (model.Song, error) {
	log.Printf("Update(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	filter := bson.M{"_id": u.ID}
	update := bson.M{
		"$set": bson.M{
			"title":       	u.Title,
			"artist":	 	u.Artist,
			"album":       	u.Album,
			"duration": 	u.Duration,
			"link": 		u.Link,
		},
	}

	var song model.Song
	err := r.col.FindOneAndUpdate(ctx, filter, update).Decode(&song)
	if err != nil {
		log.Printf("ERR 115 %v", err)
		return song, err
	}

	return song, nil
}

// Delete deletes a song from the database by its ID.
// It takes a string representing the song ID as input and returns a boolean indicating the deletion success along with any error encountered.
func (r *SongRepository) Delete(id string) (bool, error) {
	log.Printf("Delete(%s) \n", id)
	ctx, cancel := timeoutContext()
	defer cancel()

	var song model.Song
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&song)
	if err != nil {
		log.Printf("Fail to delete song: %v \n", err)
		return false, err
	}
	log.Printf("Deleted_song(%v) \n", song)
	return true, nil
}

// timeoutContext creates a context with a timeout of 60 seconds.
func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
}