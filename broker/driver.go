package broker

import (
	"errors"
	"time"

	"github.com/xuzhuoxi/infra-go/eventx"
)

var (
	EventOnTime = "TimeSlice.OnTime"
)

var (
	ErrDriverStarted = errors.New("Broker: Driver started. ")
	ErrDriverStopped = errors.New("Broker: Driver stopped. ")
)

func NewTimeSliceDriver(duration time.Duration) ITimeSliceDriver {
	return &timeSlice{duration: duration}
}

type ITimeSliceDriver interface {
	eventx.IEventDispatcher
	DriverStart() error
	DriverStop() error
}

type timeSlice struct {
	eventx.EventDispatcher
	duration time.Duration
	on       bool
}

func (s *timeSlice) DriverStart() error {
	if s.on {
		return ErrDriverStarted
	}
	go s.start()
	return nil
}

func (s *timeSlice) DriverStop() error {
	if !s.on {
		return ErrDriverStopped
	}
	s.stop()
	return nil
}

func (s *timeSlice) start() {
	s.on = true
	if s.duration <= 0 {
		for s.on {
			s.DispatchEvent(EventOnTime, s, nil)
		}
	} else {
		for s.on {
			s.DispatchEvent(EventOnTime, s, nil)
			time.Sleep(s.duration)
		}
	}
	s.on = false
}

func (s *timeSlice) stop() {
	s.on = false
}
