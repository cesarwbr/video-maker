package robots

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"video-maker/types"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

// ImageRobot create a new robot
func ImageRobot() {
	content := Load()

	fetchImagesOfAllSentences(content)
	downloadAllImages(content)
	convertAllImages(content)
	createAllSentenceImages(content)
	createYouTubeThumbnail()

	Save(content)
}

func createYouTubeThumbnail() {
	cmd := "convert"
	args := []string{
		"./content/0-converted.png",
		"./content/youtube-thumbnail.jpg"}

	err := exec.Command(cmd, args...).Run()

	if err != nil {
		panic(err)
	}
}

func createAllSentenceImages(content *types.Content) {
	for sentenceIndex, sentence := range content.Sentences {
		createSentenceImage(sentenceIndex, sentence.Text)
	}
}

type templateSetting struct {
	size    string
	gravity string
}

func createSentenceImage(sentenceIndex int, sentenceText string) {
	outputFile := "./content/" + strconv.Itoa(sentenceIndex) + "-sentence.png"

	templateSettings := []templateSetting{
		{"1920x400", "center"},
		{"1920x1080", "center"},
		{"800x1080", "west"},
		{"1920x400", "center"},
		{"1920x1080", "center"},
		{"800x1080", "west"},
		{"1920x400", "center"}}

	cmd := "convert"
	args := []string{
		"-size", templateSettings[sentenceIndex].size,
		"-gravity", templateSettings[sentenceIndex].gravity,
		"-background", "transparent",
		"-fill", "white",
		"-kerning", "-1",
		"caption:" + sentenceText,
		outputFile}

	err := exec.Command(cmd, args...).Run()

	if err != nil {
		panic(err)
	}
}

func convertAllImages(content *types.Content) {
	for sentenceIndex := range content.Sentences {
		convertImage(sentenceIndex)
	}
}

func convertImage(sentenceIndex int) {
	inputFile := "./content/" + strconv.Itoa(sentenceIndex) + "-original.png[0]"
	outputFile := "./content/" + strconv.Itoa(sentenceIndex) + "-converted.png"

	newSize := "1920x1080"

	cmd := "convert"
	args := []string{
		inputFile,
		"(", "-clone", "0", "-background", "white", "-blur", "0x9", "-resize", newSize + "^", ")",
		"(", "-clone", "0", "-background", "white", "-resize", newSize, ")",
		"-delete", "0", "-gravity", "center", "-compose", "over", "-composite", "-extent", newSize,
		outputFile}

	err := exec.Command(cmd, args...).Run()

	if err != nil {
		panic(err)
	}
}

func downloadAllImages(content *types.Content) {
	var downloadedImages []string
	for index, sentence := range content.Sentences {
		images := sentence.Images

		for _, image := range images {
			if contains(downloadedImages, image) {
				fmt.Println("imagem já foi baixada")
				continue
			}

			_, err := downloadAndSave(image, strconv.Itoa(index)+"-original.png")

			if err != nil {
				fmt.Println("Error ao baixar " + image)
				continue
			}

			downloadedImages = append(downloadedImages, image)
			fmt.Println("Baixou imagem com sucesso " + image)
			break
		}
	}
}

func contains(items []string, str string) bool {
	for _, item := range items {
		if item == str {
			return true
		}
	}

	return false
}

func downloadAndSave(url string, fileName string) (string, error) {
	response, e := http.Get(url)

	if e != nil {
		return "", e
	}

	defer response.Body.Close()

	finalFileName := "./content/" + fileName

	file, err := os.Create(finalFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return finalFileName, nil
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
