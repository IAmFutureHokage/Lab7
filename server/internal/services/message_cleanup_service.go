package services

import (
	"time"
)

type MessageStorer interface {
	RemoveOldMessages(olderThan time.Time)
}

type MessageCleanupService struct {
	store  MessageStorer
	ticker *time.Ticker
	done   chan struct{}
}

func NewMessageCleanupService(store MessageStorer) *MessageCleanupService {
	return &MessageCleanupService{
		store:  store,
		ticker: time.NewTicker(1 * time.Hour),
		done:   make(chan struct{}),
	}
}

func (mcs *MessageCleanupService) Start() {
	go func() {
		for {
			select {
			case <-mcs.ticker.C:
				mcs.store.RemoveOldMessages(time.Now().Add(-24 * time.Hour))
			case <-mcs.done:
				return
			}
		}
	}()
}
func (mcs *MessageCleanupService) Stop() {
	close(mcs.done)
	mcs.ticker.Stop()
}
