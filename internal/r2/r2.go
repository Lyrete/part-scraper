package r2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type Item struct {
	name          string
	price         string
	vatPercentage int
}

func GetPriceFromPage(page *rod.Page) string {
	price_wrapper := page.MustElement("#product-offer > div.product-info.col-xs-12.col-sm-6.col-lg-6 > div > div.product-offer > div.row > div > div")
	itemPriceElement := price_wrapper.MustElement("strong")

	return itemPriceElement.MustElementR("span", "€").MustText()
}

func GetVatFromPage(page *rod.Page) int {
	vatElementText := page.MustElement(".vat_info").MustText()

	//Replace nbsp from the text as it's very likely we'll have some in there
	vatElementText = strings.Replace(vatElementText, "\u00a0", " ", -1)

	vatPercentage := strings.Replace(strings.Split(vatElementText, " ")[1], "%", "", 1)

	vatPercentageInt, err := strconv.Atoi(vatPercentage)

	if err != nil {
		panic(err)
	}

	return vatPercentageInt
}

func GetEanFromPage(page *rod.Page) string {
	eanText := page.MustElement("#tab-description").MustElementR(".attr-characteristic", "EAN:").MustText()

	eanNumber := strings.Split(eanText, "\n")[1]

	return eanNumber
}

func SearchItemByName(itemName string) Item {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	fmt.Println("Built browser, connecting to r2-bike.com")
	page := browser.MustPage("https://www.r2-bike.com/").MustWaitStable()

	fmt.Println("Searching for " + itemName)
	page.MustElement("input[name='qs']").MustInput(itemName).MustType(input.Enter)

	fmt.Println("Opening first search result")
	page.MustElement("#product-list").MustElement(".image-wrapper").MustClick()

	fmt.Println("Grabbing price")
	item_price := getPriceFromPage(page)
	vatPercentage := getVatFromPage(page)

	fmt.Println("Grabbing EAN")
	eanNumber := getEanFromPage(page)
	fmt.Println("EAN: " + eanNumber)

	return Item{name: itemName, price: item_price, vatPercentage: vatPercentage}

}
