package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type latlong struct {
	Lat, Long float64
}

func latLongByName(name string) (latlong, bool) {
	var (
		resp, err = http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(name)))
		jsonData  []map[string]any
	)
	if err != nil {
		println("error in get", err.Error())
		return latlong{}, false
	}
	json.NewDecoder(resp.Body).Decode(&jsonData)
	if len(jsonData) == 0 {
		println("did not find anything for", name)
		return latlong{}, false
	}
	var (
		lat, _  = strconv.ParseFloat(jsonData[0]["lat"].(string), 64)
		long, _ = strconv.ParseFloat(jsonData[0]["lon"].(string), 64)
	)
	return latlong{
		Lat:  lat,
		Long: long,
	}, true
}
