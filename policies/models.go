package policies

import (
	"github.com/jinzhu/gorm"
	"github.com/miknikif/vault-auto-unseal/common"
)

type PolicyModel struct {
	gorm.Model
	Name string `gorm:"column:name,unique_index"`
	Text string `gorm:"column:text,size:4096"`
}

func FindManyPolicies() ([]PolicyModel, int64, error) {
	var models []PolicyModel
	var count int64
	l, err := common.GetLogger()
	if err != nil {
		return models, count, err
	}
	l.Debug("Starting retrieval of the all PolicyModels from the DB")
	db, err := common.GetDB()
	if err != nil {
		return models, count, err
	}
	res := db.Find(&models)
	count = res.RowsAffected
	err = res.Error
	return models, count, err
}

func FindOnePolicy(condition interface{}) (PolicyModel, error) {
	var model PolicyModel
	l, err := common.GetLogger()
	if err != nil {
		return model, err
	}
	l.Debug("Starting retrieval of the PolicyModel from the DB", condition)
	db, err := common.GetDB()
	if err != nil {
		return model, err
	}
	err = db.Where(condition).First(&model).Error
	if err != nil {
		return model, err
	}
	l.Debug("Finished retrieval of the PolicyModel from the DB", condition)
	return model, err
}

func SaveOne(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting saving the PolicyModel to the DB", data)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Save(data).Error
	l.Debug("Finished saving the PolicyModel to the DB", data)
	return err
}

func (p *PolicyModel) Update(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting update of the PolicyModel to the DB", data)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Model(p).Update(data).Error
	l.Debug("Finished update of the PolicyModel to the DB", data)
	return err
}

func DeletePolicyModel(condition interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting delete the PolicyModel from the DB", condition)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Where(condition).Delete(PolicyModel{}).Error
	l.Debug("Finished delete the PolicyModel from the DB", condition)
	return err
}
