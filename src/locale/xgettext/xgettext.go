// Copyright 2015 Luke Shumaker

package xgettext

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type KeywordSpec struct {
	ID        string
	ArgNum1   int
	ArgNum2   int
	TotalArgs int
	XComment  string
}

type Message struct {
	XComment    string
	Reference   []token.Position
	Flags       []string
	MsgCtxt     *string
	MsgID       string
	MsgIDPlural *string
}

func (msg Message) String() string {
	var str string
	if msg.XComment != "" {
		for _, line := range strings.Split(msg.XComment, "\n") {
			str += fmt.Sprintf("#. %s\n", line)
		}
	}
	if msg.Reference != nil && len(msg.Reference) > 0 {
		strs := make([]string, len(msg.Reference))
		for i, ref := range msg.Reference {
			strs[i] = ref.String()
		}
		str += fmt.Sprintf("#: %s\n", strings.Join(strs, ", "))
	}
	if msg.Flags != nil && len(msg.Flags) > 0 {
		str += fmt.Sprintf("#, %s\n", strings.Join(msg.Flags, ", "))
	}
	if msg.MsgCtxt != nil {
		str += fmt.Sprintf("msgctxt %q\n", *msg.MsgCtxt)
	}
	str += fmt.Sprintf("msgid %q\n", msg.MsgID)
	if msg.MsgIDPlural == nil {
		str += "msgstr \"\"\n"
	} else {
		str += fmt.Sprintf("msgid_plural %q\n", *msg.MsgIDPlural)
		str += "msgstr[0] \"\"\n"
		str += "msgstr[1] \"\"\n"
	}
	return str
}

type File struct {
	AST  *ast.File
	FSet *token.FileSet
	cmap cmap
}

type cmapel struct {
	Comment *ast.CommentGroup
	End     token.Position
}

type cmap struct {
	GoCMap ast.CommentMap
	MyCMap map[int]cmapel
}

func (cmap cmap) Get(node ast.Node, pos token.Position) string {
	if cgroups, hascgroups := cmap.GoCMap[node]; hascgroups {
		strs := make([]string, len(cgroups))
		for i, cgroup := range cgroups {
			strs[i] = cgroup.Text()
		}
		return strings.TrimSpace(strings.Join(strs, "\n"))
	} else if comment, ok := cmap.MyCMap[pos.Line]; ok && comment.End.Column <= pos.Column {
		return strings.TrimSpace(comment.Comment.Text())
	} else if comment, ok := cmap.MyCMap[pos.Line-1]; ok {
		return strings.TrimSpace(comment.Comment.Text())
	}
	return ""
}

func (file *File) initCMap() {
	goCMap := ast.NewCommentMap(file.FSet, file.AST, file.AST.Comments)

	myCMap := make(map[int]cmapel, len(file.AST.Comments))
	for _, comment := range file.AST.Comments {
		end := file.FSet.Position(comment.End())
		myCMap[end.Line] = cmapel{
			Comment: comment,
			End:     end,
		}
	}
	file.cmap = cmap{
		GoCMap: goCMap,
		MyCMap: myCMap,
	}
}

func toBasicLitString(n ast.Node) *ast.BasicLit {
	switch node := n.(type) {
	case *ast.BasicLit:
		if node.Kind == token.STRING {
			return node
		}
	case *ast.BinaryExpr:
		if node.Op != token.ADD {
			return nil
		}
		X := toBasicLitString(node.X)
		Y := toBasicLitString(node.Y)
		if X == nil || Y == nil {
			return nil
		}
		x, _ := strconv.Unquote(X.Value)
		y, _ := strconv.Unquote(Y.Value)
		return &ast.BasicLit{
			ValuePos: node.Pos(),
			Kind:     token.STRING,
			Value:    strconv.Quote(x + y),
		}
	case *ast.ParenExpr:
		x := toBasicLitString(node.X)
		if x == nil {
			return nil
		}
		x.ValuePos = node.Pos()
		return x
	}
	return nil
}

func (file *File) basicLitToMessage(node *ast.BasicLit) Message {
	if node.Kind != token.STRING {
		panic("only call basicLitToMessage with a string node")
	}
	val, _ := strconv.Unquote(node.Value)
	pos := file.FSet.Position(node.Pos())
	xcomment := file.cmap.Get(node, pos)

	return Message{
		XComment:    xcomment,
		Reference:   []token.Position{pos},
		Flags:       []string{},
		MsgCtxt:     nil,
		MsgID:       val,
		MsgIDPlural: nil,
	}
}

func (file *File) ExtractStrings(all bool, keywords []KeywordSpec) <-chan Message {
	file.initCMap()
	ret := make(chan Message)
	go func() {
		defer close(ret)
		ast.Inspect(file.AST, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ImportSpec:
				return false
			case *ast.BasicLit, *ast.BinaryExpr, *ast.ParenExpr:
				str := toBasicLitString(node)
				if str == nil {
					return true
				}
				if val, _ := strconv.Unquote(str.Value); val == "" {
					return true
				}
				ret <- file.basicLitToMessage(str)
				return false
			}
			return true
		})
	}()
	return ret
}
