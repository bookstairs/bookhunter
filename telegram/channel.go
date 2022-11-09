package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/pkg/log"
)

// channelInfo will return the last available book id, channel id and access hash.
func channelInfo(ctx context.Context, client *telegram.Client, config *Config) (*tgChannelInfo, error) {
	var channelID int64
	var accessHash int64
	var err error
	if strings.HasPrefix(config.ChannelID, "joinchat/") {
		channelID, accessHash, err = privateChannelInfo(ctx, client, strings.TrimPrefix(config.ChannelID, "joinchat/"))
	} else {
		channelID, accessHash, err = publicChannelInfo(ctx, client, config.ChannelID)
	}

	if err != nil {
		return nil, err
	}

	last, err := queryLatestMsgID(ctx, client, channelID, accessHash)
	if err != nil {
		return nil, err
	}

	return &tgChannelInfo{
		id:         channelID,
		accessHash: accessHash,
		lastMsgID:  last,
	}, nil
}

// Query access hash for private channel.
func privateChannelInfo(ctx context.Context, client *telegram.Client, hash string) (channelID int64, accessHash int64, err error) {
	invite, err := client.API().MessagesCheckChatInvite(ctx, hash)
	if err != nil {
		return
	}

	switch v := invite.(type) {
	case *tg.ChatInviteAlready:
		if channel, ok := v.GetChat().(*tg.Channel); ok {
			channelID = channel.ID
			accessHash = channel.AccessHash
			return
		}
	case *tg.ChatInvitePeek:
		if channel, ok := v.GetChat().(*tg.Channel); ok {
			channelID = channel.ID
			accessHash = channel.AccessHash
			return
		}
	case *tg.ChatInvite:
		log.Warn("You haven't join this private channel, plz join it manually.")
	}

	err = errors.New("couldn't find access hash")
	return
}

// Query public channel by its name.
func publicChannelInfo(ctx context.Context, client *telegram.Client, channelName string) (channelID, accessHash int64, err error) {
	username, err := client.API().ContactsResolveUsername(ctx, channelName)
	if err != nil {
		return
	}

	if len(username.Chats) == 0 {
		err = fmt.Errorf("you are not belong to channel: %s", channelName)
		return
	}

	for _, chat := range username.Chats {
		// Try to find the related channel.
		if channel, ok := chat.(*tg.Channel); ok {
			channelID = channel.ID
			accessHash = channel.AccessHash
			return
		}
	}

	err = fmt.Errorf("couldn't find channel id and hash for channel: %s", channelName)
	return
}

// queryLatestMsgID from the given channel info.
func queryLatestMsgID(ctx context.Context, client *telegram.Client, channelID, accessHash int64) (int64, error) {
	request := &tg.MessagesSearchRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  channelID,
			AccessHash: accessHash,
		},
		Filter:   &tg.InputMessagesFilterEmpty{},
		Q:        "",
		OffsetID: -1,
		Limit:    1,
	}

	last := -1
	search, err := client.API().MessagesSearch(ctx, request)
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
