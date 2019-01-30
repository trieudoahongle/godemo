package httpMethod

import (
	"fmt"
	"net/http"
)

type FuncPointer func(int, int) int

var mapfunction = make(map[string]FuncPointer)

func Testfuncpointer() {

	calculate(Plus)
	calculate(Minus)
	calculate(Multiply)

	mapfunction[http.MethodGet] = Plus
	mapfunction[http.MethodPost] = Minus
	mapfunction[http.MethodPut] = Multiply
	mapfunction[http.MethodDelete] = Multiply
	mapfunction[http.MethodPatch] = Multiply
	fmt.Println(len(mapfunction))

	calculate(mapfunction["PLUS"])
	calculate(mapfunction["MINUS"])
	calculate(mapfunction["MULTIPLY"])
}

func calculate(fp func(int, int) int) {
	ans := fp(3, 2)
	fmt.Printf("%v\n", ans)
}

// This is the same function but uses the type/fp defined above
//
// func calculate (fp ArithOp) {
//     ans := fp(3,2)
//     fmt.Printf("\n%v\n", ans)
// }

func Plus(a, b int) int {
	return a + b
}

func Minus(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}
