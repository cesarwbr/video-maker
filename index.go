package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Content store the search preferences
type Content struct {
	searchTerm string
	prefix     string
}

func main() {
	var searchTerm = askAndReturnSearchTerm()
	var prefix = askAndReturnPrefix()

	if prefix != "Cancel" {
		content := Content{searchTerm, prefix}
		fmt.Printf("\n{searchTerm: %s, prefix: %s}", content.searchTerm, content.prefix)
	}
}

func askAndReturnSearchTerm() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Type a Wikipedia search term: ")
	text, _ := reader.ReadString('\n')

	text = strings.TrimSuffix(text, "\n")

	return text
}

func askAndReturnPrefix() string {
	fmt.Printf("\n")
	var prefixes = [3]string{"Who is", "What is", "The history of"}
	var prefix int
	fmt.Println("[1] Who is")
	fmt.Println("[2] What is")
	fmt.Println("[3] The history of")
	fmt.Printf("[0] Cancel\n\n")
	fmt.Printf("Which type: ")
	fmt.Scanf("%d", &prefix)

	if (prefix > 3) || (prefix < 1) {
		return "Cancel"
	}

	return prefixes[prefix-1]
}
