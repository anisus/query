package query

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
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
	</body>
</html>
`

func TestFindByClass(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(htmlData))
	root := Nodes{node}
	nodes := root.
		Find(ByClass("container"))

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node but found %d", len(nodes))
	}
}
