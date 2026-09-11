package main

import (
	"fmt"

	"github.com/rituraj12797/cats-db/components/skiplist"
)

func main() {
	x := skiplist.NewSkipList[int, int]()

	x.Insert(5, 5)
	x.Insert(6, 6)

	fmt.Print(*(x.Search(5)))
}
