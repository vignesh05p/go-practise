//write s simple logging middleware that prints the requests method and URL before passing the rquest to the next handler

package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	
)

func loggingMiddelware(next http.Handler) http.Handler{
    
    return http.HandlerFunc(func(w  http.ResponseWriter, r *http.Request){
         log.Printf("Request received : %s %s", r.Method, r.URL.PATH)
         next .serverHTTP(w,r)
    
    })
}


func homeHandler(w http.ResponseWriter, r *http.Rquest){
    fmt.Fprintln(w. "hello, world! THis is th main handler,")
}


func main(){
    finalHanler :=loggingMiddleware(http.HandlerFunc(homeHandler))
    http.Handle("/", finalHandler)
    
    
    fmt.Println("server is running on ssserver")
    log.Fatal(http.ListenAndServer("8080",nil))
}

