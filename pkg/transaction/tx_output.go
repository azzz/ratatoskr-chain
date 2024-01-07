package transaction

type TxOutput struct {
	Value        uint
	ScriptPubKey string
}

func (tx TxOutput) CanUnclockWith(data string) bool {
	return tx.ScriptPubKey == data
}
