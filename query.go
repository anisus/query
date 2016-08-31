package query

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Nodes []*html.Node

type Selector func(*html.Node) bool

func (n Nodes) Find(selectors ...Selector) Nodes {
	matched := Nodes{}

	for _, node := range n {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, false)
		}
	}

	return matched
}

func (n Nodes) FindAll(selectors ...Selector) Nodes {
	matched := Nodes{}

	for _, node := range n {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, true)
		}
	}

	return matched
}

func (n Nodes) Children(selectors ...Selector) Nodes {
	var matched Nodes

	for _, node := range n {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if match(c, selectors) {
				matched = append(matched, c)
			}
		}
	}

	return matched
}

// ByTag returns a Selector which selects all nodes of the provided tag type.
func ByTag(a atom.Atom) Selector {
	return func(node *html.Node) bool {
		return node.DataAtom == a
	}
}

// ById returns a Selector which selects all nodes with the provided id.
func ById(id string) Selector {
	return func(node *html.Node) bool {
		return attr(node, "id") == id
	}
}

// ByClass returns a Selector which selects all nodes with the provided class.
func ByClass(class string) Selector {
	return func(node *html.Node) bool {
		cls := strings.Fields(attr(node, "class"))
		for _, cl := range cls {
			if cl == class {
				return true
			}
		}
		return false
	}
}

func find(n *Nodes, node *html.Node, selectors []Selector, nested bool) {
	if match(node, selectors) {
		(*n) = append(*n, node)

		if !nested {
			return
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		find(n, c, selectors, nested)
	}
}

func match(node *html.Node, selectors []Selector) bool {
	for _, s := range selectors {
		if s(node) {
			return true
		}
	}
	return false
}

func attr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
