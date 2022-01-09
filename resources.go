package main

import "strings"

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
func (res *Resource) isDir() bool {
	return res.Type == "dir"
}

// Удаляет префикс disk у пути
func (res *Resource) CorrectPath() string {
	return strings.Replace(res.Path, "disk:/", "", 1)
}

type ResourceList struct {
	Sort      string     `json:"sort"`
	PublicKey string     `json:"public_key"`
	Items     []Resource `json:"items"`
	Path      string     `json:"path"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	Total     int        `json:"total"`
}

type FilesResourceList struct {
	Items []struct {
		Name     string `json:"name"`
		Created  string `json:"created"`
		Modified string `json:"modified"`
		Path     string `json:"path"`
		Type     string `json:"type"`
		MimeType string `json:"mime_type"`
		Size     int    `json:"size"`
	} `json:"items"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LastUploadedResourceList struct {
	Items []Resource `json:"items"`
	Limit int        `json:"limit"`
}

type PublicResourcesList struct {
	Items  []Resource `json:"items"`
	Type   string     `json:"type"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type Disk struct {
	User          map[string]string `json:"user"`
	TrashSize     float64           `json:"trash_size"`
	TotalSpace    float64           `json:"total_space"`
	UsedSpace     float64           `json:"used_space"`
	SystemFolders map[string]string `json:"system_folders"`
}

type ErrorApi struct {
	Description string `json:"description"`
	Error       string `json:"error"`
	Message     string `json:"message"`
}

type Link struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}
