package query

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const htmlData = `
<html>
	<body>
		<div class="container">
			<li class="list">
				<ul><strong>First</strong> value</ul>
				<ul>Second value</ul>
			</li>
		</div>
		<span class="foo">Span
			<div class="foo">Div
				<div class="bar">Bar</div>
			</div>
		</span>
	</body>
</html>
`

type HtmlNodes struct {
	Doc                 *html.Node
	Html                *html.Node
	Head                *html.Node
	Body                *html.Node
	Text1               *html.Node
	Div                 *html.Node
	DivText1            *html.Node
	DivLi               *html.Node
	DivLiText1          *html.Node
	DivLiUl1            *html.Node
	DivLiUl1Strong      *html.Node
	DivLiUl1StrongText1 *html.Node
	DivLiUl1Text1       *html.Node
	DivLiText2          *html.Node
	DivLiUl2            *html.Node
	DivLiUl2Text1       *html.Node
	DivLiText3          *html.Node
	DivText2            *html.Node
	Text2               *html.Node
	Span                *html.Node
	SpanText1           *html.Node
	SpanDiv             *html.Node
	SpanDivText1        *html.Node
	SpanDivDiv          *html.Node
	SpanDivDivText1     *html.Node
	SpanDivText2        *html.Node
	SpanText2           *html.Node
	Text3               *html.Node
}

func TestFindByClass(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Body}.Find(ByClass("container")),
		nodes.Div)
}

func TestChildren(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.DivLi}.Children(),
		nodes.DivLiUl1, nodes.DivLiUl1)
}

func TestFindByClassChildrenText(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	str := root.
		Find(ByClass("list")).
		Children().
		Text()

	if str != "First valueSecond value" {
		t.Errorf("Expected \"First valueSecond value\" but found \"%s\"", str)
	}
}

func TestFindByClassContents(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	set := root.
		Find(ByClass("container")).
		Contents()

	if len(set) != 3 {
		t.Errorf("Expected 3 nodes but found %d", len(set))
	}
}

func TestFindByClassDuplicates(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	set := root.
		Find(ByClass("foo")).
		Find(ByClass("bar"))

	if len(set) != 1 {
		t.Errorf("Expected 1 node but found %d", len(set))
	}
}

func TestFirstByTag(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	set := root.
		First(ByTag(atom.Div))

	if len(set) != 1 {
		t.Errorf("Expected 1 node but found %d", len(set))
	}

	if set.Attr("class") != "container" {
		t.Errorf("Expected \"container\" but found %s", set.Attr("class"))
	}
}

func TestFirstByClassNext(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	set := root.
		First(ByClass("list")).
		First().
		Next()

	if set.Text() != "Second value" {
		t.Errorf("Expected \"Second value\" but found \"%s\"", set.Text())
	}

	set = set.Prev()
	if set.Text() != "First value" {
		t.Errorf("Expected \"First value\" but found \"%s\"", set.Text())
	}
}

func traverse(ch chan *html.Node, n *html.Node) {
	ch <- n
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(ch, c)
	}
}

func htmlNodes() HtmlNodes {
	n, _ := html.Parse(strings.NewReader(htmlData))
	ch := make(chan *html.Node)

	go traverse(ch, n)

	return HtmlNodes{
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
		<-ch,
	}
}

func nodeString(n *html.Node) string {
	var nodeStr string
	switch n.Type {
	case html.ErrorNode:
		nodeStr = "[ErrorNode]"
	case html.TextNode:
		nodeStr = fmt.Sprintf("[TextNode] %#v", n.Data)
	case html.DocumentNode:
		nodeStr = "[DocumentNode]"
	case html.ElementNode:
		nodeStr = fmt.Sprintf("<%s>", n.DataAtom)
	case html.CommentNode:
		nodeStr = "[CommentNode] " + n.Data
	case html.DoctypeNode:
		nodeStr = "[DoctypeNode]"
	}

	return nodeStr
}
func printNodes(n *html.Node, level int) {
	fmt.Printf("% *s%s\n", 2*level, "", nodeString(n))

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printNodes(c, level+1)
	}
}

func assertNodes(t *testing.T, set Set, nodes ...*html.Node) {
	if len(set) != len(nodes) {
		t.Errorf("Expected %d node(s) but found %d", len(nodes), len(set))
		return
	}

	for i, n := range nodes {
		if set[i] != n {
			t.Errorf("Expected node %d to be %s but found %s", i, nodeString(n), nodeString(set[i]))
			return
		}
	}
}
