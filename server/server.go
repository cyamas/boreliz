package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cyamas/boreliz/internal/wall"
	"github.com/cyamas/boreliz/server/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var Pool *pgxpool.Pool

type username string

type request string

var ActiveUsers map[username]request

func Run(dbPool *pgxpool.Pool) {
	Pool = dbPool
	router := http.NewServeMux()
	router.HandleFunc("/", home)
	router.HandleFunc("/close", closeDiv)
	router.HandleFunc("/select-wall", func(w http.ResponseWriter, r *http.Request) {
		selectWall(w, r, dbPool)
	})
	router.HandleFunc("/signup-form", func(w http.ResponseWriter, r *http.Request) {
		signupForm(w, r, dbPool)
	})
	router.HandleFunc("/signin-form", func(w http.ResponseWriter, r *http.Request) {
		signinForm(w, r, dbPool)
	})
	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		signup(w, r, dbPool)
	})
	router.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		signin(w, r, dbPool)
	})
	router.HandleFunc("/edit-meas-form", func(w http.ResponseWriter, r *http.Request) {
		editMeasurementsForm(w, r, dbPool)
	})

	protectedEditMeasurements := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		editMeasurements(w, r, dbPool)
	})
	router.Handle("/edit-meas", middleware.Authorization()(protectedEditMeasurements))

	protectedProfile := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile(w, r, dbPool)
	})
	router.Handle("/profile", middleware.Authorization()(protectedProfile))

	protectedLogbook := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logbook(w, r, dbPool)
	})
	router.Handle("/logbook", middleware.Authorization()(protectedLogbook))

	protectedSignout := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signout(w, r)
	})
	router.Handle("/signout", middleware.Authorization()(protectedSignout))

	router.HandleFunc("/create-hold-form", createHoldForm)

	protectedCreateHold := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createHold(w, r, dbPool)
	})
	router.Handle("/create-hold", middleware.Authorization()(protectedCreateHold))

	staticDir := http.FileServer(http.Dir("static"))
	router.Handle("/static/", http.StripPrefix("/static/", staticDir))

	server := http.Server{
		Addr:    ":6969",
		Handler: middleware.Logging(router),
	}
	log.Println("starting server on port :6969")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		IsAuthed bool
	}{
		false,
	}
	cookie, err := r.Cookie("jwt_token")
	if err == nil && cookie.Value != "" {
		jwt := cookie.Value
		_, err := middleware.DecodeJWT(jwt)
		if err == nil {
			data.IsAuthed = true
		}
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func selectWall(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	wallName := r.FormValue("wall-name")
	wall := wall.New()
	wall.ID = wallName

	switch wallName {
	case "yama":
		wall.Rows, wall.Cols = 24, 24
		wall.Angle = 37
		wall.CreateGrid()
	case "marjoram":
		wall.Rows, wall.Cols = 20, 16
		wall.Angle = 20
		wall.CreateGrid()
	}
	tmpl, err := template.ParseFiles("templates/wall.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, wall); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signinForm(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	tmpl, err := template.ParseFiles("templates/signin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signupForm(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func signup(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
	}
	userName := r.FormValue("username")
	pw := r.FormValue("password")
	confirmPW := r.FormValue("confirm-password")
	if pw != confirmPW {
		w.Write([]byte("passwords did not match. Please try again"))
		return
	}
	salt := generateSalt(w)
	hash := generateHash(w, pw, salt)
	saltStr := base64.StdEncoding.EncodeToString(salt)
	hashStr := base64.StdEncoding.EncodeToString(hash)
	insertQuery := "INSERT INTO users (username, password, salt) VALUES ($1, $2, $3)"
	_, err = dbPool.Exec(context.Background(), insertQuery, userName, hashStr, saltStr)
	if err != nil {
		log.Println("error inserting into database: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Welcome " + userName + ". Sign in to create profile."))
}

func generateSalt(w http.ResponseWriter) []byte {
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		http.Error(w, "failed to read in salt", http.StatusInternalServerError)
	}
	return salt
}

func generateHash(w http.ResponseWriter, pw string, salt []byte) []byte {
	saltedPW := append([]byte(pw), salt...)

	hash, err := bcrypt.GenerateFromPassword(saltedPW, 0)
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return hash
}

func signin(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	pw := r.FormValue("password")

	err = authenticateUser(username, pw, dbPool, w)
	if err != nil {
		w.Write([]byte("Invalid username or password"))
		return
	}
	tmpl, err := template.ParseFiles("templates/profile-btn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "jwt_token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
		MaxAge:  -1,
	}
	http.SetCookie(w, &cookie)

	tmpl, err := template.ParseFiles("templates/signout.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func authenticateUser(username, pw string, dbPool *pgxpool.Pool, w http.ResponseWriter) error {
	var id int
	var hashStr string
	var saltStr string

	query := "SELECT id, password, salt FROM users WHERE username = $1"
	err := dbPool.QueryRow(context.Background(), query, username).Scan(&id, &hashStr, &saltStr)
	if err != nil {
		log.Println("could not fetch data from database")
		http.Error(w, "could not fetch data with given username", http.StatusBadRequest)
		return errors.New("could not authenticate user")
	}
	hash, err := base64.StdEncoding.DecodeString(hashStr)
	if err != nil {
		http.Error(w, "could not validate credentials", http.StatusInternalServerError)
	}
	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		http.Error(w, "could not validate credentials", http.StatusInternalServerError)
	}
	preHash := append([]byte(pw), salt...)
	err = bcrypt.CompareHashAndPassword(hash, preHash)
	if err != nil {
		return errors.New("could not authenticate user")
	}
	jwt := middleware.CreateJWT(username, id)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    jwt,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

func profile(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	tmpl, err := template.ParseFiles("templates/profile.html", "templates/measurements.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Println(http.StatusInternalServerError, "User ID not found in context")
		http.Error(w, "User ID not found", http.StatusInternalServerError)
	}

	query := `SELECT height, wingspan, vert_reach FROM users WHERE id = $1`
	data := struct {
		Height    float64
		Wingspan  float64
		VertReach float64
	}{}
	err = dbPool.QueryRow(context.Background(), query, userID).Scan(&data.Height, &data.Wingspan, &data.VertReach)
	if err != nil {
		log.Println("error querying database.", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("error executing template: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func editMeasurementsForm(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	tmpl, err := template.ParseFiles("templates/edit-meas-form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("error executing template: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editMeasurements(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Println(http.StatusInternalServerError, "could not get userID from context")
		http.Error(w, "could not retrieve userID from context", http.StatusInternalServerError)
	}
	heightStr := r.FormValue("height")
	wingspanStr := r.FormValue("wingspan")
	vertReachStr := r.FormValue("vert-reach")

	height, err := strconv.ParseFloat(heightStr, 64)
	if err != nil {
		http.Error(w, "Invalid height value", http.StatusBadRequest)
		return
	}
	wingspan, err := strconv.ParseFloat(wingspanStr, 64)
	if err != nil {
		http.Error(w, "Invalid wingspan value", http.StatusBadRequest)
		return
	}
	vertReach, err := strconv.ParseFloat(vertReachStr, 64)
	if err != nil {
		http.Error(w, "Invalid vertical reach value", http.StatusBadRequest)
		return
	}

	query := `UPDATE users
	SET height = $1,
	wingspan = $2,
	vert_reach = $3
	WHERE id = $4`
	_, err = dbPool.Exec(context.Background(), query, height, wingspan, vertReach, userID)
	if err != nil {
		log.Println(http.StatusInternalServerError, "error executing sql query.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Height    float64
		Wingspan  float64
		VertReach float64
	}{
		height,
		wingspan,
		vertReach,
	}
	tmpl, err := template.ParseFiles("templates/profile.html", "templates/measurements.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("error executing template: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logbook(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	tmpl, err := template.ParseFiles("templates/logbook.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("error executing template: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func closeDiv(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(""))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createHoldForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create-hold-form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("error executing template: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func createHold(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(http.StatusBadRequest, "error parsing form")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	os.MkdirAll("./static/hold-images", os.ModePerm)
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(http.StatusBadRequest, "Could not convert id to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	angle, err := strconv.Atoi(r.FormValue("angle"))
	if err != nil {
		log.Println(http.StatusBadRequest, "Could not convert angle to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	texture, err := strconv.Atoi(r.FormValue("texture"))
	if err != nil {
		log.Println(http.StatusBadRequest, "Could not convert texture to int")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imagePath := fmt.Sprintf("./static/hold-images/%d.jpeg", id)
	dst, err := os.Create(imagePath)
	if err != nil {
		log.Println("Could not create file")
		return
	}
	defer dst.Close()
	image, _, err := r.FormFile("image")
	if err != nil {
		log.Println(http.StatusBadRequest, "Error retrieving formfile from request")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if _, err := io.Copy(dst, image); err != nil {
		log.Println("Could not copy file to directory")
		return
	}
	hold := struct {
		ID           int
		Manufacturer string
		Model        string
		Type         string
		Color        string
		Angle        int
		Texture      int
		ImagePath    string
		ImageLength  int
		ImageWidth   int
	}{
		id,
		r.FormValue("manufacturer"),
		r.FormValue("model"),
		r.FormValue("type"),
		r.FormValue("color"),
		angle,
		texture,
		r.FormValue("imagePath"),
		0,
		0,
	}
	fmt.Println("hold", hold)

}
