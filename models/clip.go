package models

import "time"

type Clip struct {
	Id        int
	Name      string
	Url       string
	Slug      string
	CreatedBy string
	CreatedAt time.Time
}
