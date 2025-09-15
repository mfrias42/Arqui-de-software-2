package dao

import (
	"gorm.io/gorm"
)

type InscriptionModel struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	UserID   uint `gorm:"not null;index"`
	CourseID uint `gorm:"not null;index"`
}

type InscriptionDAO struct {
	db *gorm.DB
}

func NewInscriptionDAO(db *gorm.DB) *InscriptionDAO {
	return &InscriptionDAO{db: db}
}

func (dao *InscriptionDAO) DB() *gorm.DB {
	return dao.db
}
