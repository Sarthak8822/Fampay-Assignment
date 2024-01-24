package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"Fampay_Backend_Assignment/model"
)

const (
	apiKey         = "AIzaSyB_IW1vLGuT_Oyco7cT6VwArK1Ff1OCc3o"
	mongoURI       = "mongodb+srv://sarthakpravin08:fampayassignment@cluster0.82iwxhm.mongodb.net/"
	databaseName   = "All_Videos"
	collectionName = "videos"
)

var mongoClient *mongo.Client
var youtubeService *youtube.Service

// fetchInterval is the time interval between each fetch operation
const fetchInterval = 10 * time.Second

func init() {
	// Initialize MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}

	// Ping the MongoDB server to ensure it's reachable
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB server: %v", err)
	}
	mongoClient = client

	// Initialize YouTube service
	youtubeService, err = youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}
}

func FetchAndStoreVideos(query string) {
	fmt.Println("Inside function")

	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	for {
		videos, err := performFetchAndStore(query)
		if err != nil {
			log.Printf("Error fetching and storing videos: %v", err)
		} else {
			log.Printf("Fetched and stored %d videos", len(videos))
		}

		// Wait for the next tick
		<-ticker.C
	}
}

func performFetchAndStore(query string) ([]model.Video, error) {
	fmt.Println("Inside function")

	call := youtubeService.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(20)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var videos []model.Video
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return nil, err
		}

		video := model.Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: publishedAt.Format(time.RFC3339),
			Thumbnails: model.Thumbnails{
				Default: item.Snippet.Thumbnails.Default.Url,
				Medium:  item.Snippet.Thumbnails.Medium.Url,
				High:    item.Snippet.Thumbnails.High.Url,
			},
		}
		videos = append(videos, video)
	}

	err = StoreVideosInMongoDB(videos)
	if err != nil {
		return nil, err
	}

	fmt.Println("after youtube service called", videos)

	return videos, nil
}

func StoreVideosInMongoDB(videos []model.Video) error {
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	for _, video := range videos {
		// Create a filter for finding existing documents with the same ID
		filter := bson.D{{Key: "id", Value: video.ID}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "title", Value: video.Title},
				{Key: "description", Value: video.Description},
				{Key: "publishedat", Value: video.PublishedAt},
				{Key: "thumbnails", Value: video.Thumbnails},
			}},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetLatestVideos(page, pageSize int) ([]model.Video, error) {
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	// Define options for pagination and sorting
	options := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{"publishedat", -1}}) // Sort by publishedAt in descending order

	// Execute MongoDB query
	cursor, err := collection.Find(context.Background(), bson.D{}, options)
	if err != nil {
		// Handle error (e.g., log it)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var videos []model.Video
	for cursor.Next(context.Background()) {
		var video model.Video
		if err := cursor.Decode(&video); err != nil {
			// Handle error (e.g., log it)
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}
