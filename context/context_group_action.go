package context

import (
	"fmt"
	"github.com/go-whatsapp/whatsmeow"
	waTypes "github.com/go-whatsapp/whatsmeow/types"
)

// WARN: UNSAFE COMMANDS! use it wisely

// RemoveMemberFromGroup kick members from group
func (context *Ctx) RemoveMemberFromGroup(jids []waTypes.JID) error {
	if ok, err := context.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeRemove
	}

	_, err := context.Client().UpdateGroupParticipants(context.ChatJID(), changes)
	return err
}

// AddMemberToGroup add members from group
func (context *Ctx) AddMemberToGroup(jids []waTypes.JID) error {
	if ok, err := context.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeAdd
	}

	_, err := context.Client().UpdateGroupParticipants(context.ChatJID(), changes)
	return err
}

// PromoteMemberInGroup kick a member from group
func (context *Ctx) PromoteMemberInGroup(jids []waTypes.JID) error {
	if ok, err := context.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangePromote
	}

	_, err := context.Client().UpdateGroupParticipants(context.ChatJID(), changes)
	return err
}

// DemoteMemberInGroup kick a member from group
func (context *Ctx) DemoteMemberInGroup(jids []waTypes.JID) error {
	if ok, err := context.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeDemote
	}

	_, err := context.Client().UpdateGroupParticipants(context.ChatJID(), changes)
	return err
}

// RevokeGroupInvite revoke current group invite link
func (context *Ctx) RevokeGroupInvite() (string, error) {
	if ok, err := context.IsClientGroupAdmin(); !ok || err != nil {
		return "", err
	}

	return context.Client().GetGroupInviteLink(context.ChatJID(), true)
}

// LeaveFromGroup leave group from given jid
func (context *Ctx) LeaveFromGroup(jid waTypes.JID) error {
	if jid.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	err := context.Client().LeaveGroup(jid)
	return err
}

// LeaveFromThisGroup leave from this group chat
func (context *Ctx) LeaveFromThisGroup() error {
	if context.ChatJID().Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	err := context.Client().LeaveGroup(context.ChatJID())
	return err
}
