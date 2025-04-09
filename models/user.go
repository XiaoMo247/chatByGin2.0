package models

import (
	"errors"
	"html"
	"strings"

	"chatByGin5.0/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Token    string `gorm:"unique"`
}

func (u *User) SaveUser() (*User, error) {
	tx := DB.Begin()
	err := tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// 使用gorm的hook在保存密码前对密码进行hash
func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username, password string) (string, error) {
	var err error
	u := User{}
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)
	if err != nil {
		return "", err
	}

	// 生成新 Token
	newToken, err := utils.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}

	// 更新数据库里的 Token（使旧 Token 失效）
	err = DB.Model(&u).Update("token", newToken).Error
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// 返回前将用户密码置空
func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserByID(uid uint) (User, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()
	return u, nil
}
