package main

import (
	"crdt-go/crdt"
	"time"
)

func example() {

	// Initialization, specifying the bias
	var SetS = crdt.NewLWWSet("ADD")
	var SetT = crdt.NewLWWSet("REMOVE")

	// Adding Value
	SetS.Add("Hello")
	SetT.Add(123)

	// Getting contents
	SetS.GetAddSet() // [Hello:2019-08-18 19:51:41.2206383 +0800 CST m=+0.004000201]
	SetS.GetRemoveSet() // map[]
	SetS.GetBias() // ADD

	// Querying Value
	SetS.Query("Hello")   // true
	SetS.Query("Hello213")  // false
	SetT.Query(123)  // true
	SetT.Query(1234)  // false

	// Removing value
	SetS.Remove("Hello")
	SetT.Remove(123)
	SetT.Remove(1234)

	// Querying Value Again
	SetS.Query("Hello")   // true <- because it happens on the same time, but the bias is add
	SetS.Query("Hello213")  // false
	SetT.Query(123)  // false <- because it happens on the same time, but the bias is remove
	SetT.Query(1234)  // false

	// Lets wait 1 ms before removing again
	time.Sleep(1 * time.Microsecond)
	SetS.Remove("Hello")
	SetS.Query("Hello") // now it is false

	// Lets wait 1 ms before reinserting again
	time.Sleep(1 * time.Microsecond)
	SetS.Add("Hello")
	SetS.Query("Hello") // now it is true again

	// Getting the total contents (values that exist in the set)
	SetS.GetContent() // map[Hello:2019-08-18 19:58:59.347824 +0800 CST m=+0.007056501]
	SetT.GetContent() // map[]

	// Setting up for comparison, we have multiple way to compare values
	var SetU = crdt.NewLWWSet("ADD")
	var SetV = crdt.NewLWWSet("ADD")
	var SetW = crdt.NewLWWSet("REMOVE")

	SetU.Add("Hello")
	SetV.Add(123)
	SetW.Add(123)

	// Comparing full history (including comparing their time of insertion/deletion)
	// the second parameter is whether to compare the bias or not
	SetV.Compare(SetU, false) // false, obviously
	SetV.Compare(SetW, false) // true, obviously
	SetV.Compare(SetW, true) // false, because they have different bias

	// Lets do something with SetU
	time.Sleep(1 * time.Microsecond)
	SetU.Remove("Hello")
	SetU.Add(123)

	SetU.GetContent() // map[123:2019-08-18 20:06:25.1233151 +0800 CST m=+0.008996701]
	SetW.GetContent() // map[123:2019-08-18 20:06:25.1213148 +0800 CST m=+0.006996401]

	SetU.Compare(SetW, false) // false, despite having same content
	// because their history is different (and that they dont have same insertion time)

	// If we want to ignore the operation time and just compare their content,
	// Similarly, the second parameter is whether to compare the bias or not
	SetU.CompareContent(SetW, false) // true, comparing content, ignoring time
	SetU.CompareContent(SetW, true) // false, because they have different bias

	// Merging
	// Setting up for comparison, we have multiple way to compare values
	var SetX = crdt.NewLWWSet("ADD")
	var SetY = crdt.NewLWWSet("ADD")
	SetX.Add("Hello")
	SetY.Add(123)

	// There are 2 flavor of merging
	// functional programming style:
	var SetZ = crdt.Merge(SetX, SetY)
	SetX.GetContent() // map[Hello:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001]
	SetY.GetContent() // map[123:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001]
	SetZ.GetContent() // map[123:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001 Hello:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001]
	// note that the original SetX and SetY are left untouched

	// but if performance is more important, past state is not needed then
	SetX.MergeWith(SetY) // SetS is now equal to union of SetX and SetY after this line, but SetT is unchanged
	SetX.GetContent() //map[123:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001 Hello:2019-08-18 20:17:34.9821332 +0800 CST m=+0.009075001]

}

func main(){
	example()
}
