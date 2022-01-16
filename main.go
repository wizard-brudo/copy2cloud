package main

import (
	"copy2cloud/oauth2"
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var token string
	// Если нет аргументов выдаём ошибку
	if len(os.Args) < 2 {
		fmt.Println(utils.NewError(utils.ERROR_ARGUMENTS))
		os.Exit(1)
	}
	if utils.FlagExists("--get-token") {
		// Если шаблоны существуют
		if utils.SystemResourceExists(oauth2.Wd+"/templates/index.html") == true && utils.SystemResourceExists(oauth2.Wd+"/templates/token.html") == true {
			// То получаем токен и выходим
			oauth2.GetToken()
			os.Exit(0)
		} else {
			fmt.Println(utils.NewError("Отсутствуют файлы шаблонов"))
			os.Exit(1)
		}
	}
	configFlag := utils.GetValueFlag("--config", "config.json")
	confFile, confErr := utils.GetConfigFile(configFlag)
	tokenFlag := utils.GetValueFlag("--token", "")

	// Если нет конфиг файла
	if confErr != nil {
		// То создаём конфиг
		config := map[string]string{
			"token":    tokenFlag,
			"log-file": "copy2cloud.log",
		}
		// Создаём лог файл
		logFile, err := os.Create("copy2cloud.log")
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
		}
		logFile.Close()
		// Создаём конфиг
		token = tokenFlag
		newConfFile, err := os.Create(configFlag)
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
			os.Exit(1)
		}
		data, _ := json.MarshalIndent(config, "", "\t")
		newConfFile.Write(data)
		newConfFile.Close()
	} else if tokenFlag == "" && confFile["token"] != "" {
		fmt.Println(2)
		// Если флаг токена не устоновлен и в конфиге есть токен то будем пользоваться им
		token = confFile["token"]
	} else if tokenFlag != "" && confFile["token"] == "" {
		// Если флаг токена устоновлен и в конфиге нет токен то будем пользоваться токеном из флага
		token = tokenFlag
	} else {
		fmt.Println(utils.NewError(confErr.Error()))
		os.Exit(1)
	}
	diskClient := NewDiskClient(token)
	// Потом смотрим что надо пользователю
	switch os.Args[1] {
	case "find":
		diskClient.FindCommand()
	case "move":
		diskClient.MoveCommand()
	case "trash":
		diskClient.TrashCommand()
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
		fmt.Println("0.4.1")
	case "help":
		fmt.Print(`Доступные команды:
	info [Путь к файлу/папке] - Выводит информацию о файле/папке, 
	если путь не задан,будет отображаться информация о диске.
	help - Показать сообщение этого сообщения.
	list - список файлов на диске.
	download [Путь к скачиваемому файлу] - скачать файл с диска.
	upload [Путь к загружаемому файлу] - загрузить файл на диск.
	copy [Путь к копируемому ресурсу] [Путь к создаваемой копии ресурса] - копирование файла/папки на диске
	move [Путь к перемещаемому ресурсу] [Путь к новому положению ресурса] - Перемещение файла/папки на диске
	delete [Путь к удаляемому ресурсу] - Удаление файла/папки на диске
	trash - Команды для работы с корзиной
		clear [Путь к удаляемому ресурсу] - Удаление ресурса из корзины
		restore [Путь к ресурсу который нужно восстоновить] - Восстоновить ресурс из корзины
	version - Версия программы copy2cloud.
Флаги:
	token - токен,который нужен для работы с вашим яндекс диском.
	style - стиль вывода списка файлов
	overwrite [true,false] - флаг перезаписи ресурса.
	permanently [true,false] - флаг безвозвратного удаления
	by-type - флаг указывает на то что нужно искать по типу
	get-token - Получить токен
	config - имя конфигурационого файла.
`)

	default:
		fmt.Println(utils.NewError(utils.ERROR_UNKNOWN_ARGUMENT))
	}
}
