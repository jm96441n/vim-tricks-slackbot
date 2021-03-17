package slack

import (
	"regexp"
)

type MarkdownFormatter struct {
	converter markdownConverter
}

type markdownConverter interface {
	ConvertString(string) (string, error)
}

var pattern = regexp.MustCompile(`\[(?P<title>([^\]])*)\]\((?P<url>https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*))\)`)
var replacer = regexp.MustCompile(`\[[^\]]*\]\(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)\)`)
var template = []byte("<$url|$title>")

func NewMarkdownFormatter(converter markdownConverter) MarkdownFormatter {
	return MarkdownFormatter{converter: converter}
}

func (m MarkdownFormatter) FormatHTMLToMarkdown(text string) (string, error) {
	markdown, err := m.toMarkdown(text)
	if err != nil {
		return "", err
	}
	markdown = m.formatLinks(markdown)

	return markdown, nil
}

func (m MarkdownFormatter) toMarkdown(text string) (string, error) {

	markdown, err := m.converter.ConvertString(text)
	if err != nil {
		return "", err
	}
	return markdown, nil
}

func (m MarkdownFormatter) formatLinks(text string) string {
	results := make([]string, 0)
	for _, submatches := range pattern.FindAllSubmatchIndex([]byte(text), -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		res := []byte{}
		results = append(results, string(pattern.Expand(res, template, []byte(text), submatches)))
	}
	i := 0
	replacerFn := func(_ string) string {
		res := results[i]
		i++
		return res
	}

	return replacer.ReplaceAllStringFunc(text, replacerFn)
}
