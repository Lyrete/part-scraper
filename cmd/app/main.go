package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/lyrete/part-scraper/internal/bdiscount"
	"github.com/lyrete/part-scraper/internal/r2"
)

func main() {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	r2Tab := r2.Connect(browser)
	bdTab := bdiscount.Connect(browser)

	fmt.Println("Loop search for a few items...")

	itemsToSearch := []string{"6009685090003"}

	for _, itemName := range itemsToSearch {
		item := r2.SearchForItem(r2Tab, itemName)
		bd_item := bdiscount.SearchForItem(bdTab, itemName)
		fmt.Println("B-D:", bd_item)
		fmt.Println("R2:", item)
	}

	browser.Close()

}
