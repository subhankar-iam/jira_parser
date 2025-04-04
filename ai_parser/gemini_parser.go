package ai_parser

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	constants "jira-parser/constants"
	"log"
	"net/http"
	"strings"
)

type GeminiParser struct {
	FileSaveLocation FileSaveLocation `json:"file_save_location"`
	Background       Background       `json:"background"`
	Scenarios        []Scenario       `json:"scenarios"`
}

type Background struct {
	BaseURL                string   `json:"base_url"`
	Header                 string   `json:"header"`
	Annotations            []string `json:"annotations"`
	Config                 []string `json:"config"`
	AdditionalInstructions []string `json:"additional_instructions"`
}

type Scenario struct {
	EndPoint               string   `json:"endPoint"`
	Request                string   `json:"request"`
	Method                 string   `json:"method"`
	StatusCode             int      `json:"statuscode"`
	AdditionalInstructions []string `json:"additional_instructions"`
}

type FileSaveLocation struct {
	RequestFileSaveLocation string `json:"requestfileSaveLocation"`
	HeaderFileSaveLocation  string `json:"headerFileSaveLocation"`
}

var outputFormat GeminiParser

func init() {
	outputFormat = GeminiParser{
		FileSaveLocation: FileSaveLocation{
			RequestFileSaveLocation: "location at which the request file would be saved. default: bdd/test/resources/features/headers/ if the header file exists", //dont assume anything by yourself if its null then let it be null
			HeaderFileSaveLocation:  "location at which the header file would be saved. default: bdd/test/resources/features/",                                    //dont assume anything by yourself if its null then let it be null
		},
		Background: Background{
			BaseURL: "the base url from goes here",
			Header: "the header from goes here it might be header file location or header in json format" +
				"if only the file name is given the header_file_save_path you add and get the full path " +
				"for example if headerFileSaveLocation is feature/headers  " +
				"and headerfile name is headers.json the answer would be feature/headers/headers.json",
			Config:                 []string{"any additional configuration like creating mongodb connection or creating kafka connection goes here"},
			Annotations:            []string{"any potential annotations goes here, look for labels input"},
			AdditionalInstructions: []string{"list of any other misc instruction"},
		},
		Scenarios: []Scenario{
			{
				EndPoint: "the end point goes here",
				Request: `the request  body goes here it might be request file location or request body in json format ` +
					`if its a file and only the file name is given the request_file_save_path you add and get the full path like outputFormat.RequestFileSaveLocation, outputFormat.Scenarios[0].Request. for example if requestFileSaveLocation is feature/headers and headerfile name is response.json the answer would be feature/headers/response.json. please use the full filepath`,
				Method:                 "the request method POST,GET,PUT,DELETE",
				StatusCode:             http.StatusOK, //the status code goes here
				AdditionalInstructions: []string{"Any other steps needed to be done can be mentioned here"},
			},
		},
	}
}

func ParseGemini(json_input string, file_content map[string]string) (string, error) {
	fileNames := func(m map[string]string) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys
	}(file_content)

	output_json_format, err := json2.Marshal(outputFormat)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}
	gemini_prompt := fmt.Sprintf(constants.Parser_prompt, json_input, "["+strings.Join(fileNames, ",")+"]", string(output_json_format))
	gemini_prompt = strings.ReplaceAll(gemini_prompt, `"`, `''`)
	json_resp, err := Ask_Gemini(gemini_prompt)
	if err != nil {
		return "", err
	}
	var header_file string
	var req_file string
	var parser map[string]interface{}
	json2.Unmarshal([]byte(json_resp), &parser)
	//request_file_path := parser.FileSaveLocation.RequestFileSaveLocation
	//request_header_path := parser.FileSaveLocation.HeaderFileSaveLocation
	request_file_path, ok := parser["file_save_location"].(map[string]interface{})
	if !ok {
		return "", errors.New("requestfileSaveLocation could not be parsed")
	} else {
		if val, ok := request_file_path["requestfileSaveLocation"].(string); ok {
			req_file = val
		} else {
			req_file = "/features/requests/"
		}
	}

	request_header_path, ok := parser["file_save_location"].(map[string]interface{})
	if !ok {
		return "", errors.New("requestfileSaveLocation could not be parsed")
	} else {
		if val, ok := request_header_path["headerFileSaveLocation"].(string); ok {
			header_file = val
		} else {
			header_file = "/features/headers/"
		}
	}

	for k, v := range file_content {
		var err error = nil
		if strings.Contains(k, "header") {
			err = SaveInDrive(header_file, k, v)
		} else {
			err = SaveInDrive(req_file, k, v)
		}
		if err != nil {
			return "", err
		}

	}
	delete(parser, "file_save_location")
	final_json, err := json2.Marshal(parser)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}
	return string(final_json), nil
}
