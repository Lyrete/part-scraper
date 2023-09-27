package bdiscount

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type Item struct {
	name  string
	price float64
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetPriceFromPage(page *rod.Page, vat int) float64 {
	itemPrice := page.MustElementR("#netz-price", "€").MustText()
	fmt.Println(itemPrice)

	itemPrice = strings.Replace(itemPrice, "€", "", 1)

	itemPriceFloat, err := strconv.ParseFloat(itemPrice, 64)

	if err != nil {
		panic(err)
	}

	itemPriceVatFree := itemPriceFloat / (1 + (float64(vat) / 100))
	itemPriceVatFree = RoundFloat(itemPriceVatFree, 2)

	return itemPriceVatFree
}

func GetVatFromPage(page *rod.Page) int {
	fmt.Println("Getting Bike-Discount VAT")
	vatElement, err := page.Element(".product--tax")

	if err != nil {
		panic("Couldn't find item tax element")
	}

	vatElementText := vatElement.MustText()

	//Replace nbsp from the text as it's very likely we'll have some in there
	vatElementText = strings.Replace(vatElementText, "\u00a0", " ", -1)

	vatPercentage := strings.Replace(strings.Split(vatElementText, " ")[1], "%", "", 1)

	vatPercentageInt, err := strconv.Atoi(vatPercentage)

	if err != nil {
		panic(err) //TODO: Actually handle the error, (probably before this make sure it is a number or try to search for it elsewhere otherwise)
	}

	return vatPercentageInt
}

func SearchItemByName(page *rod.Page, itemName string) Item {

	fmt.Println("Searching for " + itemName + " on Bike-Discount")

	wait := page.MustWaitRequestIdle()
	page.MustElement("input[name='sSearch']").MustInput(itemName).MustType(input.Enter)
	wait()

	page.MustElement(".search--results").MustElement(".product--image").MustClick()

	vatPercentage := GetVatFromPage(page)

	item_price := GetPriceFromPage(page, vatPercentage)

	return Item{name: itemName, price: item_price}

}

func Connect(browser *rod.Browser) *rod.Page {
	fmt.Println("Connecting to Bike-Discount")
	page := browser.MustPage("https://bike-discount.de/en/").MustWaitStable()

	return page
}