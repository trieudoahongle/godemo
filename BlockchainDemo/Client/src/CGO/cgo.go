package main

// #include "foo.cpp"
// typedef int (*intFunc) ();
//
// int bridge_int_func(intFunc f)
// {
//		return f();
// }
//
// int fortytwo()
// {
//	    return 42;
// }
// int getFoo(){
//      Foo f(1);
//      return f.getValue()
//}
//
import "C"
import "fmt"

func main() {
	f := C.intFunc(C.fortytwo)
	fmt.Println(int(C.bridge_int_func(f)))
	fmt.Println(C.fortytwo())
	fmt.Println(C.getFoo())
	// Output: 42
}
