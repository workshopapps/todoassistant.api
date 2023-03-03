package timeSrv

import (
	"log"
	"time"
)

type TimeService interface {
	CurrentTime() time.Time
	TimeSince(time2 time.Time) time.Duration
	CheckFor339Format(time string) error
	CalcEndTime() time.Time
	ScheduleDate() time.Time
	TimeBefore(time1 time.Time) bool
	TimeAfter(time1 time.Time) bool
}

type timeStruct struct{}

func (t timeStruct) CheckFor339Format(timeStr string) error {
	_, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}
	// add check if it is greater than current time

	return nil
}

func (t timeStruct) CurrentTime() time.Time {
	return time.Now().Local()
}

func (t timeStruct) TimeSince(time2 time.Time) time.Duration {
	return time.Since(time2)
}

func (t timeStruct) CalcEndTime() time.Time {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	return endOfDay
}

func (t timeStruct) ScheduleDate() time.Time {
	now := t.CurrentTime()
	schdeduleDate := time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 00, 00, time.Local)
	return schdeduleDate
}

func (t timeStruct) TimeBefore(time1 time.Time) bool {
	return t.CurrentTime().After(time1)
}

func (t timeStruct) TimeAfter(time1 time.Time) bool {
	log.Println(t.ScheduleDate())
	return t.ScheduleDate().Before(time1)
}

func NewTimeStruct() TimeService {
	return &timeStruct{}
}
