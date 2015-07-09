package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func websocketHandle(res http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{HandshakeTimeout: 10 * time.Second, ReadBufferSize: 10 * 1024, WriteBufferSize: 10 * 1024}
	if socket, err := upgrader.Upgrade(res, req, nil); err != nil {
		log.Println("upgrate error", err)
		return
	} else {
		log.Print("new websocket connection")
		socket.SetReadLimit(1024)
		socket.SetReadDeadline(time.Now().Add(5 * time.Minute))
		socket.SetPongHandler(func(string) error { socket.SetReadDeadline(time.Now().Add(5 * time.Minute)); return nil })
		for {
			if messagetype, message, err := socket.ReadMessage(); err == nil {
				log.Print(messagetype, message, err)
				socket.WriteMessage(websocket.BinaryMessage, message)
			} else {
				log.Print("connection closed, err:", err)
				break
			}

		}
	}
}
func handleGoogleLogin(res http.ResponseWriter, req *http.Request) {
	status, code := req.FormValue("state"), req.FormValue("code")
	if status == "login" && code != "" {
		log.Print("grante privilege success, setting cookie")
		cookie := &http.Cookie{Expires: time.Now().AddDate(0, 0, 1)}
		cookie.Name, cookie.Value = "oauthcode", code
		http.SetCookie(res, cookie)
		cookie.Name, cookie.Value = "oauthparty", "google"
		http.SetCookie(res, cookie)
		http.Redirect(res, req, "index.html", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(res, req, "loginfail.html", http.StatusTemporaryRedirect)
	}
}
func main() {
	log.Print("starting")
	http.HandleFunc("/ws", websocketHandle)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/oauth2callback", handleGoogleLogin)

	http.ListenAndServe(":8088", nil)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)
}
