package main

import (
	"downgrid/beatmapurl"
	"downgrid/downloadmanager"
	"downgrid/registrymanager"
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func download(setid string) {
	counter := &downloadmanager.WriteCounter{}

	go downloadmanager.DownloadFile(setid, counter)
	ui := app.New()
	window := ui.NewWindow("Downgrid") // Window has to be borderless, need help on this
	progressbar := widget.NewProgressBar()

	go func() {
		for {
			if counter.Error != nil {
				cmd := exec.Command(registrymanager.GetBrowserPath(), os.Args[1])
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				break
			} else if counter.Done {
				cmd := exec.Command("explorer", counter.FilePath)
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				break
			} else if counter.Downloading {
				progressbar.Max = float64(counter.Total)
				progressbar.SetValue(float64(counter.Current))
			}
		}
		ui.Quit()
	}()

	window.CenterOnScreen()
	window.SetContent(container.NewVBox(progressbar))
	window.ShowAndRun()
}

func main() {
	registrymanager.Register()
	val := registrymanager.GetBrowserPath()

	if len(os.Args) == 1 {
		return
	}
	if !beatmapurl.IsOsuBeatmapLink(os.Args[1]) {
		cmd := exec.Command(val, os.Args[1])
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		id, exists := beatmapurl.GetSetId(os.Args[1])
		if !exists {
			cmd := exec.Command(val, os.Args[1])
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		download(id)
	}
}
