package main

import (
    "fmt"
    "log"
    "golang.org/x/net/context"
    "cloud.google.com/go/datastore"
)
type Entity struct { 
    Value string
}

func main() {
	ctx := context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	dsClient, err := datastore.NewClient(ctx, "stunning-symbol-139515")
	if err != nil {
		log.Fatal(err)
	}

	k := datastore.NewKey(ctx, "Entity", "stringID", 0, nil)
	e := new(Entity)
	// if err := dsClient.Get(ctx, k, e); err != nil {
	// 	log.Fatal(err)
	// }

	old := "" //e.Value
	e.Value = "Hello World!"

	if _, err := dsClient.Put(ctx, k, e); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated value from %q to %q\n", old, e.Value)
}