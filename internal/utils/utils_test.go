package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestSplitPreset(t *testing.T) {
	type args struct {
		preset string
	}
	tests := []struct {
		name string
		args args
		want []Exercise
	}{
		{
			name: "strength workout",
			args: args{
				preset: "7:[12*14,12*14,12*14,12*14];8:[10*100,10*100,10*100,10*100];9:[12*40,12*40,12*50,12*50];10:[12*20,12*20,12*20,12*20];11:[14*15,14*15,14*15,14*15];12:[5*15,10*10,10*10,10*10]",
			},
			want: []Exercise{
				{
					ID: 7,
					Sets: []Set{
						{Reps: 12, Weight: 14},
						{Reps: 12, Weight: 14},
						{Reps: 12, Weight: 14},
						{Reps: 12, Weight: 14},
					},
				},
				{
					ID: 8,
					Sets: []Set{
						{Reps: 10, Weight: 100},
						{Reps: 10, Weight: 100},
						{Reps: 10, Weight: 100},
						{Reps: 10, Weight: 100},
					},
				},
				{
					ID: 9,
					Sets: []Set{
						{Reps: 12, Weight: 40},
						{Reps: 12, Weight: 40},
						{Reps: 12, Weight: 50},
						{Reps: 12, Weight: 50},
					},
				},
				{
					ID: 10,
					Sets: []Set{
						{Reps: 12, Weight: 20},
						{Reps: 12, Weight: 20},
						{Reps: 12, Weight: 20},
						{Reps: 12, Weight: 20},
					},
				},
				{
					ID: 11,
					Sets: []Set{
						{Reps: 14, Weight: 15},
						{Reps: 14, Weight: 15},
						{Reps: 14, Weight: 15},
						{Reps: 14, Weight: 15},
					},
				},
				{
					ID: 12,
					Sets: []Set{
						{Reps: 5, Weight: 15},
						{Reps: 10, Weight: 10},
						{Reps: 10, Weight: 10},
						{Reps: 10, Weight: 10},
					},
				},
			},
		},
		{
			name: "cardio workout",
			args: args{
				preset: "19:[10];20:[15]",
			},
			want: []Exercise{
				{
					ID: 19,
					Sets: []Set{
						{Minutes: 10},
					},
				},
				{
					ID: 20,
					Sets: []Set{
						{Minutes: 15},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitPreset(tt.args.preset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitPreset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPreset(t *testing.T) {
	type args struct {
		preset string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "one word",
			args: args{
				preset: "12*14",
			},
			want: true,
		},
		{
			name: "two words",
			args: args{
				preset: "12*14,12*14",
			},
			want: true,
		},
		{
			name: "many words",
			args: args{
				preset: "12*14,3*5,100*200",
			},
			want: true,
		},
		{
			name: "with cardio",
			args: args{
				preset: "12*14,20,100*200",
			},
			want: true,
		},
		{
			name: "empty",
			args: args{
				preset: "",
			},
			want: false,
		},
		{
			name: "extra comma",
			args: args{
				preset: "12*14,",
			},
			want: false,
		},
		{
			name: "extra comma 2",
			args: args{
				preset: "12*14,,12*14",
			},
			want: false,
		},
		{
			name: "extra asterisk",
			args: args{
				preset: "*14,12*14",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPreset(tt.args.preset); got != tt.want {
				t.Errorf("IsValidPreset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetThisWeek(t *testing.T) {
	type args struct {
		dateStr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "29.12.25 – 04.01.26",
			args: args{
				dateStr: "2026-01-02 14:25:06 +00:00",
			},
			want: "29.12.25 – 04.01.26",
		},
		{
			name: "29.12.25 – 04.01.26",
			args: args{
				dateStr: "2026-01-04 10:23:09 +00:00",
			},
			want: "29.12.25 – 04.01.26",
		},
		{
			name: "05.01.26 – 11.01.26",
			args: args{
				dateStr: "2026-01-05 16:26:42 +00:00",
			},
			want: "05.01.26 – 11.01.26",
		},
		{
			name: "05.01.26 – 11.01.26",
			args: args{
				dateStr: "2026-01-07 10:22:52 +00:00",
			},
			want: "05.01.26 – 11.01.26",
		},
		{
			name: "05.01.26 – 11.01.26",
			args: args{
				dateStr: "2026-01-08 16:35:58 +00:00",
			},
			want: "05.01.26 – 11.01.26",
		},
		{
			name: "05.01.26 – 11.01.26",
			args: args{
				dateStr: "2026-01-08 18:06:34 +00:00",
			},
			want: "05.01.26 – 11.01.26",
		},
		{
			name: "12.01.26 – 18.01.26",
			args: args{
				dateStr: "2026-01-12 16:57:09 +00:00",
			},
			want: "12.01.26 – 18.01.26",
		},
		{
			name: "12.01.26 – 18.01.26",
			args: args{
				dateStr: "2026-01-15 17:26:13 +00:00",
			},
			want: "12.01.26 – 18.01.26",
		},
		{
			name: "12.01.26 – 18.01.26",
			args: args{
				dateStr: "2026-01-17 09:40:33 +00:00",
			},
			want: "12.01.26 – 18.01.26",
		},
		{
			name: "12.01.26 – 18.01.26",
			args: args{
				dateStr: "2026-01-17 11:01:47 +00:00",
			},
			want: "12.01.26 – 18.01.26",
		},
		{
			name: "12.01.26 – 18.01.26",
			args: args{
				dateStr: "2026-01-18 17:52:59 +00:00",
			},
			want: "12.01.26 – 18.01.26",
		},
		{
			name: "19.01.26 – 25.01.26",
			args: args{
				dateStr: "2026-01-19 16:46:47 +00:00",
			},
			want: "19.01.26 – 25.01.26",
		},
		{
			name: "19.01.26 – 25.01.26",
			args: args{
				dateStr: "2026-01-21 16:04:36 +00:00",
			},
			want: "19.01.26 – 25.01.26",
		},
		{
			name: "19.01.26 – 25.01.26",
			args: args{
				dateStr: "2026-01-22 15:51:35 +00:00",
			},
			want: "19.01.26 – 25.01.26",
		},
		{
			name: "19.01.26 – 25.01.26",
			args: args{
				dateStr: "2026-01-24 09:03:42 +00:00",
			},
			want: "19.01.26 – 25.01.26",
		},
		{
			name: "19.01.26 – 25.01.26",
			args: args{
				dateStr: "2026-01-24 15:06:31 +00:00",
			},
			want: "19.01.26 – 25.01.26",
		},
		{
			name: "26.01.26 – 01.02.26",
			args: args{
				dateStr: "2026-01-26 15:10:13 +00:00",
			},
			want: "26.01.26 – 01.02.26",
		},
		{
			name: "26.01.26 – 01.02.26",
			args: args{
				dateStr: "2026-01-28 16:00:18 +00:00",
			},
			want: "26.01.26 – 01.02.26",
		},
		{
			name: "26.01.26 – 01.02.26",
			args: args{
				dateStr: "2026-01-30 15:25:53 +00:00",
			},
			want: "26.01.26 – 01.02.26",
		},
		{
			name: "26.01.26 – 01.02.26",
			args: args{
				dateStr: "2026-01-31 14:18:46 +00:00",
			},
			want: "26.01.26 – 01.02.26",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := "2006-01-02 15:04:05 -07:00"
			date, err := time.Parse(layout, tt.args.dateStr)
			assert.NoError(t, err)
			if got := GetThisWeek(date); got != tt.want {
				t.Errorf("GetThisWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}
