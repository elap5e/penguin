package penguin

type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Photo    string `json:"photo,omitempty"`
}

type Chat struct {
	ID            int64      `json:"id"`
	Type          string     `json:"type"`
	Title         string     `json:"title"`
	Photo         *ChatPhoto `json:"photo,omitempty"`
	Description   string     `json:"description,omitempty"`
	PinnedMessage *Message   `json:"pinned_message,omitempty"`
}

type ChatPhoto struct {
	FileID string `json:"file_id"`
}

type Channel struct{}

type Contact struct{}

type User struct {
	Account *Account `json:"account"`
	Display string   `json:"display,omitempty"`
}

type Message struct {
	MessageID     int64            `json:"message_id"`
	Chat          *Chat            `json:"chat"`
	From          *User            `json:"from,omitempty"`
	Forward       *Message         `json:"forward,omitempty"`
	ReplyTo       *Message         `json:"reply_to,omitempty"`
	PinnedMessage *Message         `json:"pinned_message,omitempty"`
	Time          int64            `json:"time"`
	EditTime      int64            `json:"edit_time,omitempty"`
	Text          string           `json:"text,omitempty"`
	Entities      []*MessageEntity `json:"entities,omitempty"`
	Audio         *Audio           `json:"audio,omitempty"`
	Document      *Document        `json:"document,omitempty"`
	Photo         *Photo           `json:"photo,omitempty"`
	Sticker       *Sticker         `json:"sticker,omitempty"`
	Video         *Video           `json:"video,omitempty"`
	Voice         *Voice           `json:"voice,omitempty"`
	Contact       *Contact         `json:"contact,omitempty"`
	Dice          *Dice            `json:"dice,omitempty"`
	Poll          *Poll            `json:"poll,omitempty"`
	NewChatUsers  []*User          `json:"new_chat_users,omitempty"`
	LeftChatUser  *User            `json:"left_chat_user,omitempty"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	URL    string `json:"url,omitempty"`
	User   *User  `json:"user,omitempty"`
}

type Audio struct {
}

type Document struct {
}

type Photo struct {
}

type Sticker struct {
}

type Video struct {
}

type Voice struct {
}

type Dice struct {
}

type Poll struct {
}
