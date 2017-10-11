package main

import (
	"fmt"
	"github.com/josedonizetti/amtest"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

var (
	app    = kingpin.New("amtest", "A command-line to create alerts on Alertmanager")
	amUrls = app.Flag("amurl", "Alertmanager URLs").Short('u').Default("http://127.0.0.1:9093").String()

	create        = app.Command("create", "Create an alert")
	amName        = create.Flag("name", "Alert name").Short('n').Required().String()
	amLabels      = create.Flag("labels", "Alert labels").Short('l').String()
	amAnnotations = create.Flag("annotations", "Alert annotations").Short('a').String()
	startTime     = create.Flag("starttime", "Start time").Short('s').Bool()
	endTime       = create.Flag("endtime", "End time").Short('e').Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case create.FullCommand():
		labels := amtest.LabelSet{
			"alertname": *amName,
		}

		if *amLabels != "" {
			for _, label := range strings.Split(*amLabels, ",") {
				arr := strings.Split(label, ":")
				labels[arr[0]] = arr[1]
			}
		}

		annotations := make(amtest.LabelSet)
		if *amAnnotations != "" {
			for _, annotation := range strings.Split(*amAnnotations, ",") {
				arr := strings.Split(annotation, ":")
				annotations[arr[0]] = arr[1]
			}
		}

		alert := amtest.Alert{
			Labels:      labels,
			Annotations: annotations,
		}

		if *startTime && *endTime {
			alert.StartsAt = time.Now()
			alert.EndsAt = time.Now()
		} else if !*startTime && *endTime {
			alert.EndsAt = time.Now()
		} else {
			alert.StartsAt = time.Now()
		}

		for _, url := range strings.Split(*amUrls, ",") {
			test := amtest.NewAmTest(url)
			err := test.Create(alert)
			fmt.Printf("%v\n", err)
		}
	}

	os.Exit(0)
}
