package notifier

// Data is a struct that represents the data of a message
type Data struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}
