package downloadmanager

import (
	"io"
	"mime"
	"net/http"
	"os"
	"regexp"
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

func DownloadFile(setid string, fn *WriteCounter) {
	cli := &http.Client{}
	req, err := http.NewRequest("GET", "https://chimu.moe/d/"+setid, nil)
	d, _ := os.UserHomeDir()

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

	r := regexp.MustCompile("[^a-zA-Z0-9\\.\\-]")

	file, err := os.Create(d + "\\" + r.ReplaceAllString(filename, "$1"))
	if err != nil {
		fn.Error = err
		return
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, fn))
	if err != nil {
		fn.Error = err
	}

	defer resp.Body.Close()
	fn.FilePath = d + "\\" + r.ReplaceAllString(filename, "$1")
	fn.Done = true
	fn.Response = resp
}
