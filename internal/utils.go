package utils

import (
	"strings"

	"github.com/go-rod/rod"
)

type Item struct {
	name    string
	price   float64
	barcode string
	URL     string
}

func GetUrl(page *rod.Page) string {
	return strings.Split(page.MustInfo().URL, "?")[0]
}
