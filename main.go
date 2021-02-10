package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stytchauth/stytch-go/stytch"
)

var client = stytch.NewClient(
	stytch.EnvTest,
	os.Getenv("STYTCH_PROJECT_ID"),
	os.Getenv("STYTCH_SECRET"),
)

var (
	address      = os.Getenv("ADDRESS")
	fullAddress  = "http://" + address
	magicLinkURL = fullAddress + "/authenticate"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homepage).Methods("GET")
	r.HandleFunc("/login_or_create_user", loginOrCreateUser).Methods("POST")
	r.HandleFunc("/authenticate", authenticate).Methods("GET")
	r.HandleFunc("/logout", logout).Methods("GET")

	// Declare the static file directory
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler)

	log.Fatal(http.ListenAndServe(address, r))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	templateData := struct {
		LoginOrCreateUserPath string
	}{
		LoginOrCreateUserPath: fullAddress + "/login_or_create_user",
	}

	parseAndExecuteTemplate("templates/loginOrSignUp.html", templateData, w)
}

func loginOrCreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := client.LoginOrCreateUser(&stytch.LoginOrCreateUser{
		Email:              r.FormValue("email"),
		LoginMagicLinkURL:  magicLinkURL,
		SignUpMagicLinkURL: magicLinkURL,
	})
	if err != nil {
		log.Printf("something went wrong sending magic link: %s\n", err)
	}

	parseAndExecuteTemplate("templates/emailSent.html", nil, w)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	_, err := client.AuthenticateMagicLink(r.URL.Query().Get("token"), nil)
	if err != nil {
		log.Printf("something went wrong authenticating the magic link: %s\n", err)
	}

	templateData := struct {
		LoggedOutPath string
	}{
		LoggedOutPath: fullAddress + "/logout",
	}

	parseAndExecuteTemplate("templates/loggedIn.html", templateData, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	parseAndExecuteTemplate("templates/loggedOut.html", nil, w)
}

func parseAndExecuteTemplate(temp string, tempData interface{}, w http.ResponseWriter) {
	t, err := template.ParseFiles(temp)
	if err != nil {
		log.Printf("something went wrong parsing template: %s\n", err)
	}

	err = t.Execute(w, tempData)
	if err != nil {
		log.Printf("something went wrong executing the template: %s\n", err)
	}
}
