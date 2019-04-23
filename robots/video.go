package robots

import (
	"os"
	"os/exec"
	"strconv"
	"video-maker/types"
)

// VideoRobot create a new robot
func VideoRobot() {
	content := Load()

	convertAllImages(content)
	createAllSentenceImages(content)
	createYouTubeThumbnail()
	createAfterEffectsScript(content)
	renderVideoWithAfterEffects()
}

func renderVideoWithAfterEffects() {
	cmd := "/Applications/Adobe After Effects CC 2019/aerender"

	templateFilePath := "/Users/cesaralvarenga/dev/golang/src/video-maker/templates/1/template.aep"
	destinationFilePath := "/Users/cesaralvarenga/dev/golang/src/video-maker/content/output.mov"

	args := []string{
		"-comp", "main",
		"-project", templateFilePath,
		"-output", destinationFilePath}

	cmdExec := exec.Command(cmd, args...)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	err := cmdExec.Run()

	if err != nil {
		panic(err)
	}
}

func createAfterEffectsScript(content *types.Content) {
	SaveScript(content)
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
