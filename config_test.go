package logsystem

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				data: []byte(`{}`),
			},
			want: Config{
				Drivers: nil,
			},
			wantErr: false,
		},
		{
			name: "two drivers",
			args: args{
				data: []byte(`{"drivers":{"driver1":{"key":"value"},"driver2":{"key2":"value2"}}}`),
			},
			want: Config{
				Drivers: map[DriverID]json.RawMessage{
					"driver1": json.RawMessage(`{"key":"value"}`),
					"driver2": json.RawMessage(`{"key2":"value2"}`),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				data: []byte(`{"drivers":{"driver1":{"key":"value"}`),
			},
			want: Config{
				Drivers: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfig(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
