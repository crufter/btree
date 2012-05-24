Btree
=====

This is technically a B+Tree, but the btree name is easier on the tongoue :). Implemented in Go.

Caution
=====
This was the first project I have written in Go, thus probably not very idiomatic.

Examples
=====

import(
	"github.com/opesun/btree"
	"fmt"
)

type Int int
func (i Int) Less(c Comper) bool {
	a, ok := c.(Int)
	if !ok {
		return false
	}
	return i < a
}
func (i Int) Eq(c Comper) bool {
	a, ok := c.(Int)
	if !ok {
		return false
	}
	return i == a
}

//type Comper interface{
//	Less(Comper) bool
//	Eq(Comper) bool
//}

func main() {
	t := btree.NewBtree(50) // branching factor of the btree
	// Btree accepts Comper interface
	t.Insert(Int(8))
	fmt.Println(t.Find(Int(8)))
}

Credits
=====

- Thanks skelterjohn (John Asmuth) for pointing out that I don't need two different struct for the node and leaf, and also for helping me grasping some concepts of Go.