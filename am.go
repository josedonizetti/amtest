package amtest

import (
	"bytes"
	"encoding/json"
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

type AmTest struct {
	url string
}

func NewAmTest(url string) *AmTest {
	return &AmTest{url}
}

func (amt *AmTest) Create(alert Alert) error {
	alerts := [1]Alert{alert}
	var buf bytes.Buffer

	json.NewEncoder(&buf).Encode(alerts)
	url := fmt.Sprintf("%s/api/v1/alerts", amt.url)
	res, err := http.Post(url, "application/json; charset=utf-8", &buf)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}
