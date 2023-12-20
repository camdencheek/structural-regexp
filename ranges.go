package regexp

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type Ranges struct {
	starts []uint32
	ends   []uint32
}

func ParseJavascript(source []byte) Ranges {
	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(javascript.GetLanguage())
	tree, err := parser.ParseCtx(context.TODO(), nil, source)
	if err != nil {
		panic(err)
	}
	defer tree.Close()

	cursor := sitter.NewTreeCursor(tree.RootNode())
	defer cursor.Close()

	var ranges Ranges
	forEachPreorder(cursor, func(n *sitter.Node) {
		ranges.starts = append(ranges.starts, n.StartByte())
		ranges.ends = append(ranges.ends, n.EndByte())
	})
	return ranges
}

func forEachPreorder(cursor *sitter.TreeCursor, f func(*sitter.Node)) {
	// visit root
	f(cursor.CurrentNode())

	// visit subtrees from left to right
	for valid := cursor.GoToFirstChild(); valid; valid = cursor.GoToNextSibling() {
		forEachPreorder(cursor, f)
	}

	// return cursor to the node it started on
	cursor.GoToParent()
}
