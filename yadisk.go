package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	HttpClient http.Client
	Token      string
	BaseUrl    string
}

func (c *Client) SendRequest(api_page string) []byte {
	// Отправляем запрос
	req, err := http.NewRequest(
		"GET", c.BaseUrl+api_page, nil,
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
	return body
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

func (yad *Yadisk) ShowInfo() {
	// если нам дали только команду info то аргументов два
	// если с токеном то аргументов 4
	if len(os.Args) == 2 || len(os.Args) == 4 {
		// выводим информацию о диске
		diskInfo := Disk{}
		response := yad.Client.SendRequest("disk/")
		json.Unmarshal(response, &diskInfo)
		fmt.Printf(`Логин:%s
Общие пространство:%g
Размер корзины:%g
Использованное пространство:%g
Системные папки:
`, diskInfo.User["login"],
			diskInfo.TotalSpace,
			diskInfo.TrashSize,
			diskInfo.UsedSpace)
		for folderName, path := range diskInfo.SystemFolders {
			fmt.Printf("	%s  %s\n", folderName, path)
		}
	} else if len(os.Args) >= 3 {
		resource := Resource{}
		response := yad.Client.SendRequest("disk/resources?path=" + os.Args[2])
		json.Unmarshal(response, &resource)
		fmt.Println("Имя:", resource.Name)
		fmt.Println("Дата модификации:", resource.Modified)
		fmt.Println("Cсылка на ресурс:", resource.PublicUrl)
		fmt.Println("Путь:", resource.Path)
		fmt.Println("Тип:", resource.Type)
		fmt.Println("Мим-тип:", resource.MimeType)
		fmt.Println("Размер:", resource.Size)
		if len(resource.Embedded.Items) > 0 {
			fmt.Println("Содержит:")
			fmt.Println("	ИМЯ 						     РАЗМЕР(В байтах) 			ДАТА СОЗДАНИЯ")
			for _, item := range resource.Embedded.Items {
				if len(item.Name) <= 6 {
					// Формат для мальньких
					fmt.Printf("	%s\t\t\t\t\t\t\t\b\b\b%d\t\t\t\t%s\n", item.Name, item.Size, item.Created)
				} else if len(item.Name) <= 16 && len(item.Name) > 6 {
					// Формат для средних
					fmt.Printf("	%s\t\t\t\t\t\t\b\b\b%d\t\t\t%s\n", item.Name, item.Size, item.Created)
				} else {
					// Формат для больших
					fmt.Printf("	%s                                   %d                   %s\n", item.Name, item.Size, item.Created)
				}
			}
		}
	}
}

func (yad *Yadisk) ShowListFile() {
	listFiles := FilesResourceList{}
	response := yad.Client.SendRequest("disk/resources/files")
	json.Unmarshal(response, &listFiles)
	// Флаг стиля выведения файлов
	style := getValueFlag("--style", "1")
	if style == "1" {
		fmt.Println("ИМЯ 						     РАЗМЕР(В байтах) 			ДАТА СОЗДАНИЯ")
	}
	for _, item := range listFiles.Items {
		// 1 вариант(Таблица)
		if style == "1" {
			if len(item.Name) <= 6 {
				// Формат для мальньких
				fmt.Printf("%s\t\t\t\t\t\t\t\b\b\b%d\t\t\t\t%s\n", item.Name, item.Size, item.Created)
			} else if len(item.Name) <= 16 && len(item.Name) > 6 {
				// Формат для средних
				fmt.Printf("%s\t\t\t\t\t\t\b\b\b%d\t\t\t%s\n", item.Name, item.Size, item.Created)
			} else {
				// Формат для больших
				fmt.Printf("%s                                   %d                   %s\n", item.Name, item.Size, item.Created)
			}
		} else if style == "2" {
			// 2 Вариант
			fmt.Println("-----------------------------------------------")
			fmt.Println("Имя:", item.Name)
			fmt.Println("Дата создания:", item.Created)
			fmt.Println("Дата модификации:", item.Modified)
			fmt.Println("Путь(на диске):", item.Path)
			fmt.Println("Тип:", item.Type)
			fmt.Println("Мим-тип:", item.MimeType)
			fmt.Println("-----------------------------------------------")
		} else {
			fmt.Println(ERROR_UNKOWN_STYLE)
			os.Exit(1)
		}
	}
}

func (yad *Yadisk) DownloadFile(filename string) {
	link := Link{}
	download_url := yad.Client.SendRequest("disk/resources/download?path=" + url.QueryEscape(filename))
	json.Unmarshal(download_url, &link)
	resp, download_error := yad.Client.HttpClient.Get(link.Href)
	if download_error != nil {
		fmt.Println(setTextColor("Unable to get file:"+download_error.Error(), RED))
		os.Exit(1)
	}
	file, err := os.Create(filename)
	if err != nil { // если возникла ошибка
		fmt.Println(setTextColor("Unable to create file:"+err.Error(), RED))
		os.Exit(1) // выходим из программы
	}
	data, _ := io.ReadAll(resp.Body)
	file.Write(data)
	file.Close()

}

func (yad *Yadisk) UploadFile(filepath string) {
	link := Link{}
	fmt.Println("Идёт загрузка " + filepath)
	upload_url := yad.Client.SendRequest("disk/resources/upload?path=" + url.QueryEscape(filepath))
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
