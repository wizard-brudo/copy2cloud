package yandex

import (
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (c *ApiClient) CheckResponse(response http.Response) error {
	if response.StatusCode >= 400 {
		errApi := ErrorApi{}
		body_byte, _ := ioutil.ReadAll(response.Body)
		err := json.Unmarshal(body_byte, &errApi)
		if err != nil {
			return utils.ERROR_JSON
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
		return nil, -1, check_error
	}
	// Читаем ответ(переводим в массив байтов)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body, resp.StatusCode, nil
}

func NewYandexDisk(token string) (YandexDisk, error) {
	if token == "" {
		return YandexDisk{}, utils.ERROR_NO_TOKEN
	}
	newApiClient := ApiClient{
		HttpClient: http.Client{},
		Token:      token,
		BaseUrl:    "https://cloud-api.yandex.net/v1/",
	}
	return YandexDisk{
		ApiClient: newApiClient,
	}, nil
}

func (yadisk *YandexDisk) GetMetaInformation(resourcePath string) (Resource, error) {
	resource := Resource{}
	response, _, err := yadisk.ApiClient.SendRequest("GET", "disk/resources?path="+url.QueryEscape(resourcePath))
	if err != nil {
		return resource, err
	}
	json.Unmarshal(response, &resource)
	return resource, nil
}

func (yadisk *YandexDisk) GetListFiles() (FilesResourceList, error) {
	listFiles := FilesResourceList{}
	response, _, err := yadisk.ApiClient.SendRequest("GET", "disk/resources/files")
	if err != nil {
		return listFiles, err
	}
	json.Unmarshal(response, &listFiles)
	return listFiles, nil
}

func (yadisk *YandexDisk) GetAllResources(rootPath string) ([]Resource, error) {
	resources := []Resource{}
	yaDiskResources, err := yadisk.GetMetaInformation(rootPath)
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

func (yadisk *YandexDisk) CreateDir(dirname string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", dirname)
	createDir, err := yadisk.GetLink("PUT", "disk/resources", urlValues)
	if err != nil {
		return createDir, err // CreateDir в данном случае пустой
	}
	return createDir, nil
}

func (yadisk *YandexDisk) CopyResource(from, path, overwrite string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("from", from)
	urlValues.Add("path", path)
	urlValues.Add("overwrite", overwrite)
	copyResource, err := yadisk.GetLink("POST", "disk/resources/copy", urlValues)
	if err != nil {
		return copyResource, err // copyResource  пустой
	}
	return copyResource, nil
}

func (yadisk *YandexDisk) DeleteResource(path, permanently string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("path", path)
	urlValues.Add("permanently", permanently)
	deleteResource, err := yadisk.GetLink("DELETE", "disk/resources", urlValues)
	if err != nil {
		return deleteResource, err // copyResource  пустой
	}
	return deleteResource, nil
}

func (yadisk *YandexDisk) MoveResource(from, path, overwrite string) (Link, error) {
	urlValues := url.Values{}
	urlValues.Add("from", from)
	urlValues.Add("path", path)
	urlValues.Add("overwrite", overwrite)
	moveResource, err := yadisk.GetLink("POST", "disk/resources/move", urlValues)
	if err != nil {
		return moveResource, err // moveResource пустой
	}
	return moveResource, nil
}

func (yadisk *YandexDisk) ResourceExists(resourceName string) bool {
	_, err := yadisk.GetMetaInformation(resourceName)
	if err != nil {
		return false
	}
	return true
}

func (yadisk *YandexDisk) DownloadFile(filepath, overwrite string) error {
	urlValues := url.Values{}
	urlValues.Add("path", filepath)
	urlValues.Add("overwrite", overwrite)
	downloadLink, LinkErr := yadisk.GetLink("GET", "disk/resources/download", urlValues)
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

func (yadisk *YandexDisk) UploadFile(filepath, overwrite string) error {
	fmt.Println("Идёт загрузка " + filepath)
	urlValues := url.Values{}
	urlValues.Add("overwrite", overwrite)
	urlValues.Add("path", filepath)
	uploadLink, uploadErr := yadisk.GetLink("GET", "disk/resources/upload", urlValues)
	if uploadErr != nil {
		return uploadErr
	}
	file, fileErr := os.Open(filepath)
	if fileErr != nil {
		return fileErr
	}
	rq, errRq := http.NewRequest("PUT", uploadLink.Href, file)
	if errRq != nil {
		return errRq
	}
	resp, errRqDo := yadisk.ApiClient.HttpClient.Do(rq)
	if errRqDo != nil {
		return errRqDo
	}
	// Обробатываем ошибки
	switch resp.StatusCode {
	case 412:
		return utils.NewError("при дозагрузке файла был передан неверный диапазон в заголовке Content-Range.")
	case 413:
		return utils.NewError("размер файла превышает 10 ГБ.")
	case 500:
	case 503:
		return utils.NewError("ошибка сервера, попробуйте повторить загрузку.")
	case 507:
		return utils.NewError("для загрузки файла не хватает места на Диске")
	}
	file.Close()
	return nil
}
