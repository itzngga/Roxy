package util

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"strings"
)

func SetUserStatus(c *whatsmeow.Client, status string) {
	if status == "" {
		fmt.Println("error: failed to blank status string")
		return
	}

	err := c.SetStatusMessage(status)
	if err != nil {
		fmt.Printf("error: failed to change status : %v\n", err)
	}
	return
}

func JoinInviteLink(c *whatsmeow.Client, link string) {
	if link == "" {
		fmt.Println("error: blank link string")
		return
	}
	// formatting group link
	link = strings.Replace(link, "https://chat.whatsapp.com/", "", -1)

	groupId, err := c.JoinGroupWithLink(link)
	if err != nil {
		fmt.Printf("error: failed to join group with invite link : %v\n", err)
		return
	}

	fmt.Printf("success: joined to group %s\n", groupId)
	return
}

func GetGroupInfoFromInviteLink(c *whatsmeow.Client, link string) *waTypes.GroupInfo {
	if link == "" {
		fmt.Println("error: blank link string")
		return nil
	}
	// formatting group link
	link = strings.Replace(link, "https://chat.whatsapp.com/", "", -1)

	groupInfo, err := c.GetGroupInfoFromLink(link)
	if err != nil {
		fmt.Printf("error: failed to get group info with invite link : %v\n", err)
		return nil
	}

	return groupInfo
}

func GetGroupInviteLink(c *whatsmeow.Client, jid any, reset bool) string {
	jids := ParseGroupJid(jid)

	link, err := c.GetGroupInviteLink(jids, reset)
	if err != nil {
		fmt.Printf("error: failed to get group invite link : %v\n", err)
		return ""
	}

	return "https://chat.whatsapp.com/" + link
}

func GetJoinedGroups(c *whatsmeow.Client) []*waTypes.GroupInfo {
	groups, err := c.GetJoinedGroups()
	if err != nil {
		fmt.Printf("error: failed to get joined group : %v\n", err)
		return nil
	}
	return groups
}

func GetGroupInfo(c *whatsmeow.Client, jid any) *waTypes.GroupInfo {
	jids := ParseGroupJid(jid)

	group, err := c.GetGroupInfo(jids)
	if err != nil {
		fmt.Printf("error: failed to get group info : %v\n", err)
		return nil
	}

	return group
}

func GetProfilePicture(c *whatsmeow.Client, jid any) string {
	jids := ParseAllJid(jid)

	pic, err := c.GetProfilePictureInfo(jids, &whatsmeow.GetProfilePictureParams{})
	if err != nil {
		fmt.Printf("error: failed to get profile link : %v\n", err)
		return ""
	}
	return pic.URL
}

func GetUser(c *whatsmeow.Client, jid any) *waTypes.UserInfo {
	jids := ParseUserJid(jid)

	user, err := c.GetUserInfo([]waTypes.JID{jids})
	if err != nil {
		if err != nil {
			fmt.Printf("error: failed to get user info : %v\n", err)
			return nil
		}
	}

	res := user[jids]
	return &res
}
