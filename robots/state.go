package robots

import (
	"encoding/json"
	"io/ioutil"
	"video-maker/types"
)

var contentFilePath = "content.json"
var scriptFilePath = "content/after-effects-script.js"

// Save persist the content struct to the filesystem
func Save(content *types.Content) {
	file, _ := json.MarshalIndent(content, "", "")

	ioutil.WriteFile(contentFilePath, file, 0644)
}

// SaveScript save after effects script
func SaveScript(content *types.Content) {
	contentJSON, err := json.Marshal(content)

	if err != nil {
		panic(err)
	}

	contentString := string(contentJSON)
	scriptFile := []byte("var content = " + contentString)

	ioutil.WriteFile(scriptFilePath, scriptFile, 0644)
}

// Load load content struct from the filesystem
func Load() *types.Content {
	file, _ := ioutil.ReadFile(contentFilePath)

	content := types.Content{}

	json.Unmarshal([]byte(file), &content)

	return &content
}
