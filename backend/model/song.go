package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SongCollection is the name of the MongoDB collection where song documents are stored.
const SongCollection = "song"

// Song represents a song in the music playlist.
type Song struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` // Unique identifier for the song
	Title    string             `bson:"title"`          // Title of the song
	Artist   string             `bson:"artist"`         // Artist of the song
	Album    string             `bson:"album"`          // Album of the song
	Duration string             `bson:"duration"`       // Duration of the song
	Link     string             `bson:"link"`           // Link to the song (e.g., SoundCloud track number)
}
