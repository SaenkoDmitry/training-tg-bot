package utils

import (
	"reflect"
	"testing"
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
			name: "",
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
