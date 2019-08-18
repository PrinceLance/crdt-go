package main

import (
	"CRDT/crdt"
	"fmt"
	"time"
)

func example() {

	// Initialization, specifying the bias
	var SetS = crdt.NewLWWSet("ADD")
	var SetT = crdt.NewLWWSet("REMOVE")

	// Adding Value
	SetS.Add("Hello")
	SetS.Add("Hello")
	SetT.Add(123)

	// Querying Value
	fmt.Println(SetS.Query("Hello"))   // true
	fmt.Println(SetS.Query("Hello213"))  // false
	fmt.Println(SetT.Query(123))  // true
	fmt.Println(SetT.Query(1234))  // false

	// Removing value
	SetS.Remove("Hello")
	SetT.Remove(123)
	SetT.Remove(1234)

	// Printing Contents for debugging
	SetS.Print()
	//Contents of add setmap[Hello:2019-08-18 16:21:24.9182782 +0800 CST m=+0.011999801]
	//Contents of remove setmap[Hello:2019-08-18 16:21:24.9192719 +0800 CST m=+0.012993401]
	//Contents of biasADD

	SetT.Print()
	//Contents of add setmap[123:2019-08-18 16:21:24.9182782 +0800 CST m=+0.011999801]
	//Contents of remove setmap[123:2019-08-18 16:21:24.9192719 +0800 CST m=+0.012993401 1234:2019-08-18 16:21:24.9192719 +0800 CST m=+0.012993401]
	//Contents of biasREMOVE

	// Querying Value
	fmt.Println(SetS.Query("Hello"))   // true <- because it happens on the same time, but the bias is add
	fmt.Println(SetS.Query("Hello213"))  // false
	fmt.Println(SetT.Query(123))  // false <- because it happens on the same time, but the remove is add
	fmt.Println(SetT.Query(1234))  // false

	// Lets wait 1 ms before removing again
	time.Sleep(1 * time.Microsecond)
	SetS.Remove("Hello")
	fmt.Println(SetS.Query("Hello")) // now it is false
	SetS.Remove("Hello")
	fmt.Println(SetS.Query("Hello")) // now it is true again

	// Resinsertion
	time.Sleep(1 * time.Microsecond)
	SetS.Add("Hello")
	fmt.Println(SetS.Query("Hello")) // now it is true
}

func main(){
	example()
}
