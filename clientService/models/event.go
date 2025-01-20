package models


type Event struct{
Criticality int `json:"criticality"`
Timestamp  string  `json:"timestamp"`
EventMessage string    `json:"eventMessage"`
}