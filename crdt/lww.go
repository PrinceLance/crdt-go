package crdt

import (
	"time"
)

// Implementation of state based LWW set
// According to https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type#LWW-Element-Set_(Last-Write-Wins-Element-Set)
// The set should provide those implementation:
//	payload set A, set R
//	query lookup(element e) : boolean b
//	update add(element e)
//	update remove(element e)
//	compare (S, T) : boolean b
//	merge (S, T) : payload U
type LWWSet struct {
	// Addition Set of LWW Set
	addSet map[interface{}]time.Time

	// Removal Set of LWW Set
	removeSet  map[interface{}]time.Time

	// Bias, possible value is "ADD", "REMOVE". If nothing is given
	// or other value is given, default to "ADD"
	bias string
}

// payload set A, set R
// NewLWWSet initialize the addset, removeset and assign the bias
func NewLWWSet(bias string) *LWWSet {
	// Bias, possible value is "ADD", "REMOVE". If nothing is given
	// or other value is given, default to "ADD"
	if bias != "REMOVE" {
		bias = "ADD"
	}

	return &LWWSet{
		addSet: make(map[interface{}]time.Time),
		removeSet:  make(map[interface{}]time.Time),
		bias: bias,
	}
}

// Content return a map containing the values
func (setS *LWWSet) GetContent() map[interface{}]time.Time {
	var content = make(map[interface{}]time.Time)

	// check element by element
	for addKey := range setS.addSet {
		var SAddTime = setS.addSet[addKey]
		var SRemoveTime, keyExistInRemoveSet  = setS.removeSet[addKey]
		// there are 4 cases
		// keyExistInRemoveSet == false --> value exist
		// SAddTime > SRemoveTime  --> value exist
		// SAddTime == SRemoveTime  --> value exist if bias is ADD
		// SAddTime < SRemoveTime  --> value do not exist / deleted
		if !keyExistInRemoveSet || SAddTime.After(SRemoveTime) ||
			( SAddTime.Equal(SRemoveTime) && setS.bias == "ADD"){
			content[addKey] = SAddTime
		}
	}

	return content
}

// GetAddSet return the add set of the set
func (setS *LWWSet) GetAddSet() map[interface{}]time.Time {
	var content = make(map[interface{}]time.Time)
	// Deep Copy
	for addKey := range setS.addSet {
		var SAddTime = setS.addSet[addKey]
		content[addKey] = SAddTime
	}
	return content
}

// GetRemoveSet return the remove set of the set
func (setS *LWWSet) GetRemoveSet() map[interface{}]time.Time {
	var content = make(map[interface{}]time.Time)
	// Deep Copy
	for removeKey := range setS.removeSet {
		var SRemoveTime = setS.removeSet[removeKey]
		content[removeKey] = SRemoveTime
	}
	return content
}

// GetBias return the bias of the set
func (setS *LWWSet) GetBias() string {
	return setS.bias
}

// update add(element e)
// Add adds a value to the add Set
func (setS *LWWSet) Add(value interface{}) {
	setS.addSet[value] = time.Now()
	// Optional, can also remove the same value from remove set
}

// update remove(element e)
// Remove adds a value to the remove Set
func (setS *LWWSet) Remove(value interface{}) {
	setS.removeSet[value] = time.Now()
	// Optional, can also remove the same value from add set
}

// query lookup(element e) : boolean b
// Query checks if the value exist in the lww set
func (setS *LWWSet) Query(value interface{}) bool {

	var addTime, addOk = setS.addSet[value]
	// If the value is not present in add set
	// then it is not in the set
	if !addOk {
		return false
	}

	var removeTime, removeOk = setS.removeSet[value]
	// If the value is present in add set but not in remove set
	// then it exist in the set
	if !removeOk {
		return true
	}

	// If the value exist in both add and remove set.
	// then the latest time wins
	// if both times are equal, then we we need to check the bias
	if addTime.Equal(removeTime) {
		// check bias
		if setS.bias == "REMOVE" {
			// biased towards remove, which mean remove wins
			return false
		} else {
			// biased towards add, which mean add wins
			return true
		}
 	} else {
		if addTime.After(removeTime) {
			// add is the latest input, it wins
			return true
		} else {
			// remove is the latest input, it wins
			return false
		}
	}
}

// compare (S, T) : boolean b
// Compare compares 2 sets
// define equal here as the LWW set having the same elements in add and remove set
// AND their times are equal
// if compareBias is true then it also check if their bias are the same
func (setS *LWWSet) Compare(setT *LWWSet, compareBias bool) bool {

	if compareBias && setS.bias != setT.bias {
		return false
	}

	// if the length of the sets are different then it is sure to be different
	if len(setS.addSet) != len(setT.addSet) ||
		len(setS.removeSet) != len(setT.removeSet) {
			return false
	}

	// check element by element
	for addKey := range setS.addSet {
		var SAddTime = setS.addSet[addKey]
		var TAddTime, keyExist = setT.addSet[addKey]
		// if key doesnt exist or not equal, it mean the 2 set are different
		if !keyExist || !SAddTime.Equal(TAddTime) {
			return false
		}
	}

	// check element by element
	for removeKey := range setS.removeSet {
		var SRemoveTime = setS.removeSet[removeKey]
		var TRemoveTime, keyExist = setT.removeSet[removeKey]
		// if key doesnt exist or not equal, it mean the 2 set are different
		if !keyExist || !SRemoveTime.Equal(TRemoveTime) {
			return false
		}
	}

	// "When you have eliminated the impossible, whatever remains, however improbable, must be the truth"
	return true
}

// CompareContent compares the contents 2 sets
// this is a less strict function to just compare the contents but don't care the time of addition
// if compareBias is true then it also check if their bias are the same
func (setS *LWWSet) CompareContent(setT *LWWSet, compareBias bool) bool {

	if compareBias && setS.bias != setT.bias {
		return false
	}

	var setSContent = setS.GetContent()
	var setTContent = setT.GetContent()

	// if the length of the sets are different then it is sure to be different
	if len(setSContent) != len(setTContent) {
		return false
	}

	// check element by element
	for addKey := range setSContent {
		var _, keyExist = setT.addSet[addKey] // do not care the time

		// if key doesnt exist or not equal, it mean the 2 set are different
		if !keyExist {
			return false
		}
	}

	// "When you have eliminated the impossible, whatever remains, however improbable, must be the truth"
	return true
}

// merge (S, T) : payload U
// Merge merges T and S and return the resulting LWW SET with bias same as S
func Merge(setS *LWWSet, setT *LWWSet) *LWWSet {
	// Deep Copying
	var setU = &LWWSet{
		addSet: setS.GetAddSet(),
		removeSet:  setS.GetRemoveSet(),
		bias: setS.bias,
	}

	setU.MergeWith(setT)
	return setU
}

// MergeWith merges T into S
func (setS *LWWSet) MergeWith(setT *LWWSet) {

	for addKey := range setT.addSet {
		var TAddTime = setT.addSet[addKey]
		var SAddTime, keyExist = setS.addSet[addKey]
		// the key doesnt exist in set S OR the add time in setT is after setS, we update the value in setS
		if !keyExist || TAddTime.After(SAddTime) {
			setS.addSet[addKey] = TAddTime
		}
	}

	for removeKey := range setT.removeSet {
		var TRemoveTime= setT.removeSet[removeKey]
		var SRemoveTime, keyExist= setS.removeSet[removeKey]
		// the key doesnt exist in set S OR the remove time in setT is after setS, we update the value in setS
		if !keyExist || (TRemoveTime.After(SRemoveTime)) {
			setS.removeSet[removeKey] = TRemoveTime
		}
	}
}