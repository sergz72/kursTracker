package main

import (
  "core/entities"
  "core/parsers"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "time"
)

func printKutsItems(kursItems []entities.KursItem, name string) {
  for _, item := range kursItems {
    if item.CurrencyCodeA == 840 && item.CurrencyCodeB == 980 {
      fmt.Printf("%v UAH to USD buy: %v sell: %v\n", name, item.RateBuy, item.RateSell)
    } else if item.CurrencyCodeA == 978 && item.CurrencyCodeB == 980 {
      fmt.Printf("%v UAH to EUR buy: %v sell: %v\n", name, item.RateBuy, item.RateSell)
    }
  }
}

func main() {
  l := len(os.Args)
  if l < 3 {
    fmt.Fprintf(os.Stderr, "Usage: kurs html_files_folder result_folder [-q]")
    return
  }

  quiet := l == 4 && os.Args[3] == "-q"

  var result map[string][]entities.KursItem
  result = make(map[string][]entities.KursItem)

  kursItems, err := entities.ReadMonoKursItemsFromJson(os.Args[1] + string(os.PathSeparator) + "mono_currencies.json", quiet)
  if err == nil {
    if !quiet {
      printKutsItems(kursItems, "Mono")
    }
    result["Mono"] = kursItems
  } else {
    fmt.Fprintf(os.Stderr, "Mono parse error: %v", err)
  }
  kursItems, err = parsers.ParsePrivatKurs(os.Args[1] + string(os.PathSeparator) + "pb.ua", quiet)
  if err == nil {
    if !quiet {
      printKutsItems(kursItems, "Privat")
    }
    result["Privat"] = kursItems
  } else {
    fmt.Fprintf(os.Stderr, "Privat parse error: %v", err)
  }
  kursItems, err = parsers.ParseAlfaKurs(os.Args[1] + string(os.PathSeparator) + "alfabank.ua", quiet)
  if err == nil {
    if !quiet {
      printKutsItems(kursItems, "Alfa")
    }
    result["Alfa"] = kursItems
  } else {
    fmt.Fprintf(os.Stderr, "Alfa parse error: %v", err)
  }
  kursItems, err = parsers.ParseMBKurs(os.Args[1] + string(os.PathSeparator) + "kurs.com.ua", quiet)
  if err == nil {
    if !quiet {
      printKutsItems(kursItems, "MB")
    }
    result["MB"] = kursItems
  } else {
    fmt.Fprintf(os.Stderr, "MB parse error: %v", err)
  }
  var jsonText []byte
  jsonText, err = json.MarshalIndent(result, "", "")
  if err != nil {
    fmt.Fprintf(os.Stderr, "JSON marshal error: %v", err)
    return
  }
  t := time.Now()
  folderName := os.Args[2] + string(os.PathSeparator) + fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
  if _, err = os.Stat(folderName); os.IsNotExist(err) {
    err = os.Mkdir(folderName, 0755)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Mkdir error: %v", err)
      return
    }
  }
  fileName := folderName + string(os.PathSeparator) + fmt.Sprintf("%d", t.Unix())
  err = ioutil.WriteFile(fileName, jsonText, 0644)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Result file write error: %v", err)
  }
}
