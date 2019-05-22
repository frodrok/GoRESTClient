package configHandler

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Environment struct {
	name      string
	variables map[string]string
}

type Configuration struct {
	environments []Environment
}

var BASE64_SPLIT_PATTERN = "~_~"

type HistoryEntry struct {
	Url      string
	Username string
	Password string
}

type History struct {
	Content []HistoryEntry
}

var HOME_PATH = os.Getenv("HOME")

var CONFIG_PATH = HOME_PATH + "/.gostman/config"

func LoadConfig() (*Configuration, error) {

	raw_bytes, err := ioutil.ReadFile(CONFIG_PATH)

	if err != nil {
		// panic(err)
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

		err := os.Mkdir(HOME_PATH+"/.gostman", 0755)

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

	/* Load history from file or return an empty history struct and err */

	var hist = History{}

	// Check if the folder exists
	if _, err := os.Stat(HOME_PATH + "/.gostman"); os.IsNotExist(err) {

		err := os.Mkdir(HOME_PATH+"/.gostman", 0755)

		if err != nil {
			panic(err)
		}

	}

	// If the history file doesn't exist, create it and return an empty history
	if _, err := os.Stat(HOME_PATH + "/.gostman/history"); os.IsNotExist(err) {
		err := ioutil.WriteFile(HOME_PATH+"/.gostman/history", []byte(nil), 0644)

		if err != nil {
			panic(err)
		}

		return &hist, nil

	} else {

		// If it exists unmarshal the content and return it
		byhtes, err := ioutil.ReadFile(HOME_PATH + "/.gostman/history")

		if err != nil {
			panic(err)
		}

		str := string(byhtes)

		println(len(str))

		if len(str) < 1 {
			return &hist, nil
		}
		var historyFileRows = strings.Split(str, "\n")

		println("lenrows")
		println(len(historyFileRows))

		for _, entry := range historyFileRows {
			// Base64 decode and add to hist.Content
			raw, err := b64.StdEncoding.DecodeString(entry)

			if err != nil {
				continue
			}

			var histEntry = stringToHistoryEntry(string(raw))

			hist.Content = append(hist.Content, histEntry)

		}

		//hist.Content = ff

		return &hist, nil
	}

	return &hist, nil
}

func historyEntryToString(hist HistoryEntry) string {

	// Base 64 encode url:username:password and store that
	toEncode := hist.Url + BASE64_SPLIT_PATTERN + hist.Username + BASE64_SPLIT_PATTERN + hist.Password
	encoded := b64.StdEncoding.EncodeToString([]byte(toEncode))
	return encoded
}

func stringToHistoryEntry(s string) HistoryEntry {

	histEntry := &HistoryEntry{}

	println("received string:")
	println(s)

	splat := strings.Split(s, BASE64_SPLIT_PATTERN)

	histEntry.Url = splat[0]
	histEntry.Username = splat[1]
	histEntry.Password = splat[2]

	return *histEntry
}

func SaveToHistory(url string,
	username string,
	password string) bool {

	// hist is a pointer to a History struct
	var hist, err = GetHistoryOrCreateHistoryFile()

	if err != nil {
		panic(err)
	}

	entry := &HistoryEntry{
		url,
		username,
		password,
	}

	hist.Content = append(hist.Content, *entry)

	var asStrings []string

	for _, content := range hist.Content {
		asString := historyEntryToString(content)
		asStrings = append(asStrings, asString)
	}

	asFinalString := strings.Join(asStrings, "\n")

	if err := ioutil.WriteFile(HOME_PATH+"/.gostman/history", []byte(asFinalString), 0644); err != nil {
		panic(err)
		return false
	}

	return true

}
