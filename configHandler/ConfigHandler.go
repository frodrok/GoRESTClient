package configHandler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Environment struct {
	name string
	variables map[string]string
}

type Configuration struct {
	environments []Environment
}

type History struct {
	Content []string
}

var HOME_PATH = os.Getenv("HOME")

var CONFIG_PATH = HOME_PATH + "/.gostman/config"

func LoadConfig() (*Configuration, error) {

	raw_bytes, err := ioutil.ReadFile(CONFIG_PATH)

	if err != nil {
		panic(err)
		return nil, err
	}

	res := Configuration{}

	smut := json.Unmarshal(raw_bytes, &res)

	if smut != nil {
		panic(smut)
		return &res, err
	}

	return &res, nil

}

func SaveConfig(conf *Configuration) bool {

	var actual = *conf
	var asJson, _ = json.Marshal(actual)

	// Check if the folder exists
	if _, err := os.Stat(HOME_PATH + "/.gostman"); os.IsNotExist(err) {

		err := os.Mkdir(HOME_PATH + "/.gostman", 0755)

		if err != nil {
			panic(err)
		}

	}

	err := ioutil.WriteFile(CONFIG_PATH, []byte(asJson), 0644)

	if err != nil {
		panic(err)
		return false
	}

	return true

}

func GetHistoryOrCreateHistoryFile() (*History, error) {

	var hist = History{}

	// Check if the folder exists
	if _, err := os.Stat(HOME_PATH + "/.gostman"); os.IsNotExist(err) {

		err := os.Mkdir(HOME_PATH + "/.gostman", 0755)

		if err != nil {
			panic(err)
		}

	}

	// If the history file doesn't exist, create it and return an empty history
	if _, err := os.Stat(HOME_PATH + "/.gostman/history"); os.IsNotExist(err) {
		err := ioutil.WriteFile(HOME_PATH + "/.gostman/history", []byte(nil), 0644)

		if err != nil {
			panic(err)
		}

		return &hist, nil

	} else {

		// If it exists unmarshal the json and return it
		byhtes, err := ioutil.ReadFile(HOME_PATH + "/.gostman/history")

		if err != nil {
			panic(err)
		}

		str := string(byhtes)
		var ff = strings.Split(str, "\n")

		hist.Content = ff

		return &hist, nil
	}


	return &hist, nil
}


func SaveToHistory(url string) bool {

	// hist is pointer to a History struct
	var hist, err = GetHistoryOrCreateHistoryFile()

	if err != nil {
		panic(err)
	}

	hist.Content = append(hist.Content, url)

	asString := strings.Join(hist.Content, "\n")

	if err := ioutil.WriteFile(HOME_PATH + "/.gostman/history", []byte(asString), 0644); err != nil {
		panic(err)
		return false
	}

	return true

}