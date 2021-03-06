package util

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	testTime, err := time.Parse("Mon Jan 02 2006 15:04:05 MST", "Wed May 22 2019 23:01:51 GMT")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		ts      string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "triple timestamp",
			ts:      "Wed May 22 2019 23:01:51 GMT+0000 (UTC)",
			want:    testTime,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
