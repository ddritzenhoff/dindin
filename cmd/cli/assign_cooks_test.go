package main

import (
	"testing"
	"time"
)

func Test_getDayDifference(t *testing.T) {
	type args struct {
		now  time.Weekday
		then time.Weekday
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sunday-sunday", args{time.Sunday, time.Sunday}, 0},
		{"sunday-monday", args{time.Sunday, time.Monday}, 1},
		{"sunday-tuesday", args{time.Sunday, time.Tuesday}, 2},
		{"sunday-wednesday", args{time.Sunday, time.Wednesday}, 3},
		{"sunday-thursday", args{time.Sunday, time.Thursday}, 4},
		{"sunday-friday", args{time.Sunday, time.Friday}, 5},
		{"sunday-saturday", args{time.Sunday, time.Saturday}, 6},
		{"monday-sunday", args{time.Monday, time.Sunday}, 6},
		{"monday-monday", args{time.Monday, time.Monday}, 0},
		{"monday-tuesday", args{time.Monday, time.Tuesday}, 1},
		{"monday-wednesday", args{time.Monday, time.Wednesday}, 2},
		{"monday-thursday", args{time.Monday, time.Thursday}, 3},
		{"monday-friday", args{time.Monday, time.Friday}, 4},
		{"monday-saturday", args{time.Monday, time.Saturday}, 5},
		{"tuesday-sunday", args{time.Tuesday, time.Sunday}, 5},
		{"tuesday-monday", args{time.Tuesday, time.Monday}, 6},
		{"tuesday-tuesday", args{time.Tuesday, time.Tuesday}, 0},
		{"tuesday-wednesday", args{time.Tuesday, time.Wednesday}, 1},
		{"tuesday-thursday", args{time.Tuesday, time.Thursday}, 2},
		{"tuesday-friday", args{time.Tuesday, time.Friday}, 3},
		{"tuesday-saturday", args{time.Tuesday, time.Saturday}, 4},
		{"wednesday-sunday", args{time.Wednesday, time.Sunday}, 4},
		{"wednesday-monday", args{time.Wednesday, time.Monday}, 5},
		{"wednesday-tuesday", args{time.Wednesday, time.Tuesday}, 6},
		{"wednesday-wednesday", args{time.Wednesday, time.Wednesday}, 0},
		{"wednesday-thursday", args{time.Wednesday, time.Thursday}, 1},
		{"wednesday-friday", args{time.Wednesday, time.Friday}, 2},
		{"wednesday-saturday", args{time.Wednesday, time.Saturday}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDayDifference(tt.args.now, tt.args.then); got != tt.want {
				t.Errorf("getDayDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}
