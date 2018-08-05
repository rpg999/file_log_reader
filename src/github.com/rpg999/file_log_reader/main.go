package main

import (
	"context"
	"flag"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/rpg999/file_log_reader/log_reader"
	"log"
)

var fileList = flag.String("file_list", "bin/file_list.json", "List of files to track in json format")

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(*fileList) == 0 {
		log.Fatal("no file list is specified")
	}

	files, err := log_reader.ParseFileList(*fileList)
	if err != nil {
		log.Fatalf("error while parsing file list: %s", err)
	}

	if len(files) == 0 {
		log.Fatal("no valid files for tracking")
	}

	// For simplicity we will use the default db connection
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("bitsane_test").Collection("log")

	// Done channel can make all goroutines exist safely
	done := make(chan struct{})
	logStream, err := log_reader.TrackFiles(done, files)
	if err != nil {
		log.Println(err)
		close(done)

		return
	}

	for l := range logStream {
		// Here we can do something with errors, for example, if there are more than 10 errors then we close done channel
		if l.Err != nil {
			log.Println(l.Err)
			continue
		}

		res, err := collection.InsertOne(context.Background(), l.Log)
		if err != nil {
			log.Println(err)
		}
		id := res.InsertedID
		log.Println("New entry is added with id: ", id)
	}

	log.Println("End of program")
}
