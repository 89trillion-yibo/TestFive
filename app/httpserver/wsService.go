package httpserver

import (
	"awesomeProject/Testfive/internal/ws"
	"awesomeProject/Testfive/internal/route"
	"net/http"
)

func InitRun() {
	management := ws.NewClientManagement()
	go management.StartRun()
	route.Routers(management)
	http.ListenAndServe(":8080",nil)
}


