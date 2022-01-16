package oauth2

import (
	"copy2cloud/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	file, _               = os.Executable()
	Wd                    = filepath.Dir(file)
	clientID              = "b15c740e54a84c3ab4dd30ba087e96d0"
	encryptedClientSecret = "Th0QkghYdt6BWMqvTrzNpknFVv62WMvSRc4PIW-7xz5JgfCdV0cC-XKVgq3P9CNs"
)

func index(w http.ResponseWriter, r *http.Request) {
	indexTemplate, _ := template.ParseFiles(Wd + "/templates/index.html")
	url := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s", clientID)
	indexTemplate.ExecuteTemplate(w, "index.html", url)
}

func token(w http.ResponseWriter, r *http.Request) {
	tokenTemplate, _ := template.ParseFiles(Wd + "/templates/token.html")
	client_secret, err := utils.MustReveal(encryptedClientSecret)
	if err != nil {
		fmt.Println(utils.NewError(err.Error()))
		os.Exit(1)
	}
	urlValues := url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{r.URL.Query().Get("code")},
		"client_id":     []string{clientID},
		"client_secret": []string{client_secret},
	}
	resp, err := http.PostForm("https://oauth.yandex.ru/token", urlValues)
	if err != nil {
		fmt.Println(err)
	}
	bytes, _ := io.ReadAll(resp.Body)
	token := struct {
		AccessToken string `json:"access_token"`
	}{}
	json.Unmarshal(bytes, &token)
	fmt.Println(token.AccessToken)
	tokenTemplate.ExecuteTemplate(w, "token.html", token.AccessToken)
}

func GetToken() {
	// Запускаем параллельно чтобы веб сервер запустился и не ждал поиска браузера
	go func() {
		browsers := []string{"firefox", "brave-browser-stable", "start-tor-browser", "google-chrome-stable", "yandex-browser-beta"}
		for _, browser := range browsers {
			fmt.Println("Запуск " + browser)
			_, err := exec.Command(browser, "http://localhost:8080/").Output()
			if err == nil {
				break
			}
		}
	}()
	http.HandleFunc("/", index)
	http.HandleFunc("/token", token)
	fmt.Println("Пажалуйста зайдите на  localhost:8080 чтобы получить токен")
	fmt.Print(http.ListenAndServe(":8080", nil))
}
