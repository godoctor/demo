package main

import "fmt"

func foo() { //<<<<<adderror,5,1,5,1,pass
	fmt.Println("in foo")
}

func main() {
	foo()
	fmt.Println("Done")
}
