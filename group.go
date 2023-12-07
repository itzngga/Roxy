package roxy

import (
	"fmt"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
	"github.com/go-whatsapp/whatsmeow/types/events"
	"github.com/itzngga/Roxy/util"
)

func (muxer *Muxer) FindGroupByJid(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error) {
	groups, ok := muxer.GroupCache.Load(groupJid.ToNonAD().String())
	if !ok {
		var client = muxer.AppMethods.Client()
		groups, err = client.GetJoinedGroups()
		if err != nil {
			return nil, err
		}
		muxer.GroupCache.Store(groupJid.ToNonAD().String(), groups)
	}
	for _, groupz := range groups {
		if groupz.JID.ToNonAD() == groupJid {
			group = groupz
			break
		}
	}
	if group == nil {
		return nil, fmt.Errorf("error: invalid group jid : %v", groupJid.String())
	}
	return group, nil
}

func (muxer *Muxer) GetAllGroups() (group []*waTypes.GroupInfo, err error) {
	var client = muxer.AppMethods.Client()
	groups, ok := muxer.GroupCache.Load(muxer.AppMethods.ClientJID().ToNonAD().String())
	if !ok {
		groups, err = client.GetJoinedGroups()
		if err != nil {
			return nil, err
		}
		muxer.GroupCache.Store(muxer.AppMethods.ClientJID().ToNonAD().String(), groups)
	}
	return groups, nil
}

func (muxer *Muxer) CacheAllGroup() {
	var client = muxer.AppMethods.Client()

	groups, err := client.GetJoinedGroups()
	if err != nil {
		return
	}
	muxer.GroupCache.Store(muxer.AppMethods.ClientJID().ToNonAD().String(), groups)
}

func (muxer *Muxer) UnCacheOneGroup(info *events.GroupInfo, joined *events.JoinedGroup) {
	var err error
	var client = muxer.AppMethods.Client()
	clientJID := muxer.AppMethods.ClientJID().ToNonAD()
	if info != nil {
		groups, ok := muxer.GroupCache.Load(clientJID.String())
		if !ok {
			groups, err = client.GetJoinedGroups()
			if err != nil {
				return
			}
			muxer.GroupCache.Store(clientJID.String(), groups)
		} else {
			for i, group := range groups {
				if group.JID.ToNonAD() == info.JID.ToNonAD() {
					info, _ := client.GetGroupInfo(group.JID.ToNonAD())
					groups[i] = info
					break
				}
			}
		}
	}
	if joined != nil {
		groups, ok := muxer.GroupCache.Load(clientJID.String())
		if !ok {
			groups, err = client.GetJoinedGroups()
			if err != nil {
				return
			}
			muxer.GroupCache.Store(clientJID.String(), groups)
		} else {
			for i, group := range groups {
				if group.JID.ToNonAD() == joined.JID.ToNonAD() {
					info, _ := client.GetGroupInfo(group.JID.ToNonAD())
					groups[i] = info
					break
				}
			}
		}
	}
}
func (muxer *Muxer) IsGroupAdmin(chat waTypes.JID, jid any) (bool, error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return false, err
	}

	group, err := muxer.FindGroupByJid(chat)
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

func (muxer *Muxer) IsClientGroupAdmin(chat waTypes.JID) (bool, error) {
	if chat.Server != waTypes.GroupServer {
		return false, fmt.Errorf("error: chat is not a group")
	}

	group, err := muxer.FindGroupByJid(chat)
	if err != nil {
		return false, err
	}

	var isAdmin bool
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == muxer.AppMethods.ClientJID().ToNonAD() {
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
