package regexp

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type Range struct {
	Start uint32
	End   uint32
}

func ParseJavascript(source []byte) []Range {
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

	var ranges []Range
	forEachPreorder(cursor, func(n *sitter.Node) {
		ranges = append(ranges, Range{n.StartByte(), n.EndByte()})
	})
	return ranges
}

func forEachPreorder(cursor *sitter.TreeCursor, f func(*sitter.Node)) {
	// visit root
	f(cursor.CurrentNode())

	// visit subtrees from left to right
	onChild := false
	for valid := cursor.GoToFirstChild(); valid; valid = cursor.GoToNextSibling() {
		onChild = true
		forEachPreorder(cursor, f)
	}

	if onChild {
		cursor.GoToParent()
	}
}
