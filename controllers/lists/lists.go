package lists

import (
	"net/http"
	"net/url"
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"encoding/json"
)

// Handler is the handler for the lists endpoints
func Handler(w http.ResponseWriter, r *http.Request, store stores.Store) error {
	listID := getListIDFromURL(r.URL)

	switch r.Method {
	case http.MethodGet:
		if listID == "" {
			r, err := store.GetLists()
			if err != nil {
				return err
			}
			controllers.WriteOkResponse(w, http.StatusOK, r)
		} else {
			l, err := store.GetSingleList(listID)
			if err != nil {
				return err
			}
			controllers.WriteOkResponse(w, http.StatusOK, l)
		}
	case http.MethodPost:
		l, err := parseListBody(r)
		if err != nil {
			return err
		}
		err = store.AddList(&l)
		if err != nil {
			return err
		}
		controllers.WriteOkResponse(w, http.StatusCreated, l)
	case http.MethodDelete:
		err := store.RemoveList(listID)
		if err != nil {
			return err
		}
		controllers.WriteOkResponse(w, http.StatusNoContent, nil)
	case http.MethodPut:
		l, err := parseListBody(r)
		if err != nil {
			return err
		}
		err = store.UpdateList(listID, &l)
		if err != nil {
			return err
		}
		controllers.WriteOkResponse(w, http.StatusOK, l)
	default:
		controllers.WriteOkResponse(w, http.StatusMethodNotAllowed, nil)
	}

	return nil
}

func getListIDFromURL(u *url.URL) string {
	var listID string

	if len(u.Path) > len("/lists") {
		listID = u.Path[len("/lists/"):]
	}

	return listID
}

func parseListBody(r *http.Request) (models.List, error) {
	if r.Body == nil {
		return models.List{}, &controllers.NoBodyError{}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.List{}, &controllers.InvalidBodyError{InternalError: err}
	}

	l := dto.ToList()

	return l, nil
}