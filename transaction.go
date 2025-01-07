package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"KNIRVCHAIN-DEV/constants"
)

type Transaction struct {
	From            string       `json:"from"`
	To              string       `json:"to"`
	Value           uint64       `json:"value"`
	Data            []byte       `json:"data"`
	Status          string       `json:"status"`
	Timestamp       int64        `json:"timestamp"`
	TransactionHash string       `json:"transaction_hash"`
	PublicKey       string       `json:"public_key,omitempty"`
	Signature       []byte       `json:"Signature"`
	DeepCopy        *Transaction `json:"deep_copy,omitempty"`
}

func NewTransaction(from, to string, value uint64, data []byte) *Transaction {
	transaction := new(Transaction)

	transaction.From = from
	transaction.To = to
	transaction.Value = value
	transaction.Timestamp = time.Now().Unix()
	transaction.Data = data
	transaction.Status = constants.PENDING // for now by default
	transaction.TransactionHash = generateHash([]byte{}, transaction.Timestamp, fmt.Sprintf("%s%s%d", from, to, value), 0)

	return transaction

}

func (t Transaction) ToJson() string {
	nb, err := json.Marshal(t)

	if err != nil {
		return err.Error()
	} else {
		return string(nb)
	}
}

func (t Transaction) VerifyTxn() bool {
	if t.Value <= 0 {
		return false
	}

	if t.Value > math.MaxUint32 {
		return false
	}

	if t.From == t.To {
		return false
	}

	valid := t.VerifySignature()

	return valid

}

func (t *Transaction) VerifySignature() bool {

	if t.Signature == nil {
		return false
	}

	if t.PublicKey == "" {
		return false
	}

	signature := t.Signature
	publicKeyHex := t.PublicKey
	t.Signature = []byte{}
	t.PublicKey = ""
	publicKeyEcdsa := GetPublicKeyFromHex(publicKeyHex)

	bs, _ := json.Marshal(t)
	hash := sha256.Sum256(bs)

	valid := ecdsa.VerifyASN1(publicKeyEcdsa, hash[:], signature)
	t.Signature = signature
	return valid
}

func (t Transaction) Hash() string {

	bs, _ := json.Marshal(t)
	sum := sha256.Sum256(bs)
	hexRep := hex.EncodeToString(sum[:32])
	formattedHexRep := constants.HEX_PREFIX + hexRep

	return formattedHexRep
}

func GetPublicKeyFromHex(publicKeyHex string) *ecdsa.PublicKey {
	rpk := publicKeyHex[2:]
	xHex := rpk[:64]
	yHex := rpk[64:]
	x := new(big.Int)
	y := new(big.Int)
	x.SetString(xHex, 16)
	y.SetString(yHex, 16)

	var npk ecdsa.PublicKey
	npk.Curve = elliptic.P256()
	npk.X = x
	npk.Y = y

	return &npk
}

func (tx *Transaction) DeepCopying() *Transaction {

	copiedSignature := make([]byte, len(tx.Signature))
	copy(copiedSignature, tx.Signature)

	copiedPublicKey := tx.PublicKey

	copiedData := make([]byte, len(tx.Data))
	copy(copiedData, tx.Data)

	return &Transaction{ // Create a completely new Transaction instance
		From:            tx.From,
		To:              tx.To,
		Value:           tx.Value,
		Data:            copiedData, // Copy the byte slice
		Status:          tx.Status,
		Timestamp:       tx.Timestamp,
		TransactionHash: tx.TransactionHash,

		PublicKey: copiedPublicKey,
		Signature: copiedSignature,
	}
}
