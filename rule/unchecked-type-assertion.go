package rule

import (
	"fmt"
	"go/ast"
	"regexp"

	"github.com/mgechev/revive/lint"
)

// UncheckedTypeAssertion lints for unchecked type assertions that could panic.
type UncheckedTypeAssertionRule struct{}

// Apply applies the rule to given file.
func (r *UncheckedTypeAssertionRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	if len(arguments) > 1 {
		panic(`invalid configuration for "unchecked-type-assertion"`)
	}
	excludeRegexpStr := `^.+_test.go$`
	if len(arguments) == 1 {
		var ok bool
		excludeRegexpStr, ok = arguments[0].(string)
		if !ok {
			panic(`invalid value passed as argument string to the "unchecked-type-assertion" rule`)
		}
	}
	excludeRegexp, err := regexp.Compile(excludeRegexpStr)
	if err != nil {
		panic(fmt.Sprintf("regex failed to compile: %q", err))
	}
	if excludeRegexp.MatchString(file.Name) {
		return nil
	}

	var failures []lint.Failure

	fileAst := file.AST
	walker := lintUncheckedTypeAssertion{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *UncheckedTypeAssertionRule) Name() string {
	return "unchecked-type-assertion"
}

type lintUncheckedTypeAssertion struct {
	file                   *lint.File
	fileAst                *ast.File
	onFailure              func(lint.Failure)
	inFuncDecl             bool
	typeSwitchStmtAssign   ast.Stmt
	inTypeSwitchStmtAssign bool
	inCheckedAssignStmt    bool
}

func (w lintUncheckedTypeAssertion) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if n == w.typeSwitchStmtAssign {
		w.inTypeSwitchStmtAssign = true
		w.typeSwitchStmtAssign = nil
	}
	if tss, ok := n.(*ast.TypeSwitchStmt); ok {
		w.typeSwitchStmtAssign = tss.Assign
	} else if _, ok := n.(*ast.FuncDecl); ok {
		w.inFuncDecl = true
	} else if as, ok := n.(*ast.AssignStmt); ok && len(as.Lhs) == 2 {
		// also reject blank identifiers
		if ident, ok := as.Lhs[1].(*ast.Ident); ok && ident.Name != "_" {
			w.inCheckedAssignStmt = true
		}
	} else if _, ok := n.(*ast.TypeAssertExpr); ok {
		// no need to check module-level type assertions since that will get
		// caught on startup and simplifies expressiveness.
		if w.inFuncDecl && !w.inCheckedAssignStmt && !w.inTypeSwitchStmtAssign {
			w.onFailure(lint.Failure{
				Category:   "errors",
				Node:       n,
				Confidence: 1,
				Failure:    "unchecked type assertion",
			})
		}
	} else if _, ok := n.(*ast.CallExpr); ok {
		w.inCheckedAssignStmt = false
	}
	return w
}
