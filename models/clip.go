package models

import "time"

type Clip struct {
	Id        int
	Name      string
	Url       string
	Slug      string
	IsReady   bool
	CreatedBy string
	CreatedAt time.Time
}
