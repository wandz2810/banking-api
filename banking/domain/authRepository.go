package domain

import (
	"encoding/json"
	"fmt"
	"github.com/wandz2810/banking-lib/logger"
	"net/http"
	"net/url"
)

type AuthRepository interface {
	IsAuthorized(token string, routeName string, vars map[string]string) bool
}

type RemoteAuthRepository struct {
}

func (r RemoteAuthRepository) IsAuthorized(token string, routeName string, vars map[string]string) bool {

	u := buildVerifyURL(token, routeName, vars)

	if response, err := http.Get(u); err != nil {
		fmt.Println("Error while sending..." + err.Error())
		return false
	} else {
		m := map[string]bool{}
		if err = json.NewDecoder(response.Body).Decode(&m); err != nil {
			logger.Error("Error while decoding response from auth server:" + err.Error())
			return false
		}
		return m["isAuthorized"]
	}
}

// Sample: http://localhost:8181/auth/verify?token:aaaa.bbbb.cccc&routeName=MakeTransaction&customer_id=2000&account_id=12345
func buildVerifyURL(token string, routeName string, vars map[string]string) string {
	u := url.URL{Host: "host.docker.internal:8181", Path: "/auth/verify", Scheme: "http"}
	//u := url.URL{Host: "localhost:8181", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	q.Add("routeName", routeName)
	for k, v := range vars {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func NewAuthRepository() RemoteAuthRepository {
	return RemoteAuthRepository{}
}
