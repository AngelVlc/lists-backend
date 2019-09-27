package users

import (
	"net/http"
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"encoding/json"
)

// Handler is the handler for the users endpoints
func Handler(w http.ResponseWriter, r *http.Request, store stores.Store) error {
	switch r.Method {
	case http.MethodPost:
		u, err := parseUserBody(r)
		if err != nil {
			return err
		}
		err = store.AddUser(&u)
		if err != nil {
			return err
		}
		controllers.WriteOkResponse(w, http.StatusCreated, u)
	default:
		controllers.WriteOkResponse(w, http.StatusMethodNotAllowed, nil)
	}

	return nil
}

func parseUserBody(r *http.Request) (models.User, error) {
	if r.Body == nil {
		return models.User{}, &controllers.NoBodyError{}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.User{}, &controllers.InvalidBodyError{InternalError: err}
	}

	l := dto.ToUser()

	return l, nil
}