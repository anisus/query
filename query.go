package query

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Nodes []*html.Node

type Selector func(*html.Node) bool

// Find gets the descendants of each element in the current set of matched elements, filtered by Selectors.
// After discovering a match, it will not attempt finding matches among the descendants of that node.
func (n Nodes) Find(selectors ...Selector) Nodes {
	matched := Nodes{}

	for _, node := range n {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, false)
		}
	}

	return matched
}

// Find gets the descendants of each element in the current set of matched elements, filtered by Selectors.
// After discovering a match, it will continue to search for matches among the descendants of that node.
func (n Nodes) FindAll(selectors ...Selector) Nodes {
	matched := Nodes{}

	for _, node := range n {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, true)
		}
	}

	return matched
}

// Children gets the children of each element in the set of matched element, optionally filtered by Selectors.
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

// Filter reduces the set of matched elements to those that match the Selectors.
func (n Nodes) Filter(selectors ...Selector) Nodes {
	var matched Nodes

	for _, node := range n {
		if match(node, selectors) {
			matched = append(matched, node)
		}
	}

	return matched
}

// ByTag returns a Selector which matches all nodes of the provided tag type.
func ByTag(a atom.Atom) Selector {
	return func(node *html.Node) bool {
		return node.DataAtom == a
	}
}

// ById returns a Selector which matches all nodes with the provided id.
func ById(id string) Selector {
	return ByAttr("id", id)
}

// ByClass returns a Selector which matches all nodes with the provided class.
func ByClass(class string) Selector {
	return func(node *html.Node) bool {
		cls := strings.Fields(getAttr(node, "class"))
		for _, cl := range cls {
			if cl == class {
				return true
			}
		}
		return false
	}
}

// ByAttr returns a Selector which matches all nodes with the provided attribute equal to the provided value.
func ByAttr(attr, value string) Selector {
	return func(node *html.Node) bool {
		return getAttr(node, attr) == value
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

func getAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
