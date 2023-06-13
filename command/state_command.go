package command

import (
	"strconv"
)

func (runFunc *RunFuncContext) SetUserState(stateName string) {
	number := strconv.FormatUint(runFunc.MessageInfo.Sender.UserInt(), 10)

	runFunc.UserStateChan <- []string{stateName, number}
}
