package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stytchauth/stytch-go/stytch"
)

// define the stytch client using your stytch project id & secret
// use stytch.EnvLive if you want to hit the live api
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

// struct to hold the values to be passed to the html templates
type templateVariables struct {
	LoginOrCreateUserPath string
	LoggedOutPath         string
}

func main() {
	r := mux.NewRouter()
	// routes
	r.HandleFunc("/", homepage).Methods("GET")
	r.HandleFunc("/login_or_create_user", loginOrCreateUser).Methods("POST")
	r.HandleFunc("/authenticate", authenticate).Methods("GET")
	r.HandleFunc("/logout", logout).Methods("GET")

	// Declare the static file directory
	// this is to ensure our static assets & css are accessible & rendered
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler)

	log.Fatal(http.ListenAndServe(address, r))
}

// handles the homepage for Hello Socks
func homepage(w http.ResponseWriter, r *http.Request) {
	parseAndExecuteTemplate(
		"templates/loginOrSignUp.html",
		&templateVariables{LoginOrCreateUserPath: fullAddress + "/login_or_create_user"},
		w,
	)
}

// takes the email entered on the homepage and hits the stytch
// loginOrCreateUser endpoint to send the user a magic link
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

// this is the endpoint the link in the magic link hits takes the token from the
// link's query params and hits the stytch authenticate endpoint to verify the token is valid
func authenticate(w http.ResponseWriter, r *http.Request) {
	_, err := client.AuthenticateMagicLink(r.URL.Query().Get("token"), nil)
	if err != nil {
		log.Printf("something went wrong authenticating the magic link: %s\n", err)
	}

	parseAndExecuteTemplate(
		"templates/loggedIn.html",
		&templateVariables{LoggedOutPath: fullAddress + "/logout"},
		w,
	)
}

// handles the logout endpoint
func logout(w http.ResponseWriter, r *http.Request) {
	parseAndExecuteTemplate("templates/loggedOut.html", nil, w)
}

// helper function to parse the template & render it with any provided data
func parseAndExecuteTemplate(temp string, templateVars *templateVariables, w http.ResponseWriter) {
	t, err := template.ParseFiles(temp)
	if err != nil {
		log.Printf("something went wrong parsing template: %s\n", err)
	}

	err = t.Execute(w, templateVars)
	if err != nil {
		log.Printf("something went wrong executing the template: %s\n", err)
	}
}
