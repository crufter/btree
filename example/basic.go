package main

import(
    "github.com/opesun/btree"
    "fmt"
)

// Here we declare a type what the btree will happily accept.
// Note this design does not really allow more than one type in a btree, but thats the point.
// Of course, if you want, you can implement your very own crazy types which compare ints to string or things like that.
type Int int
func (i Int) Less(c btree.Comper) bool {
    a, ok := c.(Int)
    if !ok {
        return false
    }
    return i < a
}
func (i Int) Eq(c btree.Comper) bool {
    a, ok := c.(Int)
    if !ok {
        return false
    }
    return i == a
}

func main() {
    t, err := btree.New(50) // branching factor of the btree. For high performance 100 is optimal.
	if err != nil {
		panic(err)
	}
    t.Insert(Int(8))		// Btree accepts btee.Comper interface
    fmt.Println(
		t.Find(Int(8)),
		t.Delete(Int(8)),
		t.Find(Int(8)),
	)
}