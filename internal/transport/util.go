package transport

import (
	"github.com/CloudyKit/jet/v6"
	"net/http"
	"strings"
)

func Render(w http.ResponseWriter, pageName string, data jet.VarMap) error {
	view, err := views.GetTemplate(pageName)
	if err != nil {
		return err
	}

	err = view.Execute(w, data, nil)
	return err
}

func GetUserList() string {
	var userList []string
	for _, user := range users {
		if user != "" {
			userList = append(userList, user)
		}
	}
	return strings.Join(userList, ",")
}

func usernameExists(username string) bool {
	for _, user := range users {
		if strings.ToLower(user) == username {
			return true
		}
	}
	return false
}
