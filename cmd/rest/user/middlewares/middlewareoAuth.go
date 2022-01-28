package user

import (
	"fmt"
	exec "golang.org/x/sys/execabs"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		token = strings.Replace(token, "Bearer ", "", -1)
		fmt.Println("token received in request-->", token)

		if !validateToken(token) {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusCreated)
			jsonData := []byte(`{"status":"FAIL","Error":"Invalid Token"}`)
			rw.Write(jsonData)
		} else {
			next.ServeHTTP(rw, r)
		}

	})
}

func validateToken(stringToken string) bool {

	//Cobra CLI app
	out, err := exec.Command("tkn-gb", "valTkn", stringToken).Output()
	if err != nil {
		fmt.Println("could not validate token, err:", err)
		return false
	}
	output := string(out)
	fmt.Println("token output-->", output)
	if strings.Contains(output, "invalid token") {
		return false
	}
	return true

}
