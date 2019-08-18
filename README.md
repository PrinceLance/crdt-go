# crdt-go

## Introduction
This is a Golang implementation of crdt Last-Writer-Wins element set.
Currently it have 100% test coverage

## Getting Started
First, install the library using below command in the console
``` bash
go get github.com/PrinceLance/crdt-go
```

In your code
``` go
import "github.com/PrinceLance/crdt-go/crdt"
```

See example.go  for quick copy and paste code examples

## Functions
### Last-write-wins element set
Last-write-wins element set tracks element addition and removal.
Each element operation is attached with a timestamp. 

#### type LWWSet struct
``` go
type LWWSet struct {
	// Addition Set of LWW Set
	addSet map[interface{}]time.Time

	// Removal Set of LWW Set
	removeSet  map[interface{}]time.Time

	// Bias, possible value is "ADD", "REMOVE". If nothing is given
	// or other value is given, default to "ADD"
	bias string
}
```

#### func NewLWWSet(bias string) *LWWSet
Given a string bias, return a LWWSet object with empty add set, empty removal set and the indicated bias. 
``` go
var SetS = crdt.NewLWWSet("ADD")
var SetT = crdt.NewLWWSet("REMOVE")

// defaulted to "ADD"
var SetS = crdt.NewLWWSet("")
```

#### func (setS *LWWSet) GetContent() map[interface{}]time.Time
Return a map of (element, timestamp) containing the element that exist in the LWW element set.
If an element exist in both add and removal set, then the latest one determine it's status.
If the time is equal, then the bias will determine it.
``` go
SetS.GetContent() // map[Hello:2019-08-18 19:58:59.347824 +0800 CST m=+0.007056501]
```

#### func (setS *LWWSet) GetAddSet() map[interface{}]time.Time
Return a map of (element, timestamp) containing the element that exist in the add set.
``` go
SetS.GetAddSet() // map[Hello:2019-08-18 19:58:59.347824 +0800 CST m=+0.007056501]
```

#### func (setS *LWWSet) GetRemoveSet() map[interface{}]time.Time
Return a map of (element, timestamp) containing the element that exist in the removal set.
``` go
SetS.GetRemoveSet() // map[Hello:2019-08-18 19:58:59.347824 +0800 CST m=+0.007056501]
```

#### func (setS *LWWSet) GetBias() string
Return the bias
``` go
SetS.GetBias() // "ADD"
```

#### func (setS *LWWSet) Add(value interface{})
Add an element to the set
``` go
SetS.Add("Hello")
SetT.Add(123)
```

#### func (setS *LWWSet) Remove(value interface{})
Remove an element from the set
``` go
SetS.Remove("Hello")
SetT.Remove(123)
```

#### func (setS *LWWSet) Query(value interface{}) bool
Find if an element is in the set
``` go
SetS.Add("Hello")
SetS.Query("Hello") // true

SetS.Remove("Hello")
SetS.Query("Hello") // false
```

#### func (setS *LWWSet) Compare(setT *LWWSet, compareBias bool) bool
Compare 2 LWW sets and determine if it is equal or not
The second parameter is for to compare the bias or not

It perform a full history search. Thus it check if both set have to have the same
(element, timestamp) values in both add and removal set.
``` go
SetV.Compare(SetU, false) // true/false
```

#### func (setS *LWWSet) CompareContent(setT *LWWSet, compareBias bool) bool
Compare 2 LWW sets and determine if the content is equal or not
The second parameter is for to compare the bias or not

It is not as strict as Compare, it simply compare if bot set have the same
element or not, regardless of the insertion time, or history of removed elements.
``` go
SetV.CompareContent(SetU, false) // true/false
```

#### func Merge(setS *LWWSet, setT *LWWSet) *LWWSet
Merge 2 sets, returning their union.
This is an idempotent operation.
``` go
var SetZ = crdt.Merge(SetX, SetY)
``` 

#### func (setS *LWWSet) MergeWith(setT *LWWSet)
Merges the set in parameter into another set.
This is NOT an idempotent operation, as it changes the value of the setS
``` go
SetX.MergeWith(SetY)
``` 