// service/youtube_service.go
package service

import (
	"Fampay_Backend_Assignment/model"
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var mongoClient *mongo.Client
var youtubeService *youtube.Service

const fetchInterval = 10 * time.Second

func init() {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Error pinging MongoDB server: %v", err)
	}
	mongoClient = client
	apiKey := os.Getenv("API_KEY")
	youtubeService, err = youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}
}

func FetchAndStoreVideos(query string) error {
	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	for {
		if err := performFetchAndStore(query); err != nil {
			log.Printf("Error fetching and storing videos: %v", err)
		} else {
			log.Print("Fetched and stored videos successfully")
		}
		<-ticker.C
	}
}

func performFetchAndStore(query string) error {
	call := youtubeService.Search.List([]string{"snippet"}).
		Q(query).
		MaxResults(20)

	response, err := call.Do()
	if err != nil {
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

	if err := StoreVideosInMongoDB(videos); err != nil {
		return err
	}

	return nil
}

func StoreVideosInMongoDB(videos []model.Video) error {
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

func GetLatestVideos(page, pageSize int) ([]model.Video, error) {
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	options := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{"publishedat", -1}})

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

func GetPaginatedVideos(page, pageSize int) ([]model.Video, error) {
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	collection := mongoClient.Database(databaseName).Collection(collectionName)

	options := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{"publishedat", -1}})

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
