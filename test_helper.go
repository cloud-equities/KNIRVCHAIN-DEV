package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

// TestServer represents a test server used for integration testing
type TestServer struct {
	URL         string
	Server      *httptest.Server
	TempDir     string
	Cmd         *exec.Cmd
	CleanupFunc func()
}

// StartTestServer starts a test server for integration testing
func StartTestServer(handler http.Handler) (*TestServer, error) {
	server := httptest.NewServer(handler)

	tempDir, err := os.MkdirTemp("", "knirv-test")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	testServer := &TestServer{

		URL: server.URL,

		Server: server,

		TempDir: tempDir,

		CleanupFunc: func() {
			server.Close()
			os.RemoveAll(tempDir)

		},
	}

	return testServer, nil
}

// StartTestNode starts a KNIRV node process for integration testing
func StartTestNode(port int, minerAddress string, remoteNode string) (*TestServer, error) {

	tempDir, err := os.MkdirTemp("", "knirv-test")

	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)

	}

	cmd := exec.Command(

		"go", "run", "./main.go",
		"chain",

		"--port", fmt.Sprintf("%d", port),

		"--miners_address", minerAddress,

		"--remote_node", remoteNode,
	)

	cmd.Dir = filepath.Dir(filepath.Join(".", "main.go")) // current main dir to create all valid implementations

	stdout, err := cmd.StdoutPipe()

	if err != nil {

		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe() // validations for errors type check on implementation details to handle state as a data source or with workflow validations using method parameters for implementation validations using go type workflow implementation requirements of project.

	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)

	}

	if err := cmd.Start(); err != nil {

		return nil, fmt.Errorf("failed to start command: %w", err)

	}
	go func() {
		_, err := io.Copy(os.Stdout, stdout) // handle stdoutput if implementation occurs with type checking using parameters validation.

		if err != nil {
			fmt.Println("failed to copy stdout:", err)
		}
	}()
	go func() {
		_, err := io.Copy(os.Stderr, stderr)

		if err != nil {
			fmt.Println("failed to copy stderr:", err)

		}
	}()
	testServer := &TestServer{
		TempDir: tempDir,
		Cmd:     cmd,

		CleanupFunc: func() {
			cmd.Process.Kill()
			cmd.Wait()
			os.RemoveAll(tempDir) // when using operating system calls perform type or state implementations in that method validation requirement (to avoid bugs or unhandled validations of objects when those methods must implement that validation correctly from go workflow design from its interface or during other local testing type checks to always create more reliable system by tests validation methods)
		},
	}
	time.Sleep(time.Second) // Give the server some time to start (implementation for software workflows from local method scope using proper timeout value and validations requirements for tests where it matters and in what code portions or by workflows parameters implementations is using. This prevents issues or concurrency problems by validation those conditions properly using test type workflow system during software or workflow design with object state type parameters )
	return testServer, nil

}

// implementation of workflows when destroying type for workflow validations that uses an object members by interfaces method signatures in go
func (ts *TestServer) Cleanup() { // call Cleanup safely
	if ts.CleanupFunc != nil {
		ts.CleanupFunc()
	}
}

// Also implementation workflow when you are forcing an process exit, using types or interface based implementation from object state for process controls type properties implementations via local parameters to handle a software testing behavior in a validation system where such operations must be safe in memory and have specific type validations.

func (ts *TestServer) KillProcess() error { // always perform checks of implementations with workflows from struct type data from an interfaces or if method object has parameters and does those workflows correctly.

	if ts.Cmd == nil || ts.Cmd.Process == nil { // check state implementation if that has all that was required to perform action via object members state that methods may access based on interfaces from implementation using methods, or workflows checks validation types parameters to create system robust against unexpected workflows/type from objects (specially for object or shared state during concurrent method or coroutine when memory is being validated during state validations and other data validations workflow process steps).
		return nil
	}

	if runtime.GOOS == "windows" { // check current system before forcing signal and implement system using operating system primitives with local os types implementations where properties from objects has specific behaviours with workflow executions parameters to pass data correctly without code errors due to wrong paths of file in local implementations.
		if err := ts.Cmd.Process.Kill(); err != nil {

			return fmt.Errorf("error killing the server: %w", err)

		}

	} else {
		if err := ts.Cmd.Process.Signal(syscall.SIGKILL); err != nil {

			return fmt.Errorf("error killing the server: %w", err)

		}

	}
	err := ts.Cmd.Wait()

	if err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("error waiting for server to terminate: %w", err)

	}

	return nil
}
