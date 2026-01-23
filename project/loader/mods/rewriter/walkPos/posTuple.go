package walkPos

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type PosTuple struct {
	// Node is the ast.Node that has this position as a field.
	Node ast.Node

	// Pos it the pointer to the position field.
	// If this value is modified, it will modify the node.
	Pos *token.Pos

	// Width is the width of the token for this position.
	// This may not always be the length of the Text.
	Width int

	// Text is the token at this position.
	Text string

	// Id is an indicator of what position this is in the document.
	// This should mostly used for debugging.
	Id string

	// Pseudo indicates if this tuple is a pseudo position that
	// does not actually exist as a field on the given node.
	// If pseudo, then setting the position will have no effect on the node.
	Pseudo bool
}

const (
	nodeStartId = `NodeStart`
	nodeEndId   = `NodeEnd`
)

// End is the first position after this position tuple
// that is not used by this tuple.
func (pt *PosTuple) End() token.Pos {
	return *(pt.Pos) + token.Pos(pt.Width)
}

func (pt *PosTuple) IsNodeEdge() bool {
	return pt.IsNodeStart() || pt.IsNodeEnd()
}

func (pt *PosTuple) IsNodeStart() bool {
	return pt.Id == nodeStartId
}

func (pt *PosTuple) IsNodeEnd() bool {
	return pt.Id == nodeEndId
}

func (pt *PosTuple) String() string {
	id := pt.Id
	if !strings.Contains(id, `.`) {
		typStr := fmt.Sprintf("%T", pt.Node)
		typStr = strings.TrimPrefix(typStr, `*ast.`)
		id = typStr + `.` + id
	}
	text := ``
	if len(pt.Text) > 0 {
		text = strconv.Quote(pt.Text)
	}
	if pt.Pseudo {
		text += `(P)`
	}
	return fmt.Sprintf("%d:%v:%d%s", *pt.Pos, id, pt.Width, text)
}

func newPosTuple(n ast.Node, pos *token.Pos, text string, width int, id string) *PosTuple {
	if pos.IsValid() {
		return &PosTuple{
			Node:  n,
			Pos:   pos,
			Width: width,
			Text:  text,
			Id:    id,
		}
	}
	return nil
}

func commentTuple(c *ast.Comment, id string) *PosTuple {
	return newPosTuple(c, &c.Slash, c.Text, len(c.Text), id)
}

func tokTuple(n ast.Node, pos *token.Pos, tok token.Token, id string) *PosTuple {
	return newPosTuple(n, pos, tok.String(), tokWidth(tok), id)
}

func appendTuple(frame []*PosTuple, pt *PosTuple) []*PosTuple {
	if pt != nil {
		return append(frame, pt)
	}
	return frame
}
