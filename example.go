package main

import (
	"crdt-go/crdt"
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

	// Querying Value
	fmt.Println(SetS.Query("Hello"))   // true <- because it happens on the same time, but the bias is add
	fmt.Println(SetS.Query("Hello213"))  // false
	fmt.Println(SetT.Query(123))  // false <- because it happens on the same time, but the bias is remove
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

	// Comparison
	fmt.Println(SetS.Compare(SetT, false)) // false

	// Merging
	SetS.Merge(SetT) // SetS is equal to union of SetS and SetT after this line, but SetT is unchanged
	SetT.Merge(SetS)

	fmt.Println(SetS.Compare(SetT, false)) // true
}

func main(){
	example()
}
