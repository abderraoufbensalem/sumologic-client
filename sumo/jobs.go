package sumo

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Jobs struct {
	executor *ClientExecutor
}

type Job struct {
	State            string            `json:"state"`
	HistogramBuckets []HistogramBucket `json:"histogramBuckets,omitempty"`
	MessageCount     int               `json:"messageCount"`
	RecordCount      int               `json:"recordCount"`
}

type HistogramBucket struct {
	StartTimestamp int `json:"startTimestamp"`
	Length         int `json:"length"`
	Count          int `json:"count"`
}

type JobCreate struct {
	Query    string    `json:"query"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
	TimeZone string    `json:"timeZone"`
}

type JobResponse struct {
	ID   string  `json:"id"`
	Link JobLink `json:"link"`
}

type JobLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type JobRecords struct {
	Fields  []JobField     `json:"fields"`
	Records []JobRecordMap `json:"records"`
}

type JobField struct {
	Name      string `json:"name"`
	FieldType string `json:"fieldType"`
	KeyField  bool   `json:"keyField"`
}

type JobRecordMap struct {
	Map JobMap `json:"map"`
}

type JobMap struct {
	App   string `json:"app"`
	Count string `json:"_count"`
}

func NewJobs(executor *ClientExecutor) *Jobs {
	return &Jobs{
		executor: executor,
	}
}

func (j *Jobs) Create(job *JobCreate) (*JobResponse, error) {
	req, err := j.executor.NewRequest()
	if err != nil {
		return nil, err
	}

	req.SetEndpoint("search/jobs")

	req.SetJSONBody(job)

	res, err := req.Post()

	if err != nil {
		return nil, err
	}

	item := &JobResponse{}
	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (j *Jobs) Get(id string) (*Job, error) {
	req, err := j.executor.NewRequest()
	if err != nil {
		return nil, err
	}
	req.SetEndpoint(fmt.Sprintf("/search/jobs/%v", id))

	res, err := req.Get()

	if err != nil {
		return nil, err
	}

	type getResponse struct {
		Job *Job `json:"job"`
	}

	item := &getResponse{}

	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}

	return item.Job, nil
}

func (j *Jobs) GetRecords(id string, offset int, limit int) (*JobRecords, error) {
	req, err := j.executor.NewRequest()
	if err != nil {
		return nil, err
	}

	req.SetEndpoint(fmt.Sprintf("/search/jobs/%v/records", id))
	req.SetQuery(url.Values{
		"offset": []string{strconv.Itoa(offset)},
		"limit":  []string{strconv.Itoa(limit)},
	})

	time.Sleep(5 * time.Second)
	res, err := req.Get()

	found, err := IsObjectFound(res, err)
	retries := 0
	maxReties := 5
	if !found {
		for {
			time.Sleep(5 * time.Second)
			res, err := req.Get()
			found, err := IsObjectFound(res, err)
			if !found {
				retries++
			} else {
				break
			}
			if retries > maxReties {
				break
			}
		}
	}

	if err != nil {
		fmt.Println("req.Get() ==> ERROR: ", err)
		return nil, err
	}

	item := &JobRecords{}

	if err := res.BodyJSON(item); err != nil {
		return nil, err
	}

	if item != nil {
		fmt.Println(" Item is Not NULL")
	}

	return item, nil

}

func (j *Jobs) Delete(ID string) error {
	req, err := j.executor.NewRequest()
	if err != nil {
		return err
	}
	req.SetEndpoint(fmt.Sprintf("/seach/jobs/%s", ID))

	type deleteResponse struct {
		ID string `json:"id"`
	}
	res, err := req.Delete()
	if err != nil {
		fmt.Println("Delete id: ", ID)
		return err
	}

	item := &deleteResponse{}
	if err := res.BodyJSON(item); err != nil {
		return err
	}
	if item.ID != ID {
		return errors.New("Could not delete the Job")
	}

	return nil
}
