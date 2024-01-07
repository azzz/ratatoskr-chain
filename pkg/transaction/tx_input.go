package transaction

// TxInput is a mapping on a specific TxOutput of a specific Transaction.
// Each transaction might have multiple Outputs indexed relatively to the Transaction.
// Vout is the index of a TxOutput, and TxID is the Transaction.ID
type TxInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

func (tx TxInput) CanUnlockWith(data string) bool {
	return tx.ScriptSig == data
}
