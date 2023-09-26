package main

import (
	"fmt"

	"github.com/lyrete/part-scraper/internal/r2"
)

func main() {

	fmt.Println("Loop search for a few items...")
	fmt.Println()

	itemsToSearch := []string{"SRAM XG-1275", "SRAM XG-1295", "GX Eagle Chain", "DT Swiss GR1600"}

	for _, itemName := range itemsToSearch {
		item := r2.SearchItemByName(itemName)
		fmt.Println(item)
	}

}
