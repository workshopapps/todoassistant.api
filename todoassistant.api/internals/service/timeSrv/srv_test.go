package timeSrv

import (
	"reflect"
	"testing"
	"time"
)

func Test_TimeStruct_CalcEndTime(t *testing.T) {
	tests := []struct {
		name string
		tr   timeStruct
		want time.Time
	}{
		// TODO: Add test cases.
		{"First test", timeStruct{}, NewTimeStruct().CalcEndTime()},
		{"Second test", timeStruct{}, NewTimeStruct().ScheduleDate()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.CalcEndTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeStruct.CalcEndTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TimeStruct_TimeBefore(t *testing.T) {
	type args struct {
		time1 time.Time
	}

	tests := []struct {
		name string
		tr   timeStruct
		args args
		want bool
	}{
		// TODO: Add test cases.
		// {"First test", timeStruct{}, args{time.Now()}, false},
		{"Second test", timeStruct{}, args{time.Date(2023, time.December, 21, 0, 0, 0, 0, time.UTC)}, false},
		{"Third test", timeStruct{}, args{time.Date(2022, time.December, 21, 0, 0, 0, 0, time.UTC)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TimeBefore(tt.args.time1); got != tt.want {
				t.Errorf("timeStruct.TimeBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}
