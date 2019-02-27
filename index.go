package main

import (
	"bufio"
	"fmt"
	"os"
)

type Content struct {
	searchTerm string
	prefix     string
}

func main() {
	var searchTerm = askAndReturnSearchTerm()
	var prefix = askAndReturnPrefix()

	content := Content{searchTerm, prefix}

	fmt.Println(content)
}

func askAndReturnSearchTerm() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a Wikipedia search term:")
	text, _ := reader.ReadString('\n')

	return text
}

func askAndReturnPrefix() string {
	var prefixes = [3]string{"Who is", "What is", "The history of"}
	var prefix int
	fmt.Println("[1] Who is")
	fmt.Println("[2] What is")
	fmt.Println("[3] The history of")
	fmt.Println("[0] Cancel")
	fmt.Scanf("%d", &prefix)

	if (prefix > 3) || (prefix < 1) {
		return "Cancel"
	}

	return prefixes[prefix-1]
}
