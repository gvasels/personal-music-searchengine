package metadata

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"

	"github.com/gvasels/personal-music-searchengine/internal/models"
)

// Extractor extracts metadata from audio files
type Extractor struct{}

// NewExtractor creates a new metadata extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts metadata from an audio file
func (e *Extractor) Extract(reader io.ReadSeeker, filename string) (*models.UploadMetadata, error) {
	// Try to read metadata using tag library
	m, err := tag.ReadFrom(reader)
	if err != nil {
		// If we can't read tags, return metadata based on filename
		return e.metadataFromFilename(filename), nil
	}

	// Calculate duration for MP3 files
	duration := 0
	if m.FileType() == tag.MP3 {
		if _, err := reader.Seek(0, io.SeekStart); err == nil {
			duration = e.calculateMP3Duration(reader)
		}
	}

	metadata := &models.UploadMetadata{
		Title:       e.parseTitle(m.Title(), filename),
		Artist:      e.parseArtist(m.Artist()),
		AlbumArtist: m.AlbumArtist(),
		Album:       m.Album(),
		Genre:       m.Genre(),
		Year:        m.Year(),
		Duration:    duration,
		Format:      e.formatToString(m.FileType()),
		HasCoverArt: m.Picture() != nil,
		Composer:    m.Composer(),
		Comment:     m.Comment(),
	}

	// Extract track and disc numbers
	track, _ := m.Track()
	disc, _ := m.Disc()
	metadata.TrackNumber = track
	metadata.DiscNumber = disc

	// Try to get additional metadata from raw tags
	if raw := m.Raw(); raw != nil {
		if lyrics, ok := raw["lyrics"].(string); ok {
			metadata.Lyrics = lyrics
		}
		if bitrate, ok := raw["bitrate"].(int); ok {
			metadata.Bitrate = bitrate
		}
	}

	return metadata, nil
}

// ExtractCoverArt extracts embedded cover art from an audio file
func (e *Extractor) ExtractCoverArt(reader io.ReadSeeker) ([]byte, string, error) {
	// Reset reader to beginning
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, "", fmt.Errorf("failed to seek: %w", err)
	}

	m, err := tag.ReadFrom(reader)
	if err != nil {
		return nil, "", nil // No error, just no cover art
	}

	picture := m.Picture()
	if picture == nil {
		return nil, "", nil
	}

	return picture.Data, picture.MIMEType, nil
}

// DetectFormat detects the audio format from a reader
func (e *Extractor) DetectFormat(reader io.ReadSeeker) (models.AudioFormat, error) {
	// Reset reader to beginning
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to seek: %w", err)
	}

	m, err := tag.ReadFrom(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read tags: %w", err)
	}

	return e.fileTypeToAudioFormat(m.FileType()), nil
}

// parseTitle returns the title from metadata or falls back to filename
func (e *Extractor) parseTitle(title, filename string) string {
	if title != "" {
		return title
	}
	// Use filename without extension as title
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// parseArtist returns the artist from metadata or "Unknown Artist"
func (e *Extractor) parseArtist(artist string) string {
	if artist != "" {
		return artist
	}
	return "Unknown Artist"
}

// formatToString converts tag.FileType to string
func (e *Extractor) formatToString(ft tag.FileType) string {
	switch ft {
	case tag.MP3:
		return "MP3"
	case tag.FLAC:
		return "FLAC"
	case tag.OGG:
		return "OGG"
	case tag.M4A, tag.M4B, tag.M4P, tag.ALAC:
		return "AAC"
	default:
		return "UNKNOWN"
	}
}

// fileTypeToAudioFormat converts tag.FileType to models.AudioFormat
func (e *Extractor) fileTypeToAudioFormat(ft tag.FileType) models.AudioFormat {
	switch ft {
	case tag.MP3:
		return models.AudioFormatMP3
	case tag.FLAC:
		return models.AudioFormatFLAC
	case tag.OGG:
		return models.AudioFormatOGG
	case tag.M4A, tag.M4B, tag.M4P, tag.ALAC:
		return models.AudioFormatAAC
	default:
		return models.AudioFormatMP3 // Default to MP3
	}
}

// metadataFromFilename creates metadata based only on filename
func (e *Extractor) metadataFromFilename(filename string) *models.UploadMetadata {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	title := strings.TrimSuffix(base, ext)

	// Try to parse format from extension
	format := e.extensionToFormat(ext)

	return &models.UploadMetadata{
		Title:       title,
		Artist:      "Unknown Artist",
		Format:      string(format),
		HasCoverArt: false,
	}
}

// extensionToFormat converts file extension to format string
func (e *Extractor) extensionToFormat(ext string) models.AudioFormat {
	switch strings.ToLower(ext) {
	case ".mp3":
		return models.AudioFormatMP3
	case ".flac":
		return models.AudioFormatFLAC
	case ".wav":
		return models.AudioFormatWAV
	case ".m4a", ".aac":
		return models.AudioFormatAAC
	case ".ogg":
		return models.AudioFormatOGG
	default:
		return models.AudioFormatMP3
	}
}

// calculateMP3Duration calculates the duration of an MP3 file by parsing frames
func (e *Extractor) calculateMP3Duration(reader io.ReadSeeker) int {
	decoder := mp3.NewDecoder(reader)
	var totalDuration time.Duration
	var frame mp3.Frame
	skipped := 0

	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			break
		}
		totalDuration += frame.Duration()
	}

	return int(totalDuration.Seconds())
}
