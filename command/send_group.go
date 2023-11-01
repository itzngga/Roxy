package command

import (
	"fmt"
	"github.com/itzngga/Roxy/types"
	"github.com/itzngga/Roxy/util"
	waTypes "go.mau.fi/whatsmeow/types"
)

// FindGroupByJid find group by jid from cache
func (runFunc *RunFuncContext) FindGroupByJid(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error) {
	FindGroupByJid := types.GetContext[types.FindGroupByJid](runFunc.Ctx, "FindGroupByJid")
	return FindGroupByJid(groupJid)
}

// GetAllGroups get all groups from cache
func (runFunc *RunFuncContext) GetAllGroups() (group []*waTypes.GroupInfo, err error) {
	GetAllGroups := types.GetContext[types.GetAllGroups](runFunc.Ctx, "GetAllGroups")
	return GetAllGroups()
}

// IsGroupAdmin check is target jid are group admin
func (runFunc *RunFuncContext) IsGroupAdmin(jid any) (bool, error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return false, err
	}

	group, err := runFunc.FindGroupByJid(runFunc.MessageInfo.Chat)
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

// IsClientGroupAdmin check if client is a group admin
func (runFunc *RunFuncContext) IsClientGroupAdmin() (bool, error) {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return false, fmt.Errorf("error: chat is not a group")
	}

	group, err := runFunc.FindGroupByJid(runFunc.MessageInfo.Chat)
	if err != nil {
		return false, err
	}

	var isAdmin bool
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == runFunc.ClientJID.ToNonAD() {
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

// SetGroupName set current group name
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

// SetGroupPhoto set current group photo
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

// SetGroupAnnounce set group announce info
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

// SetGroupLocked set current group to locked
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
