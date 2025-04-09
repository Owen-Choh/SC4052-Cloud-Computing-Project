package user

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/config"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/chatbot/types"
	"github.com/Owen-Choh/SC4052-Cloud-Computing-Assignment-2/chatbot-backend/utils"
)

type UserStore struct {
	store *sql.DB
}

func NewStore(db *sql.DB) types.UserStoreInterface {
	return &UserStore{store: db}
}

func (s *UserStore) GetUserByID(id int) (*types.User, error) {
	rows, err := s.store.Query("SELECT * FROM users WHERE userid = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Userid == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserStore) CreateUser(newUser types.RegisterUserPayload) error {
	username := newUser.Username
	password := newUser.Password
	createdDate, err := utils.GetCurrentTime()
	if err != nil {
		log.Println("unable to obtain formatted date for creating user")
		createdDate = config.Envs.Default_Time // default time if time util fails
	}
	lastLogin := createdDate

	_, dberr := s.store.Exec(
		"INSERT INTO users (username, password, createddate, lastlogin) VALUES (?, ?, ?, ?)",
		username,
		password,
		createdDate,
		lastLogin,
	)

	if dberr != nil {
		log.Println(dberr)
	}

	return nil
}

func (s *UserStore) UpdateUserLastlogin(userid int) error {
	currentTime, err := utils.GetCurrentTime()
	if err != nil {
		log.Println("unable to obtain formatted date for updating user lastlogin")
	}

	_, dberr := s.store.Exec("UPDATE users SET lastlogin=? WHERE userid=?", currentTime, userid)
	if dberr != nil {
		log.Println(dberr)
	}

	return nil
}

func (s *UserStore) GetUserByName(username string) (*types.User, error) {
	rows, err := s.store.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Userid == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.Userid,
		&user.Username,
		&user.Password,
		&user.Createddate,
		&user.Lastlogin,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
