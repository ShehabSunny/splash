package main

import (
	"flag"
	"fmt"
	"gosplash/resplash"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/reujab/wallpaper"
)

// ImageVendor interface
type ImageVendor interface {
	Download(URL string) (string, error)
	DownloadRandom() (string, error)
	DownloadRandomTopic(topic string) (string, error)
}

func main() {
	// get user home dir
	usr, err := user.Current()
	if err != nil {
		log.Println("could not get current user")
		return
	}

	gosplashHome := filepath.Join(usr.HomeDir, ".gosplash", ".env")
	err = godotenv.Load(gosplashHome)
	if err != nil {
		fmt.Println("error loading .env file in $HOME/.gosplash directory")
		fmt.Printf("1. Signup for an unsplash developer account. \n")
		fmt.Printf("2. Create an app, get credentials and put ACCESS_KEY=<acckess_key> in a .env file in $HOME/.gosplash directory. \n\n")
		return
	}

	api := "https://api.unsplash.com"
	clientID := os.Getenv("ACCESS_KEY")
	if clientID == "" {
		fmt.Println("ACCESS_KEY not found in .env file in $HOME/.gosplash directory")
		fmt.Printf("1. Signup for an unsplash developer account. \n")
		fmt.Printf("2. Create an app, get credentials and put ACCESS_KEY=<acckess_key> in a .env file in $HOME/.gosplash directory. \n\n")

		return
	}

	flag.Usage = usage
	// parse topic from flags
	topic := flag.String("t", "", "an image with this topic will be set as wallpaper.")
	flag.Parse()

	// initial resplash
	u := resplash.NewUnsplash(api, clientID)

	fmt.Println("starting..")

	if *topic != "" {
		fmt.Printf("topic: %s\n", *topic)
		err = downloadSetWallpaperByTopic(u, *topic)
	} else {
		err = downloadSetWallpaper(u)
	}
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("wallpaper changed successfully.")
}

func usage() {
	fmt.Printf("gosplash: a cli tool to change wallpaper from unsplash api.\n\n")
	fmt.Printf("Usage: no arguments will set a random wallpaper. \n[OPTIONS] argument:\n")
	flag.PrintDefaults()
}

func downloadSetWallpaper(imageVendor ImageVendor) error {
	filePath, err := imageVendor.DownloadRandom()
	if err != nil {
		return err
	}

	err = setWallpaper(filePath)
	if err != nil {
		return err
	}

	return nil
}

func downloadSetWallpaperByTopic(imageVendor ImageVendor, topic string) error {
	filePath, err := imageVendor.DownloadRandomTopic(topic)
	if err != nil {
		return err
	}

	err = setWallpaper(filePath)
	if err != nil {
		return err
	}

	return nil
}

// setWallpaper sets the wallpaper
func setWallpaper(path string) error {
	err := wallpaper.SetFromFile(path)
	if err != nil {
		return err
	}
	return nil
}
