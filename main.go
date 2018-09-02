package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	if len(os.Args) == 0 {
		log.Fatal(errors.New("args is too short"))
	}
	fname := os.Args[1]
	_, err := os.Stat(fname)
	if err != nil {
		log.Fatal(errors.New("file not exists"))
	}
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
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
	if err := p.Fprint(os.Stdout, fset, f); err != nil {
		log.Fatal(err)
	}
}

func checkComment(com, funcName string) bool {
	if com == "" {
		return true
	}
	return !strings.HasPrefix(com, funcName)
}
