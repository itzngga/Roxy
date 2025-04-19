package context

import (
	"fmt"
	"github.com/itzngga/Roxy/util"
	waTypes "go.mau.fi/whatsmeow/types"
)

// FindGroupByJid find group by jid from cache
func (context *Ctx) FindGroupByJid(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error) {
	return context.Methods().FindGroupByJid(groupJid)
}

// GetAllGroups get all groups from cache
func (context *Ctx) GetAllGroups() (group []*waTypes.GroupInfo, err error) {
	return context.Methods().GetAllGroups()
}

// IsGroupAdmin check is target jid are group admin
func (context *Ctx) IsGroupAdmin(jid any) (bool, error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return false, err
	}

	group, err := context.FindGroupByJid(context.MessageInfo().Chat)
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
func (context *Ctx) IsClientGroupAdmin() (bool, error) {
	if context.ChatJID().Server != waTypes.GroupServer {
		return false, fmt.Errorf("error: chat is not a group")
	}

	group, err := context.FindGroupByJid(context.ChatJID())
	if err != nil {
		return false, err
	}

	var isAdmin bool
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == context.ClientJID().ToNonAD() {
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
func (context *Ctx) SetGroupName(name string) error {
	if context.ChatJID().Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := context.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = context.Client().SetGroupName(context.ChatJID(), name)
	return err
}

// SetGroupPhoto set current group photo
func (context *Ctx) SetGroupPhoto(data []byte) error {
	if context.ChatJID().Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := context.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	_, err = context.Client().SetGroupPhoto(context.ChatJID(), data)
	return err
}

// SetGroupAnnounce set group announce info
func (context *Ctx) SetGroupAnnounce(announce bool) error {
	if context.ChatJID().Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := context.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = context.Client().SetGroupAnnounce(context.ChatJID(), announce)
	return err
}

// SetGroupLocked set current group to locked
func (context *Ctx) SetGroupLocked(locked bool) error {
	if context.ChatJID().Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	ok, err := context.IsClientGroupAdmin()
	if !ok || err != nil {
		return fmt.Errorf("error: client is not a group admin")
	}

	err = context.Client().SetGroupLocked(context.ChatJID(), locked)
	return err
}
