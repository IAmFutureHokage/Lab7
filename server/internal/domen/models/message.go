package models

import "time"

type Message struct {
	Date     time.Time
	Username string
	Text     string
}
