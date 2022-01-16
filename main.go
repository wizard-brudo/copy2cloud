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
		if utils.SystemResourceExists(utils.Wd+"/templates/index.html") == true && utils.SystemResourceExists(oauth2.Wd+"/templates/token.html") == true {
			// То получаем токен и выходим
			oauth2.GetToken()
			os.Exit(0)
		} else {
			fmt.Println(utils.NewError("Отсутствуют файлы шаблонов"))
			os.Exit(1)
		}
	}
	configFlag := utils.GetValueFlag("--config", utils.Wd+"/config.json")
	mainConfig, confErr := utils.GetConfigFile(configFlag)
	tokenFlag := utils.GetValueFlag("--token", "")
	// Если нет конфиг файла
	if confErr != nil {
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Создаётся конфигурационный файл")
		}
		token = tokenFlag
		// То создаём конфиг
		config := map[string]string{
			"token":    tokenFlag,
			"log-file": utils.Wd + "/copy2cloud.log",
		}
		// Создаём лог файл
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Создаётся файл логгирования")
		}
		logFile, err := os.Create(utils.Wd + "/copy2cloud.log")
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
		}
		logFile.Close()
		newConfFile, err := os.Create(configFlag)
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
			os.Exit(1)
		}
		bytes, _ := json.MarshalIndent(config, "", "\t")
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Запись конфигурации в файл")
		}
		newConfFile.Write(bytes)
		newConfFile.Close()
	} else if tokenFlag == "" && mainConfig["token"] != "" {
		// Если флаг токена не устоновлен и в конфиге есть токен то будем пользоваться им
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Устоновка токена из конфигурационного файла")
		}
		token = mainConfig["token"]
	} else if tokenFlag != "" && mainConfig["token"] == "" {
		// Если флаг токена устоновлен и в конфиге нет токен то будем пользоваться токеном из флага
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Устоновка токена из флага")
		}
		token = tokenFlag
		config, _ := utils.GetConfigFile(configFlag)
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Устоновка токена из флага в конфиг")
		}
		config["token"] = tokenFlag
		// Открываем файл о очищаем его
		configFile, err := os.OpenFile(configFlag, os.O_TRUNC|os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(utils.NewError(err.Error()))
		}
		bytes, _ := json.MarshalIndent(config, "", "\t")
		configFile.Write(bytes)
		configFile.Close()
	} else if tokenFlag != "" && mainConfig["token"] != "" {
		// Если есть токены и в и в конфиге есть что-то
		// То приоритет отдаём флагу
		if utils.FlagExists("--verbose") == true {
			fmt.Println("Устоновка токена из флага и обновление токена в конфигурации")
		}
		config, _ := utils.GetConfigFile(configFlag)
		token = tokenFlag
		if config["token"] != tokenFlag {
			config["token"] = tokenFlag

			// Очищаем и открываем файл
			configFile, err := os.OpenFile(configFlag, os.O_TRUNC|os.O_RDWR, os.ModePerm)
			if err != nil {
				fmt.Println(utils.NewError(err.Error()))
			}
			bytes, _ := json.MarshalIndent(config, "", "\t")
			configFile.Write(bytes)
			configFile.Close()

		}
	}
	if utils.FlagExists("--verbose") == true {
		fmt.Println("Создание клиента")
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
	verbose - выводить детально
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
