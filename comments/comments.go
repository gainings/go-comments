package comments

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"sort"
	"strings"
)

//Process is adding comment templates to should-commented variable names, structure definitions, and function names
func Process(filename string, src []byte) ([]byte, error) {
	if src == nil {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		src = b
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	comments := []*ast.CommentGroup{}
	for _, cg := range f.Comments {
		comments = append(comments, cg)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch nt := n.(type) {
		case *ast.FuncDecl:
			if nt.Name.IsExported() && checkComment(nt.Doc.Text(), nt.Name.Name) {
				comment := &ast.Comment{
					Text:  fmt.Sprintf("//%s is TODO: need to enter a comment", nt.Name),
					Slash: nt.Pos() - 1,
				}
				cg := &ast.CommentGroup{
					List: []*ast.Comment{comment},
				}
				nt.Doc = cg
				comments = append(comments, cg)
			}
		case *ast.GenDecl:
			for _, spec := range nt.Specs {
				switch spt := spec.(type) {
				case *ast.TypeSpec:
					if spt.Name.IsExported() && checkComment(spt.Doc.Text(), spt.Name.Name) {
						comment := &ast.Comment{
							Text:  fmt.Sprintf("//%s is TODO: need to enter a comment", spt.Name),
							Slash: nt.TokPos - 1,
						}
						cg := &ast.CommentGroup{
							List: []*ast.Comment{comment},
						}
						nt.Doc = cg
						comments = append(comments, cg)
					}
				case *ast.ValueSpec:
					for _, name := range spt.Names {
						if name.IsExported() && checkComment(nt.Doc.Text(), name.Name) {
							var pos token.Pos
							if nt.Lparen == 0 {
								pos = nt.Pos()
							} else {
								pos = name.Pos()
							}
							comment := &ast.Comment{
								Text:  fmt.Sprintf("//%s is TODO: need to enter a comment", name.Name),
								Slash: pos - 1,
							}
							cg := &ast.CommentGroup{
								List: []*ast.Comment{comment},
							}
							nt.Doc = cg
							comments = append(comments, cg)
						}
					}
				}
			}
		}
		return true
	})
	sort.Slice(comments, func(i, j int) bool { return comments[i].Pos() < comments[j].Pos() })
	f.Comments = comments
	p := printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 8,
	}
	var buf bytes.Buffer
	err = p.Fprint(&buf, fset, f)
	if err != nil {
		return nil, err
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return out, nil
}

func checkComment(com, funcName string) bool {
	if com == "" {
		return true
	}
	return !strings.HasPrefix(com, funcName)
}
