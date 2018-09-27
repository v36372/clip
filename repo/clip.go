package repo

import (
	"clip/infra"
	"clip/models"

	"github.com/jinzhu/gorm"
)

type clip struct {
	base
}

var Clip IClip

func init() {
	Clip = clip{}
}

type IClip interface {
	GetById(id int) (*models.Clip, error)
	GetLatest(offset, limit int) (clips []models.Clip, err error)
	Create(*models.Clip) (*models.Clip, error)
	Delete(*models.Clip) error
	Update(*models.Clip) error
}

func (c clip) Create(clip *models.Clip) (*models.Clip, error) {
	value, err := c.create(clip)
	return value.(*models.Clip), err
}

func (c clip) Delete(clip *models.Clip) error {
	return c.delete(clip)
}

func (c clip) Update(clip *models.Clip) error {
	return c.save(clip)
}

func (c clip) GetById(id int) (*models.Clip, error) {
	var clip models.Clip
	err := infra.PostgreSql.Model(models.Clip{}).
		Where("id = ?", id).
		Find(&clip).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &clip, err
}

func (c clip) GetLatest(offset, limit int) (clips []models.Clip, err error) {
	err = infra.PostgreSql.Model(models.Clip{}).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&clips).
		Error

	return clips, err
}
