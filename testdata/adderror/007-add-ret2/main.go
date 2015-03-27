package main

import "fmt"

func bar() { //<<<<<adderror,9,1,9,1,pass
	if 3 < 5 {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}

func main() {
	bar()
	fmt.Println("Done")
}
