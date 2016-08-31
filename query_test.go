package query

import (
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

func TestFindByClass(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Set{node}
	set := root.
		Find(ByClass("container"))

	if len(set) != 1 {
		t.Errorf("Expected 1 node but found %d", len(set))
	}
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
