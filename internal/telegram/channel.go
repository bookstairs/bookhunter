package telegram

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/log"
)

func (t *Telegram) ChannelInfo() (*ChannelInfo, error) {
	var channelID int64
	var accessHash int64
	var err error

	// Get the real channelID and accessHash.
	if strings.HasPrefix(t.channelID, "joinchat/") {
		channelID, accessHash, err = t.privateChannelInfo(strings.TrimPrefix(t.channelID, "joinchat/"))
	} else {
		channelID, accessHash, err = t.publicChannelInfo(t.channelID)
	}
	if err != nil {
		return nil, err
	}

	// Query the last message ID.
	lastMsgID, err := t.queryLastMsgID(channelID, accessHash)
	if err != nil {
		return nil, err
	}

	return &ChannelInfo{
		ID:         channelID,
		AccessHash: accessHash,
		LastMsgID:  lastMsgID,
	}, nil
}

// privateChannelInfo queries access hash for the private channel.
func (t *Telegram) privateChannelInfo(hash string) (id int64, access int64, err error) {
	invite, err := t.client.API().MessagesCheckChatInvite(t.ctx, hash)
	if err != nil {
		return
	}

	switch v := invite.(type) {
	case *tg.ChatInviteAlready:
		if channel, ok := v.GetChat().(*tg.Channel); ok {
			id = channel.ID
			access = channel.AccessHash
			return
		}
	case *tg.ChatInvitePeek:
		if channel, ok := v.GetChat().(*tg.Channel); ok {
			id = channel.ID
			access = channel.AccessHash
			return
		}
	case *tg.ChatInvite:
		log.Warn("You haven't join this private channel, plz join it manually.")
	}

	err = errors.New("couldn't find access hash")
	return
}

// publicChannelInfo queries the public channel by its name.
func (t *Telegram) publicChannelInfo(name string) (id, access int64, err error) {
	username, err := t.client.API().ContactsResolveUsername(t.ctx, &tg.ContactsResolveUsernameRequest{Username: name})
	if err != nil {
		return
	}

	if len(username.Chats) == 0 {
		err = fmt.Errorf("you are not belong to channel: %s", name)
		return
	}

	for _, chat := range username.Chats {
		// Try to find the related channel.
		if channel, ok := chat.(*tg.Channel); ok {
			id = channel.ID
			access = channel.AccessHash
			return
		}
	}

	err = fmt.Errorf("couldn't find channel id and hash for channel: %s", name)
	return
}

// queryLastMsgID from the given channel info.
func (t *Telegram) queryLastMsgID(channelID, access int64) (int64, error) {
	request := &tg.MessagesSearchRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  channelID,
			AccessHash: access,
		},
		Filter:   &tg.InputMessagesFilterEmpty{},
		Q:        "",
		OffsetID: -1,
		Limit:    1,
	}

	last := -1
	search, err := t.client.API().MessagesSearch(t.ctx, request)
	if err != nil {
		return 0, err
	}

	channelInfo, ok := search.(*tg.MessagesChannelMessages)
	if !ok {
		return 0, err
	}

	for _, msg := range channelInfo.Messages {
		if msg != nil {
			last = msg.GetID()
			break
		}
	}

	if last <= 0 {
		return 0, errors.New("couldn't find last message id")
	}

	return int64(last), nil
}
