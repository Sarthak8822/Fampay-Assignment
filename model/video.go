// model/video.go
package model

// Video struct to represent video data
type Video struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PublishedAt string     `json:"published_at"` // Keep it as a string for now
	Thumbnails  Thumbnails `json:"thumbnails"`
	// Add other fields as needed
}

// Thumbnails struct to represent video thumbnail URLs
type Thumbnails struct {
	Default string `json:"default"`
	Medium  string `json:"medium"`
	High    string `json:"high"`
}
