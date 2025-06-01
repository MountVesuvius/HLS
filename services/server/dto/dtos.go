package dto

import (
	"time"
	"gorm.io/gorm"
)

// not sure if it's worth having as mega struct, but makes sense
type File struct {
	gorm.Model // includes id, createdat, updatedat, deletedat

	Name           string
	Description    string
	UploaderId     string
	State          string
	Tags           []string `gorm:"type:text[]"`
	ThumbnailURL   string
	DurationBucket int // enum that assigns length
	StreamedCount  int

	UploadRequestTime   time.Time
	UploadCompletion    time.Time
	ParseRequestTime    time.Time
	ParseCompletionTime time.Time
	ParseDuration       time.Duration `gorm:"type:interval"`

	FileName    string
	Duration    time.Duration `gorm:"type:interval"`
	Resolution  string
	AspectRatio string
	Size        int    // in bytes?
	Format      string // should make this an enum??
	Bitrate     int    // need to choose a unit

	MasterPlaylist string // base64
	// This needs to be an array of the playlists with their resolution listed..
	Resolutions     string // need to come up with a good system for this
	BaseStoragePath string // to pass to file service

}
