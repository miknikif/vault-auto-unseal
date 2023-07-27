package tokens

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/miknikif/vault-auto-unseal/common"
	"github.com/miknikif/vault-auto-unseal/policies"
)

type TokenModel struct {
	gorm.Model
	TokenID                   string `gorm:"unique_index"`
	Accessor                  string `gorm:"unique_index"`
	CreationTime              time.Time
	CreationTTL               int
	DisplayName               string
	EntityID                  string
	ExpireTime                time.Time
	ExplicitMaxTTL            int
	LastRenewalTime           time.Time
	ExternalNamespacePolicies string
	IdentityPolicies          []policies.PolicyModel `gorm:"many2many:token_policy_ip"`
	Meta                      *string
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
	// Len in bytes
	TOKEN_TYPE_SERVICE_LEN = 70
	TOKEN_TYPE_BATCH_LEN   = 92
	TOKEN_ACCESSOR_LEN     = 20
)

var tokenTypeToLen = map[string]int{
	TOKEN_TYPE_SERVICE: TOKEN_TYPE_SERVICE_LEN,
	TOKEN_TYPE_BATCH:   TOKEN_TYPE_BATCH_LEN,
	TOKEN_ACCESSOR:     TOKEN_ACCESSOR_LEN,
}

func (p *TokenModel) Update(data interface{}) error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting update of the TokenModel", "token", data)
	db, err := common.GetDB()
	if err != nil {
		return err
	}
	err = db.Model(p).Update(data).Error
	l.Debug("Finished update of the TokenModel", "token", data)
	return err
}

func (s *TokenModel) Renew() error {
	l, err := common.GetLogger()
	if err != nil {
		return err
	}
	l.Debug("Starting renew of the TokenModel", "token", s)
	db, err := common.GetDB()
	if err != nil {
		return err
	}

	if s.CreationTTL == 0 {
		return nil
	}

	if s.ExplicitMaxTTL == 0 && s.Period == 0 {
		return nil
	}

	endTime := s.CreationTime.Add(time.Second * time.Duration(s.ExplicitMaxTTL))

	remainingTTL := s.ExpireTime.Unix() - time.Now().Unix()
	// Check that token is not expiried yet
	if remainingTTL < 0 {
		if err := DeleteTokenModel(s); err != nil {
			return fmt.Errorf("error occurred during token removal")
		}
		return fmt.Errorf("the token is expired")
	}

	l.Trace("Token data", "creationTTL", s.CreationTTL, "endTime", endTime, "explicitMaxTTL", s.ExplicitMaxTTL, "period", s.Period)

	if s.ExplicitMaxTTL > 0 && s.Period == 0 {
		l.Trace("Handling usual token")
		if int(time.Now().Unix())+s.CreationTTL < int(endTime.Unix()) {
			s.LastRenewalTime = time.Now()
			s.ExpireTime = time.Now().Add(time.Second * time.Duration(s.CreationTTL))
		} else if time.Now().Unix() < endTime.Unix() {
			s.LastRenewalTime = time.Now()
			s.ExpireTime = endTime
		}
	} else {
		l.Trace("Handling periodic token")
		s.LastRenewalTime = time.Now()
		s.ExpireTime = time.Now().Add(time.Second * time.Duration(s.CreationTTL))
	}

	err = db.Model(s).Update(s).Error
	l.Debug("Finished renew of the TokenModel", "expireTime", s.ExpireTime, "token", s)
	return err

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

func DeleteTokenModel(condition interface{}) error {
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

func SeedDB(c *common.Config) error {
	rootToken := NewRootToken()
	if err := SaveOne(rootToken); err != nil {
		return err
	}
	c.Logger.Warn("New root token was created. Please save it in the safe place!", "tokenID", rootToken.TokenID, "accessor", rootToken.Accessor)
	return nil
}
