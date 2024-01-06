package transaction

type TxInput struct {
	TxID      []byte
	Vout      int
	ScriptSig string
}

func (tx *TxInput) CanUnlockOutpuWith(data string) bool {
	return tx.ScriptSig == data
}
