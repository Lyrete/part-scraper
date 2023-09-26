package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/lyrete/part-scraper/internal/r2"
)

func main() {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	fmt.Println("Loop search for a few items...")

	itemsToSearch := []string{"SRAM XG-1275", "SRAM XG-1295", "GX Eagle Chain", "DT Swiss GR1600"}

	for _, itemName := range itemsToSearch {
		item := r2.SearchItemByName(browser, itemName)
		fmt.Println(item)
	}

	browser.Close()

}
