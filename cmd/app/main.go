package main

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/lyrete/part-scraper/internal/bike24"
)

func main() {

	url := launcher.New().
		Headless(true).
		Devtools(true).
		MustLaunch()

	browser := rod.New().Timeout(time.Minute).
		ControlURL(url).
		Trace(true).SlowMotion(2 * time.Second).MustConnect()
	defer browser.MustClose()

	//r2Tab := r2.Connect(browser)
	//bdTab := bdiscount.Connect(browser)
	bike24Tab := bike24.Connect(browser)

	fmt.Println("Loop search for a few items...")

	itemsToSearch := []string{"Squirt Lube"}

	for _, itemName := range itemsToSearch {
		//item := r2.SearchForItem(r2Tab, itemName)
		//bd_item := bdiscount.SearchForItem(bdTab, itemName)
		bike24_item := bike24.SearchForItem(bike24Tab, itemName)
		//fmt.Println("B-D:", bd_item)
		//fmt.Println("R2:", item)
		fmt.Println("Bike24:", bike24_item)
	}

	browser.Close()

}
