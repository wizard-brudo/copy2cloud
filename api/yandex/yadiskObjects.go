package yandex

import "strings"

// Описание ресурса, мета-информация о файле или папке. Включается в ответ на запрос метаинформации.
type Resource struct {
	Publickey        string            `json:"public_key"`
	Embedded         ResourceList      `json:"_embedded"`
	Name             string            `json:"name"`
	Created          string            `json:"created"`
	CustomProperties map[string]string `json:"custom_properties"`
	PublicUrl        string            `json:"public_url"`
	OriginPath       string            `json:"origin_path"`
	Modified         string            `json:"modified"`
	Path             string            `json:"path"`
	Md5              string            `json:"md5"`
	Type             string            `json:"type"`
	MimeType         string            `json:"mime_type"`
	Size             int               `json:"size"`
}

// Проверяет папка ли это
func (res *Resource) IsDir() bool {
	return res.Type == "dir"
}

// Удаляет префикс disk или trash у пути
func (res *Resource) CorrectPath() string {
	if strings.Contains(res.Path, "disk") {
		return strings.Replace(res.Path, "disk:", "", 1)
	} else {
		return strings.Replace(res.Path, "trash:", "", 1)
	}
}

// Список ресурсов, содержащихся в папке. Содержит объекты Resource и свойства списка.
type ResourceList struct {
	Sort      string     `json:"sort"`
	PublicKey string     `json:"public_key"`
	Items     []Resource `json:"items"`
	Path      string     `json:"path"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	Total     int        `json:"total"`
}

// Плоский список всех файлов на Диске в алфавитном порядке.
type FilesResourceList struct {
	Items  []Resource `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// Список последних добавленных на Диск файлов, отсортированных по дате загрузки (от поздних к ранним).
type LastUploadedResourceList struct {
	Items []Resource `json:"items"`
	Limit int        `json:"limit"`
}

// Список опубликованных файлов на Диске.
type PublicResourcesList struct {
	Items  []Resource `json:"items"`
	Type   string     `json:"type"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// Данные о свободном и занятом пространстве на Диске
type Disk struct {
	User          map[string]string `json:"user"`
	TrashSize     float64           `json:"trash_size"`
	TotalSpace    float64           `json:"total_space"`
	UsedSpace     float64           `json:"used_space"`
	SystemFolders map[string]string `json:"system_folders"`
}

// Объект содержит URL для запроса метаданных ресурса.
type Link struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

// Ошибка при обработке запроса
type ErrorApi struct {
	Description string `json:"description"`
	Error       string `json:"error"`
	Message     string `json:"message"`
}
