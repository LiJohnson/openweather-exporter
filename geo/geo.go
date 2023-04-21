// Copyright 2023 Billy Wooten
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type Geo struct {
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Name       string            `json:"name"`
	State      string            `json:"state"`
	Country    string            `json:"country"`
	LocalNames map[string]string `json:"local_names"`
}

func GetCoords(apiKey string, city string) (float64, float64, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	var rq = url.Values{}
	rq.Add("limit", "1")
	rq.Add("appid", apiKey)
	rq.Add("q", city)
	geoUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?%s", rq.Encode())
	log.Infof("geo url %s", geoUrl)
	response, err := client.Get(geoUrl)
	if err != nil {
		return 0, 0, err
	}
	defer response.Body.Close()
	var data []Geo
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0, 0, err
	}
	if len(data) == 0 {
		return 0, 0, errors.New(fmt.Sprintf("city[%s] not found ", city))
	}

	var location = data[0]
	log.Debugf("Latitude: %f Longitude: %f for %s found : Name : %s , Country : %s , State : %s  LocalName : %s", location.Lat, location.Lon, city, location.Name, location.Country, location.State, location.LocalNames)

	return location.Lat, location.Lon, nil
}
