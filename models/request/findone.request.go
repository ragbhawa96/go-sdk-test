package models

type FindOneRequest struct {
	ID int `json:"id"`
}

type AccountOneRequest struct {
	ID   string `json:"id"`
	File string `json:"file`
}
