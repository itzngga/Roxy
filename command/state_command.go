package command

func (runFunc *RunFuncContext) SetUserState(stateName string, data map[string]interface{}) {
	runFunc.UserStateChan <- []interface{}{stateName, runFunc.MessageInfo.Sender.ToNonAD().String(), data}
}
