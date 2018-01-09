package push

type Message struct {
	Token string
	Alert string
	Badge int
	Sound string
}

func (this *Message) CreateMessage(token, alert, sound string, badge int) *Message {
	model := new(Message)
	model.Token = token
	model.Alert = alert
	model.Sound = "default"
	if sound != "" {
		model.Sound = sound
	}
	model.Badge = 1
	if badge > 0 {
		model.Badge = badge
	}
	return model
}
