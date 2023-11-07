package services

import (
	"flag"
	"log"

	"net/http"
	"os"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/services/handlers/auth"
	"social-network/pkg/services/handlers/groups"
	"social-network/pkg/services/handlers/notifications"
	"social-network/pkg/services/handlers/posts"
	"social-network/pkg/services/handlers/users"
)

func Server() {
	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes(),
	}

	sqlite.DataBase()
	defer sqlite.DB.Close()

	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	mux.HandleFunc("/post/", posts.Post) // info about multiple routes inside --> /post/{id} /post/{id}/comment /post/{id}/like /post/{id}/dislike
	mux.HandleFunc("/post/all", posts.Posts)
	mux.HandleFunc("/post/create", posts.CreatePost)
	mux.HandleFunc("/categories", posts.Categories)
	mux.HandleFunc("/category/", posts.PostsByCategoryId)
	mux.HandleFunc("/signup", auth.SignUp)
	mux.HandleFunc("/signin", auth.SignIn)
	mux.HandleFunc("/signout", auth.SignOut)
	mux.HandleFunc("/me", users.GetMyProfile)
	mux.HandleFunc("/me/settings/", users.AccountSettings)
	mux.HandleFunc("/user/", users.GetProfile)
	mux.HandleFunc("/user/all", users.AllUsers)
	mux.HandleFunc("/changeavatar", users.ChangeAccountAvatar)
	mux.HandleFunc("/groups", groups.ShowGroupList)
	mux.HandleFunc("/group/", groups.Group) // info about multiple routes inside
	mux.HandleFunc("/group/create", groups.CreateGroup)
	mux.HandleFunc("/group/post/", groups.GroupPost)
	mux.HandleFunc("/notifications", notifications.UserNotifications)
	mux.HandleFunc("/notifications/clear", notifications.ClearAllGroupPostCommentsNotifications)

	wsServer := NewWebsocketServer()
	go wsServer.Run()
	mux.HandleFunc("/chat/", UserMessages)
	mux.HandleFunc("/group-chat/", GetGroupChatMessages)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	return mux
}
