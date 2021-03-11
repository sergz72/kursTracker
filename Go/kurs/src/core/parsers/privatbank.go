package parsers

import (
	"core"
	"core/entities"
	"fmt"
	"io/ioutil"
	"strconv"
)

func processPBTokens(tokens []core.Token) ([]entities.KursItem, error) {
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
					if t.Typ == core.TEXT {
						//fmt.Print("e")
						t = core.NextToken(tokens)
						if core.IsName(t, "data-cource_type") {
							//fmt.Print("f")
							core.NextToken(tokens)
							t = core.NextToken(tokens)
							if core.IsText(t, "cards_course") {
								//fmt.Print("g")
								t = core.NextToken(tokens)
								valuta1 := ""
								valuta2 := ""
								var rateBuy float64 = 0
								for t.Typ != core.EOF {
									if core.IsName(t, "td") {
										t = core.NextToken(tokens)
										if core.IsSymbol(t, '>') {
											t = core.NextToken(tokens)
											if valuta1 == "" || valuta2 == "" {
												if core.IsName(t, "EUR") {
													valuta1 = "EUR"
												}
												if core.IsName(t, "USD") {
													valuta1 = "USD"
												}
												if core.IsName(t, "UAH") {
													valuta2 = "UAH"
												}
											} else if t.Typ == core.NUMBER {
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
													valuta2 = ""
													rateBuy = 0
												}
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

func ParsePrivatKurs(path string, quiet bool) ([]entities.KursItem, error) {
	if !quiet {
		fmt.Printf("Processing file %v\n", path)
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tokens := core.Parse(dat)
	return processPBTokens(tokens)
}
