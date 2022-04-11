package beatmapurl

import (
	"net/http"
	"strings"
	"unicode"
)

func IsOsuBeatmapLink(url string) bool {
	if strings.Contains(url, "https://osu.ppy.sh/") || strings.Contains(url, "http://osu.ppy.sh/") {
		if strings.Contains(url, "/s/") || strings.Contains(url, "/beatmapse") || strings.Contains(url, "/b/") || strings.Contains(url, "/beatmaps/") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
func IsSet(url string) bool {
	if IsOsuBeatmapLink(url) && (strings.Contains(url, "s/") || strings.Contains(url, "beatmapse")) && (!strings.Contains(url, "b/") && !strings.Contains(url, "beatmaps/")) {
		return true
	} else {
		return false
	}
}

func GetId(url string) string {
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
	return _id
}

func GetSetId(url string) (string, bool) {
	resp, err := http.Head(url)
	if err != nil {
		return "", false
	}
	return GetId(resp.Request.URL.String()), true
}
