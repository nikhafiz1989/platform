package tests_test

import (
	"bytes"
	"context"
	"github.com/influxdata/flux"
	"github.com/influxdata/platform/cmd/influxd/launcher"
	"io"
	"io/ioutil"
	nethttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/influxdata/platform"
	"github.com/influxdata/platform/http"

	_ "github.com/influxdata/flux/functions" // Import the built-in functions
	_ "github.com/influxdata/flux/functions/inputs"
	_ "github.com/influxdata/flux/functions/outputs"
	_ "github.com/influxdata/flux/functions/tests"
	_ "github.com/influxdata/flux/functions/transformations"
	_ "github.com/influxdata/flux/options"             // Import the built-in options
	_ "github.com/influxdata/platform/query/functions" // Import the built-in functions
	_ "github.com/influxdata/platform/query/functions/inputs"
	_ "github.com/influxdata/platform/query/functions/outputs"
	_ "github.com/influxdata/platform/query/options" // Import the built-in options
)

// Default context.
var ctx = context.Background()

func init() {
	flux.RegisterBuiltIn("loadTest", loadTestBuiltin)
	flux.FinalizeBuiltIns()
}

var loadTestBuiltin = `
// loadData is a function that's referenced in all the transformation tests.  
// it's registered here so that we can register a different loadData function for 
// each platform/binary.  
testLoadData = (file) => {fromCSV(file: file) |> to(bucket: "BUCKET", org: "ORG")  return from(bucket: "BUCKET") |> range(start:-500000h)}
//testLoadData = (file) => {return fromCSV(file:file)}
`

//
//return from(bucket: "BUCKET") |> range(start:-50000d)

var skipTests = map[string]string{
	"string_max":                  "error: invalid use of function: *functions.MaxSelector has no implementation for type string (https://github.com/influxdata/platform/issues/224)",
	"null_as_value":               "null not supported as value in influxql (https://github.com/influxdata/platform/issues/353)",
	"string_interp":               "string interpolation not working as expected in flux (https://github.com/influxdata/platform/issues/404)",
	"to":                          "to functions are not supported in the testing framework (https://github.com/influxdata/flux/issues/77)",
	"covariance_missing_column_1": "need to support known errors in new test framework (https://github.com/influxdata/flux/issues/536)",
	"covariance_missing_column_2": "need to support known errors in new test framework (https://github.com/influxdata/flux/issues/536)",
	"drop_before_rename":          "need to support known errors in new test framework (https://github.com/influxdata/flux/issues/536)",
	"drop_referenced":             "need to support known errors in new test framework (https://github.com/influxdata/flux/issues/536)",
	"yield":                       "yield requires special test case (https://github.com/influxdata/flux/issues/535)",
}

func withEachFluxFile(t testing.TB, fn func(filename string)) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "testdata")
	path = "/Users/adam/programming/go/src/github.com/influxdata/flux/functions/tests/testdata"
	os.Chdir(path)

	fluxFiles, err := filepath.Glob("*.flux")
	if err != nil {
		t.Fatalf("error searching for Flux files: %s", err)
	}

	for _, fluxFile := range fluxFiles {
		fn(fluxFile)
	}
}

func Test_QueryEndToEnd(t *testing.T) {
	l := RunMainOrFail(t, ctx)
	l.SetupOrFail(t)
	defer l.ShutdownOrFail(t, ctx)

	//// Execute single write against the server.
	//if resp, err := nethttp.DefaultClient.Do(l.MustNewHTTPRequest("POST", fmt.Sprintf("/api/v2/write?org=%s&bucket=%s", l.Org.ID, l.Bucket.ID), `m,k=v f=100i 946684800000000000`)); err != nil {
	//	t.Fatal(err)
	//} else if err := resp.Body.Close(); err != nil {
	//	t.Fatal(err)
	//} else if resp.StatusCode != nethttp.StatusNoContent {
	//	t.Fatalf("unexpected status code: %d", resp.StatusCode)
	//}

	withEachFluxFile(t, func(filename string) {
		reason, skip := skipTests[filename[:len(filename)-len(".flux")]]
		t.Run(filename, func(t *testing.T) {
			if skip {
				t.Skip(reason)
			}
			if qs, err := ioutil.ReadFile(filename); err != nil {
				return
			} else {
				// Query server to ensure write persists.

				//qy := `from(bucket:"BUCKET") |> range(start:2000-01-01T00:00:00Z,stop:2000-01-02T00:00:00Z)`
				qy := string(qs)
				//qy = strings.Replace(qy, "\n", " ", -1)
				//qy = strings.Replace(qy, "\r", " ", -1)
				t.Log(qy)

				var buf bytes.Buffer
				req := (http.QueryRequest{Query: qy, Org: l.Org}).WithDefaults()
				if preq, err := req.ProxyRequest(); err != nil {
					t.Fatal(err)
				} else if _, err := l.FluxService().Query(ctx, &buf, preq); err != nil {
					t.Fatal(err)
				}

				t.Log(buf.String())
			}
		})
	})
}

// Launcher is a test wrapper for main.Launcher.
type Launcher struct {
	*launcher.Launcher

	// Root temporary directory for all data.
	Path string

	// Initialized after calling the Setup() helper.
	User   *platform.User
	Org    *platform.Organization
	Bucket *platform.Bucket
	Auth   *platform.Authorization

	// Standard in/out/err buffers.
	Stdin  bytes.Buffer
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

// NewLauncher returns a new instance of Launcher.
func NewLauncher() *Launcher {
	l := &Launcher{Launcher: launcher.NewLauncher()}
	l.Launcher.Stdin = &l.Stdin
	l.Launcher.Stdout = &l.Stdout
	l.Launcher.Stderr = &l.Stderr
	if testing.Verbose() {
		l.Launcher.Stdout = io.MultiWriter(l.Launcher.Stdout, os.Stdout)
		l.Launcher.Stderr = io.MultiWriter(l.Launcher.Stderr, os.Stderr)
	}

	path, err := ioutil.TempDir("/Users/adam/programming/go/src/github.com/influxdata/flux/functions/tests/testdata", "")
	if err != nil {
		panic(err)
	}
	l.Path = path
	return l
}

// RunMainOrFail initializes and starts the server.
func RunMainOrFail(tb testing.TB, ctx context.Context, args ...string) *Launcher {
	tb.Helper()
	l := NewLauncher()
	if err := l.Run(ctx, args...); err != nil {
		tb.Fatal(err)
	}
	return l
}

// Run executes the program with additional arguments to set paths and ports.
func (l *Launcher) Run(ctx context.Context, args ...string) error {
	args = append(args, "--bolt-path", filepath.Join(l.Path, "influxd.bolt"))
	args = append(args, "--engine-path", filepath.Join(l.Path, "engine"))
	args = append(args, "--nats-path", filepath.Join(l.Path, "nats"))
	args = append(args, "--http-bind-address", "127.0.0.1:0")
	args = append(args, "--log-level", "debug")
	return l.Launcher.Run(ctx, args...)
}

// Shutdown stops the program and cleans up temporary paths.
func (l *Launcher) Shutdown(ctx context.Context) error {
	l.Cancel()
	l.Launcher.Shutdown(ctx)
	return os.RemoveAll(l.Path)
}

// ShutdownOrFail stops the program and cleans up temporary paths. Fail on error.
func (l *Launcher) ShutdownOrFail(tb testing.TB, ctx context.Context) {
	tb.Helper()
	if err := l.Shutdown(ctx); err != nil {
		tb.Fatal(err)
	}
}

// SetupOrFail creates a new user, bucket, org, and auth token. Fail on error.
func (l *Launcher) SetupOrFail(tb testing.TB) {
	svc := &http.SetupService{Addr: l.URL()}
	results, err := svc.Generate(ctx, &platform.OnboardingRequest{
		User:     "USER",
		Password: "PASSWORD",
		Org:      "ORG",
		Bucket:   "BUCKET",
	})
	if err != nil {
		tb.Fatal(err)
	}

	l.User = results.User
	l.Org = results.Org
	l.Bucket = results.Bucket
	l.Auth = results.Auth
}

func (l *Launcher) FluxService() *http.FluxService {
	return &http.FluxService{Addr: l.URL(), Token: l.Auth.Token}
}

// MustNewHTTPRequest returns a new nethttp.Request with base URL and auth attached. Fail on error.
func (l *Launcher) MustNewHTTPRequest(method, rawurl, body string) *nethttp.Request {
	req, err := nethttp.NewRequest(method, l.URL()+rawurl, strings.NewReader(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Token "+l.Auth.Token)
	return req
}
