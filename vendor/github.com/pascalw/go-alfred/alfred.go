package alfred

import (
	"strings"
	"fmt"
	"encoding/xml"
)

type AlfredResponse struct {
	Items []AlfredResponseItem
	XMLName struct{} `xml:"items"`
}

type AlfredResponseItem struct {
	Valid bool `xml:"valid,attr"`
	Arg string `xml:"arg,attr,omitempty"`
	Uid string `xml:"uid,attr,omitempty"`
	Title string `xml:"title"`
	Subtitle string `xml:"subtitle"`
	Icon string `xml:"icon"`

	XMLName struct{} `xml:"item"`
}

const xmlHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"

func NewResponse() *AlfredResponse {
	return new(AlfredResponse).Init()
}

func (response *AlfredResponse) Init() *AlfredResponse {
	response.Items = []AlfredResponseItem{}
	return response
}

func (response *AlfredResponse) AddItem(item *AlfredResponseItem) {
	response.Items = append(response.Items, *item)
}

func (response *AlfredResponse) Print() {
	var xmlOutput, _ = xml.Marshal(response)
	fmt.Print(xmlHeader, string(xmlOutput))
}

func InitTerms(params []string) {
	for index, term := range params {
		params[index] = strings.ToLower(term)
	}
}

func MatchesTerms(queryTerms []string, itemName string) bool {
	nameLower := strings.ToLower(itemName)

	for _, term := range queryTerms {
		if ! strings.Contains(nameLower, term) { return false }
	}

	return true
}
