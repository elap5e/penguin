// Copyright 2022 Elapse and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package penguin

type AccountType string

const (
	AccountTypeDefault    AccountType = "DEFAULT"
	AccountTypeAnonymous  AccountType = "ANONYMOUS"
	AccountTypeChannel    AccountType = "CHANNEL"
	AccountTypeChannelBot AccountType = "CHANNEL_BOT"
)

type Account struct {
	ID       int64       `json:"id"`
	Type     AccountType `json:"type"`
	Username string      `json:"username,omitempty"`
	Photo    string      `json:"photo,omitempty"`
}

type ChatType string

const (
	// default
	ChatTypeDiscuss        ChatType = "DISCUSS"
	ChatTypeDiscussPrivate ChatType = "DISCUSS_PRIVATE"
	ChatTypeGroup          ChatType = "GROUP"
	ChatTypeGroupPrivate   ChatType = "GROUP_PRIVATE"
	ChatTypePrivate        ChatType = "PRIVATE"

	// channel
	ChatTypeChannel        ChatType = "CHANNEL"
	ChatTypeChannelPrivate ChatType = "CHANNEL_PRIVATE"

	// channel room
	ChatTypeRoomText    ChatType = "ROOM_TEXT"
	ChatTypeRoomVoice   ChatType = "ROOM_VOICE"
	ChatTypeRoomGroup   ChatType = "ROOM_GROUP"
	ChatTypeRoomLive    ChatType = "ROOM_LIVE"
	ChatTypeRoomApp     ChatType = "ROOM_APP"
	ChatTypeRoomForum   ChatType = "ROOM_FORUM"
	ChatTypeRoomPrivate ChatType = "ROOM_PRIVATE"
)

type Chat struct {
	ID      int64      `json:"id"`
	Type    ChatType   `json:"type"`
	Title   string     `json:"title,omitempty"`
	Photo   *ChatPhoto `json:"photo,omitempty"`
	Display string     `json:"display,omitempty"`

	PinnedMessage *Message `json:"pinned_message,omitempty"`

	Chat    *Chat `json:"chat,omitempty"`
	User    *User `json:"user,omitempty"`
	Channel *Chat `json:"channel,omitempty"`
}

type ChatPhoto struct {
	FileID string `json:"file_id"`
}

type Contact struct {
	Account *Account `json:"account"`
	Display string   `json:"display,omitempty"`
}

type Role struct {
	ID    int64  `json:"id"`
	Color uint32 `json:"color,omitempty"`
	Hoist bool   `json:"hoist,omitempty"`
}

type User struct {
	Account *Account `json:"account"`
	Display string   `json:"display,omitempty"`
}

type Message struct {
	MessageID int64            `json:"message_id"`
	Chat      *Chat            `json:"chat"`
	From      *User            `json:"from,omitempty"`
	Forward   *Message         `json:"forward,omitempty"`
	ReplyTo   *Message         `json:"reply_to,omitempty"`
	Time      int64            `json:"time"`
	EditTime  int64            `json:"edit_time,omitempty"`
	Text      string           `json:"text,omitempty"`
	Entities  []*MessageEntity `json:"entities,omitempty"`

	Animation *Animation `json:"animation,omitempty"`
	Audio     *Audio     `json:"audio,omitempty"`
	Document  *Document  `json:"document,omitempty"`
	Photo     *Photo     `json:"photo,omitempty"`
	Sticker   *Sticker   `json:"sticker,omitempty"`
	Video     *Video     `json:"video,omitempty"`
	Voice     *Voice     `json:"voice,omitempty"`
	Contact   *Contact   `json:"contact,omitempty"`
	Dice      *Dice      `json:"dice,omitempty"`
	Poll      *Poll      `json:"poll,omitempty"`
	Venue     *Venue     `json:"venue,omitempty"`

	NewChatUsers []*User `json:"new_chat_users,omitempty"`
	LeftChatUser *User   `json:"left_chat_user,omitempty"`

	NewChatTitle string `json:"new_chat_title,omitempty"`
	NewChatPhoto *Photo `json:"new_chat_photo,omitempty"`

	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	URL    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

type Animation struct{}

type Audio struct{}

type Document struct{}

type Photo struct {
	ID     int64   `json:"id"`
	Path   string  `json:"path"`
	Name   string  `json:"name"`
	Size   int64   `json:"size"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Digest *Digest `json:"digest,omitempty"`
}

type Sticker struct{}

type Video struct {
	ID     int64   `json:"id"`
	UUID   string  `json:"uuid"`
	Name   string  `json:"name"`
	Size   int64   `json:"size"`
	Digest *Digest `json:"digest,omitempty"`

	Thumbnail *Photo `json:"thumbnail,omitempty"`
}

type Voice struct {
	ID     int64   `json:"id"`
	Path   string  `json:"path"`
	Name   string  `json:"name"`
	Size   int64   `json:"size"`
	Digest *Digest `json:"digest,omitempty"`
}

type Dice struct {
	Emoji string `json:"emoji"`
	Value int8   `json:"value"`
}

type PollOption struct {
	Text  string `json:"text"`
	Count int64  `json:"count"`

	Users []*User `json:"users,omitempty"`
}

type Poll struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`

	Multiple bool          `json:"multiple"`
	Options  []*PollOption `json:"options"`

	Count int64 `json:"count"`

	IsAnonymous bool `json:"is_anonymous"`
	IsClosed    bool `json:"is_closed"`
}

type Venue struct{}

type Digest struct {
	MD5    []byte `json:"md5,omitempty"`
	SHA1   []byte `json:"sha1,omitempty"`
	SHA256 []byte `json:"sha256,omitempty"`
	SHA512 []byte `json:"sha512,omitempty"`
}
