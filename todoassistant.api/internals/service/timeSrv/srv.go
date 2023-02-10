package timeSrv

import "time"

type TimeService interface {
	CurrentTime() time.Time
	TimeSince(time2 time.Time) time.Duration
	CheckFor339Format(time string) error
	CalcEndTime() time.Time
	TimeBefore(time1 time.Time) bool
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

func NewTimeStruct() TimeService {
	return &timeStruct{}
}

func (t timeStruct) CurrentTime() time.Time {
	return time.Now()
}

func (t timeStruct) TimeSince(time2 time.Time) time.Duration {
	return time.Since(time2)
}

func (t timeStruct) CalcEndTime() time.Time {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	return endOfDay
}

func (t timeStruct) TimeBefore(time1 time.Time) bool {
	return t.CurrentTime().Local().Before(time1)
}
