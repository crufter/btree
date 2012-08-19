Btree
=====
This is technically a B+Tree, but the btree name is easier on the tongoue :). Implemented in Go.
The leaves are chained bidirectionally.

What's the point?
=====
I needed a sorted set, which internals I know and I can extend to my own database related needs.

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
Transform this into a counted b+tree.

Credits
=====
- Thanks skelterjohn (John Asmuth) for pointing out that I don't need two different struct for the node and leaf, and also for helping me grasping some concepts of Go.