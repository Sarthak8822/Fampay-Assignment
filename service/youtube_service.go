// service/youtube_service.go
package service

import (
	"Fampay_Backend_Assignment/model"
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	mongoClient   *mongo.Client
	youtubeClient *youtube.Service
	apiKeys       []string
	apiKeyIndex   int
	apiKeyMutex   sync.Mutex
)

const fetchInterval = 10 * time.Second

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, Please add a .env file if not Added")
	}

	mongoURI, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	apiKeys = getAPIKeys()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Error pinging MongoDB server: %v", err)
	}
	mongoClient = client

	// Initialize the YouTube client with the first API key
	initYouTubeClient(apiKeys[0])
}

func initYouTubeClient(apiKey string) {
	var err error
	youtubeClient, err = youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}
}

func getAPIKeys() []string {
	apiKeyStr, exists := os.LookupEnv("API_KEYS")
	if !exists {
		log.Fatal("API_KEYS environment variable is not set")
	}
	keys := strings.Split(apiKeyStr, ",")
	if len(keys) == 0 {
		log.Fatal("No API keys provided")
	}
	if len(keys) < 2 {
		log.Fatal("Provide More API keys to handle 403 Errors")
	}
	return keys
}

func FetchAndStoreVideos(query string) error {
	for {
		if err := performFetchAndStore(query); err != nil {
			log.Printf("Error fetching and storing videos: %v", err)
		} else {
			log.Print("Fetched and stored videos successfully")
		}
		time.Sleep(fetchInterval)
	}
}

// should be private
func performFetchAndStore(query string) error {
	call := youtubeClient.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(20)

	response, err := call.Do()
	if err != nil {
		// Check if the error is due to quota exhaustion
		if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == 403 {
			log.Printf("YouTube API Error: %v", apiErr)
			// If Yes then switch it to another API_KEY
			switchAPIKey()
			return performFetchAndStore(query)
		}
		return err
	}

	var videos []model.Video
	for _, item := range response.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			return err
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

	if err := storeVideosInMongoDB(videos); err != nil {
		return err
	}

	return nil
}

// The switchAPIKey function is introduced to manage the rotation of API keys atomically using a mutex.
func switchAPIKey() {
	apiKeyMutex.Lock()
	defer apiKeyMutex.Unlock()

	apiKeyIndex = (apiKeyIndex + 1) % len(apiKeys)
	initYouTubeClient(apiKeys[apiKeyIndex])
}

// should be private
func storeVideosInMongoDB(videos []model.Video) error {
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	for _, video := range videos {
		filter := bson.D{{Key: "id", Value: video.ID}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "title", Value: video.Title},
				{Key: "description", Value: video.Description},
				{Key: "publishedat", Value: video.PublishedAt},
				{Key: "thumbnails", Value: video.Thumbnails},
			}},
		}

		if _, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true)); err != nil {
			return err
		}
	}

	return nil
}

func GetPaginatedVideos(page, pageSize int) ([]model.Video, error) {
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	options := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "publishedat", Value: -1}}) // To retrieve the stored videos in reverse chronological order.

	cursor, err := collection.Find(context.Background(), bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var videos []model.Video
	for cursor.Next(context.Background()) {
		var video model.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}
