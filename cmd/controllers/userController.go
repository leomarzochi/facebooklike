package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/leomarzochi/facebooklike/cmd/auth"
	"github.com/leomarzochi/facebooklike/cmd/crypt"
	"github.com/leomarzochi/facebooklike/cmd/db"
	"github.com/leomarzochi/facebooklike/cmd/helpers"
	"github.com/leomarzochi/facebooklike/cmd/models"
	"github.com/leomarzochi/facebooklike/cmd/repositories"
)

// Create a user on the database
func UserCreate(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := helpers.ReadJSON(w, r, &user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := user.Prepare(models.USER_STATUS_CREATING); err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	ID, err := userRepo.Create(user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user.ID = ID

	payload := helpers.JsonResponse{
		Data: user,
	}

	helpers.WriteJSON(w, http.StatusCreated, payload)
}

// List one user on the database
func UserList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	repo := repositories.NewUserRepository(db)

	user, err := repo.ListOne(userID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, user)
}

// List all users from the database
func UserListAll(w http.ResponseWriter, r *http.Request) {
	nameORUsername := strings.ToLower(r.URL.Query().Get("q"))
	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	users, err := repo.ListAll(nameORUsername)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, users)
}

// Updates a user informations on the database
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user models.User

	err := helpers.ReadJSON(w, r, &user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err = user.Prepare(models.USER_STATUS_EDITING)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	IDfromToken, err := auth.GetIDFromToken(r)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
	}

	if userID != IDfromToken {
		helpers.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	err = repo.Update(userID, user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}

// Delete a user from the database
func UserDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}

	IDfromToken, err := auth.GetIDFromToken(r)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
	}

	if userID != IDfromToken {
		helpers.ErrorJSON(w, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}

	defer db.Close()

	repo := repositories.NewUserRepository(db)

	err = repo.Delete(userID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}

func UserFollow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := auth.GetIDFromToken(r)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	followID, err := strconv.ParseUint(params["followID"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if followID == userID {
		helpers.ErrorJSON(w, errors.New("you cannot follow your self"), http.StatusForbidden)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	newId, err := repo.FollowUser(userID, followID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var response helpers.JsonResponse

	response.Data = newId
	response.Message = "success"

	helpers.WriteJSON(w, http.StatusOK, response)

}

func UserUnfollow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := auth.GetIDFromToken(r)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	followID, err := strconv.ParseUint(params["followID"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if followID == userID {
		helpers.ErrorJSON(w, errors.New("you cannot stop following your self"), http.StatusForbidden)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	err = repo.UnfollowUser(userID, followID)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var response helpers.JsonResponse

	response.Message = "success"

	helpers.WriteJSON(w, http.StatusOK, response)

}

func UserFollows(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	usersFollowed, err := repo.Followed(userId)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var response helpers.JsonResponse

	response.Message = "success"
	response.Data = usersFollowed

	helpers.WriteJSON(w, http.StatusOK, response)
}

func UserFollowers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	followers, err := repo.Followers(userId)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	var response helpers.JsonResponse

	response.Message = "success"
	response.Data = followers

	helpers.WriteJSON(w, http.StatusOK, response)
}

func UserChangePassword(w http.ResponseWriter, r *http.Request) {
	IDFromToken, err := auth.GetIDFromToken(r)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var passwords models.Password

	err = helpers.ReadJSON(w, r, &passwords)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repo := repositories.NewUserRepository(db)

	password, err := repo.FetchUserPassword(IDFromToken)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = crypt.VerifyPassword(password, passwords.Password)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid password"), http.StatusForbidden)
		return
	}

	newPassword, err := crypt.HashPassword(passwords.NewPassword)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = repo.UpdatePassword(IDFromToken, string(newPassword))
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var response helpers.JsonResponse

	response.Message = "success"

	helpers.WriteJSON(w, http.StatusOK, response)
}
