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
	confFile, err := utils.GetConfigFile(configFlag)

	// Если конфиг файла нет или в конфиг файле нет токена
	if err != nil || confFile["token"] == "" {
		// То получаем флаг токена
		tokenFlag := utils.GetValueFlag("--token", "")
		// Если флаг токена устоновлен
		if tokenFlag != "" {
			// То создаём карту с токеном
			config := map[string]string{
				"token": tokenFlag,
				"log":   "copy2cloud.log",
			}
			// Создаём лог файл
			logFile, err := os.Create("copy2cloud.log")
			if err != nil {
				fmt.Println(utils.NewError(err.Error()))
			}
			logFile.Close()
			// Создаём конфиг
			token = tokenFlag
			newConfFile, _ := os.Create(configFlag)
			data, _ := json.Marshal(config)
			newConfFile.Write(data)
			newConfFile.Close()
		}
	} else {
		token = confFile["token"]
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
