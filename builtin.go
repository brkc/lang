package main

import (
	"fmt"
)

func builtin(name string, args []*expression) (*expression, error) {
	switch name {
	case "print":
		print(args)
	default:
		return nil, fmt.Errorf("could not find fn: '%s'", name)
	}
	return nil, nil
}

func print(args []*expression) {
	for _, arg := range args {
		switch arg.typeValue {
		case stringType:
			fmt.Printf("%s\n", arg.value.(string))
		case numberType:
			fmt.Printf("%d\n", arg.value.(int))
		case booleanType:
			fmt.Printf("%t\n", arg.value.(bool))
		}
	}
}
