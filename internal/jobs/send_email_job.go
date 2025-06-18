package jobs

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
)

// SendEmailArgs คือ struct ที่กำหนดอาร์กิวเมนต์สำหรับ Job ส่งอีเมล
type SendEmailArgs struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Kind คือฟังก์ชันที่คืนค่าชื่อเฉพาะของ Job ประเภทนี้
func (SendEmailArgs) Kind() string {
	return "send_email"
}

// SendEmailWorker คือ Worker ที่จะประมวลผล Job ประเภท SendEmailArgs
// โดยจะเก็บสิ่งที่ต้องพึ่งพา (Dependencies) เช่น repository และค่าตั้งค่า SMTP
type SendEmailWorker struct {
	river.WorkerDefaults[SendEmailArgs]

	emailLogRepo ports.EmailLogRepository
	from         string
	addr         string
	auth         smtp.Auth
}

// NewSendEmailWorker ใช้สำหรับสร้าง Email Worker ใหม่พร้อมกับ Dependencies
func NewSendEmailWorker(emailLogRepo ports.EmailLogRepository) *SendEmailWorker {
	// โหลดค่าตั้งค่า SMTP ที่นี่ เพื่อให้ Worker ทำงานได้ด้วยตัวเอง
	_ = godotenv.Load()
	from := os.Getenv("SMTP_SENDER")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	addr := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	return &SendEmailWorker{
		emailLogRepo: emailLogRepo,
		from:         from,
		addr:         addr,
		auth:         auth,
	}
}

// Work คือฟังก์ชันที่จะถูกเรียกโดย River framework เพื่อประมวลผล Job
func (w *SendEmailWorker) Work(ctx context.Context, job *river.Job[SendEmailArgs]) error {
	args := job.Args
	log.Printf("กำลังส่งอีเมลไปที่ %s ด้วยหัวข้อ '%s'", args.To, args.Subject)

	msg := []byte("From: " + w.from + "\r\n" +
		"To: " + args.To + "\r\n" +
		"Subject: " + args.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		args.Body + "\r\n")

	// ส่งอีเมล
	err := smtp.SendMail(w.addr, w.auth, w.from, []string{args.To}, msg)

	// บันทึก Log การส่งอีเมลเสมอ ไม่ว่าจะสำเร็จหรือล้มเหลว (เหมือนโค้ดเดิม)
	logErr := w.emailLogRepo.Save(&domain.EmailLog{
		EventID: uuid.New(),
		Subject: args.Subject,
		From:    w.from,
		To:      args.To,
		Body:    args.Body,
	})
	if logErr != nil {
		// หากการบันทึก Log ล้มเหลว ให้แสดงผลออกทาง standard log
		log.Printf("critical: ไม่สามารถบันทึก email log ได้: %v", logErr)
	}

	// หากการส่งล้มเหลว ให้ return error ออกไปเพื่อให้ River ลองทำงานนี้ใหม่
	if err != nil {
		return fmt.Errorf("ไม่สามารถส่งอีเมลได้: %w", err)
	}

	log.Printf("ส่งอีเมลไปที่ %s สำเร็จ", args.To)
	return nil
}
