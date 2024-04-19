package storage

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fileOrPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{{
		name: "wrong path", args: args{path: "00000"}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileOrPathExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileOrPathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fileOrPathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewProducer(t *testing.T) {
	type args struct {
		fileName string
	}

	tests := []struct {
		name    string
		args    args
		want    *Producer
		wantErr bool
	}{
		{name: "file", args: args{fileName: "test"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProducer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewProducer() = %v", got)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    *Consumer
		wantErr bool
	}{
		{name: "file", args: args{fileName: "test"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsumer(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewConsumer() = %v", got)
			}
		})
	}
}

func TestConsumer_ReadEvent(t *testing.T) {

	value := float64(10)
	metrics := []Metrics{{ID: "test", MType: GuageType, Value: &value}}
	filename := "test"
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
	WriteToFile(filename, metrics[0])
	consumer, err := NewConsumer(filename)
	require.NoError(t, err)
	tests := []struct {
		name    string
		c       *Consumer
		want    *[]Metrics
		wantErr bool
	}{
		{name: "normal file", c: consumer, want: &metrics, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.ReadEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Consumer.ReadEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Consumer.ReadEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
