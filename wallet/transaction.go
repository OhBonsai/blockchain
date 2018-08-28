
package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"strings"
	"crypto/elliptic"
)

const subsidy = 10

type Transaction struct {
	ID    []byte
	Vin   []TXInput
	Vout  []TXOutput
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) Serialize()[]byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash()[]byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}


func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTXs := prevTXs[hex.EncodeToString(vin.Txid)]

		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTXs.Vout[vin.Vout].PubKeyHash

		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s , err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("---- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}


func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx * Transaction) Verify(prevTXs map[string]Transaction) bool {

	return false
}

func NewCoinBaseTransaction(to, data string) *Transaction{
	if data == "" {
		data = fmt.Sprintf("Sealed with a kiss")
	}

	txin := TXInput{
		[]byte{},
		-1,
	data,
	}
	txout := TXOutput{
		subsidy,
		to,
	}

	t := Transaction{
		nil,
		[]TXInput{txin},
		[]TXOutput{txout},
	}
	t.SetID()
	return &t
}


func (t *Transaction)SetID() {
	var encoded bytes.Buffer
	var hash [32]byte


	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(t)

	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	t.ID = hash[:]
}




func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction{
	var inputs []TXInput
	var outputs []TXOutput


	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: Not engough funds")
	}

	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount{
		outputs = append(outputs, TXOutput{acc-amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
