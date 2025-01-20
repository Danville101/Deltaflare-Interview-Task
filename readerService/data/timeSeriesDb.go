package data

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	readApi "github.com/influxdata/influxdb-client-go/v2/api"
)

type Event struct {
	Criticality  int    `json:"criticality"`
	Timestamp    string `json:"timestamp"`
	EventMessage string `json:"eventMessage"`
}

type Reader interface {
	ReadDb(limit int, criticalityLevel int) ([]Event, error)
}

type InfluxDBClient struct {
	client   influxdb2.Client
	queryAPI readApi.QueryAPI
}

func NewInfluxDBClient(url, token, org, bucket string) *InfluxDBClient {
	client := influxdb2.NewClient(url, token)
	queryAPI := client.QueryAPI(org)

	return &InfluxDBClient{
		client:   client,
		queryAPI: queryAPI,
	}
}

func (db *InfluxDBClient) ReadDb(limit int, criticalityLevel int) ([]Event, error) {
	
	

	query := fmt.Sprintf(`
		from(bucket: "my-bucket")
		|> range(start: -1h)
		|> filter(fn: (r) => r["_measurement"] == "event")
		|> filter(fn: (r) => r["_field"] == "criticality" and r["_value"] >= %d or r["_field"] == "eventMessage")
		|> limit(n: %d)
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`, criticalityLevel, limit)

	result, err := db.queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}

	var events []Event
	for result.Next() {
		record := result.Record()

		criticality, ok := record.ValueByKey("criticality").(int64)
		if !ok {
			fmt.Println("Invalid type for criticality")
			continue
		}

		eventMessage, ok := record.ValueByKey("eventMessage").(string)
		if !ok {
			eventMessage = "No message"
		}

		event := Event{
			Criticality:  int(criticality),
			Timestamp:    record.Time().Format("2006-01-02T15:04:05Z"),
			EventMessage: eventMessage,
		}
		events = append(events, event)
	}

	if result.Err() != nil {
		fmt.Println("Query error:", result.Err())
		return nil, result.Err()
	}

	return events, nil
}
