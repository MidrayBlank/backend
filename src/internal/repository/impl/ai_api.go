package impl

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"backend/src/internal/db/abstract"
	"backend/src/internal/domain"
	"backend/src/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AiApiRepository struct{}

func NewAiApiRepository() *AiApiRepository {
	return &AiApiRepository{}
}

func (r *AiApiRepository) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (r *AiApiRepository) Insert(conn abstract.IDBConnection, token string) error {
	db := conn.Get().(*gorm.DB)
	hash := r.hashToken(token)

	dao := &model.AiApi{}
	dao, err := dao.ToModel(hash)
	if err != nil {
		return err
	}

	return db.Where("hash = ?", hash).
		FirstOrCreate(dao).Error
}

func (r *AiApiRepository) IncreaseRequests(conn abstract.IDBConnection, token string) error {
	db := conn.Get().(*gorm.DB)
	hash := r.hashToken(token)

	result := db.Model(&model.AiApi{}).
		Where("hash = ?", hash).
		Update("requests", gorm.Expr("requests + 1"))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (r *AiApiRepository) GetAllRequestsCount(conn abstract.IDBConnection, tokens []string) ([]*domain.AiApi, error) {
	db := conn.Get().(*gorm.DB)

	tokensHashMap := make(map[string]string, len(tokens))
	hashedTokens := make([]string, len(tokens))

	for i, token := range tokens {
		hash := r.hashToken(token)
		tokensHashMap[hash] = token

		hashedTokens[i] = r.hashToken(token)
	}

	var AiApiDaos []model.AiApi
	err := db.Where("hash IN ?", hashedTokens).
		Find(&AiApiDaos).Error
	if err != nil {
		return nil, err
	}

	var modelObj model.AiApi
	return modelObj.ToDomainSlice(AiApiDaos, tokensHashMap)
}

func (r *AiApiRepository) ResetAllRequestsCount(conn abstract.IDBConnection) error {
	db := conn.Get().(*gorm.DB)

	return db.Model(&model.AiApi{}).
		Update("requests", 0).Error
}

func (r *AiApiRepository) InsertIfNotExist(conn abstract.IDBConnection, tokens []string) error {
	db := conn.Get().(*gorm.DB)

	hashedTokens := make([]string, 0, len(tokens))
	for _, token := range tokens {
		hash := r.hashToken(token)
		hashedTokens = append(hashedTokens, hash)
	}

	records := make([]model.AiApi, 0, len(tokens))
	for i := range tokens {
		hash := hashedTokens[i]
		dao := &model.AiApi{}
		dao, err := dao.ToModel(hash)
		if err != nil {
			return err
		}
		records = append(records, *dao)
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "hash"}},
		DoNothing: true,
	}).Create(&records).Error
}
