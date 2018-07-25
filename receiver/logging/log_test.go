package logging

import (
	"reflect"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	tests := []struct {
		name string
		want *Opts
	}{
		{
			name: "Test Calling",
			want: &Opts{
				Level:     1,
				DebugPath: "",
				InfoPath:  "",
				ErrorPath: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		opts Opts
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Default",
			args: args{
				Opts{
					Level:     1,
					DebugPath: "",
					InfoPath:  "",
					ErrorPath: "",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Level",
			args: args{
				Opts{
					Level:     4,
					DebugPath: "",
					InfoPath:  "",
					ErrorPath: "",
				},
			},
			wantErr: true,
		},
		{
			name: "NoWritePermDebug",
			args: args{
				Opts{
					Level:     1,
					DebugPath: "/debug.not.writeable",
					InfoPath:  "",
					ErrorPath: "",
				},
			},
			wantErr: true,
		},
		{
			name: "NonexistentPath",
			args: args{
				Opts{
					Level:     1,
					DebugPath: "/home/mattu/should/not/exist",
					InfoPath:  "",
					ErrorPath: "",
				},
			},
			wantErr: true,
		},
		{
			name: "PathInHome",
			args: args{
				Opts{
					Level:     1,
					DebugPath: "/home/mattu/test_debug.log",
					InfoPath:  "",
					ErrorPath: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Init(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
