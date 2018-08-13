package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

const subsidy = 10

type Transaction struct {
	ID    []byte
	Vin   []TXInput
	Vout  []TXOutput
}


type TXOutput struct {
	Value int
	ScriptPubKey string
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
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

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (in *TXInput) CanUnlockOuputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlocked(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
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
