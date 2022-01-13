package utils

import (
	"errors"
)

var (
	ERROR_ARGUMENTS            = NewError("нет аргументов запуска")
	ERROR_LOGIN                = NewError("неверное имя пользователя и/или пароль ")
	ERROR_UNKNOWN_ARGUMENT     = NewError("неизвестный аргумент ")
	ERROR_JSON                 = NewError("неправильный json ")
	ERROR_NOT_ENOUGH_ARGUMENTS = NewError("недостаточно аргументов ")
	ERROR_UNKOWN_STYLE         = NewError("неизвестный стиль")
	ERROR_WRITE_BIT_NOT_SET    = NewError("в этом файле не установлен бит разрешения на запись  ")
	ERROR_STAT                 = NewError("не удалось получить статистику")
	ERROR_NO_PERMISSION        = NewError("у вас нет прав на запись в этот каталог")
	ERROR_NOT_DIRECTORY        = NewError("путь не является каталогом ")
	ERROR_PATH_NOT_EXISTS      = NewError("путь не существует ")
	ERROR_NO_TOKEN             = NewError("нет токена")
	ERROR_TOO_MANY_ARGUMENTS   = NewError("слишком много аргументов")
	ERROR_NO_RESOURCES         = NewError("нет ресурсов")
	ERROR_RESOURCE_EXISTS      = NewError("Ресурс уже есть на вашем устройстве(диске),используйте флаг --overwrite true чтобы перезаписать ресурс")
)

func NewError(errorMessage string) error {
	return errors.New(setTextColor("Ошибка: "+errorMessage, RED))
}
