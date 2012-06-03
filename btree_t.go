package main

import (
	"fmt"
	"github.com/opesun/btree"
	"math/rand"
	"time"
)

func leafTest(t *btree.Btree, arr *[][]btree.Comper) bool {
	var last btree.Comper
	for _, v1 := range *arr {
		if len(v1) < t.NodeSize()/2 {
			panic("Btree nodesize is too small.")
		}
		for _, v2 := range v1 {
			if last != nil && v2 != nil {
				if v2.Less(last) {
					panic("It's in descending order instead of ascending.")
				}
			} else {
				last = v2
			}
		}
	}
	return true
}

const tnum = 5000
const testnum = tnum
const perfnum = tnum * 20

func TestSmallOrder() {
	for i := 5; i < 10; i++ {
		s := btree.NewBtree(i)
		for j := 0; j < testnum/10; j++ {
			s.Insert(btree.Int(j))
		}
		for j := 0; j < testnum/10; j++ {
			if s.Find(btree.Int(j)) == false {
				panic("Shit happened.")
			}
		}
		s2 := btree.NewBtree(i)
		for j := testnum / 10; j >= 0; j-- {
			s2.Insert(btree.Int(j))
		}
		for j := testnum / 10; j >= 0; j-- {
			if s2.Find(btree.Int(j)) == false {
				panic("Shit happened.")
			}
		}
	}
}

func BenchmarkInsert() {
	order := 100
	p := btree.NewBtree(order)
	fmt.Println("Measuring performance. Order is ", order, ", iteration count is ", perfnum, ":") // (1 million inserts in C++ version with ints only takes 0.1 secs if order is 100.)
	tim := time.Now()
	for i := 0; i < perfnum; i++ {
		p.Insert(btree.Int(i))
	}
	fmt.Println("Took ", time.Since(tim))
	fmt.Println("Measuring find of every value:")
	tim = time.Now()
	for i := 0; i < perfnum; i++ {
		p.Find(btree.Int(i))
	}
	fmt.Println("Took ", time.Since(tim))
}

func TestInsert() {

}

func TestControlled() {
}

func TestBtree() {
	order := 100
	// Small order test
	a := btree.NewBtree(order)
	for i := 0; i < testnum; i++ {
		a.Insert(btree.Int(i))
	}
	for i := 0; i < testnum; i++ {
		if a.Find(btree.Int(i)) == false {
			fmt.Println(i)
			panic("We cant find something which we should definitely find.")
		}
	}
	//arr := a.GetAll()
	//leafTest(a, arr)
	b := btree.NewBtree(order)
	for i := testnum; i > 0; i-- {
		b.Insert(btree.Int(i))
	}
	for i := testnum; i > 0; i-- {
		if b.Find(btree.Int(i)) == false {
			fmt.Println(i)
			panic("We cant find something which we should definitely find.")
		}
	}
	//leafTest(b, arr)
	if a.TreeSize() != testnum {
		panic("Tree \"a\" size is not correct.")
	}
	if b.TreeSize() != testnum {
		panic("Tree \"b\" size is not correct.")
	}
	u := btree.NewBtree(order)
	for i := 0; i < testnum; i++ {
		u.Insert(btree.Int(i % 2))
	}
	for i := 0; i < 2; i++ {
		if u.Find(btree.Int(i)) == false {
			panic("Uniq stuff is not working...")
		}
	}
	fmt.Println("Insert duplicates, then delete all of them and cry if any not found...")
	for ord := 5; ord < 43; ord++ {
		dup_inc := btree.NewBtree(ord)
		dup_dec := btree.NewBtree(ord)
		for mod := 2; mod < 47; mod++ {
			for k := 0; k <= testnum/5; k++ {
				dup_inc.Insert(btree.Int(k % mod))
			}
			for k := 0; k <= testnum/5; k++ {
				if dup_inc.Delete(btree.Int(k%mod)) == false {
					panic("Dupe test failed.")
				}
			}
			for k := testnum / 5; k >= 0; k-- {
				dup_dec.Insert(btree.Int(k % mod))
			}
			for k := testnum / 5; k >= 0; k-- {
				if dup_dec.Delete(btree.Int(k%mod)) == false {
					panic("Dupe test failed.")
				}
			}
		}
	}
	fmt.Println("Doing some slightly controlled stress test.")
	for ord := 5; ord <= 31; ord++ {
		m := make(map[int]int)
		tes := btree.NewBtree(5)
		deletion_count := 0
		for i := 0; i <= tnum*200; i++ {
			c := rand.Int() % (tnum * 15)
			v, ok := m[c]
			if ok != tes.Find(btree.Int(c)) {
				fmt.Println(v, ok)
				panic("Real bad, containment does not match with map.")
			}
			if ok == true {
				m[c]++
			} else {
				m[c] = 1
			}
			tes.Insert(btree.Int(c))

			c = rand.Int() % (tnum * 15)
			v, ok = m[c]
			if ok != tes.Find(btree.Int(c)) {
				panic("Bad.")
			}
			if ok == true {
				if v == 1 {
					delete(m, c)
				} else {
					m[c]--
				}
				if tes.Delete(btree.Int(c)) == false {
					panic("Can not found.")
				}
				deletion_count++
			}
		}
		fmt.Println("	Ran with order ", ord, " deleted ", deletion_count)
	}
	fmt.Println("Test ran successfully.")
}

func main() {
	fmt.Println("\nThis test will take approx 5 minutes.\n")
	TestSmallOrder()
	BenchmarkInsert()
	TestBtree()
}
