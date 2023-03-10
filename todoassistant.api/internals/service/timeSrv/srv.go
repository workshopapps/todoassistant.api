package timeSrv

import (
	"time"
)

type TimeService interface {
	CurrentTime() time.Time
	CurrentTimeString() string
	TimeSince(time2 time.Time) time.Duration
	CheckFor339Format(timeStr string) (string, error)
	CalcEndTime() time.Time
	CalcEndTimeString() string
	Parse(time2 string) (time.Time, error)
	CalcScheduleEndTime(schedule time.Time) time.Time
	CalcScheduleEndTimeString(schedule time.Time) string
	ScheduleDate() time.Time
	TimeBefore(time1 time.Time) bool
	TimeAfter(time1 time.Time) bool
	ScheduleTimeAfter(time1 time.Time) bool
}

type timeStruct struct{}

func (t timeStruct) CheckFor339Format(timeStr string) (string, error) {
	ti, err := t.Parse(timeStr)
	if err != nil {
		return "", err
	}
	// add check if it is greater than current time

	return ti.Format(time.RFC3339), nil
}

func (t timeStruct) CurrentTime() time.Time {
	return time.Now().UTC()
}

func (t timeStruct) CurrentTimeString() string {
	return t.CurrentTime().Format(time.RFC3339)
}

func (t timeStruct) TimeSince(time2 time.Time) time.Duration {
	return time.Since(time2)
}

func (t timeStruct) Parse(time2 string) (time.Time, error) {
	schedule, err := time.ParseInLocation(time.RFC3339, time2, time.UTC)
	if err != nil {
		return time.Time{}, err
	}

	return schedule, nil
}

func (t timeStruct) CalcEndTime() time.Time {
	now := t.CurrentTime()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	return endOfDay
}

func (t timeStruct) CalcEndTimeString() string {
	return t.CalcEndTime().Format(time.RFC3339)
}

func (t timeStruct) CalcScheduleEndTime(schedule time.Time) time.Time {
	endOfDay := time.Date(schedule.Year(), schedule.Month(), schedule.Day(), 23, 59, 59, 0, time.UTC)
	return endOfDay
}

func (t timeStruct) CalcScheduleEndTimeString(schedule time.Time) string {
	return t.CalcScheduleEndTime(schedule).Format(time.RFC3339)
}

func (t timeStruct) TimeBefore(time1 time.Time) bool {
	return t.CurrentTime().After(time1)
}

func (t timeStruct) TimeAfter(time1 time.Time) bool {
	return t.CurrentTime().Before(time1)
}

func (t timeStruct) ScheduleDate() time.Time {
	now := t.CurrentTime()
	schdeduleDate := time.Date(now.Year(), now.Month(), now.Day(), 00, 00, 00, 00, time.UTC)
	return schdeduleDate
}

func (t timeStruct) ScheduleTimeAfter(time1 time.Time) bool {
	return t.ScheduleDate().Before(time1)
}

func NewTimeStruct() TimeService {
	return &timeStruct{}
}
