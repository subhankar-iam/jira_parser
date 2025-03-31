package description_parser

import "strings"

type Description struct {
	Background string   `json:"background"`
	Scenario   []string `json:"scenario"`
	Labels     []string `json:"labels"`
}

func ParseDescription(descriptionString string, labels []string) *Description {
	parts := strings.Split(strings.ToLower(descriptionString), "scenario")
	background := parts[0]
	var scenarios []string
	for _, scenario := range parts[1:] {
		scenarios = append(scenarios, scenario)
	}

	return &Description{background, scenarios, labels}

}
