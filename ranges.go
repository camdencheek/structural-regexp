package regexp

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

// TODO: compress/decompress range.
// https://github.com/ronanh/intcomp looks nice.
type Range struct {
	Start uint32
	End   uint32
}

// ColumnRanges is a column-oriented representation of []Range
type ColumnRanges struct {
	Starts []uint32
	Ends   []uint32
}

func (c *ColumnRanges) Slice(start, end int) ColumnRanges {
	return ColumnRanges{
		Starts: c.Starts[start:end],
		Ends:   c.Ends[start:end],
	}
}

func ColumnRangesFromRanges(input []Range) ColumnRanges {
	output := ColumnRanges{
		Starts: make([]uint32, 0, len(input)),
		Ends:   make([]uint32, 0, len(input)),
	}
	for _, r := range input {
		output.Starts = append(output.Starts, r.Start)
		output.Ends = append(output.Ends, r.End)
	}
	return output
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
