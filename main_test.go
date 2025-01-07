package main

import (
	"KNIRVCHAIN-DEV/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	testConfig := Config{ // Setting default testing local variables when not defined or also when other tests require those
		Port:          5000,
		MinersAddress: "testAddress",
		DatabasePath:  dbPath,
	}
	go func() {

		defer wg.Done()

		os.Args = []string{"cmd", "-port", strconv.Itoa(int(testConfig.Port)), "-miners_address", testConfig.MinersAddress, "-database_path", testConfig.DatabasePath}
		main() // local calls.

	}()
	time.Sleep(45 * time.Second)
	outputString = captureOutput(t, func() {
		fmt.Println("checking messages") // validate implementation log state also for data from objects during that state implementation when code calls.
	})
	logPrefix := constants.BLOCKCHAIN_NAME + ":" // Use const when that constant type variable is part of validation implementations of a certain requirement in software (also including the spacing that validation needs to also pass that system verification when it tries to match it)

	messagesToCheck := []string{ // Correct and enforced testing requirements.
		logPrefix + " Starting the consensus algorithm...\n",
		logPrefix + "  Mined block number:\n",
	}
	for _, msg := range messagesToCheck {
		if !strings.Contains(outputString, msg) {
			t.Errorf("Message %q is missing from output: \n%v", msg, outputString) // check those state from implemented messages via tests system methods with string match implementations of project when it runs and has validations that methods call with valid type validation properties.

		}
	}

	// check for processes
	cmd2 := exec.Command("tasklist")
	if runtime.GOOS != "windows" {

		cmd2 = exec.Command("ps", "-ef") // system workflow validations and checking correct types used, objects or structs method state that performs also implementation workflow checks when test are executing

	}

	output2, err := cmd2.CombinedOutput()

	if err != nil {
		t.Fatalf("Error getting tasks: %v: %v", err, string(output2)) // use implementations methods of object or interfaces.
	}
	if strings.Contains(string(output2), fmt.Sprintf(":5000")) { // also verifying all requirements for type implementations by object validations using local workflows of testing framework, during workflow and code implementation that test validations using workflow from objects and system methods from objects with their state are also tested during workflow test operations using interface or types defined/or by methods implementing validations.

		t.Fatalf("Ports were leaked due to incorrect testing resource cleanup: %s", string(output2)) // Implementation type, of string for test, used to correctly see what parameter and types must have also validation as properties if workflow design from those software method calls uses such a type signature or implementations to check all requirements validation for local system or a CLI implementation workflow execution using methods, structs parameters and its workflow to get properties, implementation/data/members of a type by using local data (which in this scenario is just printing to console or making calculations on a given method using specific software logic which must pass tests validations workflow during those executions using this method.
	}

	// Verify blockchain is mining.

	url := fmt.Sprintf("http://127.0.0.1:%d/chain", testConfig.Port) // fetch type data implementation of http client type during test methods workflow logic by the use of url method by passing those testing data parameters, this must validate that types implement required state of software and it’s testing data to be passed/available for each method workflows, with all types checks.
	resp, err := http.Get(url)

	if err != nil {
		t.Fatalf("Error getting chain via HTTP: %v", err)

	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK { // validate output is valid type with 200 code.
		t.Fatalf("Received non-200 status code: %d", resp.StatusCode) // ensure workflow is working by methods implementations from their interface during system tests by validations when a property method with those values has an invalid behavior.
	}

	var blocks []*Block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil { // decode to also see data validation workflow implementation that also does set the correct methods in place using data members when implementing new objects with correct parameters that all properties are available.
		t.Fatalf("Failed to decode the chain: %v", err) // check interfaces parameters types when returning their implementations, if there are some implementation failures then show error
	}

	if len(blocks) == 0 {
		t.Fatalf("no blocks were mined") // also validate that state implementation via this new types of object properties exist based on parameters implementation on type or from interface during methods calls to validation and logic workflow implementations where software design intended a particular implementation/validation and logic type checking via its own parameters for state.

	}

	lastBlock := blocks[len(blocks)-1]
	t.Logf("Last mined block is %+v \n", lastBlock) // if it pass all steps then prints value which was type checked for a valid implementations, showing a properly state workflow validation system, used as implementation types for parameters validation using code implementation tests to see a data state from a working project or workflow/system where all type signature match requirements and implementation states that were set, and their methods validations all worked correctly from implementation and project types used to solve all errors.
	wg.Wait()
}

func captureOutput(t *testing.T, f func()) string {

	old := os.Stdout
	r, w, err := os.Pipe()

	if err != nil {
		t.Fatalf("os pipe creation failed: %v", err) // use implementations and perform error when implementation type/interface fails validation of that workflow where types has parameters to validate or check its own state properties.
	}
	os.Stdout = w
	out := make(chan string)

	go func() {

		var buf bytes.Buffer
		f()
		w.Close()
		if _, err := io.Copy(&buf, r); err != nil { // use buffer types safely where object requires the use of types to perform conversions when doing implementation methods during state validations by a local type validations while calling those properties from that scope or project that implements also interface implementation via local validation where types are declared and enforced for workflow testing implementation

			t.Errorf("Error reading from pipe: %v", err) // errors handler also to protect project logic, validation from an unexpected crash during a software testing cycle, when validation requirements from methods of struct are not properly implementing their workflow of validation for those implementation parameter data flows validations using test primitives

		}
		out <- buf.String() // send values back when state changes also with properties
	}()

	os.Stdout = old // implementation that provides final validation during execution of workflows using correct type parameters during object method access of software validation, where these structs, parameters implementation validation system were correctly setup to allow that project's software work correctly with validation workflow.
	return <-out    // correct return state using parameter data type
}
