package utils

import (
	"encoding/json"
	"os"
	"strings"
	"syscall"
)

func IndexOf[T comparable](DesiredValue T, array []T) int {
	for index, value := range array {
		if value == DesiredValue {
			return index
		}
	}
	return -1
}

// Кодирует карту в строку
func UrlEncode(parameters map[string]string) string {
	var urlParametersString string
	if len(parameters) > 0 {
		for key, value := range parameters {
			urlParametersString += key + "=" + value + "&"
		}
	}
	return urlParametersString
}

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
func GetValueFlag(Flag string, value string) string {
	IndexFlag := IndexOf(Flag, os.Args)
	// Если не нашли значение флага
	// Возращаем значение по умолчанию
	if IndexFlag == -1 {
		return value
	}
	return os.Args[IndexFlag+1]
}

func IsFlag(arg string) bool {
	if strings.Contains(arg, "--") || strings.Contains(arg, "-") {
		return true
	}
	return false
}

func IsWritableDir(path string) (isWritable bool, err error) {
	var stat syscall.Stat_t
	info, err := os.Stat(path)
	if err != nil {
		return false, ERROR_PATH_NOT_EXISTS
	}
	if !info.IsDir() {
		return false, ERROR_NOT_DIRECTORY
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, ERROR_WRITE_BIT_NOT_SET
	}

	if err = syscall.Stat(path, &stat); err != nil {
		return false, ERROR_STAT
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, ERROR_NO_PERMISSION
	}

	return true, nil
}

func GetConfigFile(configName string) (map[string]string, error) {
	configFile := map[string]string{}
	dataConfFile, fileErr := os.ReadFile(configName)
	if fileErr != nil {
		return nil, NewError(fileErr.Error())
	}
	json.Unmarshal(dataConfFile, &configFile)
	return configFile, nil
}
