package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	// "cloud.google.com/go/datastore"
	// "golang.org/x/net/context"
	// "google.golang.org/api/iterator"
)

type Post struct {
	Title             string
	PubTime           time.Time
	Content           []byte // the result of parsing and executing the templates, for easy serving
	Path              string // the path - last part of the url - that is the address of this post. Might want to make this an [] ?
	StorageBucketPath string
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World?"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)

	// ctx := context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	// dsClient, err := datastore.NewClient(ctx, "stunning-symbol-139515")
	// if err != nil {
	// log.Fatal(err)
	// }

	// q := datastore.NewQuery("Post")

	// t := dsClient.Run(ctx, q)
	// for {
	// 	var p Post
	// 	_, err := t.Next(&p)
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Fprintf(w, "Post:\t%q\n", p.Title)
	// }
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
