package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	HttpClient http.Client
	Token      string
	BaseUrl    string
}

func (c *Client) SendRequest(method, apiPage string) ([]byte, int) {
	// Отправляем запрос
	req, err := http.NewRequest(
		method, c.BaseUrl+apiPage, nil,
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "OAuth "+c.Token)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	check_error := CheckResponse(*resp)
	if check_error != nil {
		fmt.Println(check_error)
		os.Exit(1)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body, resp.StatusCode
}

type Yadisk struct {
	Client Client
}

func NewYaDisk(token string) (Yadisk, error) {
	if token == "" {
		return Yadisk{}, ERROR_NO_TOKEN
	}
	newClient := Client{
		HttpClient: http.Client{},
		Token:      token,
		BaseUrl:    "https://cloud-api.yandex.net/v1/",
	}
	return Yadisk{
		Client: newClient,
	}, nil
}

func (yad *Yadisk) getMetaInformation(resourcePath string) Resource {
	resource := Resource{}
	response, _ := yad.Client.SendRequest("GET", "disk/resources?path="+url.QueryEscape(resourcePath))
	json.Unmarshal(response, &resource)
	return resource
}

func (yad *Yadisk) getListFiles() FilesResourceList {
	listFiles := FilesResourceList{}
	response, _ := yad.Client.SendRequest("GET", "disk/resources/files")
	json.Unmarshal(response, &listFiles)
	return listFiles
}

type Item struct {
	ItemPath string
	Info     os.FileInfo
	Err      error
}

func (yad *Yadisk) getAllResources(rootPath string) []Resource {
	resources := []Resource{}
	items := yad.getMetaInformation(rootPath).Embedded.Items
	if len(items) > 0 {
		for _, item := range items {
			if item.isDir() {
				resources = append(resources, item)            // Добавляем папку в наш массив ресурсов
				dirResources := yad.getAllResources(item.Path) // Запускаем гуся в папку
				resources = append(resources, dirResources...) // Добавляем то что притащил гусь в наш массив ресурсов
			} else {
				resources = append(resources, item) // Если это файл то просто добовляем путь в наш массив ресурсов
			}
		}
		return resources // Возращаем массив ресурсов
	} else {
		return resources // Возращаем пустой массив ресурсов
	}

}

func (yad *Yadisk) Upload(itemName string) {
	isDir, err := isWritableDir(itemName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if isDir == true {
		items := []Item{}
		// Получаем все пути(пути папок и файлов в директории)
		filepath.Walk(itemName, func(path string, info os.FileInfo, err error) error {
			items = append(items, Item{ItemPath: path, Info: info, Err: err})
			return nil
		})
		// Если в папке что то есть то загрузжаем эти файлы на диск
		// Иначе просто грузим папку
		if len(items) > 0 {
			for _, item := range items {
				if item.Info.IsDir() {
					yad.CreateDir(item.ItemPath)
				} else {
					yad.UploadFile(item.ItemPath)
				}
			}
		} else {
			dirExists, link := yad.CreateDir(itemName)
			if dirExists {
				fmt.Println("Папка создана")
				fmt.Println("Ссылка на папку:" + link.Href)
			}
		}

	} else {
		yad.UploadFile(itemName)
	}
}

func (yad *Yadisk) CreateDir(dirname string) (bool, Link) {
	link := Link{}
	createDir, StatusCode := yad.Client.SendRequest("PUT", "disk/resources?path="+url.QueryEscape(dirname))
	if StatusCode == 201 {
		json.Unmarshal(createDir, &link)
		return true, link
	}
	return false, link
}

func (yad *Yadisk) Download(ItemPath string) {
	resource := yad.getMetaInformation(ItemPath)
	if resource.isDir() {
		fmt.Println("Создание главной папки", resource.CorrectPath())
		os.Mkdir(resource.CorrectPath(), os.ModePerm)
		items := yad.getAllResources(resource.Path)
		for _, item := range items {
			if item.isDir() {
				fmt.Println("Создание папки", item.CorrectPath())
				os.Mkdir(item.CorrectPath(), os.ModePerm)
			} else {
				fmt.Println("Загрузка ", item.CorrectPath())
				yad.DownloadFile(item.Path)
			}
		}
	} else {
		fmt.Println("Загрузка " + ItemPath)
		yad.DownloadFile(ItemPath)
	}
}

func (yad *Yadisk) DownloadFile(filename string) {
	link := Link{}
	download_url, _ := yad.Client.SendRequest("GET", "disk/resources/download?path="+url.QueryEscape(filename))
	json.Unmarshal(download_url, &link)
	resp, download_error := yad.Client.HttpClient.Get(link.Href)
	if download_error != nil {
		fmt.Println(setTextColor("Unable to get file:"+download_error.Error(), RED))
		os.Exit(1)
	}
	file, err := os.Create(strings.Replace(filename, "disk:/", "", 1))
	if err != nil { // если возникла ошибка
		fmt.Println(setTextColor(err.Error(), RED))
		os.Exit(1) // выходим из программы
	}
	data, _ := io.ReadAll(resp.Body)
	file.Write(data)
	file.Close()

}

func (yad *Yadisk) UploadFile(filepath string) {
	link := Link{}
	fmt.Println("Идёт загрузка " + filepath)
	upload_url, _ := yad.Client.SendRequest("GET", "disk/resources/upload?path="+url.QueryEscape(filepath))
	json.Unmarshal(upload_url, &link)
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(setTextColor(err.Error(), RED))
		os.Exit(1)
	}
	rq, errRq := http.NewRequest("PUT", link.Href, file)
	if errRq != nil {
		fmt.Println(setTextColor("Unable to create request:"+errRq.Error(), RED))
		os.Exit(1)
	}
	resp, errRqDo := yad.Client.HttpClient.Do(rq)
	if errRqDo != nil {
		fmt.Println(setTextColor("Unable to make request:"+errRq.Error(), RED))
		os.Exit(1)
	}
	switch resp.StatusCode {
	case 201:
		fmt.Println("Файл создан")
	case 202:
		fmt.Println("Файл принят сервером, но еще не был перенесен непосредственно в Яндекс.Диск.")
	case 412:
		fmt.Println("При дозагрузке файла был передан неверный диапазон в заголовке Content-Range.")
	case 413:
		fmt.Println("Размер файла превышает 10 ГБ.")
	case 500:
	case 503:
		fmt.Println("Ошибка сервера, попробуйте повторить загрузку.")
	case 507:
		fmt.Println("Для загрузки файла не хватает места на Диске")
	}
	file.Close()

}

func getConfigFile() (map[string]string, error) {
	configFile := map[string]string{}
	dataConfFile, fileErr := os.ReadFile("config.json")
	if fileErr != nil {
		return nil, errors.New(setTextColor(fileErr.Error(), RED))
	}
	json.Unmarshal(dataConfFile, &configFile)
	return configFile, nil
}
