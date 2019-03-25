package robots

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"video-maker/types"
)

// InputRobot define
func InputRobot() {
	content := types.Content{MaximumSentences: 7}

	searchTerm := askAndReturnSearchTerm()

	chosePrefix := false

	for !chosePrefix {
		prefix, err := askAndReturnPrefix()

		if err == nil {
			chosePrefix = true

			if prefix != "Cancel" {
				content.SearchTerm = searchTerm
				content.Prefix = prefix
			}
		}
	}

	Save(&content)
}

func askAndReturnSearchTerm() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Type a Wikipedia search term: ")
	text, _ := reader.ReadString('\n')

	text = strings.TrimSuffix(text, "\n")

	return text
}

func askAndReturnPrefix() (string, error) {
	fmt.Printf("\n")
	prefixes := [...]string{"Cancel", "Who is", "What is", "The history of"}

	var prefix uint
	fmt.Println("[1] Who is")
	fmt.Println("[2] What is")
	fmt.Println("[3] The history of")
	fmt.Printf("[0] Cancel\n\n")
	fmt.Printf("Which type: ")
	fmt.Scanf("%d", &prefix)

	if (prefix > 3) || (prefix < 0) {
		return "", errors.New("Wrong type")
	}

	return prefixes[prefix], nil
}
