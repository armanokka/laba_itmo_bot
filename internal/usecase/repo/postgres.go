package repo

import (
	"github.com/armanokka/translobot/internal/usecase/entity"
	"gorm.io/gorm"
)

type TranslationRepo struct {
	client *gorm.DB
}

func New(client *gorm.DB) (TranslationRepo, error) {
	if err := client.AutoMigrate(&entity.Translations{}); err != nil {
		return TranslationRepo{}, err
	}
	return TranslationRepo{client: client}, nil
}
