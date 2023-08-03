package core

import (
	"fmt"
	"github.com/itzngga/Roxy/util"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (muxer *Muxer) findGroupByJid(groupJid waTypes.JID) (group *waTypes.GroupInfo, err error) {
	groups, ok := muxer.GroupCache.Load(groupJid.ToNonAD().String())
	if !ok {
		var client = muxer.getCurrentClient()
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

func (muxer *Muxer) getAllGroups() (group []*waTypes.GroupInfo, err error) {
	var client = muxer.getCurrentClient()
	clientJID := client.Store.ID.ToNonAD()
	groups, ok := muxer.GroupCache.Load(clientJID.String())
	if !ok {
		groups, err = client.GetJoinedGroups()
		if err != nil {
			return nil, err
		}
		muxer.GroupCache.Store(clientJID.String(), groups)
	}
	return groups, nil
}

func (muxer *Muxer) cacheAllGroup() {
	var client = muxer.getCurrentClient()

	clientJID := client.Store.ID.ToNonAD()
	groups, err := client.GetJoinedGroups()
	if err != nil {
		return
	}
	muxer.GroupCache.Store(clientJID.String(), groups)
}

func (muxer *Muxer) unCacheOneGroup(info *events.GroupInfo, joined *events.JoinedGroup) {
	var err error
	var client = muxer.getCurrentClient()
	clientJID := client.Store.ID.ToNonAD()
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
func (muxer *Muxer) isGroupAdmin(chat waTypes.JID, jid any) (bool, error) {
	jids, err := util.ParseUserJid(jid)
	if err != nil {
		return false, err
	}

	group, err := muxer.findGroupByJid(chat)
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

func (muxer *Muxer) isClientGroupAdmin(chat waTypes.JID) (bool, error) {
	if chat.Server != waTypes.GroupServer {
		return false, fmt.Errorf("error: chat is not a group")
	}

	group, err := muxer.findGroupByJid(chat)
	if err != nil {
		return false, err
	}

	var isAdmin bool
	var client = muxer.getCurrentClient()
	for _, participant := range group.Participants {
		if participant.JID.ToNonAD() == client.Store.ID.ToNonAD() {
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
