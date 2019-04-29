package customtime_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/burubur/go-microservice/internal/platform/customtime"
	"github.com/stretchr/testify/assert"
)

var (
	testTimezoneDefault *time.Location
	testTimezoneBangkok *time.Location
)

func init() {
	testTimezoneDefault, _ = time.LoadLocation("Asia/Jakarta")
	testTimezoneBangkok, _ = time.LoadLocation("Asia/Bangkok")
}

func TestSetLocation(t *testing.T) {
	tests := []struct {
		name    string
		tz      string
		wantErr bool
	}{
		{
			name:    "0. Negative case - set location with undefined timezone location",
			tz:      "Asia",
			wantErr: true,
		},
		{
			name:    "1. Positive case - set location with empty string",
			tz:      "",
			wantErr: false,
		},
		{
			name:    "2. Positive case - set location with valid timezone location",
			tz:      "Asia/Bangkok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := customtime.SetLocation(tt.tz); (err != nil) != tt.wantErr {
				t.Errorf("SetLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocation(t *testing.T) {

	t.Run("set location to Asia/Bangkok", func(t *testing.T) {
		_ = customtime.SetLocation("Asia/Bangkok")
		if got := customtime.Location(); !reflect.DeepEqual(got, testTimezoneBangkok) {
			t.Errorf("Location() = %v, want %v", got, testTimezoneBangkok)
		}
	})
	t.Run("set location to default timezone", func(t *testing.T) {
		_ = customtime.SetLocation("Asia/Jakarta")
		if got := customtime.Location(); !reflect.DeepEqual(got, testTimezoneDefault) {
			t.Errorf("Location() = %v, want %v", got, testTimezoneDefault)
		}
	})
}

func TestNow(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "0. time case 1",
			want: time.Now(),
		},
		{
			name: "1. time case 2",
			want: time.Now(),
		},
		{
			name: "2. time case 3",
			want: time.Now(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := customtime.Now()
			assert.WithinDuration(t, tt.want, got, time.Minute, "at least return the same minute")
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		layout string
		value  string
	}
	want, _ := time.Parse("2006-01-02 15:04:05 +0700 WIB", "2018-01-01 00:00:00 +0700 WIB")
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "parse year to time",
			args: args{
				layout: "2006",
				value:  "2018",
			},
			want:    want.Format("2006-01-02 15:04:05.999999 +0700 WIB"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := customtime.Parse(tt.args.layout, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Format("2006-01-02 15:04:05.999999 +0700 WIB"), tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
