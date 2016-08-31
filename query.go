package query

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Set []*html.Node

type Selector func(*html.Node) bool

// Find gets the descending ElementNodes of each element in the Set, filtered by Selectors.
func (s Set) Find(selectors ...Selector) Set {
	matched := Set{}

	for _, node := range s {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, true)
		}
	}

	return matched
}

// FindShallow gets the descending ElementNodes of each element in the Set, filtered by Selectors.
// After discovering a match, it will not attempt finding matches among the descendants of that node.
func (s Set) FindShallow(selectors ...Selector) Set {
	matched := Set{}

	for _, node := range s {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			find(&matched, c, selectors, false)
		}
	}

	return matched
}

// Children gets the child ElementNodes of each element in the Set, filtered by Selectors.
func (s Set) Children(selectors ...Selector) Set {
	var matched Set

	for _, node := range s {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && match(c, selectors) {
				matched = append(matched, c)
			}
		}
	}

	return matched
}

// Contents gets the children of each element in the Set, including TextNodes, CommentNodes, and DoctypeNodes, filtered by Selectors.
func (s Set) Contents(selectors ...Selector) Set {
	var matched Set

	for _, node := range s {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if match(c, selectors) {
				matched = append(matched, c)
			}
		}
	}

	return matched
}

// Filter reduces the set of matched elements to those that match the Selectors.
func (s Set) Filter(selectors ...Selector) Set {
	var matched Set

	for _, node := range s {
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

// Attr returns the value of an attribute for the first element in the set of matched elements.
// If the set contains no elements, or the first misses the attribute, an empty string is returned.
func (s Set) Attr(key string) string {
	if len(s) == 0 {
		return ""
	}

	for _, a := range s[0].Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// Text returns the combined text contents of each element in the set of matched elements, including their descendants.
func (s Set) Text() string {
	str := make([]string, 0, len(s))

	for _, node := range s {
		text(&str, node)
	}

	return strings.Join(str, "")
}

// text traverses the node tree and adds any TextNode's Data to the provided string slice.
func text(str *[]string, n *html.Node) {
	if n.Type == html.TextNode {
		(*str) = append(*str, n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text(str, c)
	}
}

// find traverses the node tree and append any matching descending ElementNode to the Set.
// If nested is true, the function will continue to search subnodes on matched nodes.
func find(s *Set, n *html.Node, selectors []Selector, nested bool) {
	if n.Type == html.ElementNode && match(n, selectors) {
		appendNode(s, n)

		if !nested {
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		find(s, c, selectors, nested)
	}
}

// appendNode adds the node to the Set unless it already exists.
func appendNode(s *Set, n *html.Node) {
	for _, node := range *s {
		if node == n {
			return
		}
	}
	(*s) = append(*s, n)
}

// match returns true if the node matches all of the Selectors or if the slice has 0 Selectors.
func match(node *html.Node, selectors []Selector) bool {
	if len(selectors) == 0 {
		return true
	}

	for _, s := range selectors {
		if s(node) {
			return true
		}
	}
	return false
}

// getAttr returns the provided attribute of the node, or an empty
// string if the attribute was missing.
func getAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}
