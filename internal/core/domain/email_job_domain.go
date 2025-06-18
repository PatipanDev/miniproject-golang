package domain

type SendEmailArgs struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Kind คือฟังก์ชันที่คืนค่าชื่อเฉพาะของ Job ประเภทนี้
func (SendEmailArgs) Kind() string {
	return "send_email"
}
