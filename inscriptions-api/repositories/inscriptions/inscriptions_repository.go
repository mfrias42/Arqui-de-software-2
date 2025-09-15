package repositories

import (
	"context"
	"errors"
	"fmt"
	dao "inscriptions-api/DAOs/inscriptions"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Connect() (*gorm.DB, error) {
	dbHost := getEnv("DB_HOST", "mysql")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "inscriptions")

	fmt.Printf("Connecting to MySQL: Host=%s, Port=%s, User=%s, DBName=%s\n", dbHost, dbPort, dbUser, dbName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=30s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %v", err)
	}

	err = db.AutoMigrate(&dao.InscriptionModel{})
	if err != nil {
		return nil, fmt.Errorf("error migrating database: %v", err)
	}

	return db, nil
}

type InscriptionRepository struct {
	dao *dao.InscriptionDAO
}

func NewInscriptionRepository(dao *dao.InscriptionDAO) *InscriptionRepository {
	return &InscriptionRepository{dao: dao}
}

func (r *InscriptionRepository) CreateInscription(ctx context.Context, userID, courseID uint) (*dao.InscriptionModel, error) {
	var inscription dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&inscription).Error; err == nil {
		return nil, errors.New("inscription already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	newInscription := dao.InscriptionModel{UserID: userID, CourseID: courseID}
	if err := r.dao.DB().WithContext(ctx).Create(&newInscription).Error; err != nil {
		return nil, err
	}

	return &newInscription, nil
}

func (r *InscriptionRepository) GetInscriptions(ctx context.Context) ([]dao.InscriptionModel, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return inscriptionsModel, nil
}

func (r *InscriptionRepository) GetInscriptionsByUser(ctx context.Context, userID uint) ([]dao.InscriptionModel, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("user_id = ?", userID).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return inscriptionsModel, nil
}

func (r *InscriptionRepository) GetInscriptionsByCourse(ctx context.Context, courseID uint) ([]dao.InscriptionModel, error) {
	var inscriptionsModel []dao.InscriptionModel
	if err := r.dao.DB().WithContext(ctx).Where("course_id = ?", courseID).Find(&inscriptionsModel).Error; err != nil {
		return nil, err
	}

	return inscriptionsModel, nil
}
