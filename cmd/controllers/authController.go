package controllers

import (
	"errors"
	"net/http"

	"github.com/leomarzochi/facebooklike/cmd/auth"
	"github.com/leomarzochi/facebooklike/cmd/crypt"
	"github.com/leomarzochi/facebooklike/cmd/db"
	"github.com/leomarzochi/facebooklike/cmd/helpers"
	"github.com/leomarzochi/facebooklike/cmd/models"
	"github.com/leomarzochi/facebooklike/cmd/repositories"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//Recebe o email e senha
	var user models.User

	err := helpers.ReadJSON(w, r, &user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	//conecta com o banco

	db, err := db.ConnectDB()
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// instancia o repositorio

	repo := repositories.NewUserRepository(db)

	//Busca no bd se o email existe, se existe retorna junto a senha
	userOnDB, err := repo.GetUserByEmail(user.Email)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	//Compara se a senha é correta

	err = crypt.VerifyPassword(userOnDB.Password, user.Password)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("invalid email or password"), http.StatusUnauthorized)
		return
	}

	//Retorna token
	token, err := auth.CreateToken(uint(userOnDB.ID))
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", token)

	helpers.WriteJSON(w, http.StatusNoContent, nil)
}
