package router

import (
	"avito/controllers"
	mdw "avito/middlewares"
	"avito/responses"
	"avito/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRoutes() *mux.Router {
	Router := mux.NewRouter()

	Router.HandleFunc("/subscribe", mdw.SetJSONmdw(controllers.Subscribe)).Methods("POST")
	Router.HandleFunc("/unsubscribe", mdw.SetJSONmdw(controllers.Unsubscribe)).Methods("POST")
	Router.HandleFunc("/get_subscribes", mdw.SetJSONmdw(controllers.GetSubscribes)).Methods("GET")
	Router.HandleFunc("/confirm_email", mdw.SetJSONmdw(controllers.EmailConfirmation)).Methods("GET")
	Router.HandleFunc("/manual_check_price", mdw.SetJSONmdw(controllers.ManualCheckPrice)).Methods("GET")

	Router.NotFoundHandler = http.HandlerFunc(notFound)
	return Router
}

func notFound(response http.ResponseWriter, request *http.Request) {
	utils.LogRequest(request)
	responses.Status(response, 404)
}
