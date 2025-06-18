package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImp struct {
	userRepo  ports.UserRepository
	fileRepo  ports.FileStorageRepository
	emailRepo ports.EmailLogRepository
	roleRepo  ports.RoleRepository
}

func NewUserService(repo ports.UserRepository, fileRepo ports.FileStorageRepository, roleRepo ports.RoleRepository, emailRepo ports.EmailLogRepository) *UserServiceImp {
	return &UserServiceImp{
		userRepo:  repo,
		fileRepo:  fileRepo,
		roleRepo:  roleRepo,
		emailRepo: emailRepo,
	}
}

func (s *UserServiceImp) RegisterUser(user *domain.User) error {
	existing, _ := s.userRepo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// กำหนด Role เริ่มต้นให้เป็น employee
	employeeRole := domain.Role{}
	if err := s.roleRepo.FindByName("employee", &employeeRole); err != nil {
		// create if not exists
		employeeRole = domain.Role{Name: domain.USER_ROLE_EMPLOYEE}
		if err := s.roleRepo.Create(&employeeRole); err != nil {
			return errors.New("cannot create default role 'employee'")
		}
	}

	newUser := &domain.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Password:  string(hashedPassword),
		Status:    domain.USER_STATUS_ACTIVE,
		Roles:     []domain.Role{employeeRole},
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return err
	}

	// --- ส่วนของตรรกะการส่งอีเมลถูกแทนที่ ---
	// Goroutine หายไป ตอนนี้เราจะนำ Job เข้าคิวแทน

	subject := "ยินดีต้อนรับสู่ระบบของเรา!"
	// เกร็ดความรู้: ลิงก์ในโค้ดเดิมของคุณคือวิดีโอเพลง "Never Gonna Give You Up" ของ Rick Astley ครับ
	body := fmt.Sprintf(`<h2 style="background-color:powderblue; color:blue;">สวัสดีคุณ %s %s</h2>
		<p style="color:red;">ขอบคุณที่สมัครใช้งานระบบของเรา </p>
		<p color:yellow;>หากมีคำถามสามารถติดต่อทีมงานได้ตลอดเวลา <a href="https://youtu.be/dQw4w9WgXcQ?si=7f8rrTxSrYJd2aZV">คลิ๊กที่นี่</a></p>`,
		newUser.FirstName, newUser.LastName)

	// สร้างข้อมูลสำหรับ Job
	emailJobArgs := domain.SendEmailArgs{
		To:      newUser.Email,
		Subject: subject,
		Body:    body,
	}

	if err := s.emailRepo.ProcessUserEmailJob(&emailJobArgs); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImp) UpdateUser(user *domain.User, id string) error {
	err := s.userRepo.Update(user, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImp) DeleteUser(id string) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImp) GetUserById(id string) (*domain.User, error) {
	user, err := s.userRepo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImp) FindUsers() ([]domain.User, error) {
	r, err := s.userRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("error get users : %w", err)
	}
	return r, nil
}

func (s *UserServiceImp) GetPaginationUsers(page int, limit int) (*domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	total, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	users, err := s.userRepo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	totalPage := int((total + int64(limit) - 1) / int64(limit))

	users1 := domain.User{}

	var roles []string
	for _, r := range users1.Roles {
		roles = append(roles, string(r.Name))
	}

	var result []domain.ResUSerNoID
	for _, u := range users {
		result = append(result, domain.ResUSerNoID{
			FullName:   u.FirstName + " " + u.LastName,
			EmployeeID: u.EmployeeID,
			Email:      u.Email,
			Status:     string(u.Status),
			UpdatedAt:  u.UpdatedAt,
			Roles:      roles,
		})
	}

	return &domain.Pagination{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPage,
		Data:      result,
	}, nil
}

func (s *UserServiceImp) GetUsers(filter *domain.UserFilter) ([]domain.ResUSerNoID, int64, error) {
	users, total, err := s.userRepo.FindUsers(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("error get users: %w", err)
	}

	var results []domain.ResUSerNoID
	for _, u := range users {
		// Map roles to []string
		var roleNames []string
		for _, r := range u.Roles {
			roleNames = append(roleNames, string(r.Name))
		}

		results = append(results, domain.ResUSerNoID{
			FullName:   u.FirstName + " " + u.LastName,
			EmployeeID: u.EmployeeID,
			Email:      u.Email,
			Status:     string(u.Status),
			Roles:      roleNames,
		})
	}
	return results, total, nil
}

func (s *UserServiceImp) UploadProfilePicture(id string, file []byte, filename string) (string, error) {
	ext := filepath.Ext(filename)
	newFile := fmt.Sprintf("%s_%d%s", id, time.Now().UnixNano(), ext)
	folderPath := "uploads/profile_pictures"

	fileName, err := s.fileRepo.SaveFile(folderPath, newFile, file)

	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Println("oooooooooooo", fileName)

	err = s.userRepo.UpdateUserProfilePicName(id, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to uploade user profile picture URL in database:%w", err)
	}
	return fileName, nil
}
