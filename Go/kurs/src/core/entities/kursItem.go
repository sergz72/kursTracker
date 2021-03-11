package entities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type KursItem struct {
	CurrencyCodeA int
	CurrencyCodeB int
	RateBuy float64
	RateSell float64
}

func ReadMonoKursItemsFromJson(path string, quiet bool) ([]KursItem, error) {
	if !quiet {
		fmt.Printf("Processing file %v\n", path)
	}
	var items []KursItem
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dat, &items); err != nil {
		return nil, err
	}

	var result []KursItem
	for _, item := range items {
		if (item.CurrencyCodeA == 840 && item.CurrencyCodeB == 980) || (item.CurrencyCodeA == 978 && item.CurrencyCodeB == 980) {
			result = append(result, item)
		}
	}
	return result, nil
}