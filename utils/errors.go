package utils

import (
	"errors"
	"strings"
)

const (
	ERROR_ARGUMENTS            = "Нет аргументов запуска"
	ERROR_UNKNOWN_ARGUMENT     = "Неизвестный аргумент "
	ERROR_JSON                 = "Неправильный json "
	ERROR_NOT_ENOUGH_ARGUMENTS = "Недостаточно аргументов "
	ERROR_UNKOWN_STYLE         = "Неизвестный стиль"
	ERROR_WRITE_BIT_NOT_SET    = "В этом файле не установлен бит разрешения на запись  "
	ERROR_STAT                 = "Не удалось получить статистику"
	ERROR_NO_PERMISSION        = "У вас нет прав на запись в этот каталог"
	ERROR_NOT_DIRECTORY        = "Путь не является каталогом "
	ERROR_PATH_NOT_EXISTS      = "Путь не существует "
	ERROR_NO_TOKEN             = "Нет токена"
	ERROR_TOO_MANY_ARGUMENTS   = "Слишком много аргументов"
	ERROR_NO_RESOURCES         = "Нет ресурсов"
	ERROR_RESOURCE_EXISTS      = "Ресурс уже есть на вашем устройстве(диске),используйте флаг --overwrite true чтобы перезаписать ресурс"
)

func NewError(text string) error {
	errorMessage := text
	if strings.Contains(errorMessage, "Ошибка:") == true {
		errorMessage = strings.Replace(text, "Ошибка: ", "", 1)
	}
	err := WriteLog(text)
	if err != nil {
		errorMessage += "\nОшибка при записи лога: " + err.Error()
	}
	return errors.New(SetTextColor("Ошибка: "+errorMessage, RED))
}
