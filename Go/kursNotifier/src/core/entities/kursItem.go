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

func ReadKursItemsFromJson(path string, quiet bool) (map[string][]KursItem, error) {
	if !quiet {
		fmt.Printf("Processing file %v\n", path)
	}
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