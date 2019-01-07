package main

type symbol struct {
	name  string
	value interface{}
}

type scope struct {
	symbols map[string]*symbol
	parent  *scope
}

func newScope(s *scope) *scope {
	return &scope{map[string]*symbol{}, s}
}

func (scope *scope) declare(name string, value interface{}) {
	scope.symbols[name] = &symbol{name, value}
}

func (scope *scope) assign(name string, value interface{}) {
	if symbol := scope.resolve(name); symbol != nil {
		symbol.value = value
	}
}

func (scope *scope) resolve(name string) *symbol {
	for scope != nil {
		if s, ok := scope.symbols[name]; ok {
			return s
		}
		scope = scope.parent
	}
	return nil
}
