package main

import (
	"copy2cloud/api/yandex"
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

func NewDiskClient(token string) DiskClient {
	diskClient := DiskClient{}
	yaDisk, err := yandex.NewYandexDisk(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	diskClient.yadisk = yaDisk
	return diskClient
}

type DiskClient struct {
	yadisk yandex.YandexDisk
}

func (dclient *DiskClient) ListCommand() {
	listFiles, err := dclient.yadisk.GetListFiles()
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	// Флаг стиля выведения файлов
	style := utils.GetValueFlag("--style", "1")
	w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
	if style == "1" {
		fmt.Fprintln(w, "ИМЯ\t\t\tРАЗМЕР(В байтах)\t\t\tДАТА СОЗДАНИЯ")
	}
	for _, item := range listFiles.Items {
		// 1 вариант(Таблица)
		if style == "1" {
			fmt.Fprintf(w, "%s\t\t\t%d\t\t\t%s\n", item.Name, item.Size, item.Created)
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
			fmt.Println(utils.ERROR_UNKOWN_STYLE)
			os.Exit(1)
		}
	}
	if style == "1" {
		// Рисуем таблицу
		w.Flush()
	}
}

func (dclient *DiskClient) DownloadCommand() {
	if len(os.Args) < 3 {
		fmt.Println(utils.ERROR_NOT_ENOUGH_ARGUMENTS)
		os.Exit(1)
	}
	dclient.Download(os.Args[2], utils.GetValueFlag("--overwrite", "false"))
}

func (dclient *DiskClient) InfoCommand() {
	if len(os.Args) == 2 || len(os.Args) == 4 {
		// выводим информацию о диске
		diskInfo := yandex.Disk{}
		response, _, err := dclient.yadisk.ApiClient.SendRequest("GET", "disk/")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
		resource, err := dclient.yadisk.GetMetaInformation(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Имя:", resource.Name)
		fmt.Println("Дата модификации:", resource.Modified)
		fmt.Println("Cсылка на ресурс:", resource.PublicUrl)
		fmt.Println("Путь:", resource.Path)
		fmt.Println("Тип:", resource.Type)
		fmt.Println("Мим-тип:", resource.MimeType)
		fmt.Println("Размер:", resource.Size)
		if len(resource.Embedded.Items) > 0 {
			fmt.Println("Содержит:")
			w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
			fmt.Fprintln(w, "\tИМЯ\t\tРАЗМЕР(В байтах)\t\tДАТА СОЗДАНИЯ")
			for _, item := range resource.Embedded.Items {
				fmt.Fprintf(w, "\t%s\t\t%d\t\t%s\n", item.Name, item.Size, item.Created)
			}
			w.Flush()
		}
	}
}

func (dclient *DiskClient) UploadCommand() {
	if len(os.Args) < 3 {
		fmt.Println(utils.ERROR_NOT_ENOUGH_ARGUMENTS)
		os.Exit(1)
	}
	overwriteFlag := utils.GetValueFlag("--overwrite", "false")
	dclient.Upload(os.Args[2], overwriteFlag)
}

func (dclient *DiskClient) CopyCommand() {
	if len(os.Args) < 4 {
		fmt.Println(utils.ERROR_NOT_ENOUGH_ARGUMENTS)
		os.Exit(1)
	}
	_, copyErr := dclient.yadisk.CopyResource(os.Args[2], os.Args[3], utils.GetValueFlag("--overwrite", "false"))
	if copyErr != nil {
		fmt.Println(copyErr)
		os.Exit(1)
	}
	fmt.Println("Ресурс скопирован")
}

func (dclient *DiskClient) Download(resourcePath, overwrite string) {
	mainResource, err := dclient.yadisk.GetMetaInformation(resourcePath)
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	resourceExists := utils.SystemResourceExists(resourcePath)
	if resourceExists == false || overwrite == "true" {
		if mainResource.IsDir() {
			fmt.Println("Создание главной папки", mainResource.CorrectPath())
			os.Mkdir(mainResource.CorrectPath(), os.ModePerm)
			items, err := dclient.yadisk.GetAllResources(mainResource.Path)
			if err != nil {
				fmt.Println(utils.NewError(err.Error()))
				os.Exit(1)
			}
			for _, item := range items {
				if item.IsDir() {
					fmt.Println("Создание папки", item.CorrectPath())
					os.Mkdir(item.CorrectPath(), os.ModePerm)
				} else {
					fmt.Println("Загрузка ", item.CorrectPath())
					dclient.yadisk.DownloadFile(item.Path, overwrite)
				}
			}
		} else {
			fmt.Println("Загрузка " + resourcePath)
			dclient.yadisk.DownloadFile(resourcePath, overwrite)
		}
	} else {
		fmt.Println(utils.ERROR_RESOURCE_EXISTS)
		os.Exit(1)
	}
}

func (dclient *DiskClient) MoveCommand() {
	if len(os.Args) < 4 {
		fmt.Println(utils.ERROR_NOT_ENOUGH_ARGUMENTS)
		os.Exit(1)
	}
	_, moveErr := dclient.yadisk.MoveResource(os.Args[2], os.Args[3], utils.GetValueFlag("--overwrite", "false"))
	if moveErr != nil {
		fmt.Println(moveErr)
		os.Exit(1)
	}
	fmt.Println("Ресурс перемещен")
}

func (dclient *DiskClient) DeleteCommand() {
	if len(os.Args) < 3 {
		fmt.Println(utils.ERROR_NOT_ENOUGH_ARGUMENTS)
		os.Exit(1)
	}
	_, deleteErr := dclient.yadisk.DeleteResource(os.Args[2], utils.GetValueFlag("--permanently", "false"))
	if deleteErr != nil {
		fmt.Println(deleteErr)
		os.Exit(1)
	}
	fmt.Println("Ресурс удалён")
}

func (dclient *DiskClient) Upload(resourceName, overwrite string) {
	isDir, err := utils.IsWritableDir(resourceName)
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	isYaDiskResourceExists := dclient.yadisk.ResourceExists(resourceName)
	// Мы загрузжаем ресурсы,если ресурса не существует на диске или указан флаг --overwrite true
	if isYaDiskResourceExists == false || overwrite == "true" {
		if isDir == true {
			resources := []yandex.OsResource{}
			// Получаем все пути(пути папок и файлов в директории)
			filepath.Walk(resourceName, func(path string, info os.FileInfo, _ error) error {
				resources = append(resources, yandex.OsResource{ItemPath: path, Info: info})
				return nil
			})
			// Если в папке что то есть то загружаем папку
			// иначе просто создаём папку
			if len(resources) > 0 {
				// Проходимся по списку ресурсов полученные через filepath.Walk
				for _, resource := range resources {
					// Если это папка то создаём папку
					if resource.Info.IsDir() {
						_, DirError := dclient.yadisk.CreateDir(resource.ItemPath)
						if DirError != nil && overwrite == "false" {
							// Если есть ошибка и нету флага перезаписи то выдаём ошибку
							fmt.Println(utils.NewError(DirError.Error()))
							os.Exit(1)
						}
					} else {
						// Иначе это файл и загрузжаем его
						UploadError := dclient.yadisk.UploadFile(resource.ItemPath, overwrite)
						if UploadError != nil {
							fmt.Println(utils.NewError(UploadError.Error()))
							os.Exit(1)
						}
					}
				}
			}
		} else {
			// Если это просто файл то загрузжаем его
			UploadError := dclient.yadisk.UploadFile(resourceName, overwrite)
			if UploadError != nil {
				fmt.Println(utils.NewError(UploadError.Error()))
				os.Exit(1)
			}
		}
	} else {
		fmt.Println(utils.ERROR_RESOURCE_EXISTS)
		os.Exit(1)
	}

}
