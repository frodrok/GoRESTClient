package configHandler

import (
	"GoRESTClient/httpClient"
	"encoding/json"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func LOG(s string) {
	fmt.Printf("%s - %s\n", time.Now().Format("2009-09-01T15:04"), s)
}

type ConfigHandlerActions interface {
	GetConfig() Configuration
	AddRequestToHistory(h httpClient.HttpRequest) bool
	GetRequestHistory() []httpClient.HttpRequest
}

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

type ConfigHandler struct {
	filePath string
}

func NewConfigHandler(filePath string) (cHandler *ConfigHandler, err error) {

	db, err = leveldb.OpenFile(filePath, nil)

	if err != nil {
		return nil, err
	}

	return &ConfigHandler{
		filePath,
	}, nil
}

func (cHandler *ConfigHandler) GetConfig() Configuration {
	return Configuration{[]Environment{}}
}

func requestToJson(h httpClient.HttpRequest) string {

	bytes, _ := json.Marshal(&h)

	return string(bytes)
}

func requestsToJson(h *[]httpClient.HttpRequest) string {
	bytes, _ := json.Marshal(h)
	return string(bytes)
}

func (cHandler *ConfigHandler) AddRequestToHistory(h httpClient.HttpRequest) bool {

	storedRequests := cHandler.GetRequestHistory()

	storedRequests = append(storedRequests, h)

	asJson := requestsToJson(&storedRequests)

	err := db.Put([]byte("history"), []byte(asJson), nil)

	if err != nil {
		LOG(err.Error())
		return false
	}
	return true
}

func (*ConfigHandler) GetRequestHistory() []httpClient.HttpRequest {

	data, _ := db.Get([]byte("history"), nil)

	var list []httpClient.HttpRequest

	json.Unmarshal(data, &list)

	return list
}
