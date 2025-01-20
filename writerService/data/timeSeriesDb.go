package data

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	writeApi "github.com/influxdata/influxdb-client-go/v2/api"
)

type Event struct{
	Criticality int `json:"criticality"`
	Timestamp  string  `json:"timestamp"`
	EventMessage string    `json:"eventMessage"`
	}


type Writer interface{
	
	WriteEvent(event Event) error
}


type InfluxDBClient struct{
	client influxdb2.Client
	writeAPI writeApi.WriteAPI
}


func NewInfluxDBClient(url, token, org, bucket string) *InfluxDBClient{
	client := influxdb2.NewClient(url, token)
	writeAPI:=  client.WriteAPI(org, bucket)

	return &InfluxDBClient{
		client: client,
		writeAPI: writeAPI,
	}
}


func (db *InfluxDBClient) WriteEvent(event Event) error{
	parsedTime, err := time.Parse(time.RFC3339, event.Timestamp)
	if err != nil {
	    fmt.Println("Error parsing time:", err)
	    return err
	}
 
	point := influxdb2.NewPoint(
		"event",
		nil, 
    map[string]interface{}{
        "criticality": event.Criticality, 
        "eventMessage":     event.EventMessage, 
    }, 
		parsedTime,
	)

	db.writeAPI.WritePoint(point)
	db.writeAPI.Flush()

	return nil
}