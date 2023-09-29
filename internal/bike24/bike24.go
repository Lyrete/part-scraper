package bike24

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/stealth"
	utils "github.com/lyrete/part-scraper/internal"
)

type Item struct {
	name     string
	price    float64
	itemCode string
	URL      string
}

func Connect(browser *rod.Browser) *rod.Page {
	fmt.Println("Connecting to Bike24...")
	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.bike24.com/")

	return page
}

func SearchForItem(page *rod.Page, searchTerm string) Item {
	fmt.Println("Searching for " + searchTerm + " on Bike24")

	wait := page.MustWaitRequestIdle()
	page.MustElement("input[type='search']").MustInput(searchTerm).MustType(input.Enter)
	wait()

	page.MustElement("#search-results .container-lg > div > div:nth-child(3) div.row > div").MustClick()

	vatPercentage := getVatFromPage(page)

	itemPrice := getPriceFromPage(page, vatPercentage)

	eanNumber := getEanFromPage(page)

	return Item{name: getFullName(page), price: itemPrice, barcode: eanNumber, URL: utils.GetUrl(page)}
}

func getVatFromPage(page *rod.Page) int {
	return 0
}

func getPriceFromPage(page *rod.Page, vatPercentage int) float64 {
	itemPrice := page.MustElement(".product-detail-price p.price__value").MustText()
	itemPrice = strings.Replace(itemPrice, " €", "", 1)
	itemPrice = strings.Replace(itemPrice, ",", ".", 1)
	itemPrice = strings.TrimSpace(itemPrice)

	itemPriceFloat, err := strconv.ParseFloat(itemPrice, 64)

	if err != nil {
		panic(err)
	}

	return itemPriceFloat
}

func getEanFromPage(page *rod.Page) string {
	return ""
}

func getFullName(page *rod.Page) string {
	return page.MustElement(".product-detail-information-area__product-name").MustText()
}
