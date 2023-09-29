package utils

import "github.com/go-rod/rod"

func GetUrl(page *rod.Page) string {
	return page.MustInfo().URL
}
