package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	file, _ = os.Executable()
	Wd      = filepath.Dir(file)
)

// Функция проверяет существует ли ресурс в папке которой находиться прогграма
func SystemResourceExists(resourceName string) bool {
	_, err := os.Stat(resourceName)
	if err != nil {
		return false
	}
	return true
}

/*
	Аргумент Flag это искомые флаг.
	Аргумент Value это значение если flag не был найден
*/
func GetValueFlag(flag string, value string) string {

	for index, arg := range os.Args {
		if arg == flag {
			// Если значение есть
			if index+1 < len(os.Args) {
				// То возращаем его
				return os.Args[index+1]
			}
		}
	}
	// Если не нашли флаг
	// Возращаем значение по умолчанию
	return value
}

// Проверяет существует ли флаг
func FlagExists(flag string) bool {
	for _, arg := range os.Args {
		if arg == flag {
			return true
		}
	}
	return false
}

func IsWritableDir(path string) (isWritable bool, err error) {
	var stat syscall.Stat_t
	info, err := os.Stat(path)
	if err != nil {
		return false, errors.New(ERROR_PATH_NOT_EXISTS)
	}
	if !info.IsDir() {
		return false, errors.New(ERROR_NOT_DIRECTORY)
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, errors.New(ERROR_WRITE_BIT_NOT_SET)
	}

	if err = syscall.Stat(path, &stat); err != nil {
		return false, errors.New(ERROR_STAT)
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, errors.New(ERROR_NO_PERMISSION)
	}

	return true, nil
}

func GetConfigFile(configFileName string) (map[string]string, error) {
	config := map[string]string{}
	dataConfFile, fileErr := os.ReadFile(configFileName)
	if fileErr != nil {
		return nil, fileErr
	}
	json.Unmarshal(dataConfFile, &config)
	return config, nil
}

func WriteLog(text string) error {
	// Читаем конфиг
	config, err := GetConfigFile(GetValueFlag("--config", "config.json"))
	if err != nil {
		return err
	}
	if config["log-file"] == "" {
		// Мы здесь используем errors.New чтобы не уйти в рекурсию
		return errors.New("В конфигурационном файле не определенно имя файла логгирования")
	}
	// Открываем лог файл
	LogFile, err := os.OpenFile(config["log-file"], os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	_, errWrite := LogFile.WriteString(text + " " + time.Now().Format(time.RFC3339) + "\n")
	if errWrite != nil {
		return err
	}
	LogFile.Close()
	return nil

}
