package command

import (
	"fmt"
	"github.com/zhangyunhao116/skipmap"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var GroupCache *skipmap.StringMap[[]*waTypes.GroupInfo]

func init() {
	GroupCache = skipmap.NewString[[]*waTypes.GroupInfo]()
}

func FindGroupByJid(c *whatsmeow.Client, groupJid waTypes.JID) (group *waTypes.GroupInfo, err error) {
	clientJID := c.Store.ID.ToNonAD()
	groups, ok := GroupCache.Load(clientJID.String())
	if !ok {
		groups, err = c.GetJoinedGroups()
		if err != nil {
			return nil, err
		}
		GroupCache.Store(clientJID.String(), groups)
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

func GetAllGroups(c *whatsmeow.Client) (group []*waTypes.GroupInfo, err error) {
	clientJID := c.Store.ID.ToNonAD()
	groups, ok := GroupCache.Load(clientJID.String())
	if !ok {
		groups, err = c.GetJoinedGroups()
		if err != nil {
			return nil, err
		}
		GroupCache.Store(clientJID.String(), groups)
	}
	return groups, nil
}

func CacheAllGroup(c *whatsmeow.Client) {
	clientJID := c.Store.ID.ToNonAD()
	groups, err := c.GetJoinedGroups()
	if err != nil {
		return
	}
	GroupCache.Store(clientJID.String(), groups)
}

func UNCacheOneGroup(c *whatsmeow.Client, info *events.GroupInfo, joined *events.JoinedGroup) {
	var err error
	clientJID := c.Store.ID.ToNonAD()
	if info != nil {
		groups, ok := GroupCache.Load(clientJID.String())
		if !ok {
			groups, err = c.GetJoinedGroups()
			if err != nil {
				return
			}
			GroupCache.Store(clientJID.String(), groups)
		} else {
			for i, group := range groups {
				if group.JID.ToNonAD() == info.JID.ToNonAD() {
					info, _ := c.GetGroupInfo(group.JID.ToNonAD())
					groups[i] = info
					break
				}
			}
		}
	}
	if joined != nil {
		groups, ok := GroupCache.Load(clientJID.String())
		if !ok {
			groups, err = c.GetJoinedGroups()
			if err != nil {
				return
			}
			GroupCache.Store(clientJID.String(), groups)
		} else {
			for i, group := range groups {
				if group.JID.ToNonAD() == joined.JID.ToNonAD() {
					info, _ := c.GetGroupInfo(group.JID.ToNonAD())
					groups[i] = info
					break
				}
			}
		}
	}
}
