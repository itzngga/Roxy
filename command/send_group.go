package command

import (
	"fmt"
	"github.com/itzngga/Roxy/util"
	waTypes "go.mau.fi/whatsmeow/types"
)

func (runFunc *RunFuncContext) IsGroupAdmin(jid any) (bool, error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return false, err
	}

	group, err := FindGroupByJid(runFunc.Client, runFunc.MessageInfo.Chat)
	if err != nil {
		return false, err
	}

	var isAdmin bool
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == jids {
			if participant.IsSuperAdmin {
				isAdmin = true
				break
			}
			if participant.IsAdmin {
				isAdmin = true
				break
			}
		}
	}

	return isAdmin, nil
}

func (runFunc *RunFuncContext) IsClientGroupAdmin() (bool, error) {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return false, fmt.Errorf("error: chat is not a group")
	}

	group, err := FindGroupByJid(runFunc.Client, runFunc.MessageInfo.Chat)
	if err != nil {
		return false, err
	}

	var isAdmin bool
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == runFunc.Client.Store.ID.ToNonAD() {
			if participant.IsSuperAdmin {
				isAdmin = true
				break
			}
			if participant.IsAdmin {
				isAdmin = true
				break
			}
		}
	}

	return isAdmin, nil
}

func (runFunc *RunFuncContext) SetGroupName(name string) error {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := runFunc.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = runFunc.Client.SetGroupName(runFunc.MessageInfo.Chat, name)
	return err
}

func (runFunc *RunFuncContext) SetGroupPhoto(data []byte) error {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := runFunc.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	_, err = runFunc.Client.SetGroupPhoto(runFunc.MessageInfo.Chat, data)
	return err
}

func (runFunc *RunFuncContext) SetGroupAnnounce(announce bool) error {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := runFunc.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = runFunc.Client.SetGroupAnnounce(runFunc.MessageInfo.Chat, announce)
	return err
}

func (runFunc *RunFuncContext) SetGroupLocked(locked bool) error {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := runFunc.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = runFunc.Client.SetGroupLocked(runFunc.MessageInfo.Chat, locked)
	return err
}
