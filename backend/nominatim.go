package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type latlong struct {
	Lat, Long float64
}

func latLongByName(name string) (latlong, bool) {
	var (
		resp, err = http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", name))
		jsonData  []map[string]any
	)
	if err != nil {
		println("error in get", err.Error())
		return latlong{}, false
	}
	bodyInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println("error reading response", err.Error())
		return latlong{}, false
	}
	if err := json.Unmarshal(bodyInBytes, &jsonData); err != nil {
		println(err.Error())
		return latlong{}, false
	}
	if len(jsonData) == 0 {
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
