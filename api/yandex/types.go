package yandex

import (
	"net/http"
	"os"
)

type ApiClient struct {
	HttpClient http.Client
	Token      string
	BaseUrl    string
}

type YandexDisk struct {
	ApiClient ApiClient
}

type OsResource struct {
	ItemPath string
	Info     os.FileInfo
}
