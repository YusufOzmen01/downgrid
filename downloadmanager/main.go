package downloadmanager

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

type WriteCounter struct {
	Current     uint64
	Total       uint64
	Done        bool
	Downloading bool
	Error       error
	FilePath    string
	Response    *http.Response
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.Downloading = true
	n := len(p)
	wc.Current += uint64(n)
	return n, nil
}

func DownloadFile(url string, path string, fn *WriteCounter, text string) {
	log.Println("File download requested with url " + url + " and path " + path)
	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fn.Error = err
	}

	resp, err := cli.Do(req)
	if err != nil {
		fn.Error = err
	}

	if resp.StatusCode != 200 {
		fn.Error = err
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		fn.Error = err
	}

	filename := params["filename"]

	if err != nil {
		fn.Error = err
	}

	fn.Total = uint64(resp.ContentLength)

	log.Println("File created")
	file, err := os.Create(path + filename)
	if err != nil {
		fn.Error = err
	}

	_, err = io.Copy(io.MultiWriter(progressbar.DefaultBytes(resp.ContentLength, text), file), io.TeeReader(resp.Body, fn))
	if err != nil {
		fn.Error = err
	}

	log.Println("File downloaded successfully")

	defer resp.Body.Close()
	fn.Done = true
	fn.FilePath = path + filename
	fn.Response = resp
}
