package crdt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLWWSet(t *testing.T) {

	t.Run("Set creation",func(t *testing.T){
		var S = NewLWWSet("ADD")
		require.Empty(t, S.addSet, "add set is not empty")
		require.Empty(t, S.removeSet, "remove set is not empty")
		require.Equal(t, S.bias, "ADD","bias is not correct")

		var T = NewLWWSet("REMOVE")
		require.Equal(t, T.bias, "REMOVE","bias is not correct")

		var U = NewLWWSet("")
		require.Equal(t, U.bias, "ADD","default bias is not correct")
	})

	t.Run("Set empty content",func(t *testing.T){
		var S = NewLWWSet("ADD")
		require.Empty(t, S.GetAddSet(), "add set is not empty")
		require.Empty(t, S.GetRemoveSet(), "remove set is not empty")
		require.Equal(t, S.GetBias(), "ADD","bias is not correct")
		require.Empty(t, S.GetContent(), "ADD","bias is not correct")
	})

	t.Run("Set addition, deletion, query and content ",func(t *testing.T){
		// init
		var S = NewLWWSet("ADD")
		require.False(t, S.Query("HelloWorld"), "set should be empty")
		require.Equal(t, 0, len(S.GetAddSet()), "add should be empty")
		require.Equal(t, 0, len(S.GetRemoveSet()), "remove should be empty")
		require.Equal(t, 0, len(S.GetContent()), "add should be empty")

		// addition
		time.Sleep(1 * time.Microsecond)
		S.Add("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "set should have this value")
		require.Equal(t, 1, len(S.GetAddSet()), "add should have 1 value")
		require.Equal(t, 0, len(S.GetRemoveSet()), "remove should be empty")
		require.Equal(t, 1, len(S.GetContent()), "content should exist")
		require.NotNil(t, S.GetContent()["HelloWorld"], "content should exist")

		// deletion
		time.Sleep(1 * time.Microsecond)
		S.Remove("HelloWorld")
		require.False(t, S.Query("HelloWorld"), "value should be deleted from the set")
		require.Equal(t, 1, len(S.GetAddSet()), "add should have 1 value")
		require.Equal(t, 1, len(S.GetRemoveSet()), "remove have this value")
		require.Equal(t, 0, len(S.GetContent()), "content should be empty")

		// reinsertion
		time.Sleep(1 * time.Microsecond)
		S.Add("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "value should be reinserted to the set")
		require.Equal(t, 1, len(S.GetAddSet()), "add should have 1 value")
		require.Equal(t, 1, len(S.GetRemoveSet()), "remove should be empty")
		require.Equal(t, 1, len(S.GetContent()), "content should exist")
		require.NotNil(t, S.GetContent()["HelloWorld"], "content should exist")
	})

	t.Run("Set Bias testing",func(t *testing.T){
		var S = NewLWWSet("ADD")
		S.Add("HelloWorld")
		S.Remove("HelloWorld")
		require.True(t, S.Query("HelloWorld"), "set should have this value, because bias is ADD")

		S.bias = "REMOVE"
		require.False(t, S.Query("HelloWorld"), "set should have this value, because bias has been hacked to REMOVE")

		time.Sleep(1 * time.Microsecond)
		var T = NewLWWSet("REMOVE")
		T.Add("HelloWorld")
		T.Remove("HelloWorld")
		require.False(t, T.Query("HelloWorld"), "set should have this value, because bias is REMOVE")
	})

	t.Run("Set Comparison",func(t *testing.T){
		var R = NewLWWSet("ADD")
		R.Add("HelloWorld")
		R.Remove("HelloWorld")

		var S = NewLWWSet("REMOVE")
		S.Add("HelloWorld")
		S.Remove("HelloWorld")

		var T = NewLWWSet("ADD")
		T.Add("HelloWorld")
		T.Remove("HelloWorld")

		var U = NewLWWSet("REMOVE")
		U.Add("HelloWorld")
		U.Remove("HelloWorld")

		var V = NewLWWSet("REMOVE")
		V.Add("HelloWorld")

		var X = NewLWWSet("REMOVE")
		X.Remove("HelloWorld")

		var Y = NewLWWSet("REMOVE")
		Y.Add("HelloWorld")
		Y.Remove("HelloWorld2")

		var Z = NewLWWSet("REMOVE")
		Z.Add("HelloWorld2")
		Z.Remove("HelloWorld")

		time.Sleep(1 * time.Microsecond)
		T.Add("HelloWorld")
		U.Add("HelloWorld")

		// Historic comparison
		require.True(t, R.Compare(S, false) ,"they are same")
		require.False(t, R.Compare(S, true) ,"should not be same as they have different bias")

		require.True(t, T.Compare(U, false) ,"they are same")
		require.False(t, T.Compare(U, true) ,"should not be same as they have different bias")

		require.False(t, R.Compare(V, false) ,"they have different add set")
		require.False(t, R.Compare(X, false) ,"they have different remove set")
		require.False(t, R.Compare(Y, false) ,"they have different add set")
		require.False(t, R.Compare(Z, false) ,"they have different remove set")

		// Content comparison
		require.True(t, R.CompareContent(T, false) ,"they have same content")
		require.False(t, R.CompareContent(U, true) ,"should not be same as they have different bias")
		require.False(t, R.CompareContent(X, false) ,"they have different content")
		require.False(t, R.CompareContent(Z, false) ,"they have different content")
	})

	t.Run("Set Merging",func(t *testing.T){
		var R = NewLWWSet("ADD")
		R.Add("HelloWorld")
		R.Remove("HelloWorld")

		var U = NewLWWSet("REMOVE")
		U.Add(123)
		U.Remove(123)
		require.False(t, R.Compare(U, false) ,"they are completely different")

		var V = NewLWWSet("ADD")
		V.Add("HelloWorld")
		V.Remove("HelloWorld")
		V.Add(123)
		V.Remove(123)

		// Merge
		var W = Merge(R,U)
		require.True(t, V.Compare(W, true) ,"after merge, they should have same value")

		// Merge With
		R.MergeWith(U)
		require.True(t, V.Compare(R, true) ,"after merge, they should have same value")

		R.Add("yolo")
	})
}