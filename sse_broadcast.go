package main

import (
	"net/http"
	"time"

	"github.com/lonord/sse"
)

const readInterval = time.Second * 3

type SSESystemBoradcast struct {
	sse *sse.Service
}

func NewSSESystemBoradcast() *SSESystemBoradcast {
	return &SSESystemBoradcast{
		sse: sse.NewServiceWithOption(sse.Option{
			Headers: map[string]string{"X-Accel-Buffering": "no"},
		}),
	}
}

func (s *SSESystemBoradcast) handleClient(clientID interface{}, w http.ResponseWriter) error {
	c, err := s.sse.HandleClient(clientID, w)
	if err != nil {
		return err
	}
	go s.runSystemInfoReader()
	<-c
	return nil
}

func (s *SSESystemBoradcast) runSystemInfoReader() {
	timer := time.NewTimer(readInterval)
	for {
		select {
		case <-timer.C:
			if s.sse.GetClientCount() == 0 {
				return
			}
			result := readSystemInfoAll()
			s.sse.Broadcast(sse.Event{
				Data: result,
			})
			timer.Reset(readInterval)
		}
	}
}
