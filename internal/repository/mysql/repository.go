package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/vodolaz095/purser/model"
	"github.com/vodolaz095/purser/pkg"
)

type secretData struct {
	ID        string    `gorm:"primaryKey"`
	Encoded   []byte    `gorm:"type:text"`
	CreatedAt time.Time `json:"createdAt" gorm:"index"`
}

type bodyData struct {
	Body string            `json:"Body"`
	Meta map[string]string `json:"meta"`
}

// Repository реализует интерфейс SecretRepo с базой данных mysql/mariadb внутри
type Repository struct {
	DatabaseConnectionString string
	db                       *gorm.DB
}

// Ping проверяет соединение с базой данных
func (r *Repository) Ping(ctx context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.PingContext(ctx)
}

// Init настраивает соединение с базой данных
func (r *Repository) Init(ctx context.Context) error {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(r.DatabaseConnectionString), &gorm.Config{})
	if err != nil {
		return err
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)
	r.db = db
	err = r.Ping(ctx)
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).
		Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&secretData{})
	return err
}

// Close закрывает соединение с базой данных
func (r *Repository) Close(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Create создаёт новый model.Secret
func (r *Repository) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	bd := bodyData{
		Body: body,
		Meta: meta,
	}
	data, err := json.Marshal(bd)
	if err != nil {
		return model.Secret{}, err
	}

	databaseSecretData := secretData{
		ID:        pkg.UUID(),
		Encoded:   data,
		CreatedAt: time.Now(),
	}
	secret := model.Secret{
		ID:        databaseSecretData.ID,
		Body:      body,
		Meta:      meta,
		CreatedAt: databaseSecretData.CreatedAt,
		ExpireAt:  databaseSecretData.CreatedAt.Add(model.TTL),
	}

	err = r.db.
		WithContext(ctx).
		Save(&databaseSecretData).Error
	if err != nil {
		return model.Secret{}, err
	}
	return secret, nil
}

// FindByID ищет model.Secret по идентификатору
func (r *Repository) FindByID(ctx context.Context, id string) (model.Secret, error) {
	var databaseSecretData secretData
	err := r.db.WithContext(ctx).First(&databaseSecretData, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Secret{}, model.ErrSecretNotFound
		}
		return model.Secret{}, err
	}
	// expired
	if databaseSecretData.CreatedAt.Add(model.TTL).Before(time.Now()) {
		return model.Secret{}, model.ErrSecretNotFound
	}
	var params bodyData
	err = json.Unmarshal(databaseSecretData.Encoded, &params)
	if err != nil {
		return model.Secret{}, err
	}
	secret := model.Secret{
		ID:        id,
		Body:      params.Body,
		Meta:      params.Meta,
		CreatedAt: databaseSecretData.CreatedAt,
		ExpireAt:  databaseSecretData.CreatedAt.Add(model.TTL),
	}
	return secret, nil
}

// DeleteByID удаляет секрет по идентификатору
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(&secretData{}, id).Error
}

// Prune удаляет старые секреты
func (r *Repository) Prune(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", time.Now().Add(-model.TTL)).
		Delete(&secretData{}).Error
}
