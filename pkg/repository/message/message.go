package message

import (
	"github.com/go-redis/redis"
	"github.com/nlsra/message-system/pkg/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewRepository(db *gorm.DB, client *redis.Client) *Repository {
	return &Repository{
		DB:          db,
		RedisClient: client,
	}
}

func (msg *Repository) FindAllSentMessages() ([]model.Message, error) {
	messages := []model.Message{}
	err := msg.DB.Find(&messages, "status = ?", "success").Error
	return messages, err
}

func (msg *Repository) UpdateMessageStatus(message *model.Message, status string) error {
	return msg.DB.Model(&message).Update("status", status).Error
}

func (msg *Repository) Migrate() error {
	return msg.DB.AutoMigrate(&model.Message{})
}
