package parser

import (
	"github.com/robertkrimen/otto/ast"
)

type _scope struct {
	outer           *_scope
	allowIn         bool
	inIteration     bool
	inSwitch        bool
	inFunction      bool
	declarationList []ast.Declaration
	block           bool

	labels []string
}

func (self *_parser) openScope() {
	self.scope = &_scope{
		outer:   self.scope,
		allowIn: true,
	}
}

func (self *_parser) openBlockScope() {
	self.scope = &_scope{
		outer:       self.scope,
		block:       true,
		allowIn:     self.scope.allowIn,
		inIteration: self.scope.inIteration,
		inSwitch:    self.scope.inSwitch,
		inFunction:  self.scope.inFunction,
	}
}

func (self *_parser) closeScope() {
	self.scope = self.scope.outer
}

func (self *_scope) definitionScope() *_scope {
	scope := self
	for scope.block && scope.outer != nil {
		scope = scope.outer
	}
	return scope
}

func (self *_scope) declare(declaration ast.Declaration) {
	scope := self
	if declaration, ok := declaration.(*ast.VariableDeclaration); !ok || !declaration.Block {
		scope = scope.definitionScope()
	}
	scope.declarationList = append(scope.declarationList, declaration)
}

func (self *_scope) hasLabel(name string) bool {
	scope := self.definitionScope()
	for _, label := range scope.labels {
		if label == name {
			return true
		}
	}
	return false
}

func (self *_scope) pushLabel(name string) {
	scope := self.definitionScope()
	scope.labels = append(scope.labels, name)
}
