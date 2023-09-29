package r2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	utils "github.com/lyrete/part-scraper/internal"
	"github.com/lyrete/part-scraper/internal/bdiscount"
)

type Item struct {
	name    string
	price   float64
	barcode string
	URL     string
}

func GetPriceFromPage(page *rod.Page, vat int) float64 {
	price_wrapper := page.MustElement("#product-offer > div.product-info.col-xs-12.col-sm-6.col-lg-6 > div > div.product-offer > div.row > div > div") //TODO: Replace this selector, it's probably 100% gonna break at some point when they add like a random deal
	itemPriceElement := price_wrapper.MustElement("strong")

	itemPrice := itemPriceElement.MustElementR("span", "€").MustText()
	itemPrice = strings.Replace(itemPrice, " €", "", 1)
	itemPrice = strings.Replace(itemPrice, ",", ".", 1)

	itemPriceFloat, err := strconv.ParseFloat(itemPrice, 64)

	if err != nil {
		panic(err)
	}

	itemPriceVatFree := itemPriceFloat / (1 + (float64(vat) / 100))
	itemPriceVatFree = bdiscount.RoundFloat(itemPriceVatFree, 2)

	return itemPriceVatFree
}

func GetVatFromPage(page *rod.Page) int {
	vatElementText := page.MustElement(".vat_info").MustText()

	//Replace nbsp from the text as it's very likely we'll have some in there
	vatElementText = strings.Replace(vatElementText, "\u00a0", " ", -1)

	vatPercentage := strings.Replace(strings.Split(vatElementText, " ")[1], "%", "", 1)

	vatPercentageInt, err := strconv.Atoi(vatPercentage)

	if err != nil {
		panic(err) //TODO: Actually handle the error, (probably before this make sure it is a number or try to search for it elsewhere otherwise)
	}

	return vatPercentageInt
}

func GetEanFromPage(page *rod.Page) string {
	eanText := page.MustElement("#tab-description").MustElementR(".attr-characteristic", "EAN:").MustText()

	eanNumber := strings.Split(eanText, "\n")[1]

	return eanNumber
}

func SearchForItem(page *rod.Page, itemName string) Item {
	fmt.Println("Searching for " + itemName)

	wait := page.MustWaitRequestIdle()
	page.MustElement("input[name='qs']").MustInput(itemName).MustType(input.Enter)
	wait()

	page.MustElement("#product-list").MustElement(".image-wrapper").MustClick()

	vatPercentage := GetVatFromPage(page)

	item_price := GetPriceFromPage(page, vatPercentage)

	eanNumber := GetEanFromPage(page)

	return Item{name: itemName, price: item_price, barcode: eanNumber, URL: utils.GetUrl(page)}

}

func Connect(browser *rod.Browser) *rod.Page {
	fmt.Println("Connecting to r2-bike.com")
	page := browser.MustPage("https://www.r2-bike.com/").MustWaitStable()

	return page
}
