package handlers

import (
	"reflect"
	"testing"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func TestMakepath(t *testing.T) {
	type args struct {
		m storage.MemStorage
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "тест формирования url",
			args: args{storage.MemStorage{Countermem: map[string]storage.Counter{"PollCounter": 2}}},
			want: "http://localhost:8080/update/counter/PollCounter/2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := agentconfigs.Makepath(tt.args.m, "localhost:8080"); !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("Makepath() = %v, want %v", got[0], tt.want)
			}
		})
	}
}

func TestPostJSONMetrics(t *testing.T) {
	type args struct {
		b []byte
		a string
		r bool
		f agentconfigs.AgentData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no server running",
			args:    args{b: []byte{}, a: "localhost:8080", f: agentconfigs.AgentData{}, r: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostJSONMetrics(tt.args.b, tt.args.a, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("PostJSONMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgentWorkers(t *testing.T) {
	type args struct {
		data WorkerData
	}
	tests := []struct {
		name string
		args args
	}{{name: "Просто для покрытия",
		args: args{
			data: WorkerData{Batchrout: "localhost:8080"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AgentWorkers(tt.args.data)
		})
	}
}
