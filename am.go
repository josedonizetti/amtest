package amtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type LabelSet map[string]string

type Alert struct {
	Labels       LabelSet  `json:"labels"`
	Annotations  LabelSet  `json:"annotations"`
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL"`
}

type alertApiResponse struct {
	Status    string   `json:"status"`
	Data      []*Alert `json:"data"`
	ErrorType string   `json:errorType,omitempty"`
	Error     string   `json:"error,omitempty"`
}

type AmTest struct {
	url string
}

func NewAmTest(url string) *AmTest {
	return &AmTest{url}
}

func (amt *AmTest) Create(alert Alert) error {
	err := sendPost(amt.url, alert)
	if err != nil {
		return err
	}

	fmt.Println("Alert created")
	return nil
}

func (amt *AmTest) GetAlert(name string) (Alert, error) {
	url := fmt.Sprintf("%s/api/v1/alerts", amt.url)
	resp, err := http.Get(url)

	var zeroValueAlert Alert
	if err != nil {
		return zeroValueAlert, err
	}

	var apiResponse alertApiResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return zeroValueAlert, err
	}

	for _, a := range apiResponse.Data {
		if a.Labels["alertname"] == name {
			return *a, nil
		}
	}

	return zeroValueAlert, errors.New("Not found")
}

func (amt *AmTest) Resolve(alert Alert) error {
	alert.EndsAt = time.Now()
	err := sendPost(amt.url, alert)
	if err != nil {
		return err
	}

	fmt.Println("Alert resolved")
	return nil
}

func sendPost(url string, alert Alert) error {
	alerts := [1]Alert{alert}
	var buf bytes.Buffer

	json.NewEncoder(&buf).Encode(alerts)
	url = fmt.Sprintf("%s/api/v1/alerts", url)
	res, err := http.Post(url, "application/json; charset=utf-8", &buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("Server returned %d", res.StatusCode)
		return nil
	}

	return nil
}
