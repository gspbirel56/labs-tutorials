package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

	var postsCollection mongo.Collection

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		// Log the error and return(1) for the program.
		log.Fatal(err)
	}

	ctx, _:= context.WithTimeout(context.Background(), 10 * time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	// Test the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
    	log.Fatal(err)
	}
	fmt.Println("Ping successful!")

	// List the available databases on the server
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
    	log.Fatal(err)
	}
	fmt.Println(databases)

	// Get a handle to a specific collection
	tutorialDatabase := client.Database("GoTutorial")
	postsCollection = *tutorialDatabase.Collection("Posts")

	// Insert a post
	// insertPost(ctx, bson.D{
	// 	{"title", "Hello!"},
	// 	{"body", "This is another post from Go"},
	// })

	// Read all posts from the collection
	// getAllPosts(ctx, false)
	// getAllPosts(ctx, true) // better for a larger dataset
	// getSinglePost(ctx, bson.M{"title": "Hello!"}) // note that this is a filter.
		// Because it is a filter, the parameters have to be exact.

	// Update a single record
	// Note that there is an UpdateMany() function as referenced here: https://www.mongodb.com/blog/post/quick-start-golang--mongodb--how-to-update-documents.
	// This may be useful when updating records for several students in a class, for example.
	// Note that updating a single record will replace its contents with whatever is specified, while UpdateMany() will
	//	append the docuemnt or replace a record.
	updateSingelRecord(ctx, bson.M{"title": "Inserted from Go"}, bson.M{
		"title": "Inserted from Go",
		"body": "This post was appended using the ReplaceOne function.",
	})


	// Why not, update multiple records.
	updateMultipleRecords(ctx, bson.M{"title": "Hello!"}, bson.D{
		{"$set", bson.D{{"lastUpdated", time.Now()}},
		},
	})
}

func insertPost(ctx context.Context, newPost bson.D) {
	insertResult, err := postsCollection.InsertOne(ctx, newPost)
	if (err != nil) {
		log.Fatal(err)
	}

	fmt.Printf("Record with ID %v inserted into the collection!", insertResult.InsertedID)
}

func getAllPosts(ctx context.Context, batches bool) {
	cursor, err := postsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if (!batches) {
		var posts []bson.M
		if err = cursor.All(ctx, &posts); err != nil {
			log.Fatal(err)
		}
		fmt.Println(posts)
	} else {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var post bson.M
			if err = cursor.Decode(&post); err != nil {
				log.Fatal(err)
			}
			fmt.Println(post)
		}
	}
}

func getSinglePost(ctx context.Context, filter bson.M) {
	filterCursor, err := postsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	var postsFiltered []bson.M
	if err = filterCursor.All(ctx, &postsFiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(postsFiltered)
}

func updateSingelRecord(ctx context.Context, filter bson.M, replacementObject bson.M) {
	result, err := postsCollection.ReplaceOne(
		ctx,
		filter,
		replacementObject, // Just figured this out, you need a comma to continue the statement to the next line in Go.
							// Otherwise the compiler will throw an error for unexpected end of line or something.
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Replaced %v document(s)\n", result.ModifiedCount)
}

func updateMultipleRecords(ctx context.Context, filter bson.M, appendObject bson.D) {
	result, err := postsCollection.UpdateMany(
		ctx,
		filter,
		appendObject,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v document(s)\n", result.ModifiedCount)
}