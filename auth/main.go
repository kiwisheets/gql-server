package internalauth

import (
	"fmt"
	"log"
	"strings"
	"time"

	argonpass "github.com/dwin/goArgonPass"
	"github.com/emvi/hide"
	"github.com/kiwisheets/auth/permission"
	"github.com/kiwisheets/auth/token"
	"github.com/kiwisheets/gql-server/orm/model"
	"github.com/kiwisheets/util"
	"gorm.io/gorm"
)

// LoginUser generates a signed JWT token
func LoginUser(user *model.User, permissions []*permission.Permission, cfg *util.JWTConfig) (string, error) {
	return login(user, permissions, cfg, 0)
}

func LoginUserSecure(user *model.User, permissions []*permission.Permission, cfg *util.JWTConfig) (string, error) {
	return login(user, permissions, cfg, 1*time.Hour)
}

func login(user *model.User, permissions []*permission.Permission, cfg *util.JWTConfig, expires time.Duration) (string, error) {
	// put into queue and wait for queue to finish
	// this is to prevent OOM errors

	return token.BuildAndSignToken(token.UserTokenParams{
		ID:          user.ID,
		CompanyID:   user.CompanyID,
		Email:       user.Email,
		Permissions: permissions,
	}, cfg, expires)
}

// VerifyPassword verifies a password against the stored hash
func VerifyPassword(db *gorm.DB, userID hide.ID, password string) bool {
	var user model.User
	if err := db.Model(&model.User{}).Where(userID).First(&user).Error; err != nil {
		return false
	}

	start := time.Now()

	err := argonpass.Verify(password, user.Password)

	elapsed := time.Since(start)
	log.Printf("Password hash verify took %s", elapsed)

	return err == nil
}

// HashPassword attempts to hash the supplied password
func HashPassword(password string) (string, error) {
	// debug check time

	start := time.Now()

	hash, err := argonpass.Hash(password, argonpass.ArgonParams{
		Time:        15,
		Memory:      48 * 1024,
		Parallelism: 2,
		OutputSize:  1,
		Function:    "argon2id",
		SaltSize:    8,
	})

	elapsed := time.Since(start)
	log.Printf("Password hash took %s", elapsed)

	return hash, err
}

func splitToken(header string) (string, error) {
	splitToken := strings.Split(header, "Bearer")

	if len(splitToken) != 2 || len(splitToken[1]) < 2 {
		return "", fmt.Errorf("bad token format")
	}

	return strings.TrimSpace(splitToken[1]), nil
}
