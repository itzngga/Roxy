package context

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// SetUserStatus set client status
func (context *Ctx) SetUserStatus(status string) error {
	if status == "" {
		return errors.New("error: failed to blank status string")
	}

	err := context.client.SetStatusMessage(status)
	if err != nil {
		return fmt.Errorf("error: failed to change status : %v", err)
	}

	return nil
}

// JoinGroupInviteLink join a group invite link
func (context *Ctx) JoinGroupInviteLink(link string) error {
	if link == "" {
		return errors.New("error: blank link string")
	}
	// formatting group link
	link = strings.ReplaceAll(link, "https://chat.whatsapp.com/", "")

	groupId, err := context.client.JoinGroupWithLink(link)
	if err != nil {
		return fmt.Errorf("error: failed to join group with invite link : %v", err)
	}

	fmt.Printf("success: joined to group %s\n", groupId)
	return nil
}

// GetGroupInfoFromInviteLink get group info from invite link
func (context *Ctx) GetGroupInfoFromInviteLink(link string) (*waTypes.GroupInfo, error) {
	if link == "" {
		return nil, errors.New("error: blank link string")
	}
	// formatting group link
	link = strings.ReplaceAll(link, "https://chat.whatsapp.com/", "")

	groupInfo, err := context.client.GetGroupInfoFromLink(link)
	if err != nil {
		return nil, fmt.Errorf("error: failed to get group info with invite link: %v", err)
	}

	return groupInfo, nil
}

// GetGroupInviteLink get current group invite link
func (context *Ctx) GetGroupInviteLink(jid any) (string, error) {
	jids, err := util.ParseGroupJid(jid)
	if err != nil {
		return "", err
	}

	link, err := context.client.GetGroupInviteLink(jids, false)
	if err != nil {
		return "", fmt.Errorf("error: failed to get group invite link : %v", err)
	}

	return "https://chat.whatsapp.com/" + link, nil
}

// GetJoinedGroups get client joined group from cache
func (context *Ctx) GetJoinedGroups() ([]*waTypes.GroupInfo, error) {
	groups, err := context.client.GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("error: failed to get joined group : %v", err)
	}
	return groups, nil
}

// GetGroupInfo get group info from cache
func (context *Ctx) GetGroupInfo(jid any) (*waTypes.GroupInfo, error) {
	jids, err := util.ParseGroupJid(jid)
	if err != nil {
		return nil, err
	}

	group, err := context.Methods().FindGroupByJid(jids)
	if err != nil {
		return nil, fmt.Errorf("error: failed to get group info : %v", err)
	}

	return group, nil
}

// GetGroupProfilePicture get group profile picture
func (context *Ctx) GetGroupProfilePicture(jid any) (string, error) {
	jids, err := util.ParseAllJid(jid)
	if err != nil {
		return "", err
	}

	pic, err := context.client.GetProfilePictureInfo(jids, &whatsmeow.GetProfilePictureParams{})
	if err != nil {
		return "", fmt.Errorf("error: failed to get group info : %v", err)
	}
	return pic.URL, nil
}

// UpdateClientProfilePicture update client profile picture
func (context *Ctx) UpdateClientProfilePicture(data []byte) error {
	_, err := context.client.SetGroupPhoto(waTypes.JID{}, data)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGroupProfilePicture update group profile picture
func (context *Ctx) UpdateGroupProfilePicture(jid any, data []byte) error {
	val, ok := jid.(string)
	if ok && val == "" {
		_, err := context.client.SetGroupPhoto(waTypes.JID{}, data)
		if err != nil {
			return err
		}

		return nil
	}

	jids, err := util.ParseAllJid(jid)
	if err != nil {
		return err
	}

	if jids.Server == waTypes.GroupServer {
		group, err := context.Methods().FindGroupByJid(jids)
		if err != nil {
			return err
		}

		if group == nil {
			return fmt.Errorf("error: group id not valid : %v", jids.ToNonAD().String())
		}

		var isAdmin bool
		for _, participant := range group.Participants {
			if participant.JID.ToNonAD() == context.clientJid.ToNonAD() {
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

		if !isAdmin {
			return fmt.Errorf("error: client is not admin : %v", context.Methods().ClientJID().ToNonAD().String())
		}

		_, err = context.Client().SetGroupPhoto(jids, data)
		if err != nil {
			return err
		}

		return nil
	} else {
		return fmt.Errorf("error: cannot set other user profile")
	}
}

// GetUserInfo get contact user info from cache
func (context *Ctx) GetUserInfo(jid any) (result waTypes.UserInfo, err error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return result, err
	}

	user, err := context.Client().GetUserInfo([]waTypes.JID{jids})
	if err != nil {
		return result, fmt.Errorf("error: failed to get user info : %v", err)
	}

	if val, ok := user[jids]; ok {
		return val, nil
	} else {
		return waTypes.UserInfo{}, nil
	}
}
