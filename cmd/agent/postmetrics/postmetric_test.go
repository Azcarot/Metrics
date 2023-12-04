package postmetrics

import (
	"agent/measure"
	"reflect"
	"testing"
)

func TestMakepath(t *testing.T) {
	type args struct {
		m measure.MemStorage
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "тест формирования url",
			args: args{measure.MemStorage{Countermem: map[string]measure.Counter{"PollCounter": 2}}},
			want: "http://localhost:8080/update/counter/PollCounter/2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Makepath(tt.args.m); !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("Makepath() = %v, want %v", got[0], tt.want)
			}
		})
	}
}
