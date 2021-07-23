package route

import (
	"awesomeProject/Testfive/internal/ctrl"
	"awesomeProject/Testfive/internal/ws"
	"net/http"
)

func Routers(cm *ws.ClientManagement) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ctrl.ServeHttp(cm,w,r)
	})
}
