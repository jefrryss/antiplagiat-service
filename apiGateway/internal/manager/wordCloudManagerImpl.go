package manager

import (
	"net/url"
	"strings"
	"unicode"
)

type WordCloudManagerImpl struct{}

func NewWordCloudManager() *WordCloudManagerImpl {
	return &WordCloudManagerImpl{}
}

func (w *WordCloudManagerImpl) GenerateWordCloud(text string) string {
	words := extractWords(text)

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	wordText := strings.Join(words, " ")

	baseURL := "https://quickchart.io/wordcloud"
	params := url.Values{}
	params.Set("text", wordText)
	params.Set("format", "png")
	params.Set("width", "800")
	params.Set("height", "800")
	params.Set("removeStopwords", "true")
	params.Set("minWordLength", "3")

	return baseURL + "?" + params.Encode()
}

func extractWords(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
