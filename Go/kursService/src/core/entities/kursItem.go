package entities

import (
	"encoding/json"
	"io/ioutil"
)

type KursItem struct {
	CurrencyCodeA int
	CurrencyCodeB int
	RateBuy float64
	RateSell float64
}

type KursOutItem struct {
	Date int64
	BankName string
	RateBuyUSD float64
	RateSellUSD float64
	RateBuyEUR float64
	RateSellEUR float64
}

func ReadKursItemsFromJson(path string) (map[string][]KursItem, error) {
	var items map[string][]KursItem
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dat, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func ConvertKursItems(result []KursOutItem, date int64, kursItems map[string][]KursItem) []KursOutItem {
    for k, v := range kursItems {
		var rateBuyUSD float64
		var rateSellUSD float64
		var rateBuyEUR float64
		var rateSellEUR float64
		for _, item := range v {
			if item.CurrencyCodeA == 840 {
				rateBuyUSD = item.RateBuy
				rateSellUSD = item.RateSell
			} else {
				rateBuyEUR = item.RateBuy
				rateSellEUR = item.RateSell
			}
		}
		result = append(result, KursOutItem{
			Date:        date,
			BankName:    k,
			RateBuyUSD:  rateBuyUSD,
			RateSellUSD: rateSellUSD,
			RateBuyEUR:  rateBuyEUR,
			RateSellEUR: rateSellEUR,
		})
    }

    return result
}