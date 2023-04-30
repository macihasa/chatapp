package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/macihasa/chatapp/backend/pkg/models"
)

type httpErr struct {
	msg string
	err error
}

func (e httpErr) handleError(w http.ResponseWriter) {
	log.Println(e.msg, e.err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, e.msg+"%v", e.err)
}

func Landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is live!", r.RemoteAddr)
}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	var user models.UserModel

	defer func() {
		err := r.Body.Close()
		if err != nil {
			e := httpErr{"Couldnt close body", err}
			e.handleError(w)
		}
	}()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		e := httpErr{"Failed to read request body", err}
		e.handleError(w)
		return
	}

	err = json.Unmarshal(bs, &user)
	if err != nil {
		e := httpErr{"Failed to unmarshal: ", err}
		e.handleError(w)
		return
	}
	id, err := user.Register()
	if err != nil {
		e := httpErr{"Couldnt register user", err}
		e.handleError(w)
		return
	}

	w.Write([]byte(fmt.Sprintln("New user registered: ", id, user)))

}
