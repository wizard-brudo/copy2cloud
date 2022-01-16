package yandex

import (
	"copy2cloud/utils"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Проверяет ответ от  yandex api
func (c *ApiClient) CheckResponse(response http.Response) error {
	if response.StatusCode >= 400 {
		errApi := ErrorApi{}
		body_byte, _ := ioutil.ReadAll(response.Body)
		err := json.Unmarshal(body_byte, &errApi)
		if err != nil {
			return utils.NewError(utils.ERROR_JSON)
		}
		return utils.NewError(errApi.Message)
	}
	return nil
}

// Делает запрос к Yandex api
func (client *ApiClient) SendRequest(method, apiPage string) ([]byte, int, error) {
	// Отправляем запрос
	req, err := http.NewRequest(
		method, client.BaseUrl+apiPage, nil,
	)
	// Добовляем заголовки
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "OAuth "+client.Token)
	// Делаем запрос
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, -1, err
	}
	// Проверяем ответ
	check_error := client.CheckResponse(*resp)
	if check_error != nil {
		return nil, -2, check_error
	}
	// Читаем ответ(переводим в массив байтов)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body, resp.StatusCode, nil
}

// Создаёт и возращает структуру YandexDisk(нужна для работы с яндекс диском)
func NewYandexDisk(token string) (YandexDisk, error) {
	if token == "" {
		return YandexDisk{}, utils.NewError(utils.ERROR_NO_TOKEN)
	}
	newApiClient := ApiClient{
		HttpClient: http.Client{},
		Token:      token,
		BaseUrl:    "https://cloud-api.yandex.net/v1/disk/",
	}
	return YandexDisk{
		ApiClient: newApiClient,
	}, nil
}

// Возвращает мета-информация о ресурсе в корзине
func (yadisk *YandexDisk) GetTrashResource(resourcePath string) (Resource, error) {
	resource := Resource{}
	response, _, err := yadisk.ApiClient.SendRequest("GET", "trash/resources?path="+url.QueryEscape(resourcePath))
	if err != nil {
		return resource, err
	}
	json.Unmarshal(response, &resource)
	return resource, nil
}

// Возвращает мета-информация о ресурсе на диске
func (yadisk *YandexDisk) GetResource(resourcePath string) (Resource, error) {
	resource := Resource{}
	response, _, err := yadisk.ApiClient.SendRequest("GET", "resources?path="+url.QueryEscape(resourcePath))
	if err != nil {
		return resource, err
	}
	json.Unmarshal(response, &resource)
	return resource, nil
}

// Возращает список всех файлов на диске
func (yadisk *YandexDisk) GetListFiles(mediaType string) (FilesResourceList, error) {
	listFiles := FilesResourceList{}
	response, _, err := yadisk.ApiClient.SendRequest("GET", "resources/files?media_type="+mediaType)
	if err != nil {
		return listFiles, err
	}
	json.Unmarshal(response, &listFiles)
	return listFiles, nil
}

// Возращает все ресурсы в указанной директории на диске
func (yadisk *YandexDisk) GetAllResources(rootPath string) ([]Resource, error) {
	resources := []Resource{}
	yaDiskResources, err := yadisk.GetResource(rootPath)
	if err != nil {
		return nil, err
	}
	if len(yaDiskResources.Embedded.Items) > 0 {
		for _, resource := range yaDiskResources.Embedded.Items {
			if resource.IsDir() {
				resources = append(resources, resource)                  // Добавляем папку в наш массив ресурсов
				dirResources, _ := yadisk.GetAllResources(resource.Path) // Запускаем гуся в папку
				resources = append(resources, dirResources...)           // Добавляем то что притащил гусь в наш массив ресурсов
			} else {
				resources = append(resources, resource) // Если это файл то просто добовляем путь в наш массив ресурсов
			}
		}
		return resources, nil // Возращаем массив ресурсов
	} else {
		return resources, nil // Возращаем пустой массив ресурсов
	}

}

// Возрашает yandex ссылку(объект содержит ссылку для запроса метаданных)
func (yadisk *YandexDisk) GetLink(method string, apiUrl string, urlParameters url.Values) (Link, error) {
	link := Link{}
	apiResponse, StatusCode, rqError := yadisk.ApiClient.SendRequest(method, apiUrl+"?"+urlParameters.Encode())
	// Коды удачных запрос это: 200-хорошо,201-создано,202-Принято,204-Нет содержимого(Удаление ресурса)
	if StatusCode == 200 || StatusCode == 201 || StatusCode == 202 || StatusCode == 204 {
		json.Unmarshal(apiResponse, &link)
		return link, nil // Возращаем Link  и nil
	}
	return link, rqError
}

// Создаёт директорию
func (yadisk *YandexDisk) CreateDir(dirname string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", dirname)
	createDir, err := yadisk.GetLink("PUT", "resources", urlValues)
	if err != nil {
		return createDir, err // CreateDir в данном случае пустой
	}
	return createDir, nil
}

// Копирует ресурс из одного места(from) в другое(path)
func (yadisk *YandexDisk) CopyResource(from, path, overwrite string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("from", from)
	urlValues.Add("path", path)
	urlValues.Add("overwrite", overwrite)
	copyResource, err := yadisk.GetLink("POST", "resources/copy", urlValues)
	if err != nil {
		return copyResource, err // copyResource  пустой
	}
	return copyResource, nil
}

// Удаляет ресурс заданному пути(path)
func (yadisk *YandexDisk) DeleteResource(path, permanently string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", path)
	urlValues.Add("permanently", permanently)
	deleteResource, err := yadisk.GetLink("DELETE", "resources", urlValues)
	if err != nil {
		return deleteResource, err // copyResource  пустой
	}
	return deleteResource, nil
}

// Удаляет ресурс в корзине по заданному пути
func (yadisk *YandexDisk) ClearTrash(path string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", path)
	clearTrash, err := yadisk.GetLink("DELETE", "trash/resources", urlValues)
	if err != nil {
		return clearTrash, err // clearTrash  пустой
	}
	return clearTrash, nil
}

func (yadisk *YandexDisk) Find(rootPath, name string) (Resource, error) {
	listResources, err := yadisk.GetAllResources(rootPath)
	if err != nil {
		return Resource{}, err
	}
	for _, resource := range listResources {
		if resource.Name == name {
			return resource, nil
		}
	}
	return Resource{}, err // Если мы сюда дошли мы нечего не нашли и возращаем пустой resource
}

// Восстанавливает ресурс в корзине по заданному пути
func (yadisk *YandexDisk) RestoreTrash(path, overwrite string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", path)
	urlValues.Add("overwrite", overwrite)
	restoreTrash, err := yadisk.GetLink("PUT", "trash/resources/restore", urlValues)
	if err != nil {
		return restoreTrash, err // copyResource  пустой
	}
	return restoreTrash, nil
}

// Перемещает ресурс из одного места(from) в другое(path)
func (yadisk *YandexDisk) MoveResource(from, path, overwrite string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("from", from)
	urlValues.Add("path", path)
	urlValues.Add("overwrite", overwrite)
	moveResource, err := yadisk.GetLink("POST", "resources/move", urlValues)
	if err != nil {
		return moveResource, err // moveResource пустой
	}
	return moveResource, nil
}

//	Проверяет существует ли ресурс на диске
func (yadisk *YandexDisk) ResourceExists(resourceName string) bool {
	if _, err := yadisk.GetResource(resourceName); err != nil {
		return false
	}
	return true
}

// Скачивает файл с диска
func (yadisk *YandexDisk) DownloadFile(filepath, overwrite string) error {
	urlValues := url.Values{}
	urlValues.Add("path", filepath)
	urlValues.Add("overwrite", overwrite)
	downloadLink, LinkErr := yadisk.GetLink("GET", "resources/download", urlValues)
	if LinkErr != nil {
		return LinkErr
	}
	resp, download_error := yadisk.ApiClient.HttpClient.Get(downloadLink.Href)
	if download_error != nil {
		return utils.NewError("невозможно получить файл")
	}
	file, err := os.Create(strings.Replace(filepath, "disk:/", "", 1))
	if err != nil {
		return utils.NewError("при создание файла произошла ошибка")
	}
	data, _ := io.ReadAll(resp.Body)
	file.Write(data)
	file.Close()
	return nil
}

// Загрузжает системный ресурс на yandex disk
func (yadisk *YandexDisk) UploadFile(filepath, overwrite string) (int, error) {
	urlValues := url.Values{}
	urlValues.Add("overwrite", overwrite)
	urlValues.Add("path", filepath)
	uploadLink, uploadErr := yadisk.GetLink("GET", "resources/upload", urlValues)
	if uploadErr != nil {
		return -1, uploadErr
	}
	file, fileErr := os.Open(filepath)
	if fileErr != nil {
		return -2, fileErr
	}
	rq, errRq := http.NewRequest("PUT", uploadLink.Href, file)
	if errRq != nil {
		return -3, errRq
	}
	resp, errRqDo := yadisk.ApiClient.HttpClient.Do(rq)
	if errRqDo != nil {
		return -4, errRqDo
	}
	// Обробатываем ошибки
	switch resp.StatusCode {
	case 412:
		return resp.StatusCode, utils.NewError("при дозагрузке файла был передан неверный диапазон в заголовке Content-Range.")
	case 413:
		return resp.StatusCode, utils.NewError("размер файла превышает 10 ГБ.")
	case 500:
	case 503:
		return resp.StatusCode, utils.NewError("ошибка сервера, попробуйте повторить загрузку.")
	case 507:
		return resp.StatusCode, utils.NewError("для загрузки файла не хватает места на Диске")
	}
	file.Close()
	return resp.StatusCode, nil
}
