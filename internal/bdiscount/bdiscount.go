package bdiscount

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	utils "github.com/lyrete/part-scraper/internal"
)

type Item struct {
	name    string
	price   float64
	barcode string
	URL     string
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetPriceFromPage(page *rod.Page, vat int) float64 {
	itemPrice := page.MustElementR("#netz-price", "€").MustText()

	itemPrice = strings.Replace(itemPrice, "€", "", 1)
	itemPrice = strings.Replace(itemPrice, ",", ".", 1)
	itemPrice = strings.Trim(itemPrice, " ")

	itemPriceFloat, err := strconv.ParseFloat(itemPrice, 64)

	if err != nil {
		panic(err)
	}

	itemPriceVatFree := itemPriceFloat / (1 + (float64(vat) / 100))
	itemPriceVatFree = RoundFloat(itemPriceVatFree, 4)

	return itemPriceVatFree
}

func GetVatFromPage(page *rod.Page) int {
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

func getBarcodeFromPage(page *rod.Page) string {
	barcodeElement := page.MustElements(".netz-ean, .netz-upc").First()
	if barcodeElement == nil {
		return "NOT_FOUND"
	}
	return barcodeElement.MustText()
}

func isBarcode(searchTerm string) bool {
	_, err := strconv.Atoi(searchTerm)
	if err != nil {
		return false
	}
	return true
}

func getFullName(page *rod.Page) string {
	return page.MustElement("h1.product--title").MustText()
}

func SearchForItem(page *rod.Page, searchTerm string) Item {

	fmt.Println("Searching for " + searchTerm + " on Bike-Discount")

	wait := page.MustWaitRequestIdle()
	page.MustElement("input[name='sSearch']").MustInput(searchTerm).MustType(input.Enter)
	wait()

	page.MustElement(".search--results").MustElement(".product--image").MustClick()

	fullName := getFullName(page)

	if page.MustHas(".custom-variants--select--box") {
		page.MustElement(".custom-variants--select--box").MustClick()
		if isBarcode(searchTerm) {
			selector := "input[ean*=\"" + searchTerm + "\"]"
			variant := page.MustElement(".product--configurator").MustElement(selector)
			variant.MustClick()

		} else {
			page.MustElement(".product--configurator").MustElement(".variant--option").MustClick()
		}
	}

	vatPercentage := GetVatFromPage(page)

	item_price := GetPriceFromPage(page, vatPercentage)

	return Item{name: fullName, price: item_price, barcode: getBarcodeFromPage(page), URL: utils.GetUrl(page)}

}

func Connect(browser *rod.Browser) *rod.Page {
	fmt.Println("Connecting to Bike-Discount")
	page := browser.MustPage("https://bike-discount.de/en/").MustWaitStable()

	return page
}
