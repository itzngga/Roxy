package command

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// WARN: UNSAFE COMMANDS! use it wisely

// RemoveMemberFromGroup kick members from group
func (runFunc *RunFuncContext) RemoveMemberFromGroup(jids []waTypes.JID) error {
	if ok, err := runFunc.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeRemove
	}

	_, err := runFunc.Client.UpdateGroupParticipants(runFunc.MessageInfo.Chat, changes)
	return err
}

// AddMemberToGroup add members from group
func (runFunc *RunFuncContext) AddMemberToGroup(jids []waTypes.JID) error {
	if ok, err := runFunc.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeAdd
	}

	_, err := runFunc.Client.UpdateGroupParticipants(runFunc.MessageInfo.Chat, changes)
	return err
}

// PromoteMemberInGroup kick a member from group
func (runFunc *RunFuncContext) PromoteMemberInGroup(jids []waTypes.JID) error {
	if ok, err := runFunc.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangePromote
	}

	_, err := runFunc.Client.UpdateGroupParticipants(runFunc.MessageInfo.Chat, changes)
	return err
}

// DemoteMemberInGroup kick a member from group
func (runFunc *RunFuncContext) DemoteMemberInGroup(jids []waTypes.JID) error {
	if ok, err := runFunc.IsClientGroupAdmin(); !ok || err != nil {
		return err
	}

	changes := map[waTypes.JID]whatsmeow.ParticipantChange{}
	for _, jid := range jids {
		changes[jid] = whatsmeow.ParticipantChangeDemote
	}

	_, err := runFunc.Client.UpdateGroupParticipants(runFunc.MessageInfo.Chat, changes)
	return err
}

// RevokeGroupInvite revoke current group invite link
func (runFunc *RunFuncContext) RevokeGroupInvite() (string, error) {
	if ok, err := runFunc.IsClientGroupAdmin(); !ok || err != nil {
		return "", err
	}

	return runFunc.Client.GetGroupInviteLink(runFunc.MessageInfo.Chat, true)
}

// LeaveFromGroup leave group from given jid
func (runFunc *RunFuncContext) LeaveFromGroup(jid waTypes.JID) error {
	if jid.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	err := runFunc.Client.LeaveGroup(jid)
	return err
}

// LeaveFromThisGroup leave from this group chat
func (runFunc *RunFuncContext) LeaveFromThisGroup() error {
	if runFunc.MessageInfo.Chat.Server != waTypes.GroupServer {
		return fmt.Errorf("error: chat is not a group")
	}

	err := runFunc.Client.LeaveGroup(runFunc.MessageInfo.Chat)
	return err
}
