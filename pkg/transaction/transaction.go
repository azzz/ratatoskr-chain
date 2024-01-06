package transaction

import (
	"fmt"
)

const (
	subsidy = 5
)

type Transaction struct {
	ID []byte
	// Inputs
	Vin []TxInput
	// Outputs
	Vout []TxOutput
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1
}

func NewCoinbaseTx(to, data string) Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{TxID: []byte{}, Vout: -1, ScriptSig: data}
	txout := TxOutput{Value: subsidy, ScriptPubKey: to}

	tx := Transaction{ID: nil, Vin: []TxInput{txin}, Vout: []TxOutput{txout}}

	return tx
}
