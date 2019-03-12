package robots

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"video-maker/types"

	algorithmia "github.com/algorithmiaio/algorithmia-go"
	"gopkg.in/neurosnap/sentences.v1/english"
)

// TextRobot create a new robot
func TextRobot(content *types.Content) {
	fetchContentFromWikipedia(content)
	sanitizeContent(content)
	breakContentIntoSentences(content)
}

func fetchContentFromWikipedia(content *types.Content) {
	client := algorithmia.NewClient(getAPIKey(), "")
	algo, _ := client.Algo("web/WikipediaParser/0.1.2?timeout=300")
	resp, _ := algo.Pipe(content.SearchTerm)
	response := resp.(*algorithmia.AlgoResponse)

	content.SourceContentOriginal = response.Result.(map[string]interface{})["summary"].(string)
}

func sanitizeContent(content *types.Content) {
	reg := regexp.MustCompile(`\((.*?)\)`)

	withoutDatesInParentheses := reg.ReplaceAllString(content.SourceContentOriginal, "")
	withoutDatesInParentheses = strings.Replace(withoutDatesInParentheses, "  ", " ", -1)

	content.SourceContentSanitized = withoutDatesInParentheses
}

func breakContentIntoSentences(content *types.Content) {
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	sentences := tokenizer.Tokenize(content.SourceContentSanitized)

	var arrSentences []types.Sentence
	for _, s := range sentences {
		sentence := types.Sentence{Text: s.Text, Keywords: []string{}, Images: []string{}}
		arrSentences = append(arrSentences, sentence)
	}

	content.Sentences = arrSentences
}

func getAPIKey() string {
	file, _ := os.Open("credentials/credentials.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	credentials := types.Credentials{}
	err := decoder.Decode(&credentials)
	if err != nil {
		panic(err)
	}

	return credentials.APIKey
}
