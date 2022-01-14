package main

import (
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var token string
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
	}
}
