//+build !test

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type yahooSunsetInfo struct {
	Query struct {
		Count   int
		Created time.Time
		Lang    string
		Results struct {
			Channel struct {
				Astronomy struct {
					Sunset string `json:"sunset"`
				} `json:"astronomy"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}

type yahooSunsetFinder struct {
}

func parse(jsonResponse []byte) (sunsetResult, error) {

	ysi := &yahooSunsetInfo{}

	err := json.Unmarshal(jsonResponse, ysi)

	if err != nil {
		return sunsetResult{}, err
	}

	sr := sunsetResult{Sunset: ysi.Query.Results.Channel.Astronomy.Sunset, Timestamp: ysi.Query.Created}

	return sr, nil
}

func (ysf *yahooSunsetFinder) Query(location string) (sunsetResult, error) {

	yahooQuery := fmt.Sprintf("select astronomy.sunset from weather.forecast where woeid in (select woeid from geo.places(1) where text=\"%s\")", location)

	uri := fmt.Sprintf("https://query.yahooapis.com/v1/public/yql?q=%s&format=json", url.QueryEscape(yahooQuery))

	response, err := http.Get(uri)

	defer response.Body.Close()

	if err != nil {
		return sunsetResult{}, err
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return sunsetResult{}, err
	}

	ysi := &yahooSunsetInfo{}

	err = json.Unmarshal(body, ysi)

	if err != nil {
		return sunsetResult{}, err
	}

	sr := sunsetResult{Sunset: ysi.Query.Results.Channel.Astronomy.Sunset, Timestamp: ysi.Query.Created}

	if sr.Sunset == "" {
		return sr, errNotFound
	}

	return sr, nil
}
