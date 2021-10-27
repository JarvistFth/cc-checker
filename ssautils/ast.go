package ssautils

import (
	"go/ast"
	"go/token"
	"go/types"
)

func checkAst(fset *token.FileSet, astfile *ast.File, pkg *types.Info) {

	ast.Inspect(astfile, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:

			for _,expr := range stmt.Rhs{
				if callExpr, ok := expr.(*ast.CallExpr); ok{
					pos := returnsError(callExpr,pkg)

					if pos < 0 || pos >=len(stmt.Lhs){
						continue
					}


					if id, ok := stmt.Lhs[pos].(*ast.Ident); ok && id.Name == "_"{
						log.Warningf("unhandled error at %d", fset.Position(stmt.Pos()))
					}
				}
			}
		case *ast.ExprStmt:
			if callExpr, ok := stmt.X.(*ast.CallExpr); ok{
				pos := returnsError(callExpr, pkg)
				if pos >= 0{
					log.Warningf("unhandled error at %d", fset.Position(stmt.Pos()))
				}
			}
		}
		return true
	})
}

func returnsError(callExpr *ast.CallExpr, pkgInfo *types.Info) int {
	if tv := pkgInfo.TypeOf(callExpr); tv != nil{
		switch t := tv.(type) {
		case *types.Tuple:
			for pos := 0; pos < t.Len(); pos++{
				variable := t.At(pos)
				if variable != nil && variable.Type().String() == "error"{
					return pos
				}
			}
		case *types.Named:
			if t.String() == "error"{
				return 0
			}
		}
	}

	return -1
}
