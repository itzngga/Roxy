package command

import (
	"errors"
	"fmt"
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"strings"
)

func (runFunc *RunFuncContext) SetUserStatus(status string) error {
	if status == "" {
		return errors.New("error: failed to blank status string")
	}

	err := runFunc.Client.SetStatusMessage(status)
	if err != nil {
		return fmt.Errorf("error: failed to change status : %v\n", err)
	}

	return nil
}

func (runFunc *RunFuncContext) JoinInviteLink(link string) error {
	if link == "" {
		return errors.New("error: blank link string")
	}
	// formatting group link
	link = strings.Replace(link, "https://chat.whatsapp.com/", "", -1)

	groupId, err := runFunc.Client.JoinGroupWithLink(link)
	if err != nil {
		return fmt.Errorf("error: failed to join group with invite link : %v\n", err)
	}

	fmt.Printf("success: joined to group %s\n", groupId)
	return nil
}

func (runFunc *RunFuncContext) GetGroupInfoFromInviteLink(link string) (*waTypes.GroupInfo, error) {
	if link == "" {
		return nil, errors.New("error: blank link string")
	}
	// formatting group link
	link = strings.Replace(link, "https://chat.whatsapp.com/", "", -1)

	groupInfo, err := runFunc.Client.GetGroupInfoFromLink(link)
	if err != nil {
		return nil, fmt.Errorf("error: failed to get group info with invite link: %v\n", err)
	}

	return groupInfo, nil
}

func (runFunc *RunFuncContext) GetGroupInviteLink(jid any, reset bool) (string, error) {
	jids, err := util.ParseGroupJid(jid)
	if err != nil {
		return "", err
	}

	link, err := runFunc.Client.GetGroupInviteLink(jids, reset)
	if err != nil {
		return "", fmt.Errorf("error: failed to get group invite link : %v\n", err)
	}

	return "https://chat.whatsapp.com/" + link, nil
}

func (runFunc *RunFuncContext) GetJoinedGroups() ([]*waTypes.GroupInfo, error) {
	groups, err := runFunc.Client.GetJoinedGroups()
	if err != nil {
		return nil, fmt.Errorf("error: failed to get joined group : %v\n", err)
	}
	return groups, nil
}

func (runFunc *RunFuncContext) GetGroupInfo(jid any) (*waTypes.GroupInfo, error) {
	jids, err := util.ParseGroupJid(jid)
	if err != nil {
		return nil, err
	}

	group, err := FindGroupByJid(runFunc.Client, jids)
	if err != nil {
		return nil, fmt.Errorf("error: failed to get group info : %v\n", err)
	}

	return group, nil
}

func (runFunc *RunFuncContext) GetProfilePicture(jid any) (string, error) {
	jids, err := util.ParseAllJid(jid)
	if err != nil {
		return "", err
	}

	pic, err := runFunc.Client.GetProfilePictureInfo(jids, &whatsmeow.GetProfilePictureParams{})
	if err != nil {
		return "", fmt.Errorf("error: failed to get group info : %v\n", err)
	}
	return pic.URL, nil
}

func (runFunc *RunFuncContext) UpdateClientProfilePicture(data []byte) error {
	_, err := runFunc.Client.SetGroupPhoto(waTypes.JID{}, data)
	if err != nil {
		return err
	}

	return nil
}

func (runFunc *RunFuncContext) UpdateProfilePicture(jid any, data []byte) error {
	val, ok := jid.(string)
	if ok && val == "" {
		_, err := runFunc.Client.SetGroupPhoto(waTypes.JID{}, data)
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
		group, err := FindGroupByJid(runFunc.Client, jids)
		if err != nil {
			return err
		}

		if group == nil {
			return fmt.Errorf("error: group id not valid : %v", jids.ToNonAD().String())
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

		if !isAdmin {
			return fmt.Errorf("error: client is not admin : %v", runFunc.Client.Store.ID.ToNonAD().String())
		}

		_, err = runFunc.Client.SetGroupPhoto(jids, data)
		if err != nil {
			return err
		}

		return nil
	} else {
		return fmt.Errorf("error: cannot set other user profile")
	}
}

func (runFunc *RunFuncContext) GetUser(jid any) (result waTypes.UserInfo, err error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return result, err
	}

	user, err := runFunc.Client.GetUserInfo([]waTypes.JID{jids})
	if err != nil {
		if err != nil {
			return result, fmt.Errorf("error: failed to get user info : %v\n", err)
		}
	}

	if val, ok := user[jids]; ok {
		return val, nil
	} else {
		return waTypes.UserInfo{}, nil
	}
}
