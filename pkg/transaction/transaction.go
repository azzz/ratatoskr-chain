package transaction

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin) == 1 && tx.Vin[0].TxID == nil
}

func NewCoinbaseTx(to, data string) (Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TxInput{TxID: []byte{}, Vout: -1, ScriptSig: data}
	txout := TxOutput{Value: subsidy, ScriptPubKey: to}

	return newTransaction([]TxInput{txin}, []TxOutput{txout})
}

// Output is an extended TxOutput used in some business logic
type Output struct {
	TxOutput
	TxID []byte
	Vout int
}

func NewUTXOTransaction(sender, receiver string, amount uint64, availableOutputs []Output) (Transaction, error) {
	var (
		balance uint64
		inputs  []TxInput
		outputs []TxOutput
	)

	for _, out := range availableOutputs {
		balance += out.Value
		inputs = append(inputs, TxInput{out.TxID, out.Vout, sender})
	}

	outputs = append(outputs, TxOutput{
		Value:        amount,
		ScriptPubKey: receiver,
	})

	if balance < amount {
		return Transaction{}, errors.New("insufficient funds")
	}

	if amount < balance {
		// return change to the sender
		outputs = append(outputs, TxOutput{balance - amount, sender})
	}

	return newTransaction(inputs, outputs)
}

func newTransaction(Vin []TxInput, Vout []TxOutput) (Transaction, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Transaction{}, fmt.Errorf("generate id: %w", err)
	}

	return Transaction{id[:], Vin, Vout}, nil
}
