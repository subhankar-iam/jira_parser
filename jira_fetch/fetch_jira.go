package jira_fetch

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"io"
	"jira-parser/ai_parser"
	parser "jira-parser/description_parser"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var JIRA_SERVER string
var JIRA_USERNAME string
var JIRA_PASSWORD string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	JIRA_SERVER = os.Getenv("JIRA_SERVER")
	JIRA_USERNAME = os.Getenv("JIRA_USER")
	JIRA_PASSWORD = os.Getenv("JIRA_API_KEY")

}

func downloadAttatchment(attatchment *jira.Attachment, file *os.File) {
	req, err := http.NewRequest("GET", attatchment.Content, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(JIRA_USERNAME, JIRA_PASSWORD)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			ForceAttemptHTTP2: false, // Disable HTTP/2
		},
		Timeout: 300 * time.Second, // Longer timeout for large files
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Errorf("failed to read response body: %w", err)
		}
		if n == 0 {
			break
		}
		if _, writeErr := file.Write(buf[:n]); writeErr != nil {
			fmt.Errorf("failed to write to file: %w", writeErr)
		}
	}

	defer resp.Body.Close()
}

func Fetch(ticket_id string) (string, error) {
	tp := jira.BasicAuthTransport{
		Username: JIRA_USERNAME,
		Password: JIRA_PASSWORD,
	}
	client, err := jira.NewClient(tp.Client(), JIRA_SERVER)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}
	issue, _, err := client.Issue.Get(ticket_id, nil)
	if err != nil {
		log.Fatalf("Failed to fetch Jira issue: %v", err)
	}

	description := issue.Fields.Description
	labels := issue.Fields.Labels

	fmt.Printf("Ticket ID: %s\n", issue.Key)
	fmt.Printf("Summary: %s\n", issue.Fields.Summary)
	fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
	fmt.Printf("Created at: %s\n", issue.Fields.Created)
	fmt.Printf("Description: %s\n", description)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
	fmt.Printf("Lables: %s\n", labels)

	var fileName []string
	if issue.Fields.Attachments != nil {
		for _, attachment := range issue.Fields.Attachments {
			file_name := attachment.Filename
			fmt.Printf("attachment downloadig %s\n", file_name)
			filePath := filepath.Join("attachments", file_name)
			fileName = append(fileName, file_name)
			err := os.MkdirAll(filepath.Dir(filePath), 0777)
			if err != nil {
				fmt.Printf("Failed to create directory: %v", err)
			}
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Failed to create file: %v", err)
			}
			defer file.Close()
			downloadAttatchment(attachment, file)

		}
	}
	parserd_description := parser.ParseDescription(description, labels)
	json_data, err := json.Marshal(parserd_description)
	return ai_parser.ParseGemini(string(json_data), fileName)
}
