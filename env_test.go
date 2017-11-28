package cfsvcenv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestBind(t *testing.T) {
	if s := os.Getenv("GO_TEST_BIND_INDEX"); s != "" {
		index, err := strconv.Atoi(s)
		if err != nil {
			t.Fatal(err)
		}
		testBind(t, index)
		return
	}
	for index := range bindTests {
		cmd := exec.Command(os.Args[0], "-test.run=TestBind")
		cmd.Env = []string{fmt.Sprintf("GO_TEST_BIND_INDEX=%d", index)}
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			if str := stdout.String(); str != "" {
				t.Logf("%s", str)
			}
			if str := stderr.String(); str != "" {
				t.Logf("%s", str)
			}
			t.Fatalf("process ran with err %v, want none", err)
		}
	}
}

var bindTests = []struct {
	name    string
	before  func()
	wantErr bool
	wantEnv map[string]string
}{
	{"VCAP_SERVICES not set", func() {
		os.Unsetenv(vcapServicesEnv)
	}, false, nil},
	{"VCAP_SERVICES set, but empty", func() {
		os.Setenv(vcapServicesEnv, "")
	}, false, nil},
	{"VCAP_SERVICES set, but invalid", func() {
		os.Setenv(vcapServicesEnv, "blah")
	}, true, nil},
	{"VCAP_SERVICES valid, no services", func() {
		os.Setenv(vcapServicesEnv, "{}")
	}, false, nil},
	{"service with unsupported credentials format", func() {
		os.Setenv(vcapServicesEnv, `{
			"postgres": [{
				"name": "postgres",
				"credentials": ["a", "b"]
			}]
		}`)
	}, false, nil},
	{"service with supported credentials format", func() {
		os.Setenv(vcapServicesEnv, `{
			"postgres": [{
				"name": "postgres",
				"credentials": {
					"username": "u",
					"password": "p"
				}
			}]
		}`)
	}, false, map[string]string{
		"POSTGRES_USERNAME": "u",
		"POSTGRES_PASSWORD": "p",
	}},
	{"multiple services", func() {
		os.Setenv(vcapServicesEnv, `{
			"s": [{
				"name": "s",
				"credentials": {
					"a": "b"
				}
			}],
			"user-provided": [{
				"name": "t",
				"credentials": {
					"a": "b"
				}
			}, {
				"name": "u",
				"credentials": {
					"a": "c"
				}
			}]
		}`)
	}, false, map[string]string{
		"S_A": "b",
		"T_A": "b",
		"U_A": "c",
	}},
	{"service with non-string credentials", func() {
		os.Setenv(vcapServicesEnv, `{
			"s": [{
				"name": "s",
				"credentials": {
					"a": true,
					"b": 123,
					"c": 0.123
				}
			}]
		}`)
	}, false, map[string]string{
		"S_A": "true",
		"S_B": "123",
		"S_C": "0.123",
	}},
	{"service credential names transformed", func() {
		os.Setenv(vcapServicesEnv, `{
			"s": [{
				"name": "s",
				"credentials": {
					"a-b": "c",
					"d_E": "f",
					"g-h-i": 0.123
				}
			}]
		}`)
	}, false, map[string]string{
		"S_A_B":   "c",
		"S_D_E":   "f",
		"S_G_H_I": "0.123",
	}},
}

func testBind(t *testing.T, index int) {
	tt := bindTests[index]
	t.Run(tt.name, func(t *testing.T) {
		tt.before()
		got := Bind()
		switch got {
		case nil:
			if tt.wantErr {
				t.Fatal("got no error, want one")
			}
		default:
			if !tt.wantErr {
				t.Fatalf("got error %v, want none", got)
			}
		}
		for k, v := range tt.wantEnv {
			if got := os.Getenv(k); got != v {
				t.Errorf("got env %s=%v, want %s=%v", k, got, k, v)
			}
		}
	})
}
