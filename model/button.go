package model

type Keyboard struct {
	rows []*KeyboardRow
}

type KeyboardRow struct {
	cols []*Button
}

type Button struct {
	Text      string     `json:"text"`
	Type      ButtonType `json:"type"`
	Intent    Intent     `json:"intent"`
	URL       string     `json:"url,omitempty"`
	Quick     bool       `json:"quick,omitempty"`
	WebApp    string     `json:"web_app,omitempty"`
	ContactID int64      `json:"contact_id,omitempty"`
	Payload   string     `json:"payload,omitempty"`
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		rows: make([]*KeyboardRow, 0),
	}
}

func (k *Keyboard) Build() [][]*Button {
	buttons := make([][]*Button, 0)
	for _, row := range k.rows {
		buttons = append(buttons, row.Build())
	}

	return buttons
}

func (k *Keyboard) AddRow() *KeyboardRow {
	kr := &KeyboardRow{}
	k.rows = append(k.rows, kr)

	return kr
}

func (k *KeyboardRow) Build() []*Button {
	buttons := make([]*Button, 0, len(k.cols))
	buttons = append(buttons, k.cols...)

	return buttons
}

func (k *KeyboardRow) AddCallback(text string, indent Intent, payload string) *KeyboardRow {
	kr := &Button{
		Type:    ButtonCallback,
		Text:    text,
		Intent:  indent,
		Payload: payload,
	}

	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddLink(text, link string) *KeyboardRow {
	kr := &Button{

		Type: ButtonLink,
		Text: text,
		URL:  link,
	}
	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddGeoLocation(text string, quick bool) *KeyboardRow {
	kr := &Button{
		Type:  ButtonRequestGeo,
		Text:  text,
		Quick: quick,
	}
	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddContact(text string) *KeyboardRow {
	kr := &Button{
		Type: ButtonRequestContact,
		Text: text,
	}
	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddMessage(text string) *KeyboardRow {
	kr := &Button{
		Type: ButtonMessage,
		Text: text,
	}
	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddOpenApp(text string, botID int64) *KeyboardRow {
	kr := &Button{
		Type:      ButtonOpenApp,
		Text:      text,
		WebApp:    text, // ignored, required
		ContactID: botID,
	}
	k.cols = append(k.cols, kr)

	return k
}

func (k *KeyboardRow) AddClipboard(text, payload string) *KeyboardRow {
	kr := &Button{
		Type:    ButtonClipboard,
		Text:    text,
		Payload: payload,
	}
	k.cols = append(k.cols, kr)

	return k
}
