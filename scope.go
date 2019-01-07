package main

type scope struct {
	symbols map[string]interface{}
	parent  *scope
}

func newScope(s *scope) *scope {
	return &scope{map[string]interface{}{}, s}
}

func (scope *scope) resolve(name string) interface{} {
	for scope != nil {
		if s, ok := scope.symbols[name]; ok {
			return s
		}
		scope = scope.parent
	}
	return nil
}
