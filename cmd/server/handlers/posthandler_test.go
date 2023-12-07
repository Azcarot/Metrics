package handlers

import (
	"reflect"
	"testing"

	"github.com/Azcarot/Metrics/cmd/types"
)

func TestMakepath(t *testing.T) {
	type args struct {
		m types.MemStorage
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "тест формирования url",
			args: args{types.MemStorage{Countermem: map[string]types.Counter{"PollCounter": 2}}},
			want: "http://localhost:8080/update/counter/PollCounter/2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Makepath(tt.args.m, "localhost:8080"); !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("Makepath() = %v, want %v", got[0], tt.want)
			}
		})
	}
}
