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

func TestChainStartup(t *testing.T) {
	tempDir := "test_main_chain"
	defer os.RemoveAll(tempDir)
	dbPath := filepath.Join(tempDir, "knirv.db")
	var wg sync.WaitGroup
	var outputString string

	wg.Add(1)
	testConfig := Config{ // defining custom parameters instead via environment
		Port:          5000,
		MinersAddress: "testAddress",
		DatabasePath:  dbPath, // Corrected tests to always define test configs when tests perform validation steps from workflows implementations
	}
	go func() {
		defer wg.Done()
		// Set command-line arguments for test execution from this config variable that implements interface struct and its properties types properly using type variables to have implementations and testing workflow from software running with methods as are defined by each object implementation during the workflow process and those logic that test validations enforce.
		os.Args = []string{"cmd", "-port", strconv.Itoa(int(testConfig.Port)), "-miners_address", testConfig.MinersAddress, "-database_path", testConfig.DatabasePath}

		main() // Use method in tests

	}()
	time.Sleep(25 * time.Second)
	outputString = captureOutput(t, func() {
		fmt.Println("checking messages")
	})
	logPrefix := constants.BLOCKCHAIN_NAME + ":"
	messagesToCheck := []string{
		logPrefix + " Starting the consensus algorithm...\n",
		logPrefix + " Mined block number:\n",
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
		t.Fatalf("Ports were leaked due to incorrect testing resource cleanup: %s", string(output2))
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
		f()
		w.Close()
		if _, err := io.Copy(&buf, r); err != nil {
			t.Errorf("Error reading from pipe: %v", err)
		}
		out <- buf.String()
	}()
	os.Stdout = old
	return <-out
}
