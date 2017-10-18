package main

import (
	"fmt"
	"github.com/josedonizetti/amtest"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app    = kingpin.New("amtest", "A command-line to create alerts on Alertmanager")
	amUrls = app.Flag("amurl", "Alertmanager URLs").Short('u').Default("http://127.0.0.1:9093").Strings()

	create        = app.Command("create", "Create an alert")
	amName        = create.Flag("name", "Alert name").Short('n').Required().String()
	amLabels      = create.Flag("labels", "Alert labels").Short('l').StringMap()
	amAnnotations = create.Flag("annotations", "Alert annotations").Short('a').StringMap()
	generatorUrl  = create.Flag("generatorUrl", "Generator URL").Short('g').String()
	startTime     = create.Flag("starttime", "Start time").Short('s').Bool()
	endTime       = create.Flag("endtime", "End time").Short('e').Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case create.FullCommand():
		labels := amtest.LabelSet{
			"alertname": *amName,
		}

		for k, v := range *amLabels {
			labels[k] = v
		}

		annotations := make(amtest.LabelSet)
		for k, v := range *amAnnotations {
			annotations[k] = v
		}

		alert := amtest.Alert{
			Labels:       labels,
			Annotations:  annotations,
			GeneratorURL: *generatorUrl,
		}

		if *endTime {
			alert.StartsAt = time.Now()
			alert.EndsAt = time.Now()
		} else {
			alert.StartsAt = time.Now()
		}

		for _, url := range *amUrls {
			test := amtest.NewAmTest(url)
			err := test.Create(alert)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}

	os.Exit(0)
}
