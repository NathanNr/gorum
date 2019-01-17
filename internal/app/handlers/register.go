package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ltheinrich/captcha"

	"github.com/ltheinrich/gorum/internal/pkg/config"
	"github.com/ltheinrich/gorum/internal/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

// Register handler
func Register(request map[string]interface{}, username string, auth bool) interface{} {
	var err error

	// get strings from request
	password := GetString(request, "password")
	cap := GetString(request, "captcha")
	capVal := GetString(request, "captchaValue")

	// check if username and password are provided
	if username == "" || password == "" || len(username) > 32 {
		// return not provided
		return errors.New("400")
	}

	// verify captcha
	if config.Get("https", "captcha") == TRUE && !captcha.VerifyString(cap, capVal) {
		// invalid captcha
		return errors.New("403 captcha")
	}

	// query db
	var id int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username = $1;", username).Scan(&id)
	if err == sql.ErrNoRows {
		// not exists
		var passwordHash []byte
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+1)
		if err != nil {
			// return error
			return err
		}

		// insert into database
		_, err = db.DB.Exec("INSERT INTO users (username, passwordhash, registered) VALUES ($1, $2, $3);",
			username, string(passwordHash), time.Now().Format("2006-01-02T15:04:05"))
		if err != nil {
			// return error
			return err
		}

		// registered
		return map[string]interface{}{"done": true}
	} else if err != nil {
		// return error
		return err
	}

	// username exists
	return map[string]interface{}{"done": false}
}
