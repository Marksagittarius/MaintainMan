package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetDivisionByID(id uint) (*model.Division, error) {
	return TxGetDivisionByID(database.DB, id)
}

func TxGetDivisionByID(tx *gorm.DB, id uint) (*model.Division, error) {
	division := &model.Division{}
	if err := tx.Preload("Children").First(division, id).Error; err != nil {
		logger.Logger.Debugf("GetDivisionByIDErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func CreateDivision(aul *model.CreateDivisionRequest) (*model.Division, error) {
	return TxCreateDivision(database.DB, aul)
}

func TxCreateDivision(tx *gorm.DB, aul *model.CreateDivisionRequest) (*model.Division, error) {
	division := &model.Division{}
	copier.Copy(division, aul)
	if err := tx.Create(division).Error; err != nil {
		logger.Logger.Debugf("CreateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func UpdateDivision(id uint, aul *model.UpdateDivisionRequest) (*model.Division, error) {
	return TxUpdateDivision(database.DB, id, aul)
}

func TxUpdateDivision(tx *gorm.DB, id uint, aul *model.UpdateDivisionRequest) (*model.Division, error) {
	where := &model.Division{}
	where.ID = id
	division := &model.Division{}
	copier.Copy(division, aul)
	if err := tx.Model(where).Where(where).Updates(division).Error; err != nil {
		logger.Logger.Debugf("UpdateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func DeleteDivision(id uint) error {
	return TxDeleteDivision(database.DB, id)
}

func TxDeleteDivision(tx *gorm.DB, id uint) (err error) {
	if err = tx.Delete(&model.Division{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteDivisionErr: %v\n", err)
	}
	return
}