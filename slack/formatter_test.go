package slack_test

import (
	"errors"
	"testing"
	"vim-tricks-slackbot/slack"
)

func TestFormatHTMLToMarkdownConvertsHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  error
		converter      mockConverter
	}{
		{
			name:           "String contains normal elements",
			expectedOutput: "Hello there",
			input:          "<p>Hello there</p>",
			expectedError:  nil,
			converter:      mockConverter{returnFromConverter: "Hello there"},
		},
		{
			name:           "String contains links",
			expectedOutput: "<https://example.com|Examples>",
			input:          "<a href=\"https://example.com\">Examples</a>",
			expectedError:  nil,
			converter:      mockConverter{returnFromConverter: "[Examples](https://example.com)"},
		},
		{
			name:           "String contains special character",
			expectedOutput: "Insert special • characters ◆",
			input:          "<h1>Insert special • characters ◆<h1>",
			expectedError:  nil,
			converter:      mockConverter{returnFromConverter: "Insert special • characters ◆"},
		},
		{
			name:           "String contains special character in link with nested elements",
			expectedOutput: "The post <https://vimtricks.com/p/insert-special-characters|Insert special • characters ◆> appeared first on <https://vimtricks.com|VimTricks>",
			input:          "<p>The post <a rel=\"nofollow\" href=\"https://vimtricks.com/p/insert-special-characters/\" data-wpel-link=\"internal\">Insert special • characters ◆</a> appeared first on <a rel=\"nofollow\" href=\"https://vimtricks.com\" data-wpel-link=\"internal\">VimTricks</a>.</p>",
			expectedError:  nil,
			converter:      mockConverter{returnFromConverter: "The post [Insert special • characters ◆](https://vimtricks.com/p/insert-special-characters) appeared first on [VimTricks](https://vimtricks.com)"},
		},
	}
	for _, test := range tests {
		fmter := slack.NewMarkdownFormatter(test.converter)
		actual, err := fmter.FormatHTMLToMarkdown(test.input)
		if err != test.expectedError {
			t.Errorf("%s:\nExpected %q error,\nGot %q error\n", test.name, test.expectedError, err)
		}

		if actual != test.expectedOutput {
			t.Errorf("%s:\nExpected %q,\nGot %q\n", test.name, test.expectedOutput, actual)
		}
	}
}

func TestFormatHTMLToMarkdownHandlesErrors(t *testing.T) {
	converter := mockConverter{errOnConversion: true, returnFromConverter: ""}
	fmter := slack.NewMarkdownFormatter(converter)
	_, err := fmter.FormatHTMLToMarkdown("<p>Hello there</p>")
	if err == nil {
		t.Error("Expected an error,Got nil")
	}
}

type mockConverter struct {
	errOnConversion     bool
	returnFromConverter string
}

func (m mockConverter) ConvertString(s string) (string, error) {
	if m.errOnConversion {
		return "", errors.New("Failed to convert")
	}

	return m.returnFromConverter, nil
}
