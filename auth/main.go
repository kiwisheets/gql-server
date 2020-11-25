package auth

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/util"

	"git.maxtroughear.dev/max.troughear/digital-timesheet/go-server/orm/model"
	argonpass "github.com/dwin/goArgonPass"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type AuthContext struct {
	User   *model.User
	Secure bool
}

// LoginUser generates a signed JWT token
func LoginUser(user *model.User, cfg *util.JWTConfig) (string, error) {
	return login(user, cfg, 0)
}

func LoginUserSecure(user *model.User, cfg *util.JWTConfig) (string, error) {
	return login(user, cfg, 1*time.Hour)
}

func login(user *model.User, cfg *util.JWTConfig, expires time.Duration) (string, error) {
	// put into queue and wait for queue to finish
	// this is to prevent OOM errors

	return buildAndSignToken(user, cfg, expires)
}

// VerifyPassword verifies a password against the stored hash
func VerifyPassword(user *model.User, password string) bool {
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

// Middleware decodes the authorization header
func Middleware(db *gorm.DB, cfg *util.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		a := c.Request.Header.Get("authorization")
		aSecure := c.Request.Header.Get("authorizationsecure")

		// allow unauthenticated users through
		if a == "" {
			c.Next()
			return
		}

		token, err := splitToken(a)

		if err != nil {
			c.Next()
			return
		}

		// process and validate jwt token
		userID, err := validateTokenAndGetUserID(token, cfg)
		if err != nil {
			c.Next()
			return
		}
		// get user from database

		var user model.User
		err = db.Model(&model.User{}).Where(int64(userID)).First(&user).Error
		if err != nil || user.ID == 0 {
			c.Next()
			return
		}

		// check if secure token exists and is valid
		secure := false
		tokenSecure, err := splitToken(aSecure)
		if err == nil {
			userIDSecure, _ := validateTokenAndGetUserID(tokenSecure, cfg)
			if userIDSecure == user.ID {
				secure = true
			}
		}

		ctx := context.WithValue(c.Request.Context(), userCtxKey, AuthContext{
			User:   &user,
			Secure: secure,
		})

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// For find the user from the context. Middleware must have run
func For(ctx context.Context) AuthContext {
	raw, _ := ctx.Value(userCtxKey).(AuthContext)
	return raw
}

func splitToken(header string) (string, error) {
	splitToken := strings.Split(header, "Bearer")

	if len(splitToken) != 2 || len(splitToken[1]) < 2 {
		return "", fmt.Errorf("bad token format")
	}

	return strings.TrimSpace(splitToken[1]), nil
}
