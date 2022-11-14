package dashbot

type Message struct {
	Text         string   `json:"text"`
	UserId       int64    `json:"userId"`
	Intent       Intent   `json:"intent,omitempty"`
	Images       []string `json:"images,omitempty"`
	Buttons      []Button `json:"buttons,omitempty"`
	Postback     Postback `json:"postback,omitempty"`
	PlatformJson string   `json:"platformJson,omitempty"` // any json
	//PlatformUserJson map[string]interface{} `json:"platformUserJson"` // there is paid access to this field
	SessionId int64 `json:"sessionId,omitempty"` // the same as Message.UserId
}

type Postback struct {
	ButtonClick ButtonClick `json:"buttonClick"`
}

type ButtonClick struct {
	ButtonId string `json:"buttonId"`
}

type Intent struct {
	Name       string  `json:"name"`
	Inputs     []Input `json:"inputs,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type Input struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Button struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type EventType string

const (
	CustomEvent     EventType = "customEvent"
	RevenueEvent    EventType = "revenueEvent"
	ShareEvent      EventType = "shareEvent"
	PageLaunchEvent EventType = "pageLaunchEvent"
	ReferralEvent   EventType = "referralEvent"
)

type Event struct {
	Name           string `json:"name"`
	ConversationId string `json:"conversationId"`
	Type           EventType
	ExtraInfo      string `json:"extraInfo"`
}
