package main

import (
	"encoding/json"
	"fmt"
	"time"
	"video-maker/robots"

	"github.com/briandowns/spinner"
)

func main() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	robots.InputRobot()

	s.Start()
	robots.TextRobot()
	s.Stop()

	content := robots.Load()
	fmt.Println(prettyPrint(content))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
