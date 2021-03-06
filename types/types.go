package types

// Sentence define a sentence
type Sentence struct {
	Text              string
	Keywords          []string
	Images            []string
	GoogleSearchQuery string
}

// Content store the search preferences
type Content struct {
	SearchTerm             string
	Prefix                 string
	SourceContentOriginal  string
	SourceContentSanitized string
	Sentences              []Sentence
	MaximumSentences       int
}

// Credentials store the app credentials
type Credentials struct {
	APIKey         string
	SearchEngineID string
}
