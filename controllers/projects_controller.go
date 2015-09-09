package controllers
import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/zenazn/goji/web"
	projectModel "../models/project"
)

type Projects struct {
}
func (u *Projects)Create(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	current_user, result := LoginRequired(c, w, r)
	if result {
		fmt.Printf("current_user: %+v\n", *current_user)
	}
	project := projectModel.NewProject("hoge")
	encoder := json.NewEncoder(w)
	encoder.Encode(*project)
}
