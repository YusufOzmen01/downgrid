package main

import (
	"downgrid/beatmapurl"
	"downgrid/downloadmanager"
	"downgrid/registrymanager"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func download(setid string, provider int) string {
	counter := &downloadmanager.WriteCounter{}
	text := ""

	url := ""
	switch provider {
	case 1:
		url = "https://chimu.moe/d/"
		text = "Downloading beatmap from chimu.moe"
	case 2:
		url = "https://beatconnect.io/b/"
		text = "Downloading beatmap from beatconnect.io"
	}

	go downloadmanager.DownloadFile(url+setid, ".\\", counter, text)

	for {
		if counter.Error != nil {
			return ""
		}
		if counter.Done {
			fmt.Println(counter.FilePath)
			return counter.FilePath
		}
	}
}

func main() {
	file, err := os.OpenFile("logs-"+strconv.Itoa(time.Now().Nanosecond())+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("Logger initialized")
	registrymanager.Register()
	val := registrymanager.GetBrowserPath()

	log.Println("Initializing beatmap url checks")
	if len(os.Args) == 1 {
		log.Println("There is nothing to do")
		return
	}
	if !beatmapurl.IsOsuBeatmapLink(os.Args[1]) {
		cmd := exec.Command(val, os.Args[1])
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Trying to download the beatmap from chimu.moe")
		id, exists := beatmapurl.GetSetId(os.Args[1])
		if !exists {
			log.Println("Running the web browser")
			cmd := exec.Command(val, os.Args[1])
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		path := download(id, 1)

		if path == "" {
			log.Println("Trying to download the beatmap from beatconnect.io")
			path = download(id, 2)
			if path == "" {
				log.Println("Running the web browser")
				cmd := exec.Command(val, os.Args[1])
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}

		log.Println("Running osu!")
		cmd := exec.Command("explorer", path)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}
