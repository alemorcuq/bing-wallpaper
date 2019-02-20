package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/antchfx/htmlquery"
)

func main() {
	bing := "https://www.bing.com"
	homeDir := os.Getenv("HOME")
	wallpaperDir := fmt.Sprintf("%s/.wallpapers", homeDir)

	doc, err := htmlquery.LoadURL(bing)
	if err != nil {
		panic(err)
	}

	n := htmlquery.FindOne(doc, "//link[@id='bgLink']/@href")
	url := fmt.Sprintf("%s%s", bing, htmlquery.InnerText(n))

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	fileName := filepath.Base(url)
	exists, err := fileExists(wallpaperDir)
	if err != nil {
		panic(err)
	}

	if !exists {
		if err := os.Mkdir(wallpaperDir, 0755); err != nil {
			panic(err)
		}
	}

	filePath := filepath.Join(wallpaperDir, fileName)

	img, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(img, resp.Body)
	if err != nil {
		panic(err)
	}

	err = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", filePath).Run()
	if err != nil {
		panic(err)
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
