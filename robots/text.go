package robots

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"video-maker/types"

	algorithmia "github.com/algorithmiaio/algorithmia-go"
	"github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

// TextRobot create a new robot
func TextRobot() {
	content := Load()

	fetchContentFromWikipedia(content)
	sanitizeContent(content)
	breakContentIntoSentences(content)
	limitMaximumSentences(content)
	fetchKeywordsOfAllSentences(content)

	Save(content)
}

func fetchKeywordsOfAllSentences(content *types.Content) {
	for index, sentence := range content.Sentences {
		content.Sentences[index].Keywords = fetchWatsonAndReturnKeywords(sentence.Text)
	}
}

func limitMaximumSentences(content *types.Content) {
	content.Sentences = content.Sentences[0:content.MaximumSentences]
}

func fetchWatsonAndReturnKeywords(sentence string) []string {
	service, serviceErr := naturallanguageunderstandingv1.NewNaturalLanguageUnderstandingV1(&naturallanguageunderstandingv1.NaturalLanguageUnderstandingV1Options{
		URL:       "https://gateway.watsonplatform.net/natural-language-understanding/api",
		Version:   "2018-03-16",
		IAMApiKey: getWatsonAPIKey(),
	})

	if serviceErr != nil {
		panic(serviceErr)
	}

	response, responseErr := service.Analyze(
		&naturallanguageunderstandingv1.AnalyzeOptions{
			Text: &sentence,
			Features: &naturallanguageunderstandingv1.Features{
				Keywords: &naturallanguageunderstandingv1.KeywordsOptions{},
			},
		},
	)

	if responseErr != nil {
		panic(responseErr)
	}

	analyze := service.GetAnalyzeResult(response)

	var keywords []string
	for _, keyword := range analyze.Keywords {
		keywords = append(keywords, *keyword.Text)
	}

	return keywords
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
	file, _ := os.Open("credentials/algorithmia.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	credentials := types.Credentials{}
	err := decoder.Decode(&credentials)
	if err != nil {
		panic(err)
	}

	return credentials.APIKey
}

func getWatsonAPIKey() string {
	file, _ := os.Open("credentials/watson-nlu.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	credentials := types.Credentials{}
	err := decoder.Decode(&credentials)
	if err != nil {
		panic(err)
	}

	return credentials.APIKey
}
