package core

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

type Configuration struct {
	BotToken string
	ChatIDs []int64
	Sources []string
}

func LoadConfiguration(iniFileName string) (*Configuration, error) {
	var config Configuration

	iniFileContents, err := ioutil.ReadFile(iniFileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(iniFileContents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadEnabledChatIDs(fileName string) ([]int64, error) {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var result []int64
	for _, id := range strings.Split(string(fileContents), ",") {
		if id != "" {
			n, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, n)
		}
	}
	return result, nil
}

func SaveEnabledChatIDs(fileName string, ids []int64) error {
	var result strings.Builder
	first := true
	for _, id := range ids {
		s := strconv.FormatInt(id, 10)
		if first {
			first = false
		} else {
			result.WriteString(",")
		}
		result.WriteString(s)
	}
	err := ioutil.WriteFile(fileName, []byte(result.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}
