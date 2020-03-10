package fixtures

import "fmt"

var foo = 1
var bar = foo.(string)

func uncheckedTypeAssertion1() {
	foo := interface{}(1)
	_ = foo.(string) // MATCH /unchecked type assertion/
}

func uncheckedTypeAssertion2() {
	foo := interface{}(1)
	if _, ok := foo.(string); !ok {
		println("not ok")
	}
}

func uncheckedTypeAssertion3() {
	foo := interface{}(1)
	bar, _ := foo.(string) // MATCH /unchecked type assertion/
	println(bar)
}

func uncheckedTypeAssertion4() {
	foo := interface{}(1)
	switch foo.(type) {
	case int:
		println(1)
	}
}

func uncheckedTypeAssertion5() {
	foo := interface{}(1)
	switch bar := foo.(type) {
	case int:
		println(bar + 1)
	}
}

func uncheckedTypeAssertion6() {
	foo := interface{}(1)
	if _, err := fmt.Println(foo.(string)); err != nil { // MATCH /unchecked type assertion/
		return nil
	}
}
