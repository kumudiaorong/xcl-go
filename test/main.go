package main

import (
	"fmt"
	"k0e.top/xcl"
)

func main() {
	c, err := xcl.Load("./config.xcl")
	if err != nil {
		fmt.Println(err)
		return
	}
	// c.TryInsert("a'b'c", 1.0)
	// c.TryInsert("a",1)
	// c.TryInsert("a'b",1)
	c.TryInsert("a'd","asd")
	fmt.Println(c.String())
	c.Save(true)
}

// package main

// import (
//     "fmt"

//     "example.com/greetings"
// )

// func main() {
//     // Get a greeting message and print it.
//     message := greetings.Hello("Gladys")
//     fmt.Println(message)
// }
