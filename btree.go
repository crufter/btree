// B+tree implementation by Dobronszki János (John Dobronszki) @ OPESUN TECHNOLOGIES Kft. (Opesun Technologies Ltd.) 2012
package btree

import( 
	//"sort"
)

	func newLeaf(nodesize int) *Node{
		return &Node{values: make([]Comper, nodesize)}
	}
	
	func newNode(nodesize int) *Node{
		return &Node{values: make([]Comper, nodesize), pointers: make([]*Node, nodesize+1)}
	}

func NewBtree(nodesize int) *Btree {
	if nodesize < 5 {
		panic("Node size must be at least 3 because of insert, 5 because of deletion, due to the characteristics of the implementation.")
	}
	nt := new(Btree)
	nt.root = newLeaf(nodesize)
	nt.nodesize = nodesize
	nt.duplicates_allowed = true
	return nt
}

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

type Comper interface{
	Less(Comper) bool
	Eq(Comper) bool
}

type Btree struct {
	root *Node
	typ string
	duplicates_allowed bool
	nodesize, size int
}
			
			func setParents(fresh *Node) {
				l := fresh.size
				for i:=0;i<=l;i++ {
					fresh.pointers[i].parent = fresh
				}
			}
		
		// We must go over the moved ones and set the parent pointer of all.
		// We move count amount of values and pointers from the end of from to the empty to.
		func moveToNode(from *Node, to *Node, count int) {
			l := 0
			for i:=from.size - 1; i >= from.size - count; i-- {
				to.values[count - l - 1] = from.values[i]
				to.pointers[count - l] = from.pointers[i+1]
				from.values[i] = nil
				from.pointers[i+1] = nil
				l++
			}
			from.size -= count
			to.size = count
			to.pointers[0] = from.pointers[from.size]
			from.pointers[from.size] = nil
			from.values[from.size-1] = nil
			from.size--
			setParents(to)
		}

	func (tree *Btree) splitNode(parentPos int, old_node *Node, val Comper) *Node {
		new_node := newNode(tree.nodesize)
		move_count := old_node.size/2
		moveToNode(old_node, new_node, move_count)
		median := getLeftmost(new_node).values[0]
		if old_node.parent == nil {
			tree.root = newRoot(tree.nodesize, old_node, new_node, median)	// 
			old_node.parent = tree.root
			new_node.parent = tree.root
		} else {
			toNodeMiddle(old_node.parent, parentPos, median, new_node)
			new_node.parent = old_node.parent
		}
	
		if val.Less(median) {
			return old_node
		}
		return new_node
	}
	
		// We move count amount of values from the end of from to the beginning of the empty to.
		func moveToLeaf(from *Node, to *Node, count int) {
			l := 0
			for i:=from.size - 1; i >= from.size - count; i-- {
				to.values[count - l - 1] = from.values[i]
				from.values[i] = nil
				l++
			}
			from.size -= count
			to.size = count
		}
	
		func newRoot(size int, old *Node, fresh *Node, midval Comper) *Node {
			nr := newNode(size)
			nr.values[0] = midval
			nr.pointers[0] = old
			nr.pointers[1] = fresh
			nr.size = 1
			return nr
		}
	
	func (tree *Btree) splitLeaf(parentPos int, old_leaf *Node, val Comper) *Node {
		new_leaf := newLeaf(tree.nodesize)
		move_count := old_leaf.size/2
		moveToLeaf(old_leaf, new_leaf, move_count)
		
		if old_leaf.next != nil {
			old_leaf.next.prev = new_leaf
			new_leaf.next = old_leaf.next
		}
		old_leaf.next = new_leaf
		new_leaf.prev = old_leaf
		
		if old_leaf.parent == nil {			// Was a leaf root...
			tree.root = newRoot(tree.nodesize, old_leaf, new_leaf, new_leaf.values[0])
			old_leaf.parent = tree.root
			new_leaf.parent = tree.root
		} else {
			toNodeMiddle(old_leaf.parent, parentPos, new_leaf.values[0], new_leaf)
			new_leaf.parent = old_leaf.parent
		}
		// If the binary search positions the equal values after the existing ones, then the median (in the parent node) will be the right leaf's first value,
		// otherwise it will be the last value of the left leaf.
		if val.Less(new_leaf.values[0]) {
			return old_leaf
		}
		return new_leaf
	}
		
		// See *1
		func toNodeMiddle(node *Node, pos int, midval Comper, pointer *Node) {
			for c := node.size; c >= 0; c-- {
				if c > pos {
					node.values[c] = node.values[c-1]
					node.pointers[c+1] = node.pointers[c]
				} else {
					node.values[c] = midval
					node.pointers[c+1] = pointer
					break
				}
			}
			node.size++
		}
		
		// *1 If it panics here with index out of bound, then the leafsplit is not happening. God save us all.
		func toLeafMiddle(leaf *Node, pos int, val Comper) {
			for c := leaf.size; c >= 0; c-- {
				if c > pos {
					leaf.values[c] = leaf.values[c-1]
				} else {
					leaf.values[c] = val
					break
				}
			}
		}
			
	func insertToLeaf(leaf *Node, val Comper) {
		pos := findPos(leaf.size, leaf.values, val)
		if leaf.values[pos] == nil {
			leaf.values[pos] = val
		} else {
			toLeafMiddle(leaf, pos, val)
		}
		leaf.size++
	}
	
	// "Handlined" here the binary search of sort pkg because the higher order function construct was real slow.
	func findPos(n int, values []Comper, val Comper) int{
		// Define f(-1) == false and f(n) == true.
		// Invariant: f(i-1) == false, f(j) == true.
		i, j := 0, n
		for i < j {
			h := i + (j-i)/2 // avoid overflow when computing h
			// i ≤ h < j
			if !val.Less(values[h]) {
				i = h + 1 // preserves f(i-1) == false
			} else {
				j = h // preserves f(j) == true
			}
		}
		// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
		return i
	}
	//func findPos(size int, values []Comper, val Comper) int{
	//	return sort.Search(size,
	//		func(i int) bool {
	//			return val.Less(values[i])
	//		})
	//}
	
func (tree *Btree) Insert(val Comper) int {
	var leaf *Node;
	var c *Node = tree.root
	var p int
	for {
		if c.isNode() {
			if c.size >= tree.nodesize {
				c = tree.splitNode(p, c, val)
			}
			p = findPos(c.size, c.values, val)
			c = c.pointers[p]
			if c == nil {
				panic("This should definitely not happen.")
			}
		} else {	// Leaf
			leaf = c; 
			if	leaf.size >= tree.nodesize {
				leaf = tree.splitLeaf(p, leaf, val)
			}
			break
		}
	}
	insertToLeaf(leaf, val)
	tree.size++
	return 1
}
	

func (tree *Btree) Find(val Comper) bool {
	var leaf *Node;
	var c *Node = tree.root
	var p int
	for {
		if c.isNode() {
			p = findPos(c.size, c.values, val)
			c = c.pointers[p]
			if c == nil {
				panic("Woowoowoo... calm down.")
			}
		} else {
			leaf = c; 
			break
		}
	}
	p = findPos(leaf.size, leaf.values, val)
	if p > 0 {
		return val.Eq(leaf.values[p-1])
	}
	return false
}
		// p == a is the index of the pointer in the parent node, which is pointing to the right node
		func delFromNode(node *Node, p int) {
			s := node.size
			for i:=p;i<s;i++ {
				node.pointers[i] = node.pointers[i+1]
				node.values[i-1] = node.values[i]
			}
			node.pointers[s] = nil
			node.values[s-1] = nil
			node.size--
		}
		
	// p == a is the index of the pointer pointing to left
	// We must choose a new median in the merged node, between the two neighbouring pointer (if we merge nodes).
	// We must not forget to set the next and prev pointers.
	func (tree *Btree) merge(left *Node, right *Node, val Comper, p int) (*Node, int) {
		//fmt.Println("merging. isnode: ", left.isNode(), left, right)
		ls := left.size
		if left.isNode() {	// Nodes
			for i:=0;i<right.size;i++ {
				left.values[ls+1+i] = right.values[i]
				left.pointers[ls+1+i] = right.pointers[i]
			}
			left.pointers[ls+1+right.size] = right.pointers[right.size]
			left.values[ls] = getLeftmost(right).values[0]
			left.size+=right.size+1		// because of the inserted median it became bigger by one
		} else {	// Leaves
			for i:=0;i<right.size;i++{
				left.values[ls+i] = right.values[i]
			}
			left.size+=right.size
		}
		if left.parent.size == 1 && left.parent.parent == nil {
			tree.root = left; left.prev = nil; left.next = nil
			left.parent = nil
		} else {
			delFromNode(left.parent, p+1)
		}
		if left.isNode() {
			setParents(left)
		} else {
			left.next = right.next
			if right.next != nil {
				right.next.prev = left
			}
		}
		return left, p
		
	}
		
		// Create space at the beginning of the node.
		func createSpaceN(node *Node, c int) {
			copy(node.values[c:],node.values[0:])
			copy(node.pointers[c:],node.pointers[0:])
			// It is pointless to nil the others, since they will be overwritten anyway.
		}
		
		// Create space at the beginning of the leaf.
		func createSpaceL(node *Node, c int) {
			copy(node.values[c:],node.values[0:])
			// It is pointless to nil the others, since they will be overwritten anyway.
		}
		
		// Move from left to right
		// We can copy amt amount of pointers but only amt - 1 amount of values... The 1 value staying in the left branch is becoming redundant, hence we delete it.
		func moveFromEndN(left *Node, right *Node, amt int) {
			//fmt.Println(amt)
			//fmt.Println("wooo", left, right)
			copy(right.values, left.values[left.size-amt+1:left.size])
			copy(right.pointers, left.pointers[left.size-amt+1:left.size+1])
			//fmt.Println("wooo", left, right)
			right.values[amt-1] = getLeftmost(right.pointers[amt]).values[0]		// median
			// Maybe we whould nil lef here, not untirely sure yet...
		}
		
		func moveFromEndL(left *Node, right *Node, amt int) {
			copy(right.values, left.values[left.size-amt:left.size])
		}
		
	// We must put to parent.values[p] position the smallest of parent.pointers[p+1] branch.
	func (tree *Btree) spillToRight(left *Node, right *Node, val Comper, p int) (*Node, int) {
		//fmt.Println("spilling to right, isnode: ", left.isNode(), left, right)
		spill_amt := (left.size-right.size)/2
		if spill_amt <= 0 {
			panic("This is a bug :S.")
		}
		if left.isNode() {
			createSpaceN(right, spill_amt)
			moveFromEndN(left, right, spill_amt)
		} else {
			createSpaceL(right, spill_amt)
			moveFromEndL(left, right, spill_amt)
		}
		left.size-=spill_amt
		right.size+=spill_amt
		parent := left.parent
		median := getLeftmost(right).values[0]
		parent.values[p] = median
		if left.isNode() {
			setParents(right)
		}
		if val.Less(median) {
			return left, p
		}
		return right, p+1
	}
		
		func copyFromBegN(left *Node, right *Node, amt int) {
			copy(left.values[left.size+1:],right.values[0:amt-1])
			copy(left.pointers[left.size+1:],right.pointers[:amt])
			left.values[left.size] = getLeftmost(right.pointers[0]).values[0]
		}
	
		func copyFromBegL(left *Node, right *Node, amt int) {
			copy(left.values[left.size:],right.values[0:amt])
		}
		
		func delBegN(node *Node, amt int) {
			copy(node.values[0:], node.values[amt:])
			copy(node.pointers[0:], node.pointers[amt:])
		}
		
		func delBegL(node *Node, amt int) {
			copy(node.values[0:], node.values[amt:])
		}
		
	//a parent.values[p] pozícióba parent.pointers[p+1] branch legkissebbjét kell választani mediannak
	func (tree *Btree) spillToLeft(left *Node, right *Node, val Comper, p int) (*Node, int) {
		spill_amt := (right.size-left.size)/2
		//fmt.Println("spilling to left", spill_amt, right.size, left.size, left, right)
		if spill_amt<= 0 {
			panic("Please turn off your computer before it explodes.")
		}
		if left.isNode() {
			copyFromBegN(left, right, spill_amt)
			delBegN(right, spill_amt)
		} else {
			copyFromBegL(left, right, spill_amt)
			delBegL(right, spill_amt)
		}
		left.size+=spill_amt
		right.size-=spill_amt
		//fmt.Println("__", right)
		parent := left.parent
		median:= getLeftmost(right).values[0]
		parent.values[p] = median
		if left.isNode() {
			setParents(left)
		}
		if val.Less(median) {
			return left, p
		}
		return right, p+1
	}
	
	func (tree *Btree) balance(node *Node, p int, val Comper) (*Node, int) {
		to_ret := node
		s := node.size
		ps := node.parent.size
		half_o := tree.nodesize/2
		//fmt.Println(p, node, node.size, node.parent, node.parent.size)
		if p>ps {
			panic("Houston, we've got a problem here.")
		}
		if node.isNode() {
			//fmt.Println("!!!!!!!!!!!!", node)
		}
		switch {
		case p == ps:
			left_s := node.parent.pointers[p-1].size
			switch {
			case left_s > half_o:
				to_ret, p = tree.spillToRight(node.parent.pointers[p-1], node, val, p-1)
			case left_s + s + 1< tree.nodesize:							// Maybe we should fine to the number
				to_ret, p = tree.merge(node.parent.pointers[p-1], node, val, p-1)
			}
		case p == 0:
			right_s := node.parent.pointers[p+1].size
			switch {
			case right_s > half_o:
				to_ret, p = tree.spillToLeft(node, node.parent.pointers[p+1], val, p)
			case right_s + s + 1< tree.nodesize:
				to_ret, p = tree.merge(node, node.parent.pointers[p+1], val, p)
			}
		default:
			left_s := node.parent.pointers[p-1].size
			right_s := node.parent.pointers[p+1].size
			switch {
				case right_s > half_o:
					to_ret, p = tree.spillToLeft(node, node.parent.pointers[p+1], val, p)
				case right_s + s + 1< tree.nodesize:
					to_ret, p = tree.merge(node, node.parent.pointers[p+1], val, p)
				// If we get here that means it's not worth  to merge, because right == half...
				case left_s > half_o:
					to_ret, p = tree.spillToRight(node.parent.pointers[p-1], node, val, p-1)
				case left_s + s + 1< tree.nodesize:
					to_ret, p = tree.merge(node.parent.pointers[p-1], node, val, p-1)
			}
		}
		return to_ret, p
	}
	
	func delFromLeaf(leaf *Node, p int) {
		copy(leaf.values[p:],leaf.values[p+1:])
	}
	
	func runUpAndCorrigate(leaf *Node, dead_val Comper, new_val Comper, pst []int) bool {
		l := 0
		c := leaf.parent
		for {
			if c != nil {
				cp := pst[(len(pst)-1)-l]
				if cp == 0 {
					c = c.parent
				} else if c.values[cp-1].Eq(dead_val) {
					c.values[cp-1] = new_val
					return true
				} else {
					break
				}
				l++
			} else {
				break
			}
		}
		panic("Ohohoo, we should definitely find but we did not.")
		return false
	}
	
func (tree *Btree) Delete(val Comper) bool {
	var leaf *Node
	var c *Node = tree.root
	pstack := make([]int, 0, 100)
	var p, np int
	for {
		if c.isNode() {
			if c.parent != nil && c.size < tree.nodesize/2 {		// If p nil, then root // Then min 5, think about it again
				c, np = tree.balance(c, p, val)
				pstack[len(pstack)-1] = np
			}
			p = findPos(c.size, c.values, val)
			pstack = append(pstack, p)
			c = c.pointers[p]
			if c == nil {
				panic("This should frankly never happen.")
			}
		} else {
			leaf = c;
			if c.parent != nil && leaf.size < tree.nodesize/2 {
				leaf, _ = tree.balance(leaf, p, val)
			}
			break
		}
	}
	p = findPos(leaf.size, leaf.values, val)
	if p > 0 {
		if val.Eq(leaf.values[p-1]) {
			delFromLeaf(leaf, p-1)
			leaf.size--
			if tree.duplicates_allowed == true && p == 1 && leaf.prev != nil && val.Eq(leaf.prev.values[leaf.prev.size-1]) {
				runUpAndCorrigate(leaf, val, leaf.values[0], pstack)
			}
			return true
		} else {
			return false
		}
	}
	return false	// Didn't find
}

	func getLeftmost(c *Node) *Node {
		for {
			if c.isNode() {
				c = c.pointers[0]
			} else {
				break;
			}
		}
		return c
	}
	
	func getRightmost(c *Node) *Node {
		for {
			if c.isNode() {
				c = c.pointers[len(c.pointers)-1]
			} else {
				break;
			}
		}
		return c
	}

type Node struct{
	parent, prev, next 		*Node		// pointer
	size 					int
	values					[]Comper
	pointers				[]*Node
}
func (n *Node) isNode() bool {
	return len(n.pointers)!=0
}

func (t *Btree) NodeSize() int {
	return t.nodesize
}
func (t *Btree) TreeSize() int {
	return t.size
}