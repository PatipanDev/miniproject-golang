package domain

type EmailJobPayload struct {
	To      string
	Subject string
	Body    string
}

func (p EmailJobPayload) Kind() string {
	return "send_email"
}
