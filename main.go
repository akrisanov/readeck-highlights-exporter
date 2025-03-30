package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

// Highlight represents a single Readeck bookmark's highlight
type Highlight struct {
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	Href          string    `json:"href"`
	BookmarkTitle string    `json:"bookmark_title"`
	BookmarkURL   string    `json:"bookmark_url"`
	BookmarkHref  string    `json:"bookmark_href"`
	Created       time.Time `json:"created"`
}

// Helper to load required environment variables
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		slog.Error("missing required environment variable", slog.String("key", key))
		panic(fmt.Errorf("missing required environment variable: %s", key))
	}
	return value
}

// General error handler to reduce repetition
func check(err error, context string) {
	if err != nil {
		slog.Error(context, slog.String("error", err.Error()))
		panic(fmt.Errorf("%s: %w", context, err))
	}
}

// Fetch highlights from the Readeck API
func fetchHighlights(baseURL, apiPath, token string) (highlights []Highlight, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("fetchHighlights failed: %w", err)
		}
	}()

	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	parsedBaseURL.Path = path.Join(parsedBaseURL.Path, apiPath)
	apiURL := parsedBaseURL.String()

	slog.Info("Fetching highlights", slog.String("url", apiURL))

	req, err := http.NewRequest("GET", apiURL, nil)
	check(err, "Failed to create request")

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch highlights: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&highlights)
	return
}

// Export highlights to a CSV file
func exportToCSV(highlights []Highlight, outputPath string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("exportToCSV failed: %w", err)
		}
	}()

	slog.Info("Exporting highlights to CSV", slog.String("outputPath", outputPath))

	file, err := os.Create(outputPath)
	check(err, "Failed to create CSV file")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"Highlight", "Title", "URL", "Date"}
	check(writer.Write(header), "Failed to write CSV header")

	for _, h := range highlights {
		row := []string{
			h.Text,
			h.BookmarkTitle,
			h.BookmarkURL,
			h.Created.Format("2006-01-02 15:04:05"),
		}
		check(writer.Write(row), "Failed to write CSV row")
	}

	slog.Info("Exporting completed successfully 🎉")
	return nil
}

func main() {
	// Initialize logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file", slog.String("error", err.Error()))
	}

	baseURL := mustEnv("READECK_API_BASE_URL")
	apiToken := mustEnv("READECK_API_KEY")
	outputPath := mustEnv("CSV_OUTPUT_PATH")

	// Export highlights to CSV
	highlights, err := fetchHighlights(baseURL, "/bookmarks/annotations", apiToken)
	check(err, "Error fetching highlights")

	check(exportToCSV(highlights, outputPath), "Error exporting highlights to CSV")
}
