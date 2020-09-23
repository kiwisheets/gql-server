package auth

import (
	"errors"
	"fmt"
	"log"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

func GetTwoFactor(db *gorm.DB, u *model.User) (*model.TwoFactor, error) {
	var twoFactor model.TwoFactor
	if err := db.Model(&u).Association("TwoFactor").Find(&twoFactor); err != nil {
		log.Println(err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if twoFactor.Secret == "" {
		return nil, nil
	}

	return &twoFactor, nil
}

func IsTwoFactorEnabled(db *gorm.DB, u *model.User) (bool, error) {
	if db.Model(&u).Association("TwoFactor").Count() != 0 {
		return true, nil
	}
	return false, nil
}

// EnableTwoFactor will enable 2FA for the user passed to it if the token validates against the secret
func EnableTwoFactor(db *gorm.DB, u *model.User, secret string, token string) ([]string, error) {
	{
		twoFactor, _ := GetTwoFactor(db, u)
		if twoFactor != nil {
			return nil, fmt.Errorf("2FA is already enabled")
		}
	}

	if !totp.Validate(token, secret) {
		return nil, fmt.Errorf("Invalid 2FA code")
	}
	twoFactor := model.TwoFactor{
		UserID:     u.ID,
		Secret:     secret,
		BackupKeys: generateBackupKeys(),
	}
	if err := db.Create(&twoFactor).Error; err != nil {
		return nil, err
	}
	return twoFactor.BackupKeys, nil
}

func DisableTwoFactor(db *gorm.DB, u *model.User, password string) (bool, error) {
	if !VerifyPassword(u, password) {
		return false, fmt.Errorf("Password incorrect")
	}

	twoFactor, err := GetTwoFactor(db, u)

	if twoFactor == nil || err != nil {
		return true, err
	}

	if err := db.Delete(&twoFactor).Error; err != nil {
		return false, err
	}

	return true, nil
}

func VerifyTwoFactor(t *model.TwoFactor, token string) bool {
	// compare token to backup keys
	for _, k := range t.BackupKeys {
		if token == k {
			return true
		}
	}
	// validate against secret
	return totp.Validate(token, t.Secret)
}

func GetBackupKeys(db *gorm.DB, authCtx authContext) ([]string, error) {
	if !authCtx.Secure {
		return nil, fmt.Errorf("Login required")
	}
	var twoFactor model.TwoFactor
	if err := db.Model(&authCtx.User).Association("TwoFactor").Find(&twoFactor); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("2FA is not enabled")
		}
		return nil, err
	}
	return twoFactor.BackupKeys, nil
}

func generateBackupKeys() []string {
	keys := []string{}
	for i := 0; i < 10; i++ {
		keys = append(keys, util.RandString(10))
	}
	return keys
}
