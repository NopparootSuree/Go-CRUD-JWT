package models

import (
	"time"
)

type Users struct {
	ID             uint      `gorm:"primarykey;column:id;autoIncrement"`
	Username       string    `gorm:"column:username;not null"`
	HashedPassword string    `gorm:"column:hashedPassword;not null"`
	Fullname       string    `gorm:"column:fullName;not null"`
	Email          string    `gorm:"column:email;index;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

type Posts struct {
	PostID    uint      `gorm:"primarykey;column:postID;autoIncrement"`
	Title     string    `gorm:"column:title;not null"`
	Body      string    `gorm:"column:body;not null"`
	UserID    uint      `gorm:"column:userID;index;foreignkey:UserID;references:ID;not null"`
	Status    string    `gorm:"column:status;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Follows struct {
	FollowingUserID uint      `gorm:"column:followingUserID;foreignkey:FollowingUserID;references:ID;not null"`
	FollowerUserID  uint      `gorm:"column:followerUserID;foreignkey:FollowerUserID;references:ID;not null"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}
