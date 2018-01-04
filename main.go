package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/abderraoufbensalem/sumologic-client/sumo"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var appErrorsVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "app_error_vec",
	Help: "Number of errors for applications running on k8s",
},
	[]string{"app"},
)

func init() {
	sumo.InitConfigSettings()
	prometheus.MustRegister(appErrorsVec)
}

func prometheusHandler() http.Handler {
	return prometheus.Handler()
}

func main() {

	r := mux.NewRouter()
	r.Handle("/metrics", prometheusHandler())

	//prometheus.Register(histogram)

	s := &http.Server{
		Addr:           ":8001",
		ReadTimeout:    8 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	go func() {

		session := sumo.DefaultSession()
		session.Discover()
		client := sumo.NewClient(session)

		/* Create a New Job */
		createJob := &sumo.JobCreate{}
		createJob.From = time.Date(2017, 11, 21, 20, 0, 0, 0, time.UTC)
		createJob.To = time.Date(2017, 11, 22, 20, 0, 0, 0, time.UTC)
		createJob.TimeZone = "PST"
		createJob.Query = "_collector=prod ERROR | count app"
		job, err := client.Jobs().Create(createJob)

		if err != nil {
			fmt.Println("client.Jobs().Create() ==> ERROR: ", err)
		}

		/* Retrieve the Job*/
		// time.Sleep(5 * time.Second)
		// jobs, err := client.Jobs().Get(job.ID)
		// rawJobs, _ := json.Marshal(jobs)
		// fmt.Println(string(rawJobs), err)

		/*Get Job Records*/
		jobRecords, err := client.Jobs().GetRecords(job.ID, 0, 100)

		if err != nil {
			fmt.Println("client.Jobs().GetRecords() ==> ERROR: ", err)
		}

		/*Add records to Prometheus*/
		if jobRecords != nil {
			for _, record := range jobRecords.Records {
				rawRecord, _ := json.Marshal(record)
				fmt.Println(string(rawRecord), err)

				app := record.Map.App

				if app != "" {
					count, _ := strconv.ParseFloat(record.Map.Count, 64)
					appErrorsVec.WithLabelValues(app).Set(count)
				}
			}
		} else {
			fmt.Println("No Records found")
		}

		/*Delete The Job*/
		err = client.Jobs().Delete(job.ID)
		if err != nil {
			fmt.Println("client.Jobs().Delete() ==> ERROR: ", err)
		}
	}()

	// collectors, err := client.Collectors().List(0, 5)
	// raw, _ := json.Marshal(collectors)
	// fmt.Println(string(raw), err)

	// sources, err := client.Collectors().Sources(collectors[0].ID).List()
	// raw, _ = json.Marshal(sources)
	// fmt.Println(string(raw), err)

	/*
		collector, err := client.Collectors().Create(&api.CollectorCreate{
			CollectorType: "Hosted",
			Name:          "abc-collector",
			Description:   "",
			Category:      "",
		})
		fmt.Printf("%+v\n%s\n", collector, err)
	*/

	log.Fatal(s.ListenAndServe())
}
