//authentication middleware


package main


import (
	"fmt"
	"net/http"
)

func  AuthMid(next http.Handler) http.Handler{

	return  http.HandlerFunc( func( w http.RsponseWriter, r *http.Request)){

		apiKey := r.Header.Get("X-API-Key")

		if
	}
}

