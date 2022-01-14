package main

import (
<<<<<<< HEAD
	"copy2cloud/utils"
=======
>>>>>>> 6e6eb226087cd70441ee2922ecf10aec15135e77
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var token string
<<<<<<< HEAD
	var configName string = "config.json"
	// Если нет аргументов выдаём ошибку
	if len(os.Args) < 2 {
		fmt.Println(utils.ERROR_ARGUMENTS)
		os.Exit(1)
	}
	configFlag := utils.GetValueFlag("--config", "config.json")
	if configFlag != "config.json" {
		configName = configFlag
	}
	confFile, err := utils.GetConfigFile(configName)
	// Если конфигурационого файла нет или токен не устоновлен
	if err != nil || confFile["token"] == "" {
		// Если есть флаг токена то ставим его значение  как токен
		tokenFlag := utils.GetValueFlag("--token", "")
=======
	// Если нет аргументов выдаём ошибку
	if len(os.Args) < 2 {
		fmt.Println(ERROR_ARGUMENTS)
		os.Exit(1)
	}
	confFile, err := getConfigFile()
	// Если конфигурационого файла нет или токен не устоновлен
	if err != nil || confFile["token"] == "" {
		// Если есть флаг токена то ставим его значение  как токен
		tokenFlag := getValueFlag("--token", "")
>>>>>>> 6e6eb226087cd70441ee2922ecf10aec15135e77
		if tokenFlag != "" {
			config := map[string]string{
				"token": tokenFlag,
			}
			token = tokenFlag
			newConfFile, _ := os.Create("config.json")
			data, _ := json.Marshal(config)
			newConfFile.Write(data)
			newConfFile.Close()
		}
	} else {
		token = confFile["token"]
	}
<<<<<<< HEAD
	diskClient := NewDiskClient(token)
	// Потом смотрим что надо пользователю
	switch os.Args[1] {
	case "move":
		diskClient.MoveCommand()
	case "delete":
		diskClient.DeleteCommand()
	case "info":
		diskClient.InfoCommand()
	case "copy":
		diskClient.CopyCommand()
	case "list":
		diskClient.ListCommand()
	case "upload":
		diskClient.UploadCommand()
	case "download":
		diskClient.DownloadCommand()
	case "version":
		fmt.Println("0.3")
=======
	yaDisk, err := NewYaDisk(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Потом смотрим что надо пользователю
	switch os.Args[1] {
	case "info":
		if len(os.Args) == 2 || len(os.Args) == 4 {
			// выводим информацию о диске
			diskInfo := Disk{}
			response, _ := yaDisk.Client.SendRequest("GET", "disk/")
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
			resource := yaDisk.getMetaInformation(os.Args[2])
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
	case "list":
		listFiles := yaDisk.getListFiles()
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
	case "upload":
		if len(os.Args) >= 3 {
			yaDisk.Upload(os.Args[2])
			os.Exit(1)
		}
		fmt.Println(ERROR_NOT_ENOUGH_ARGUMENTS)
	case "download":
		if len(os.Args) >= 3 {
			yaDisk.Download(os.Args[2])
			os.Exit(1)
		}
		fmt.Println(ERROR_NOT_ENOUGH_ARGUMENTS)

	case "version":
		fmt.Println("0.2.1")
>>>>>>> 6e6eb226087cd70441ee2922ecf10aec15135e77
	case "help":
		fmt.Print(`Доступные команды:
	info [Путь к файлу/папке] - Выводит информацию о файле/папке, 
	если путь не задан,будет отображаться информация о диске.
	help - Показать сообщение этого сообщения.
	list - список файлов на диске.
	download [Путь к скачиваемому файлу] - скачать файл с диска.
	upload [Путь к загружаемому файлу] - загрузить файл на диск.
<<<<<<< HEAD
	copy [Путь к копируемому ресурсу] [Путь к создаваемой копии ресурса] - копирование файла/папки на диске
	move [Путь к перемещаемому ресурсу] [Путь к новому положению ресурса] - Перемещение файла/папки на диске
	delete [Путь к удаляемому ресурсу] - Удаление файла/папки на диске
	version - Версия программы copy2cloud.
Флаги:
	token - токен,который нужен для работы с вашим яндекс диском.
	style - стиль вывода списка файлов
	overwrite [true,false] - флаг перезаписи ресурса.
	permanently [true,false] - флаг безвозвратного удаления
	config - имя конфигурационого файла.
`)

	default:
		fmt.Println(utils.ERROR_UNKNOWN_ARGUMENT)
=======
	version - Версия программы copy2cloud.
`)

	default:
		fmt.Println(ERROR_UNKNOWN_ARGUMENT)
>>>>>>> 6e6eb226087cd70441ee2922ecf10aec15135e77
	}
}
