package parsers

import (
	"core"
	"core/entities"
	"fmt"
	"io/ioutil"
	"strconv"
)

func processAlfaTokens(tokens []core.Token) ([]entities.KursItem, error) {
	var result []entities.KursItem
	core.ResetTokenId()
	t := core.NextToken(tokens)
	for t.Typ != core.EOF {
		if core.IsSymbol(t, '<') {
			//fmt.Print("a")
			t = core.NextToken(tokens)
			if core.IsName(t, "div") {
				//fmt.Print("b")
				t = core.NextToken(tokens)
				//fmt.Print("c")
				if core.IsName(t, "class") {
					//fmt.Print("d")
					core.NextToken(tokens)
					t = core.NextToken(tokens)
					if core.IsText(t, "currency-tab-block") {
						//fmt.Print("e")
						t = core.NextToken(tokens)
						if core.IsName(t, "data-tab") {
							//fmt.Print("f")
							core.NextToken(tokens)
							t = core.NextToken(tokens)
							if core.IsText(t, "0") {
								//fmt.Print("g")
								t = core.NextToken(tokens)
								valuta1 := ""
								//valuta2 := ""
								var rateBuy float64 = 0
								for t.Typ != core.EOF {
									if core.IsText(t, "title") {
										core.NextToken(tokens)
										t = core.NextToken(tokens)
										valuta1 = t.StringValue
										//core.NextToken(tokens)
										//t = core.NextToken(tokens)
										//valuta2 = t.StringValue
									} else if core.IsText(t, "rate-number") {
										core.NextToken(tokens)
										t = core.NextToken(tokens)
										if t.Typ == core.NUMBER {
											if rateBuy == 0 {
												var err error
												//fmt.Println(t.StringValue)
												rateBuy, err = strconv.ParseFloat(t.StringValue, 64)
												if err != nil {
													return nil, err
												}
											} else {
												//fmt.Println(t.StringValue)
												rateSell, err := strconv.ParseFloat(t.StringValue, 64)
												if err != nil {
													return nil, err
												}
												if valuta1 == "EUR" {
													result = append(result, entities.KursItem{
														CurrencyCodeA: 978,
														CurrencyCodeB: 980,
														RateBuy: rateBuy,
														RateSell: rateSell,
													})
												} else {
													result = append(result, entities.KursItem{
														CurrencyCodeA: 840,
														CurrencyCodeB: 980,
														RateBuy: rateBuy,
														RateSell: rateSell,
													})
												}
												valuta1 = ""
												//valuta2 = ""
												rateBuy = 0
											}
										}
									}
									if len(result) == 2 {
										return result, nil
									}
									t = core.NextToken(tokens)
								}
							}
						}
					}
				}
				//fmt.Println()
			}
		}
		t = core.NextToken(tokens)
	}

	return nil, fmt.Errorf("Incorrect number of results: %v\n", result)
}

func ParseAlfaKurs(path string, quiet bool) ([]entities.KursItem, error) {
	if !quiet {
		fmt.Printf("Processing file %v\n", path)
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tokens := core.Parse(dat)
	return processAlfaTokens(tokens)
}