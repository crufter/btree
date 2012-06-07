Btree
=====
This is technically a B+Tree, but the btree name is easier on the tongoue :). Implemented in Go.
The leaves are chained bidirectionally.

Caution
=====
This was the first project I have written in Go, thus probably not very idiomatic.

What's the point?
=====
Of course, the functionality this pkg provides is almost the same as the bultin map's.
But I have written this pkg for myself mostly to be able to tweak internals, provide any type which implements a given interface,
implement this as a counted B+tree (not ready yet, see below), and be able to iterate in order.

For now use the source more like an algorithm than an all around pkg solution. (This will change soon)

Is it threadsafe?
=====
No.

Examples
=====
```
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
```

What's up with those panics in the source code, a Go pkg shall not panic!
=====
The fact is, the pkg only panics if there is a bug in the algorithm, and it's better for you, me, and the globe if you notice it.

Future
=====
There is a missing functionality in the tree what I want to implement soon... the FindXth() method, but that requires a counted B+tree.
All the building blocks are there to make this change, I have just not found the time yet.

Credits
=====
- Thanks skelterjohn (John Asmuth) for pointing out that I don't need two different struct for the node and leaf, and also for helping me grasping some concepts of Go.