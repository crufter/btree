package btree

import(
	"fmt"
)

	func (tree *Btree) GetAll() *[][]Comper {
		arr := make([][]Comper, 0)
		c := getLeftmost(tree.root)
		for {
			arr = append(arr, c.values)
			if c.next != nil {
				c = c.next
			} else {
				break
			}
		}
		return &arr
	}

func (tree *Btree) PrintAll() {
	arr := *tree.GetAll()
	for _, v1 := range arr{
		for _, v2 := range v1 {
			if v2 != nil {
				fmt.Print(v2)
			} else {
				fmt.Print("_")
			}
			fmt.Print(" ")
		}
		fmt.Print("     ")
	}
}

		func getLevStr(lev int) string{
			str := ""
			for i:=0;i<lev;i++ {
				str += "   "
			}
			return str
		}
		
	func visualize(p *Node, lev int) {
		if p.isNode() {
			fmt.Print(getLevStr(lev))
			for _, v := range p.values {
				fmt.Print(" ", v)
			}
			for _, v := range p.pointers {
				if v != nil {
					fmt.Print("/")
				}
			}
			fmt.Print(" (", p.size, ")")
			fmt.Println("")
			for _, v := range p.pointers {
				if v != nil {
					visualize(v, lev + 1)
				}
			}
		} else {
			l := p
			fmt.Print(getLevStr(lev))
			fmt.Print("-")
			for _, v := range l.values {
				fmt.Print(" ", v)
			}
			fmt.Print(" (", l.size, ")")
			fmt.Println("")
		}
	}
	
func (t *Btree) Visualize() {
	visualize(t.root, 0)
	fmt.Println(" ")
}