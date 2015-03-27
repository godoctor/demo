package main

import "fmt"

func foo() (int, string) { //<<<<<adderror,5,1,5,1,pass
	return 5, "foo"
}

func main() {
	foo()
	fmt.Println("Done")
}
