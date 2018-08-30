package main

import (
	"github.com/boltdb/bolt"
	"encoding/hex"
	"os"
	"fmt"
	"log"
	"bytes"
	"errors"
	"crypto/ecdsa"
)


const dbFile = "./blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "Sad world"


type BlockChain struct {
	tip []byte
	db *bolt.DB
}


func LoadBlockChain(address string) *BlockChain {
	if dbExist() == false {
		fmt.Println("Block is not exist. Create One Firstly")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc


}

func CreateBlockChain(address string) *BlockChain {
	if dbExist() {
		fmt.Println("Block already exist")
		os.Exit(1)
	}


	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err:= tx.CreateBucket([]byte(blocksBucket))

		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serializer())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)

		if err != nil{
			log.Panic(err)
		}


		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func dbExist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (bc *BlockChain) AddBlock(transactions []*Transaction) {

	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serializer())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}


func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _,tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash){
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}


func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions{
			txID := hex.EncodeToString(tx.ID)

			Outputs:
				for outIdx, out := range tx.Vout{

					// 这个out 是否是某个交易的 in
					if spentTXOs[txID] != nil {
						for _, spentOut := range spentTXOs[txID] {
							if spentOut == outIdx{
								continue Outputs
							}
						}
					}


					if out.IsLockedWithKey(pubKeyHash) {
						unspentTXs = append(unspentTXs, *tx)
					}
				}

				if tx.IsCoinBase() == false {
					for _, in := range tx.Vin {
						if in.UsesKey(pubKeyHash) {
							inTxID := hex.EncodeToString(in.Txid)
							spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
						}
					}
				}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}


func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int)(int, map[string][]int){

	unspentOutPuts := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)

	accumulated := 0

Work:
	for _, tx := range unspentTXs{
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout{
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutPuts[txID] = append(unspentOutPuts[txID], outIdx)
			}

			if accumulated >= amount {
				break Work
			}
		}
	}

	return accumulated, unspentOutPuts
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found .")
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin{
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTxs)
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, vin := range tx.Vin{
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTxs)
}