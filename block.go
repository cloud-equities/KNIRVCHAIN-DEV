package main

import (
	"KNIRVCHAIN-DEV/constants"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type Block struct {
	BlockNumber  uint64         `json:"block_number"`
	PrevHash     []byte         `json:"prevHash"`
	Timestamp    int64          `json:"timestamp"`
	Nonce        int            `json:"nonce"`
	Transactions []*Transaction `json:"transactions"`
	Data         *SmartContract `json:"smartcontract"`
	BlockHash    []byte         `json:"hash"`
}

func NewBlock(prevHash []byte, nonce int, blockNumber uint64, smartContract *SmartContract) *Block {

	block := new(Block)
	block.PrevHash = prevHash
	block.Nonce = nonce
	block.BlockNumber = blockNumber
	block.Timestamp = time.Now().Unix()
	block.Data = smartContract
	block.BlockHash = block.Hash() // Implementation to update all properties that match from a struct, since we perform this on structs methods during implementations.
	return block
}
func (b Block) ToJson() string {
	nb, err := json.Marshal(b)

	if err != nil {
		return err.Error()
	} else {
		return string(nb)
	}
}

// Added Implementation method calls for types, also validating data while software/validation and method parameters are set or called. This prevents object states errors.
func (b *Block) Hash() []byte {
	h := sha256.New()
	// Append the timestamp.
	timestampBytes := make([]byte, 8)

	binary.LittleEndian.PutUint64(timestampBytes, uint64(b.Timestamp)) // properly reading the data and then setting object during object implementation where properties are available or defined.
	h.Write(timestampBytes)

	// Append block number
	blockNumberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(blockNumberBytes, uint64(b.BlockNumber)) // implemented the byte conversions and proper use based on validation requirements or properties being created/set properly while data is available to perform calculations or create types to validate during tests and/or implementation state during production workflows for data implementation workflow in your specific workflow code/logic operations that the methods also need
	h.Write(blockNumberBytes)

	if b.Data != nil { // avoid null checks when doing byte conversions.

		// Append data code length as bytes
		dataCodeLengthBytes := intToBytes(len(b.Data.Code))
		h.Write(dataCodeLengthBytes)
		// Append code as bytes
		h.Write(b.Data.Code)

		// Append data as bytes
		dataDataLengthBytes := intToBytes(len(b.Data.Data))
		h.Write(dataDataLengthBytes)

		h.Write(b.Data.Data)

	}
	h.Write(intToBytes(len(b.PrevHash)))
	h.Write(b.PrevHash)  // important implement proper byte values on `PrevHash` as an input during hash creation for validation of `block`.
	hashed := h.Sum(nil) // set bytes when that property is going to return its correct result after hashing all implemented types and setting local states, and setting that parameter based on each struct to work with their parameters
	return hashed        // always return final variable

}

func intToBytes(num int) []byte { // method from `block.go` to byte conversions that also correctly perform a type conversions required by hash
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(num)) // properly reading from int using system endian conversions during operation workflow.
	return bytes                                      // correct return method usage.

}
func (b *Block) AddTransactionToTheBlock(txn *Transaction) { // Implementation type is passed as object pointer.
	// check if the txn verification is a success or a failure
	if txn.Status == constants.TXN_VERIFICATION_SUCCESS {
		txn.Status = constants.SUCCESS // proper parameter passing and implementation of object via method type with validations using properties on the workflow implementation.
	} else {
		txn.Status = constants.FAILED // set correct validation state in implementation using structs object for test workflows and logic that validates types/states in programs via tests using go validation primitives.

	}

	b.Transactions = append(b.Transactions, txn) // setting now properties after validation and types is matched correctly to object's property being updated/modified correctly using methods of structs.
}
func (b *Block) MineBlock() {
	nonce := 0                                                      // define variable
	desiredHash := strings.Repeat("0", constants.MINING_DIFFICULTY) // set from parameter from test validation.

	for {
		guessHash := b.Hash()                                                          // call proper method as it was designed on system validation workflows
		ourSolutionHash := hex.EncodeToString(guessHash[:constants.MINING_DIFFICULTY]) // set parameter based on hash byte result.

		if ourSolutionHash == desiredHash { // make state implementation changes only on success
			return // finish if hash was created based on project requirements from structs, also matching parameter with types for those objects and their related implementation when software must run/test/ validate for certain behaviors as intended by the design with a properly implemented workflow.
		}
		nonce++ // set type variable changes for nonce as they iterate during validations for loop that implement this validation, since objects where methods will also have validation implementation with parameters to test it is implementation in local object context from this `block.go` implementations code to be fully validated in method where this code also runs under `go test` environment.

		b.Nonce = nonce // after variable updates the methods local types of that `struct`, must also validate all objects methods, implementation requirements where all code workflow is performing, by setting now, that property/attribute of a struct using methods signatures type parameters being passed as objects when it performs calculations via object variables/properties to guarantee types are used properly when structs, or interfaces exist and where they are supposed to have a type of object or a property or type with members or interfaces implementations in a method signatures with data properties/types during software operation, to prevent runtime code/workflow crashes due to miss implementation of proper methods/variables parameters implementations from methods with implementation signature of struct or objects during development of go project using their implementation/types during software implementation workflows of project when method signatures require a concrete types by a object and method implementations.
	}
}
