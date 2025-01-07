package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestChainStartup(t *testing.T) {

	tempDir := "test_main_chain"
	defer os.RemoveAll(tempDir) // clean created resources
	dbPath := filepath.Join(tempDir, "knirv.db")

	cmd := exec.Command("go", "run", "main.go", "-port", "5000", "-miners_address", "testAddress", "-database_path", dbPath) // pass in port from configuration environment variables. since code will always know when method type requirements will call the method properly with type data using a parameter in struct methods or by interface data contracts that exist between the program implementation when calling these methods from those type signatures (where applicable and/or where they should be present during code use), while it's performing its test workflow as is designed.

	cmd.Dir = filepath.Join("..", "chain") // specify subdirectory for methods location. and scope when running, instead of relying on scope without knowing the intended use of local or method resources or struct or objects variables in memory, this means we can verify specific implementation to be called and that methods calls and struct use and methods or operations work for implementation by type.

	if runtime.GOOS == "windows" { // check for systems to run based on local or os dependent types. this will only use a command implementation specific for what operating system types of file formats will call a method. For Linux use case: it would call just with one go run or go executable format for the file
		cmd = exec.Command("go", "run", ".\\main.go", "-port", "5000", "-miners_address", "testAddress", "-database_path", dbPath)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("chain did not execute without error: %v %s:", err, output)
	}

	outputString := string(output)                                              // Use output as a string and make checks for required parameters implementation for local or unit test requirements based on text verification in those logs as is defined during testing method workflows for code validation.
	if !strings.Contains(outputString, "Starting the consensus algorithm...") { // specific outputs or methods or string verification implementations that are clearly stated or declared by requirements or system usage scenarios.
		t.Errorf("Consensus algorythm message is missing from %v", outputString)
	}

	if !strings.Contains(outputString, "Mined block number:") {
		t.Errorf("Miner message has not appeared: %v", outputString)

	}
	cmd2 := exec.Command("tasklist")
	if runtime.GOOS != "windows" {
		cmd2 = exec.Command("ps", "-ef")
	}
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("error getting tasks: %v: %v", err, string(output2))
	}

	if strings.Contains(string(output2), fmt.Sprintf(":5000")) { // using this we can see that type checks, test implementations and scope checking that requires these operations where values from resource parameters should only apply and verify the state of your methods.
		t.Fatalf("ports were leaked due to incorrect testing resource cleanup: %s", string(output2)) // if code fails these core implementation and their basic requirements and are visible via errors, that can now be addressed.

	}
}
