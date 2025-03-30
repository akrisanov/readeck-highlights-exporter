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
)

type Highlight struct {
	ID            string    `json:"id"`
	Text          string    `json:"text"`
	Href          string    `json:"href"`
	BookmarkTitle string    `json:"bookmark_title"`
	BookmarkURL   string    `json:"bookmark_url"`
	BookmarkHref  string    `json:"bookmark_href"`
	Created       time.Time `json:"created"`
}

func fetchHighlights(baseURL, apiPath, token string) ([]Highlight, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	parsedBaseURL.Path = path.Join(parsedBaseURL.Path, apiPath)
	apiURL := parsedBaseURL.String()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

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

	var highlights []Highlight
	if err := json.NewDecoder(resp.Body).Decode(&highlights); err != nil {
		return nil, err
	}

	return highlights, nil
}

func exportToCSV(highlights []Highlight, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{"Highlight", "Title", "URL", "Date"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, h := range highlights {
		row := []string{
			h.Text,
			h.BookmarkTitle,
			h.BookmarkURL,
			h.Created.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	baseURL := os.Getenv("READECK_API_BASE_URL")
	apiToken := os.Getenv("READECK_API_KEY")
	outputPath := os.Getenv("CSV_OUTPUT_PATH")

	if baseURL == "" || apiToken == "" || outputPath == "" {
		fmt.Println("Missing required environment variables")
		return
	}

	highlights, err := fetchHighlights(baseURL, "/bookmarks/annotations", apiToken)
	if err != nil {
		fmt.Println("Error fetching highlights:", err)
		return
	}

	if err := exportToCSV(highlights, outputPath); err != nil {
		fmt.Println("Error exporting highlights to CSV:", err)
		return
	}

	fmt.Println("Highlights exported successfully to", outputPath)
}
