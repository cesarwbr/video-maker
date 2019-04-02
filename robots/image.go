package robots

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"video-maker/types"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

// ImageRobot create a new robot
func ImageRobot() {
	content := Load()

	fetchImagesOfAllSentences(content)
	Save(content)
}

func fetchImagesOfAllSentences(content *types.Content) {
	for index, sentence := range content.Sentences {
		query := content.SearchTerm + " " + sentence.Keywords[0]
		content.Sentences[index].Images = fetchGoogleAndReturnImagesLinks(query)
		content.Sentences[index].GoogleSearchQuery = query
	}
}

func fetchGoogleAndReturnImagesLinks(query string) []string {
	apiKey, searchEngineID := getAPIKeySearchEngineID()
	hc := &http.Client{Transport: &transport.APIKey{Key: apiKey}}
	svc, err := customsearch.New(hc)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := svc.Cse.List(query).SearchType("image").ImgSize("huge").Cx(searchEngineID).Num(2).Do()
	if err != nil {
		log.Fatal(err)
	}

	var links []string
	for _, result := range resp.Items {
		links = append(links, result.Link)
	}

	return links
}

func getAPIKeySearchEngineID() (string, string) {
	file, _ := os.Open("credentials/google-search.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	credentials := types.Credentials{}
	err := decoder.Decode(&credentials)
	if err != nil {
		panic(err)
	}

	return credentials.APIKey, credentials.SearchEngineID
}
