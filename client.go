package main

import (
	"copy2cloud/api/yandex"
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	listFiles, err := dclient.yadisk.GetListFiles("")
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	// Флаг стиля выведения файлов
	style := utils.GetValueFlag("--style", "1")
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
	if style == "1" {
		fmt.Fprintln(w, "ПУТЬ\t\t\tРАЗМЕР(В байтах)\t\t\tДАТА СОЗДАНИЯ")
	}
	for _, item := range listFiles.Items {
		// 1 вариант(Таблица)
		if style == "1" {
			fmt.Fprintf(w, "%s\t\t\t%d\t\t\t%s\n", strings.Replace(item.Path, "disk:/", "", 1), item.Size, item.Created)
		} else if style == "2" {
			// 2 Вариант
			fmt.Println("-----------------------------------------------")
			// Вывод информации о ресурсе
			dclient.printResourceFields(item)
			fmt.Println("-----------------------------------------------")
		} else {
			fmt.Println(utils.NewError(utils.ERROR_UNKOWN_STYLE))
			os.Exit(1)
		}
	}
	if style == "1" {
		// Рисуем таблицу
		w.Flush()
	}
}

// Вывод Ресурса
func (dclient *DiskClient) printResourceFields(resource yandex.Resource) {
	if resource.PublicUrl == "" {
		resource.PublicUrl = "Нет"
	}
	if resource.MimeType == "" {
		resource.MimeType = "Неизвестен"
	}
	fmt.Printf(`Имя: %s
Дата модификации: %s
Cсылка на ресурс: %s
Путь: %s
Тип: %s
Мим-тип: %s
Размер: %d
`, resource.Name,
		resource.Modified,
		resource.PublicUrl,
		resource.CorrectPath(),
		resource.Type,
		resource.MimeType,
		resource.Size)
	if len(resource.Embedded.Items) > 0 {
		fmt.Println("Содержит:")
		w := tabwriter.NewWriter(os.Stdout, 10, 0, 1, ' ', 0)
		fmt.Fprintln(w, "\tПУТЬ\t\tДАТА СОЗДАНИЯ")
		for _, item := range resource.Embedded.Items {
			fmt.Fprintf(w, "\t%s\t\t%s\n", item.CorrectPath(), item.Created)
		}
		w.Flush()
	}
}

func (dclient *DiskClient) FindCommand() {
	if len(os.Args) < 4 {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
		os.Exit(1)
	}
	if utils.FlagExists("--by-type") == true {
		listFiles, errList := dclient.yadisk.GetListFiles(os.Args[2])
		if errList != nil {
			fmt.Println(utils.NewError(errList.Error()))
			os.Exit(1)
		}
		for _, resource := range listFiles.Items {
			fmt.Println("---------------------------------------------------------------")
			dclient.printResourceFields(resource)
			fmt.Println("---------------------------------------------------------------")
		}
	} else {
		desiredResource, err := dclient.yadisk.Find(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
			os.Exit(1)
		}
		dclient.printResourceFields(desiredResource)
	}
}

func (dclient *DiskClient) TrashCommand() {
	if len(os.Args) > 2 {
		if os.Args[2] == "clear" {
			_, errClear := dclient.yadisk.ClearTrash(os.Args[3])
			if errClear != nil {
				fmt.Println(utils.NewError(errClear.Error()))
				os.Exit(1)
			}
			fmt.Println("Ресурс удалён из корзины")
		} else if os.Args[2] == "restore" {
			_, errRestore := dclient.yadisk.RestoreTrash(os.Args[3], utils.GetValueFlag("--overwrite", "false"))
			if errRestore != nil {
				fmt.Println(utils.NewError(errRestore.Error()))
				os.Exit(1)
			}
			fmt.Println("Ресурс восстановлен из корзины")
		} else if os.Args[2] == "info" {
			if len(os.Args) > 3 {
				resource, errInfo := dclient.yadisk.GetTrashResource(os.Args[3])
				if errInfo != nil {
					fmt.Println(utils.NewError(errInfo.Error()))
					os.Exit(1)
				}
				// Вывод информации о ресурсе
				dclient.printResourceFields(resource)
			} else {
				fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
				os.Exit(1)
			}
		} else {
			fmt.Println(utils.NewError("неизвестная команда"))
			os.Exit(1)
		}
	} else {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
		os.Exit(1)
	}
}

func (dclient *DiskClient) DownloadCommand() {
	if len(os.Args) < 3 {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
		os.Exit(1)
	}
	dclient.Download(os.Args[2], utils.GetValueFlag("--overwrite", "false"))
}

func (dclient *DiskClient) InfoCommand() {
	if len(os.Args) == 2 || len(os.Args) == 4 || len(os.Args) == 5 {
		// выводим информацию о диске
		diskInfo := yandex.Disk{}
		response, _, err := dclient.yadisk.ApiClient.SendRequest("GET", "")
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
	} else if len(os.Args) >= 3 && utils.FlagExists("--"+os.Args[2]) == false {
		fmt.Println(utils.FlagExists("--" + os.Args[2]))
		fmt.Println(os.Args)
		resource, err := dclient.yadisk.GetResource(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dclient.printResourceFields(resource)
	}
}

func (dclient *DiskClient) UploadCommand() {
	if len(os.Args) < 3 {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
		os.Exit(1)
	}
	overwriteFlag := utils.GetValueFlag("--overwrite", "false")
	dclient.Upload(os.Args[2], overwriteFlag)
}

func (dclient *DiskClient) CopyCommand() {
	if len(os.Args) < 4 {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
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
	if utils.FlagExists("--verbose") == true {
		fmt.Println("Получение информации о ресурсе")
	}
	mainResource, err := dclient.yadisk.GetResource(resourcePath)
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	resourceExists := utils.SystemResourceExists(resourcePath)
	if resourceExists == false || overwrite == "true" {
		if mainResource.IsDir() {
			fmt.Println("Создание главной папки", mainResource.CorrectPath())
			os.Mkdir(mainResource.CorrectPath(), os.ModePerm)
			if utils.FlagExists("--verbose") == true {
				fmt.Println("Получение информации о ресурсах в папках")
			}
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
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
		os.Exit(1)
	}
}

func (dclient *DiskClient) MoveCommand() {
	if len(os.Args) < 4 {
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
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
		fmt.Println(utils.NewError(utils.ERROR_NOT_ENOUGH_ARGUMENTS))
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
			dirResources := []yandex.OsResource{}
			// Получаем все пути(пути папок и файлов в директории)
			filepath.Walk(resourceName, func(path string, info os.FileInfo, _ error) error {
				dirResources = append(dirResources, yandex.OsResource{ItemPath: path, Info: info})
				return nil
			})
			// Если в папке что то есть то загружаем папку
			// иначе просто создаём папку
			if len(dirResources) > 0 {
				// Проходимся по списку ресурсов полученные через filepath.Walk
				for _, resource := range dirResources {
					// Если это папка то создаём папку
					if resource.Info.IsDir() {
						_, DirError := dclient.yadisk.CreateDir(resource.ItemPath)
						if DirError != nil && overwrite == "false" {
							// Если есть ошибка и нету флага перезаписи то выдаём ошибку
							fmt.Println(utils.NewError(DirError.Error()))
							os.Exit(1)
						}
					} else {
						fmt.Println("Идёт загрузка ", resource.ItemPath)
						// Иначе это файл и загрузжаем его
						StatusCode, UploadError := dclient.yadisk.UploadFile(resource.ItemPath, overwrite)
						if utils.FlagExists("--verbose") == true {
							fmt.Println("Статус операции:", StatusCode)
						}
						if UploadError != nil {
							fmt.Println(utils.NewError(UploadError.Error()))
							os.Exit(1)
						}
					}
				}
			}
		} else {
			// Если это просто файл то загрузжаем его
			fmt.Println("Идёт загрузка ресурса", resourceName)
			StatusCode, UploadError := dclient.yadisk.UploadFile(resourceName, overwrite)
			if utils.FlagExists("--verbose") == true {
				fmt.Println("Статус операции:", StatusCode)
			}
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
