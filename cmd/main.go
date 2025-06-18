package main

import (
	"context"
	"log"

	"github.com/PatipanDev/mini-project-golang/internal/adapters/http/handlers"
	routers "github.com/PatipanDev/mini-project-golang/internal/adapters/http/routes"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/repositories"
	"github.com/PatipanDev/mini-project-golang/internal/adapters/worker"
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/services"
	"github.com/PatipanDev/mini-project-golang/pkg/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	param := configs.NewFiberHttpServiceParams()
	fiberServ := configs.NewFiberHTTPService(param)

	// db, err := configs.NewDatabase()
	// if err != nil {
	// 	log.Fatal("Failled to start Database:", err)

	// }

	db, err := configs.NewDatabaseGromRiver()
	if err != nil {
		log.Fatal("Failed to start Database:", err)
	}

	pgxPool, err := pgxpool.New(ctx, configs.DB_URL)
	if err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อ PGX Pool ได้: %v", err)
	}
	defer pgxPool.Close()

	// Auto migrate models
	err = db.AutoMigrate(&domain.User{}, &domain.Role{}, &domain.EmailLog{})
	sqlDB, _ := db.DB()

	// minio: เอาไว็เก็บไฟล์ต่างๆ
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "profile"
	location := "us-east-1"
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Println("Bucket already exists")
		} else {
			log.Fatal(err)
		}
	}

	// EmailLogRepository & EmailSender

	roleRepo := repositories.NewRoleRepository(db)

	userRepo := repositories.NewUserRepository(db, minioClient, bucketName)
	fileRepo := repositories.NewFileStorageRepository(bucketName, configs.MINIO_URL, minioClient)
	uploadRepo := repositories.NewUploadRepository(db)

	authServ := services.NewAurhService(userRepo, configs.SECRET_KEY)
	uploadServ := services.NewUploadService(uploadRepo)

	authHandler := handlers.NewAuthHandler(authServ)

	workers_email := river.NewWorkers()
	// **นี่คือจุดที่สำคัญที่สุด: ส่ง 'workers' เข้าไปใน 'river.Config' ตอนสร้าง 'client'**
	//ประกาศตัวแปรเอาไว้ก่อนยังไงใส่ค่า
	sendEmailWorker := &worker.SendEmailWorker{}
	if err := river.AddWorkerSafely(workers_email, sendEmailWorker); err != nil {
		log.Fatal(err)
	}

	riverClient_email, err := river.NewClient(riverdatabasesql.New(sqlDB), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers_email, // <--- ตรงนี้ต้องมี
	})
	if err != nil {
		log.Fatalf("ไม่สามารถสร้าง river client ได้: %v", err)
	}

	emailLogRepo := repositories.NewEmailLogRepository(db, riverClient_email)

	// InjectRepository
	sendEmailWorker.InjectRepository(emailLogRepo)
	userServ := services.NewUserService(userRepo, fileRepo, roleRepo, emailLogRepo)

	/// upload
	workers := river.NewWorkers()
	uploadWorker := worker.NewUploadWorker(uploadServ)
	if err := river.AddWorkerSafely(workers, uploadWorker); err != nil {
		log.Fatal(err)
	}

	riverClient, err := river.NewClient(riverdatabasesql.New(sqlDB), &river.Config{
		Workers: workers,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := riverClient.Start(ctx); err != nil {
		log.Fatal(err)
	}

	uploadHandler := handlers.NewUploadHandler(uploadServ, db, riverClient)
	userHandler := handlers.NewHttpUserHandler(userServ)

	api := fiberServ.Group("/api")
	routers.UserRoutes(api, userHandler, authHandler)
	routers.UploadRoutes(api, userHandler)
	routers.UploadProfileRoutes(api, uploadHandler)

	fiberServ.Static("uploads", "internal/adapters/storage/uploads")

	err = fiberServ.Listen(":" + configs.SERVER_HTTP_PORT)
	if err != nil {
		log.Fatal("Failed to start Fiber HTTP server:", err)
	}
}
