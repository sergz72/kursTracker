package parsers

import (
	"core"
	"core/entities"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

func processMBTokens(tokens []core.Token) ([]entities.KursItem, error) {
	var result []entities.KursItem
	core.ResetTokenId()
	t := core.NextToken(tokens)
	isBuy := false
	sectionStart := false
	var rate float64 = 0
	var rateBuyEUR float64 = 0
	var rateBuyUSD float64 = 0
	var rateSellEUR float64 = 0
	var rateSellUSD float64 = 0
	for t.Typ != core.EOF {
		if core.IsText(t, "csvw:name") {
			core.NextToken(tokens)
			t = core.NextToken(tokens)
			if core.IsText(t, "Покупка") {
				sectionStart = true
				isBuy = true
			} else if core.IsText(t, "Продажа") {
				sectionStart = true
				isBuy = false
			}
		} else if sectionStart {
			if core.IsText(t, "csvw:value") {
				core.NextToken(tokens)
				t = core.NextToken(tokens)
				if t.Typ != core.TEXT {
					return nil, errors.New("wrong csvw:value")
				}
				var err error
				rate, err = strconv.ParseFloat(t.StringValue, 64)
				if err != nil {
					return nil, err
				}
			} else if core.IsText(t, "csvw:primaryKey") {
				core.NextToken(tokens)
				t = core.NextToken(tokens)
				if t.Typ != core.TEXT {
					return nil, errors.New("wrong csvw:primaryKey")
				}
				if t.StringValue == "EUR" {
					if isBuy {
						rateBuyEUR = rate
					} else {
						rateSellEUR = rate
					}
				} else if t.StringValue == "USD" {
					if isBuy {
						rateBuyUSD = rate
					} else {
						rateSellUSD = rate
					}
				}
				if rateBuyEUR != 0 && rateSellEUR != 0 && rateBuyUSD != 0 && rateSellUSD != 0 {
					result = append(result, entities.KursItem{
						CurrencyCodeA: 978,
						CurrencyCodeB: 980,
						RateBuy: rateBuyEUR,
						RateSell: rateSellEUR,
					})
					result = append(result, entities.KursItem{
						CurrencyCodeA: 840,
						CurrencyCodeB: 980,
						RateBuy: rateBuyUSD,
						RateSell: rateSellUSD,
					})
					return result, nil
				}
			}
		}
		t = core.NextToken(tokens)
	}
	return nil, fmt.Errorf("Incorrect number of results: %v\n", result)
}

func ParseMBKurs(path string, quiet bool) ([]entities.KursItem, error) {
	if !quiet {
		fmt.Printf("Processing file %v\n", path)
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tokens := core.Parse(dat)
	return processMBTokens(tokens)
}
