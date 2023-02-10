package timeSrv

import (
	"reflect"
	"testing"
	"time"
)

func Test_timeStruct_CalcEndTime(t *testing.T) {
	tests := []struct {
		name string
		tr   timeStruct
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := timeStruct{}
			if got := tr.CalcEndTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeStruct.CalcEndTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
