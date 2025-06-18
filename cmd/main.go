/* package main

import (
	"log"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"
	routers "github.com/PatipanDev/mini-project-golang/internal/adapters/http/routes"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/repositories"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/services"
	"github.com/PatipanDev/mini-project-golang/pkg/configs"
	"github.com/PatipanDev/mini-project-golang/pkg/smtp"
)

func main() {

	param := configs.NewFiberHttpServiceParams()
	fiberServ := configs.NewFiberHTTPService(param)

	db, err := configs.NewDatabase()
	if err != nil {
		log.Fatal("Failed to start Database:", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.EmailLog{})
	if err != nil {
		log.Fatal("Failed to auto migrate models:", err)
	}

	userRepo := repositories.NewUerRepository(db)
	fileRepo := repositories.NewFileStorageRepository("internal/adapters/storage", configs.BASE_URL)

	// EmailLogRepository & EmailSender
	emailLogRepo := repositories.NewEmailLogRepository(db)
	emailSender := smtp.NewSMTPSender(emailLogRepo)

	roleRepo := repositories.NewRoleRepository(db)

	userServ := services.NewUserService(userRepo, fileRepo, emailSender, roleRepo)
	userHandler := handlers.NewHttpUserHandler(userServ)

	authServ := services.NewAurhService(userRepo, configs.SECRET_KEY)
	authHandler := handlers.NewAuthHandler(authServ)

	api := fiberServ.Group("/api")
	routers.UserRoutes(api, userHandler, authHandler)
	routers.UploadRoutes(api, userHandler)

	fiberServ.Static("uploads", "internal/adapters/storage/uploads")

	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
*/

package main

import (
	"context"
	"log"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"
	routers "github.com/PatipanDev/mini-project-golang/internal/adapters/http/routes"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/repositories"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/services"
	"github.com/PatipanDev/mini-project-golang/internal/jobs"
	"github.com/PatipanDev/mini-project-golang/pkg/configs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func main() {
	// --- การตั้งค่าพื้นฐาน (คงเดิม) ---
	param := configs.NewFiberHttpServiceParams()
	fiberServ := configs.NewFiberHTTPService(param)

	db, err := configs.NewDatabase()
	if err != nil {
		log.Fatal("Failed to start Database:", err)
	}

	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.EmailLog{})
	if err != nil {
		log.Fatal("Failed to auto migrate models:", err)
	}

	// --- การตั้งค่า River ---
	ctx := context.Background()

	// 1. สร้างการเชื่อมต่อด้วย pgx สำหรับ River
	pgxPool, err := pgxpool.New(ctx, configs.DB_URL)
	if err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อ PGX Pool ได้: %v", err)
	}
	defer pgxPool.Close()

	// --- การสร้าง Repositories ---
	userRepo := repositories.NewUerRepository(db)
	fileRepo := repositories.NewFileStorageRepository("internal/adapters/storage", configs.BASE_URL)
	emailLogRepo := repositories.NewEmailLogRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// --- สร้างและเริ่ม River Worker ---
	// 3. สร้าง Worker สำหรับส่งอีเมล
	sendEmailWorker := jobs.NewSendEmailWorker(emailLogRepo)

	// 4. ลงทะเบียน Worker (สร้าง Workers collection)
	// **ตรวจสอบให้แน่ใจว่า 'workers' นี้ถูกสร้างและใช้งานเพียงครั้งเดียว**
	workers := river.NewWorkers()
	river.AddWorker(workers, sendEmailWorker)

	// 2. สร้าง River Driver และ Client
	riverDriver := riverpgxv5.New(pgxPool)
	// **นี่คือจุดที่สำคัญที่สุด: ส่ง 'workers' เข้าไปใน 'river.Config' ตอนสร้าง 'client'**
	riverClient, err := river.NewClient(riverDriver, &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers, // <--- ตรงนี้ต้องมี
	})
	if err != nil {
		log.Fatalf("ไม่สามารถสร้าง river client ได้: %v", err)
	}

	// 5. เริ่มการทำงานของ Worker
	// **ตรวจสอบให้แน่ใจว่า 'riverClient' ที่นี่คือตัวเดียวกับที่สร้างด้านบน**

	// ใช้โค้ดส่วนนี้ได้เลย เพราะไลบรารีเป็นเวอร์ชันล่าสุดแล้ว

	// **แก้ไขเป็นแบบนี้**
	// ในเวอร์ชันเก่า client จะเป็นตัวจัดการเริ่ม worker เอง
	if err := riverClient.Start(ctx); err != nil {
		//log.Fatalf("ไม่สามารถเริ่ม river client/worker ได้: %v", err)
	}

	if err := riverClient.Stop(ctx); err != nil {
		// handle error
	}

	// --- การสร้าง Services ---
	userServ := services.NewUserService(userRepo, fileRepo, roleRepo, riverClient)
	userHandler := handlers.NewHttpUserHandler(userServ)

	authServ := services.NewAurhService(userRepo, configs.SECRET_KEY)
	authHandler := handlers.NewAuthHandler(authServ)

	// --- การตั้งค่า Routes และ Server ---
	api := fiberServ.Group("/api")
	routers.UserRoutes(api, userHandler, authHandler)
	routers.UploadRoutes(api, userHandler)

	fiberServ.Static("uploads", "internal/adapters/storage/uploads")

	log.Println("Starting Fiber HTTP server on port " + configs.SERVER_HTTP_PORT)
	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
