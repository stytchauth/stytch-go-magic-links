package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stytchauth/stytch-go/v5/stytch"
	"github.com/stytchauth/stytch-go/v5/stytch/stytchapi"
)

type config struct {
	address      string
	fullAddress  string
	magicLinkURL string
	stytchClient *stytchapi.API
}

// struct to hold the values to be passed to the html templates
type templateVariables struct {
	LoginOrCreateUserPath string
	LoggedOutPath         string
}

func main() {
	// Load .env & set config
	c, err := initializeConfig()
	if err != nil {
		log.Fatal("error initializing config")
	}

	r := mux.NewRouter()
	// routes
	r.HandleFunc("/", c.homepage).Methods("GET")
	r.HandleFunc("/login_or_create_user", c.loginOrCreateUser).Methods("POST")
	r.HandleFunc("/authenticate", c.authenticate).Methods("GET")
	r.HandleFunc("/logout", c.logout).Methods("GET")

	// Declare the static file directory
	// this is to ensure our static assets & css are accessible & rendered
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler)

	log.Fatal(http.ListenAndServe(c.address, r))
}

// handles the homepage for Hello Socks
func (c *config) homepage(w http.ResponseWriter, r *http.Request) {
	parseAndExecuteTemplate(
		"templates/loginOrSignUp.html",
		&templateVariables{LoginOrCreateUserPath: c.fullAddress + "/login_or_create_user"},
		w,
	)
}

// takes the email entered on the homepage and hits the stytch
// loginOrCreateUser endpoint to send the user a magic link
func (c *config) loginOrCreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := c.stytchClient.MagicLinks.Email.LoginOrCreate(
		&stytch.MagicLinksEmailLoginOrCreateParams{
			Email:              r.FormValue("email"),
			LoginMagicLinkURL:  c.magicLinkURL,
			SignupMagicLinkURL: c.magicLinkURL,
		})
	if err != nil {
		log.Printf("something went wrong sending magic link: %s\n", err)
	}

	parseAndExecuteTemplate("templates/emailSent.html", nil, w)
}

// this is the endpoint the link in the magic link hits takes the token from the
// link's query params and hits the stytch authenticate endpoint to verify the token is valid
func (c *config) authenticate(w http.ResponseWriter, r *http.Request) {
	_, err := c.stytchClient.MagicLinks.Authenticate(
		&stytch.MagicLinksAuthenticateParams{
			Token: r.URL.Query().Get("token"),
		})
	if err != nil {
		log.Printf("something went wrong authenticating the magic link: %s\n", err)
	}

	parseAndExecuteTemplate(
		"templates/loggedIn.html",
		&templateVariables{LoggedOutPath: c.fullAddress + "/logout"},
		w,
	)
}

// handles the logout endpoint
func (c *config) logout(w http.ResponseWriter, r *http.Request) {
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

// helper function so see if a key is in the .env file
// if so return that value, otherwise return the default value
func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if value, exists = os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// helper function to load in the .env file & set config values
func initializeConfig() (*config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("No .env file found at '%s'", ".env")
		return &config{}, errors.New("error loading .env file")
	}
	address := getEnv("ADDRESS", "localhost:3000")

	// define the stytch client using your stytch project id & secret
	// use stytch.EnvLive if you want to hit the live api
	stytchAPIClient, err := stytchapi.NewAPIClient(
		stytch.EnvTest,
		os.Getenv("STYTCH_PROJECT_ID"),
		os.Getenv("STYTCH_SECRET"),
	)
	if err != nil {
		log.Fatalf("error instantiating API client %s", err)
	}

	return &config{
		address:      address,
		fullAddress:  "http://" + address,
		magicLinkURL: "http://" + address + "/authenticate",
		stytchClient: stytchAPIClient,
	}, nil

}
