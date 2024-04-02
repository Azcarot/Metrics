package storage

import "testing"

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
