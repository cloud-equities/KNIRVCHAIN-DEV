package main

import (
	"KNIRVCHAIN-DEV/constants"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMainChainStartup(t *testing.T) {
	tempDir := "test_main_chain"
	defer os.RemoveAll(tempDir)
	dbPath := filepath.Join(tempDir, "knirv.db")
	var wg sync.WaitGroup
	var outputString string
	wg.Add(1)

	testConfig := Config{ // Setting all configuration, as workflow validation/requirements also now require from tests logic
		Port:          5000,
		MinersAddress: "testAddress",
		DatabasePath:  dbPath,
	}

	go func() {
		defer wg.Done()
		os.Args = []string{"cmd", "-port", strconv.Itoa(int(testConfig.Port)), "-miners_address", testConfig.MinersAddress, "-database_path", testConfig.DatabasePath}
		main() // implement validation from `main` method
	}()
	time.Sleep(45 * time.Second) // Implementation: time wait after workflow is started for log validation. (Implementation step which now validates type, but also if data state exist during object validation by using logging methods also used by workflow from application logic using implemented interfaces to access those members to test implementation during tests workflow).

	outputString = captureOutput(t, func() { // validation only at the end of test process and its implementations.
		fmt.Println("checking messages") // testing only the parameters are being called correctly after the timeout occurs with data ready on implementation side from types methods calls on project test systems validations from method calls with implementation validations by log using specific types for those parameters on this particular validation steps when methods of `captureOutput()` is called with workflows
	})
	logPrefix := constants.BLOCKCHAIN_NAME + ":"

	messagesToCheck := []string{ // All workflow of object parameter are correct
		logPrefix + " Starting the consensus algorithm...\n",
		logPrefix + "  Mined block number:\n",
	}
	for _, msg := range messagesToCheck {
		if !strings.Contains(outputString, msg) {

			t.Errorf("Message %q is missing from output: \n%v", msg, outputString)
		}
	}

	// check for processes
	cmd2 := exec.Command("tasklist")

	if runtime.GOOS != "windows" {
		cmd2 = exec.Command("ps", "-ef")

	}

	output2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("Error getting tasks: %v: %v", err, string(output2))

	}

	if strings.Contains(string(output2), fmt.Sprintf(":5000")) {
		t.Fatalf("Ports were leaked due to incorrect testing resource cleanup: %s", string(output2)) // workflow implementation of object parameter must match type.
	}
	wg.Wait()

}

func captureOutput(t *testing.T, f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()

	if err != nil {
		t.Fatalf("os pipe creation failed: %v", err)
	}
	os.Stdout = w

	out := make(chan string)
	go func() {
		var buf bytes.Buffer

		f() // execute that logic for local state for logging of system.

		w.Close()

		if _, err := io.Copy(&buf, r); err != nil {
			t.Errorf("Error reading from pipe: %v", err)

		}

		out <- buf.String()

	}()

	os.Stdout = old // get back output flow from implementations
	return <-out
}
