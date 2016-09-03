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
			<ul class="list">
				<li><strong>First</strong> value</li>
				<li>Second value</li>
			</ul>
		</div>
		<span class="foo">Span
			<div class="foo moo" id="divId">Div
				<div class="bar" title="text">Bar</div>
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
	DivUl               *html.Node
	DivUlText1          *html.Node
	DivUlLi1            *html.Node
	DivUlLi1Strong      *html.Node
	DivUlLi1StrongText1 *html.Node
	DivUlLi1Text1       *html.Node
	DivUlText2          *html.Node
	DivUlLi2            *html.Node
	DivUlLi2Text1       *html.Node
	DivUlText3          *html.Node
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
	assertNodes(t,
		Set{nodes.Body}.Find(ByClass("foo")),
		nodes.Span, nodes.SpanDiv)
}

func TestFindShallow(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Body}.FindShallow(ByClass("foo")),
		nodes.Span)
}

func TestChildren(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.DivUl}.Children(),
		nodes.DivUlLi1, nodes.DivUlLi2)
}

func TestText(t *testing.T) {
	nodes := htmlNodes()
	assertString(t,
		Set{nodes.DivUlLi1, nodes.DivUlLi2}.Text(),
		"First valueSecond value")
}

func TestAttr(t *testing.T) {
	nodes := htmlNodes()
	set := Set{nodes.SpanDiv, nodes.SpanDivDiv}
	assertString(t,
		set.Attr("class"),
		"foo moo")
	assertString(t,
		set.Attr(""),
		"")
	assertString(t,
		set.Attr("href"),
		"")
	assertString(t,
		Set{}.Attr("class"),
		"")
}

func TestEq(t *testing.T) {
	nodes := htmlNodes()
	set := Set{nodes.DivUlLi1, nodes.DivUlLi2}
	assertNodes(t,
		set.Eq(-1))
	assertNodes(t,
		set.Eq(0),
		nodes.DivUlLi1)
	assertNodes(t,
		set.Eq(1),
		nodes.DivUlLi2)
	assertNodes(t,
		set.Eq(2))
}

func TestContents(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Div}.Contents(),
		nodes.DivText1, nodes.DivUl, nodes.DivText2)
}

// TestFindDuplicates tests that even if a Find will match the same element twice
// only a single instance will be found in the resulting Set
func TestFindDuplicates(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Span, nodes.SpanDiv}.Find(ByClass("bar")),
		nodes.SpanDivDiv)
}

func TestFirst(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Doc}.First(ByTag(atom.Div)),
		nodes.Div)
	assertNodes(t,
		Set{nodes.DivUl}.First(),
		nodes.DivUlLi1)
	assertNodes(t,
		Set{nodes.SpanDivDiv}.First())
}

func TestNext(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.DivUlLi1}.Next(),
		nodes.DivUlLi2)
	assertNodes(t,
		Set{nodes.DivUlLi2}.Next())
}

func TestPrev(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.DivUlLi2}.Prev(),
		nodes.DivUlLi1)
	assertNodes(t,
		Set{nodes.DivUlLi1}.Prev())
}

func TestFirstChild(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Body, nodes.DivUl}.FirstChild(),
		nodes.Div, nodes.DivUlLi1)
	assertNodes(t,
		Set{nodes.Body, nodes.DivUlLi2}.FirstChild(),
		nodes.Div)
}

func TestLastChild(t *testing.T) {
	nodes := htmlNodes()
	assertNodes(t,
		Set{nodes.Body, nodes.DivUl}.LastChild(),
		nodes.Span, nodes.DivUlLi2)
	assertNodes(t,
		Set{nodes.Body, nodes.DivUlLi2}.LastChild(),
		nodes.Span)
}

func TestFilter(t *testing.T) {
	nodes := htmlNodes()
	set := Set{
		nodes.Span,
		nodes.SpanText1,
		nodes.SpanDiv,
		nodes.SpanDivText1,
		nodes.SpanDivDiv,
		nodes.SpanDivDivText1,
		nodes.SpanDivText2,
		nodes.SpanText2,
	}

	assertNodes(t,
		set.Filter(ByTag(atom.Div)),
		nodes.SpanDiv, nodes.SpanDivDiv)
	assertNodes(t,
		set.Filter(ByTag(atom.H1)))
	assertNodes(t,
		set.Filter(ById("divId")),
		nodes.SpanDiv)
	assertNodes(t,
		set.Filter(ById("divId")),
		nodes.SpanDiv)
	assertNodes(t,
		set.Filter(ByAttr("title", "text")),
		nodes.SpanDivDiv)
	assertNodes(t,
		set.Filter(ByClass("foo")),
		nodes.Span, nodes.SpanDiv)
	assertNodes(t,
		set.Filter(ByClass("moo")),
		nodes.SpanDiv)
	assertNodes(t,
		set.Filter(ByType(html.TextNode)),
		nodes.SpanText1, nodes.SpanDivText1, nodes.SpanDivDivText1, nodes.SpanDivText2, nodes.SpanText2)
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

func assertString(t *testing.T, str, exp string) {
	if str != exp {
		t.Errorf("Expected %#v but found %#v", exp, str)
	}
}
