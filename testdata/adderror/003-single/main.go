package main

import "fmt"

func foo() int { //<<<<<adderror,5,1,5,1,pass
	return 5
}

func main() {
	foo()
	fmt.Println("Done")
}
