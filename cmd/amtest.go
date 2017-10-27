package main

import (
	"context"
	"fmt"
	"github.com/josedonizetti/amtest"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	defaultAmURL    = "http://localhost:9093"
	defaultInterval = "1m"
)

var (
	app    = kingpin.New("amtest", "A command-line to create alerts on Alertmanager")
	amUrls = app.Flag("amurl", "Alertmanager URLs").Short('u').Default(defaultAmURL).Strings()

	create         = app.Command("create", "Create an alert")
	amName         = create.Flag("name", "Alert name").Short('n').Required().String()
	amLabels       = create.Flag("labels", "Alert labels").Short('l').StringMap()
	amAnnotations  = create.Flag("annotations", "Alert annotations").Short('a').StringMap()
	generatorUrl   = create.Flag("generatorUrl", "Generator URL").Short('g').String()
	startTime      = create.Flag("starttime", "Start time").Short('s').Bool()
	endTime        = create.Flag("endtime", "End time").Short('e').Bool()
	repeat         = create.Flag("repeat", "Send the alert repeatedly based on repeat-interval").Short('r').Bool()
	repeatInterval = create.Flag("repeat-interval", "Interval to repeat alerts").Default(defaultInterval).Short('i').Duration()

	resolve        = app.Command("resolve", "Create an alert")
	amResolvedName = resolve.Flag("name", "Alert name").Short('n').Required().String()
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-term
		fmt.Println("Stopping...")
		cancel()
	}()

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

		var wg sync.WaitGroup
		wg.Add(len(*amUrls))

		for _, url := range *amUrls {
			go func(ctx context.Context, u string) {
				test := amtest.NewAmTest(u)
				err := test.Create(alert)
				if err != nil {
					fmt.Printf("%v\n", err)
				}

				for *repeat {
					select {
					case <-ctx.Done():
						*repeat = false
					case <-time.After(*repeatInterval):
						err := test.Create(alert)
						if err != nil {
							fmt.Printf("%v\n", err)
						}
					}
				}

				wg.Done()
			}(ctx, url)
		}

		wg.Wait()
	case resolve.FullCommand():
		for _, url := range *amUrls {
			test := amtest.NewAmTest(url)
			alert, err := test.GetAlert(*amResolvedName)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}

			test.Resolve(alert)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}

	os.Exit(0)
}
