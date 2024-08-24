package testUtils

import (
	"regexp"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func CheckElement(t *testing.T, n *html.Node, tag string, expectedContent string) bool {
	hasMatchingHtmlTag := n.Type == html.ElementNode && n.Data == tag
	if hasMatchingHtmlTag {
		var textContent strings.Builder
		extractText(n, &textContent)
		actualContent := cleanTextContent(textContent.String())

		contentDoesNotMatch := expectedContent != "" && actualContent != expectedContent
		if contentDoesNotMatch {
			t.Errorf("Expected <%s> content to be %q, but got %q", tag, expectedContent, actualContent)
			return true
		}
		return true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if CheckElement(t, c, tag, expectedContent) {
			return true
		}
	}

	return false
}

func extractText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, sb)
	}
}

func cleanTextContent(textContent string) string {
	spaceNormalizer := regexp.MustCompile(`\s+`)
	normalizedContent := spaceNormalizer.ReplaceAllString(textContent, " ")
	return strings.TrimSpace(normalizedContent)
}
