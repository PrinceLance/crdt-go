package crdt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLWWSet(t *testing.T) {

	t.Run("Set Creation",func(t *testing.T){
		var S =NewLWWSet("ADD")
		require.Empty(t, S.addSet, "add set is not empty")
		require.Empty(t, S.removeSet, "remove set is not empty")
		require.Equal(t, S.bias, "ADD","bias is not correct")

		var T =NewLWWSet("REMOVE")
		require.Equal(t, T.bias, "REMOVE","bias is not correct")

		var U =NewLWWSet("")
		require.Equal(t, U.bias, "ADD","default bias is not correct")
	})

	t.Run("Set Addition, Deletion and Query",func(t *testing.T){
		var S =NewLWWSet("ADD")
		require.False(t, S.Query("HelloWorld"), "set should be empty")

		time.Sleep(1 * time.Microsecond)
		S.Add("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "set should have this value")

		time.Sleep(1 * time.Microsecond)
		S.Remove("HelloWorld")
		require.False(t, S.Query("HelloWorld"), "value should be deleted from the set")

		time.Sleep(1 * time.Microsecond)
		S.Add("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "value should be reinserted to the set")
	})

	t.Run("Set Bias testing",func(t *testing.T){
		var S =NewLWWSet("ADD")
		S.Add("HelloWorld")
		S.Remove("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "set should have this value, because bias is ADD")

		S.bias = "REMOVE"
		require.False(t, S.Query("HelloWorld"), "set should have this value, because bias has been hacked to REMOVE")

		time.Sleep(1 * time.Microsecond)
		var T =NewLWWSet("REMOVE")
		T.Add("HelloWorld")
		T.Remove("HelloWorld")
		require.False(t, T.Query("HelloWorld"), "set should have this value, because bias is REMOVE")
	})

	t.Run("Set Comparison",func(t *testing.T){
		var R =NewLWWSet("ADD")
		R.Add("HelloWorld")
		R.Remove("HelloWorld")

		var S =NewLWWSet("REMOVE")
		S.Add("HelloWorld")
		S.Remove("HelloWorld")
		require.True(t, R.Compare(S, false) ,"they are same")
		require.False(t, R.Compare(S, true) ,"should not be same as they have different bias")

		var T =NewLWWSet("REMOVE")
		T.Add("HelloWorld")
		require.False(t, R.Compare(T, false) ,"they are different, their remove set is different")

		var T2 =NewLWWSet("REMOVE")
		T2.Add("HelloWorld")
		T2.Remove("HelloWorld2")
		require.False(t, R.Compare(T2, false) ,"they are different, their remove set is different")

		var U =NewLWWSet("REMOVE")
		U.Remove("HelloWorld")
		require.False(t, R.Compare(U, false) ,"they are different, their add set is different")

		var U2 =NewLWWSet("REMOVE")
		U2.Add("HelloWorld2")
		U2.Remove("HelloWorld")
		require.False(t, R.Compare(U2, false) ,"they are different, their add set is different")

		var V =NewLWWSet("REMOVE")
		V.Add(123)
		V.Remove(123)
		require.False(t, R.Compare(V, false) ,"they are different, they are completely different")
	})

	t.Run("Set Merging",func(t *testing.T){
		var R =NewLWWSet("ADD")
		R.Add("HelloWorld")
		R.Remove("HelloWorld")

		var U =NewLWWSet("REMOVE")
		U.Add(123)
		U.Remove(123)
		require.False(t, R.Compare(U, false) ,"they are completely different")

		var V =NewLWWSet("ADD")
		V.Add("HelloWorld")
		V.Remove("HelloWorld")
		V.Add(123)
		V.Remove(123)

		require.False(t, R.Compare(V, true) ,"they are completely different")
		R.Merge(U)
		require.True(t, R.Compare(V, true) ,"after merge, they should have same value")
	})
}