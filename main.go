package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ImageData struct {
	DHD string `json:"dhd"`
	DSD string `json:"dsd"`
}

type Data struct {
	Data map[string]ImageData `json:"data"`
}

func main() {
	url := "https://storage.googleapis.com/panels-api/data/20240916/media-1a-i-p~s"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Error: Status code", response.StatusCode)
		return
	}

	byteValue, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var data Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	os.MkdirAll("HD", os.ModePerm)
	os.MkdirAll("Normal", os.ModePerm)

	for _, img := range data.Data {
		if img.DHD != "" {
			downloadImage(img.DHD, "HD")
		}
		if img.DSD != "" {
			downloadImage(img.DSD, "Normal")
		}
	}
}

func downloadImage(url string, folder string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Error: Status code", response.StatusCode)
		return
	}

	contentType := response.Header.Get("Content-Type")
	var fileExtension string

	switch contentType {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	default:
		fileExtension = ".jpg" // Default to jpg if unknown type
	}

	fileName := filepath.Join(folder, fmt.Sprintf("%s%s", strings.TrimSuffix(filepath.Base(url), filepath.Ext(url)), fileExtension))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Downloaded:", fileName)
}
