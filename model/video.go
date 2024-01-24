// model/video.go
// model/video.go
package model

type Video struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PublishedAt string     `json:"published_at"` // Keep it as a string for now
	Thumbnails  Thumbnails `json:"thumbnails"`
	// Add other fields as needed
}

type Thumbnails struct {
	Default string `json:"default"`
	Medium  string `json:"medium"`
	High    string `json:"high"`
}
