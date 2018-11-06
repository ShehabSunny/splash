package resplash

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// Unsplash type
type Unsplash struct {
	api      string
	clientID string
}

// Image response type
type Image struct {
	Description string   `json:"description"`
	URLS        ImageURL `json:"urls"`
}

// ImageURL response type
type ImageURL struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
}

// SearchResult response type
type SearchResult struct {
	Results []Image `json:"results"`
}

// NewUnsplash Unsplash constructor
func NewUnsplash(api string, clientID string) Unsplash {
	return Unsplash{api, clientID}
}

// Download image
func (u Unsplash) Download(url string) (string, error) {
	// get image from url
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// get user home dir
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// store image in $home/wallpapers dir
	t := time.Now().Format("20060102150405")
	dataDir := filepath.Join(usr.HomeDir, "wallpapers")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, os.ModePerm)
	}
	filePath := filepath.Join(dataDir, fmt.Sprintf("%s.jpg", t))

	//open a file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	file.Close()

	return filePath, nil
}

// DownloadRandom image
func (u Unsplash) DownloadRandom() (string, error) {
	resp, err := http.Get(u.api + "/photos/random/?client_id=" + u.clientID)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	// process response
	var resImage Image
	json.Unmarshal(body, &resImage)

	// get image url from json response
	imageURL := resImage.URLS.Full
	if imageURL == "" {
		return "", errors.New("could not get full image url")
	}

	// download from url
	path, err := u.Download(imageURL)
	if err != nil {
		return "", err
	}

	return path, nil
}

// DownloadRandomTopic image
func (u Unsplash) DownloadRandomTopic(topic string) (string, error) {
	// get by query
	resp, err := http.Get(u.api + "/search/photos?client_id=" + u.clientID + "&query=" + topic + "&per_page=100")
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)

	// process response
	var searchRes SearchResult
	json.Unmarshal(body, &searchRes)
	resLength := len(searchRes.Results)
	if resLength < 1 {
		return "", errors.New("no image found")
	}

	// get a random image
	rand.Seed(time.Now().UnixNano())
	randomIdx := random(0, resLength)

	// get image url from json response
	imageURL := searchRes.Results[randomIdx].URLS.Full
	if imageURL == "" {
		return "", errors.New("could not get full image url")
	}

	// download from url
	path, err := u.Download(imageURL)
	if err != nil {
		return "", err
	}

	return path, nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
