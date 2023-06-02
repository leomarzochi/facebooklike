package repositories

import (
	"database/sql"
	"fmt"

	"github.com/leomarzochi/facebooklike/cmd/models"
)

type users struct {
	*sql.DB
}

func NewUserRepository(DB *sql.DB) *users {
	return &users{DB}
}

func (u *users) Create(user models.User) (uint64, error) {
	statement, err := u.DB.Prepare(
		"INSERT INTO users(nome, username, email, password) values (?, ?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	r, err := statement.Exec(user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil

}

func (u *users) ListOne(id uint64) (models.User, error) {
	rows, err := u.DB.Query(
		"SELECT id, nome, username, email, createdAt FROM users WHERE id = ?",
		id,
	)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User

	if rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.CreateAt,
		)
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil

}

func (u *users) FetchUserPassword(id uint64) (string, error) {
	rows, err := u.DB.Query(
		"SELECT password FROM users WHERE id = ?",
		id,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var user models.User

	if rows.Next() {
		err = rows.Scan(
			&user.Password,
		)
	}

	if err != nil {
		return "", err
	}

	return user.Password, nil

}

func (u *users) ListAll(nameORUsername string) ([]models.User, error) {
	nameORUsername = fmt.Sprintf("%%%s%%", nameORUsername)

	rows, err := u.DB.Query(
		"SELECT id, nome, username, email, createdAt FROM users WHERE nome LIKE ? OR username LIKE ?",
		nameORUsername, nameORUsername,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var returnedUsers []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.CreateAt,
		)
		if err != nil {
			return nil, err
		}

		returnedUsers = append(returnedUsers, user)
	}

	return returnedUsers, nil
}

func (u *users) Update(id uint64, user models.User) error {
	statement, err := u.DB.Prepare("UPDATE users SET nome = ?, email = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Email, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *users) UpdatePassword(userID uint64, newPassword string) error {
	statement, err := u.DB.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(newPassword, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *users) Delete(id uint64) error {
	statement, err := u.DB.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *users) GetUserByEmail(email string) (models.User, error) {
	rows, err := u.DB.Query("SELECT id, password FROM users WHERE email = ?", email)
	if err != nil {
		return models.User{}, err
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Password,
		); err != nil {
			return models.User{}, nil
		}
	}

	return user, nil

}

func (u *users) FollowUser(userID, followID uint64) (uint64, error) {
	statement, err := u.DB.Prepare(
		"INSERT IGNORE INTO followers(user_id, follower_id) VALUES (?, ?)",
	)
	if err != nil {
		return 0, nil
	}

	result, err := statement.Exec(followID, userID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (u *users) UnfollowUser(userID, unfollowID uint64) error {
	statement, err := u.DB.Prepare(
		"DELETE FROM followers WHERE user_id = ? AND follower_id = ?",
	)
	if err != nil {
		return err
	}

	result, err := statement.Exec(unfollowID, userID)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (u *users) Followed(id uint64) ([]models.User, error) {
	rows, err := u.DB.Query(`
		SELECT u.nome, u.username, u.createdAt FROM users u
		INNER JOIN followers f ON u.id = f.user_id
		WHERE f.follower_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.Name,
			&user.Username,
			&user.CreateAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

func (u *users) Followers(id uint64) ([]models.User, error) {
	rows, err := u.DB.Query(`
		SELECT u.nome, u.username, u.createdAt FROM users u
		INNER JOIN followers f ON u.id = f.follower_id
		WHERE f.user_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.Name,
			&user.Username,
			&user.CreateAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}
