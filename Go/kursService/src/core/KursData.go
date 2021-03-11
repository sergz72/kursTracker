package core

import (
	"bytes"
	"core/entities"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
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
		result = append(result, fil.Name())
	}

	return result
}

func KursDataHandler(server *Server, w *bytes.Buffer, req string) {
	values, err := url.ParseQuery(req)
	if err != nil {
		w.Write([]byte("400 Bad request: query parsing error"))
		return
	}
	periods := values.Get("period")
	var period int
	if len(periods) > 0 {
		period, err = strconv.Atoi(periods)
		if err != nil || period <= 0 {
			w.Write([]byte("400 Bad request: " + err.Error()))
			return
		}
	} else {
		w.Write([]byte("400 Bad request: period is missing"))
		return
	}

	period--

	t := time.Now().AddDate(0, 0, -period)
	var resultItems []entities.KursOutItem
	periodv := period
	for period >= 0 {
		currentFolderName := server.kursFolder + string(os.PathSeparator) + fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
		if _, err = os.Stat(currentFolderName); err == nil {
			files := readFiles(currentFolderName)
			l := len(files)
			if l > 0 {
				var result []string
				if l > 1 {
					if periodv > 1 {
						sort.Strings(files)
						result = append(result, files[len(files)-1])
					} else {
						result = files
					}
				} else {
					result = files
				}
				for _, file := range result {
					date, err := strconv.ParseInt(file, 10, 64)
					if err != nil {
						log.Println(err)
						w.Write([]byte("500 Server error: " + err.Error()))
						return
					}
					kursItems, err := entities.ReadKursItemsFromJson(currentFolderName + string(os.PathSeparator) + file)
					if err != nil {
						log.Println(err)
						w.Write([]byte("500 Server error: " + err.Error()))
						return
					}
					resultItems = entities.ConvertKursItems(resultItems, date, kursItems)
				}
			}
		}

		t = t.AddDate(0, 0, 1)
		period--
	}

	if len(resultItems) > 0 {
		bytess, err := json.Marshal(resultItems)
		if err != nil {
			log.Println(err)
			w.Write([]byte("500 Server error: " + err.Error()))
			return
		}
		w.Write(bytess)
	} else {
		w.Write([]byte("[]"))
	}
}

