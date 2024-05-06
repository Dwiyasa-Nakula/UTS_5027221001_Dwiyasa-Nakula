package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Dwiyasa-Nakula/master/backend/genproto/musicplaylist"
	"github.com/Dwiyasa-Nakula/master/backend/repository"
	"github.com/Dwiyasa-Nakula/master/backend/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// init initializes the configuration.
// It loads the default-config file or the specified config file from command-line arguments.
func init() {
	args := os.Args[1:]
	var configname string = "default-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}

	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s.yml", err))
	}
}

func main() {
	// Log the start of the GRPC server.
	log.Println("Starting up GRPC server")

	// Create connection to database.
	log.Println("Creating connection to database")
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("app.mongodb.uri")))
	if err != nil {
		log.Fatalf("%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db := client.Database(viper.GetString("app.mongodb.database"))
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Connected to database...")

	// Create new GRPC server.
	server := grpc.NewServer()

	// Initialize repository and service.
	urepo := repository.NewSongRepo(db)
	usvc := service.NewSongService(urepo)
	musicplaylist.RegisterSongApiServer(server, usvc)

	// Get port from configuration.
	port := ":" + viper.GetString("app.grpc.port")

	// Listen on the specified port.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", port, err)
	}

	// Start the server and panic if there's an error.
	panic(server.Serve(listener))
}