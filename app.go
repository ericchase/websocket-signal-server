package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var (
	cookie_name             = "cookie-name"
	AUTHENTICATION_KEY_1, _ = hex.DecodeString(os.Getenv("AUTHENTICATION_KEY_1"))
	ENCRYPTION_KEY_1, _     = hex.DecodeString(os.Getenv("ENCRYPTION_KEY_1"))
	store                   = sessions.NewCookieStore(AUTHENTICATION_KEY_1, ENCRYPTION_KEY_1)
)

func main() {
	store.Options = &sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/login", login)
	http.HandleFunc("/secret", secret)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/echo", echo)
	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
}

func login(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	session, err_get := store.Get(r, cookie_name)
	if err_get != nil {
		log.Println(err_get)
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
		return
	}

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	err_save := session.Save(r, w)
	if err_save != nil {
		log.Println(err_save)
		fmt.Fprintln(w, "LOGIN FAILURE")
		return
	}

	fmt.Fprintln(w, "LOGIN SUCCESS")
}

func secret(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	session, err_get := store.Get(r, cookie_name)
	if err_get != nil {
		log.Println(err_get)
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func echo(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	session, err_get := store.Get(r, cookie_name)
	if err_get != nil {
		log.Println(err_get)
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
		return
	}

	ws_connection, err_upgrade := wsUpgrader.Upgrade(w, r, nil)
	if err_upgrade != nil {
		log.Print("ws upgrade:", err_upgrade)
		return
	}
	defer ws_connection.Close()
	for {
		message_type, message, err_read := ws_connection.ReadMessage()
		if err_read != nil {
			log.Println("echo r:", err_read)
			break
		}
		log.Printf("recv: %s", message)
		err_read = ws_connection.WriteMessage(message_type, message)
		if err_read != nil {
			log.Println("echo w:", err_read)
			break
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	enableCors(&w, r)
	session, err_get := store.Get(r, cookie_name)
	if err_get != nil {
		log.Println(err_get)
	}
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
		return
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	fmt.Fprintln(w, "LOGOUT SUCCESS")
}
