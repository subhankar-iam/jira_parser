package ai_parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"jira-parser/constants"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

var gemini_url string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	gemini_api_key := os.Getenv("GEMINI_API_KEY")
	gemini_url = fmt.Sprintf(constants.Gemini_url, gemini_api_key)
}

func Ask_Gemini(prompt string) (string, error) {
	request_body_json := fmt.Sprintf(constants.Gemini_Request, prompt)
	requestBody := bytes.NewBuffer([]byte(request_body_json))
	resp, err := http.Post(gemini_url, "application/json", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return "", err
	}

	var response_json string
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		response_json = response.Candidates[0].Content.Parts[0].Text
	} else {
		fmt.Println("No text found")
	}
	return response_json, nil
}
