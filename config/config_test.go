package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/claudiodangelis/qrcp/application"
)

func TestNew(t *testing.T) {
	os.Clearenv()
	_, f, _, _ := runtime.Caller(0)
	testdir := filepath.Join(filepath.Dir(f), "testdata")
	tempfile, err := ioutil.TempFile("", "qrcp*tmp.yml")
	if err != nil {
		t.Skip()
	}
	defer os.Remove(tempfile.Name())
	partialconfig, err := ioutil.TempFile("", "qrcp*partial.yml")
	if err != nil {
		panic(err)
	}
	defer os.Remove(partialconfig.Name())
	if err := os.WriteFile(partialconfig.Name(), []byte(`port: 9090`), os.ModePerm); err != nil {
		panic(err)
	}
	type args struct {
		app application.App
	}

	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"partial", args{
				app: application.App{
					Flags: application.Flags{
						Config: partialconfig.Name(),
					},
				},
			},
			Config{
				Interface: "wlo1",
				Port:      9090,
			},
		},
		{
			"init", args{
				app: application.App{
					Flags: application.Flags{
						Config: tempfile.Name(),
					},
				},
			},
			Config{
				Interface: "wlo1",
			},
		},
		{
			"#2", args{
				app: application.App{
					Flags: application.Flags{
						Config: filepath.Join(testdir, "qrcp.yml"),
					},
				},
			},
			Config{
				Interface: "wlo1",
			},
		},
		{
			"#2", args{
				app: application.App{
					Flags: application.Flags{
						Config: filepath.Join(testdir, "full.yml"),
					},
				},
			},
			Config{
				Interface: "wlo1",
				Port:      18080,
				KeepAlive: false,
				Bind:      "10.20.30.40",
				Path:      "random",
				Secure:    false,
				TlsKey:    "/path/to/key",
				TlsCert:   "/path/to/cert",
				FQDN:      "mylan.com",
				Output:    "/path/to/default/output/dir",
			},
		},
		{
			"overrides", args{
				app: application.App{
					Flags: application.Flags{
						Config: filepath.Join(testdir, "full.yml"),
						Port:   99999,
					},
				},
			},
			Config{
				Interface: "wlo1",
				Port:      99999,
				Bind:      "10.20.30.40",
				KeepAlive: false,
				Path:      "random",
				Secure:    false,
				TlsKey:    "/path/to/key",
				TlsCert:   "/path/to/cert",
				FQDN:      "mylan.com",
				Output:    "/path/to/default/output/dir",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
