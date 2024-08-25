package testUtils

import "golang.org/x/net/html"

func FindElementByID(node *html.Node, id string) *html.Node {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == id {
				return node
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		found := FindElementByID(child, id)
		if found != nil {
			return found
		}
	}
	return nil
}
