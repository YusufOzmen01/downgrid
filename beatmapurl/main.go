package beatmapurl

import (
	"log"
	"net/http"
	"strings"
	"unicode"
)

func IsOsuBeatmapLink(url string) bool {
	log.Println("Checking url " + url)
	if strings.Contains(url, "https://osu.ppy.sh/") || strings.Contains(url, "http://osu.ppy.sh/") {
		log.Println("Url is an osu! url")
		if strings.Contains(url, "/s/") || strings.Contains(url, "/beatmapse") || strings.Contains(url, "/b/") || strings.Contains(url, "/beatmaps/") {
			log.Println("Url is an osu! beatmap url")
			return true
		} else {
			log.Println("Url is not an osu! beatmap url")
			return false
		}
	} else {
		log.Println("Url is not an osu! url")
		return false
	}
}
func IsSet(url string) bool {
	log.Println("Checking url " + url)
	if IsOsuBeatmapLink(url) && (strings.Contains(url, "s/") || strings.Contains(url, "beatmapse")) && (!strings.Contains(url, "b/") && !strings.Contains(url, "beatmaps/")) {
		log.Println("Url is an osu! beatmapset url")
		return true
	} else {
		log.Println("Url is not an osu! beatmapset url")
		return false
	}
}

func GetId(url string) string {
	log.Println("Parsing url " + url)
	_id := ""
	if IsOsuBeatmapLink(url) {
		found := false
		for _, i := range url {
			if unicode.IsNumber(i) {
				_id += string(i)
				found = true
			} else if found {
				break
			}
		}
	}
	log.Println("Parsed id is " + _id)
	return _id
}

func GetSetId(url string) (string, bool) {
	log.Println("Parsing url " + url)
	resp, err := http.Head(url)
	if err != nil {
		return "", false
	}
	log.Println("Parsed set id is " + GetId(resp.Request.URL.String()))
	return GetId(resp.Request.URL.String()), true
}
