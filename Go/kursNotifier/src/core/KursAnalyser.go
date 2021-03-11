package core

import (
	"core/entities"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func readFiles(folder string) []string {
	var result []string
	d, err := os.Open(folder)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, fil := range fi {
		if fil.IsDir() {
			continue
		}
		result = append(result, folder + string(os.PathSeparator) + fil.Name())
	}

	return result
}
func BuildKursMessage(folder string, quiet bool, sources []string) string {
	t := time.Now()
	currentFolderName := folder + string(os.PathSeparator) + fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	if _, err := os.Stat(currentFolderName); err == nil {
		files := readFiles(currentFolderName)
		for i := 0; i < 3; i++ {
			t = t.AddDate(0, 0, -1)
			currentFolderName = folder + string(os.PathSeparator) + fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
			if _, err = os.Stat(currentFolderName); err == nil {
				files = append(files, readFiles(currentFolderName)...)
				break
			}
		}
		if len(files) > 1 {
			sort.Strings(files)
			previousKursItems, err := entities.ReadKursItemsFromJson(files[len(files) - 2], quiet)
            if err != nil {
            	panic(err)
			}
			//fmt.Println(previousKursItems)
			actualKursItems, err := entities.ReadKursItemsFromJson(files[len(files) - 1], quiet)
			if err != nil {
				panic(err)
			}
			//fmt.Println(actualKursItems)
			return compareKursItems(previousKursItems, actualKursItems, sources)
		}
	}

	return ""
}

func compareKursItems(previousItems map[string][]entities.KursItem, actualItems map[string][]entities.KursItem,
	                  sources []string) string {
	var b strings.Builder
	for _, source := range sources {
		prev, ok := previousItems[source]
		curr, ok2 := actualItems[source]
		if !ok || !ok2 {
			b.WriteString(source + " is missing.\n")
		} else {
			for _, message := range compareSourceKursItems(source, prev, curr) {
				b.WriteString(message)
			}
		}
	}
	return b.String()
}

func compareSourceKursItems(sourceName string, prev []entities.KursItem, curr []entities.KursItem) []string {
	var prevUSDItem *entities.KursItem = nil
	var prevEURItem *entities.KursItem = nil
	var currUSDItem *entities.KursItem = nil
	var currEURItem *entities.KursItem = nil

	for _, item := range prev {
		if item.CurrencyCodeA == 840 {
			prevUSDItem = &entities.KursItem{
				RateBuy: item.RateBuy,
				RateSell: item.RateSell,
			}
		} else if item.CurrencyCodeA == 978 {
			prevEURItem = &entities.KursItem{
				RateBuy: item.RateBuy,
				RateSell: item.RateSell,
			}
		}
	}

	for _, item := range curr {
		if item.CurrencyCodeA == 840 {
			currUSDItem = &entities.KursItem{
				RateBuy: item.RateBuy,
				RateSell: item.RateSell,
			}
		} else if item.CurrencyCodeA == 978 {
			currEURItem = &entities.KursItem{
				RateBuy: item.RateBuy,
				RateSell: item.RateSell,
			}
		}
	}

	var result []string

	if prevEURItem == nil || prevUSDItem == nil || currEURItem == nil || currUSDItem == nil {
		result = append(result, sourceName + ": Incorrect record format\n")
	} else {
		result = compareItems(result, sourceName, "EUR", prevEURItem, currEURItem)
		result = compareItems(result, sourceName, "USD", prevUSDItem, currUSDItem)
	}

	return result
}

func compareItems(result []string, sourceName string, currency string, prevItem *entities.KursItem, currItem *entities.KursItem) []string {
	if prevItem.RateBuy > currItem.RateBuy {
		result = append(result, fmt.Sprintf("%v %v BUY DOWN %v -> %v\n", sourceName, currency, prevItem.RateBuy, currItem.RateBuy))
	} else if prevItem.RateBuy < currItem.RateBuy {
		result = append(result, fmt.Sprintf("%v %v BUY UP %v -> %v\n", sourceName, currency, prevItem.RateBuy, currItem.RateBuy))
	}
	if prevItem.RateSell > currItem.RateSell {
		result = append(result, fmt.Sprintf("%v %v SELL DOWN %v -> %v\n", sourceName, currency, prevItem.RateSell, currItem.RateSell))
	} else if prevItem.RateSell < currItem.RateSell {
		result = append(result, fmt.Sprintf("%v %v SELL UP %v -> %v\n", sourceName, currency, prevItem.RateSell, currItem.RateSell))
	}
	return result
}
