package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/go-shiori/obelisk"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// downloadFile - функция для скачивания файла
func downloadFile(filepath string, url string) error {

	// отправка get запроса по ссылке
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// закрытие чтения результата запроса в defer
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(resp.Body)

	// создание файла для записи результата
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	// закрытие файла в defer
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(out)

	// копирование результата запроса в файл
	_, err = io.Copy(out, resp.Body)
	return err
}

// downloadWebpage - функция для скачивания веб страницы со всеми зависимостями
func downloadWebpage(filename, url string) error {
	// скачивание полной веб страницы с зависимостями
	req := obelisk.Request{
		URL: url,
	}

	// объявление архиватора
	arc := obelisk.Archiver{EnableLog: false}
	arc.Validate()

	// архивирование результата запроса
	result, _, err := arc.Archive(context.Background(), req)
	if err != nil {
		return err
	}

	// создание файла архива
	f, err := os.Create(filename + ".html.gz")
	if err != nil {
		return err
	}

	// закрытие файла архива в defer
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)

	// запись результата запроса в архив
	gz := gzip.NewWriter(f)
	_, err = gz.Write(result)
	if err != nil {
		return err
	}

	// закрытие архиватора
	err = gz.Close()
	if err != nil {
		return err
	}

	return nil
}

func main() {

	fileURL := os.Args[1] // получение ссылки на файл или веб страницу

	splitURL := strings.Split(fileURL, "/") // разделение ссылку на элементы пути

	// проверка на флаг --mirror
	if len(os.Args) == 3 {
		if os.Args[2] == "--mirror" {
			// если флаг присутствует, скачивается веб страница
			filename := splitURL[len(splitURL)-1]
			if len(filename) == 0 {
				if len(splitURL) >= 2 {
					filename = splitURL[len(splitURL)-2]
				}
			}
			err := downloadWebpage(splitURL[len(splitURL)-1], fileURL) // скачивание веб страницы
			if err != nil {
				fmt.Println("error while downloading file: ", err)
				return
			}
		} else {
			log.Fatalln("invalid arguments")
		}
	} else {
		// если флаг отсутствует, скачивается файл
		err := downloadFile(splitURL[len(splitURL)-1], fileURL) // скачивание файла
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("downloaded: " + fileURL)
}
