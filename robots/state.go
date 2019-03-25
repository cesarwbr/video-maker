package robots

import (
	"encoding/json"
	"io/ioutil"
	"video-maker/types"
)

var contentFilePath = "content.json"

// Save persist the content struct to the filesystem
func Save(content *types.Content) {
	file, _ := json.MarshalIndent(content, "", "")

	ioutil.WriteFile(contentFilePath, file, 0644)
}

// Load load content struct from the filesystem
func Load() *types.Content {
	file, _ := ioutil.ReadFile(contentFilePath)

	content := types.Content{}

	json.Unmarshal([]byte(file), &content)

	return &content
}
