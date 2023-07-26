package tokens

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/policies"
)

type TokenModel struct {
	gorm.Model
	TokenID                   string `gorm:"unique_index"`
	Accessor                  string `gorm:"unique_index"`
	CreationTime              int
	CreationTTL               int
	DisplayName               string
	EntityID                  string
	ExpireTime                time.Time
	ExplicitMaxTTL            int
	ExternalNamespacePolicies string
	IdentityPolicies          []policies.PolicyModel `gorm:"many2many:token_policy_ip"`
	Meta                      string
	NumUses                   int
	Orphan                    bool
	Path                      string
	Policies                  []policies.PolicyModel `gorm:"many2many:token_policy"`
	Renewable                 bool
	TTL                       int
	Period                    int
	Type                      string
}

const (
	TOKEN_TYPE_SERVICE = "service"
	TOKEN_TYPE_BATCH   = "batch"
	TOKEN_ACCESSOR     = "accessor"
)

const (
	TOKEN_TYPE_SERVICE_LEN = 96
	TOKEN_TYPE_BATCH_LEN   = 128
	TOKEN_ACCESSOR_LEN     = 32
)

var tokenTypeToLen = map[string]int{
	TOKEN_TYPE_SERVICE: TOKEN_TYPE_SERVICE_LEN,
	TOKEN_TYPE_BATCH:   TOKEN_TYPE_BATCH_LEN,
	TOKEN_ACCESSOR:     TOKEN_ACCESSOR_LEN,
}

func FindOneToken(condition interface{}) (TokenModel, error) {
	var model TokenModel
	l, err := common.GetLogger()
	if err != nil {
		return model, err
	}
	l.Debug("Searching for token in the DB: ", "search_condition", condition)
	db, err := common.GetDB()
	if err != nil {
		return model, err
	}
	err = db.Where(condition).Preload("Policies").Preload("IdentityPolicies").First(&model).Error
	if err != nil {
		return model, err
	}
	l.Debug("Token search results: ", "token", model, "err", err)
	return model, err
}

func SaveOne(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting saving the TokenModel to the DB", "token", data)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Save(data).Error
	l.Debug("Finished saving the TokenModel to the DB", "token", data)
	return err
}

func DeletePolicyModel(condition interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting delete the TokenModel from the DB", "token", condition)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Where(condition).Delete(TokenModel{}).Error
	l.Debug("Finished delete the TokenModel from the DB", "token", condition)
	return err
}
