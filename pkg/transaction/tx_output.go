package transaction

type TxOutput struct {
	Value        uint64
	ScriptPubKey string
}

func (tx TxOutput) CanUnclockWith(data string) bool {
	return tx.ScriptPubKey == data
}
