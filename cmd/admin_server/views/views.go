package views

type ContextKey string

type FlashMessages struct {
	Error []string
}

var FlashMessagesContextKey ContextKey = "flash-messages"
