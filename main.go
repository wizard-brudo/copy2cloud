package main

import (
	"fmt"
	"os"
)

func main() {
	var token string
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
		if tokenFlag != "" {
			token = tokenFlag
		}
	} else {
		token = confFile["token"]
	}
	yaDisk, err := NewYaDisk(token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Потом смотрим что надо пользователю
	switch os.Args[1] {
	case "info":
		yaDisk.ShowInfo()
	case "list":
		yaDisk.ShowListFile()
	case "upload":
		if len(os.Args) >= 3 {
			yaDisk.UploadFile(os.Args[2])
			os.Exit(1)
		}
		fmt.Println(ERROR_NOT_ENOUGH_ARGUMENTS)
	case "download":
		if len(os.Args) >= 3 {
			yaDisk.DownloadFile(os.Args[2])
			os.Exit(1)
		}
		fmt.Println(ERROR_NOT_ENOUGH_ARGUMENTS)
	case "version":
		fmt.Println("0.1")
	case "help":
		fmt.Print(`Доступные команды:
	info [Путь к файл/папке] - Вывод информации о файле/папке если путь не устоновлен будет выведена информации о диске
	help - Вывод сообщения этого сообщения
	list - список файлов на диске
	download [Путь к скачиваемому файлу] - скачать файл с диска
	upload [Путь к загружаемому файл] - загрузить файл на диск
	version - Версия прогграмы copy2cloud
`)

	default:
		fmt.Println(ERROR_UNKNOWN_ARGUMENT)
	}
}
