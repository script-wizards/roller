package main

import (
	"context"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/travis-g/dice"
	"github.com/travis-g/dice/math"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	r := chi.NewRouter()

	r.Get("/", handler)
	r.Get("/api/v1/roll", apiHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/404.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Msg string
		}{
			Msg: random404(),
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd := ""
	msg := ""
	if r.URL.Query().Get("cmd") != "" {
		cmd = r.URL.Query().Get("cmd")
	}
	if cmd != "" {
		ctx := dice.NewContextFromContext(context.Background())
		exp, err := math.EvaluateExpression(ctx, cmd)
		if err != nil {
			msg = err.Error()
		} else {
			msg = exp.String()
		}

	}
	data := struct {
		Msg string
	}{
		Msg: msg,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	cmd := ""
	msg := ""
	if r.URL.Query().Get("cmd") != "" {
		cmd = r.URL.Query().Get("cmd")
	}
	if cmd != "" {
		ctx := dice.NewContextFromContext(context.Background())
		exp, err := math.EvaluateExpression(ctx, cmd)
		if err != nil {
			msg = err.Error()
		} else {
			msg = strconv.FormatFloat(exp.Result, 'f', -1, 64)
		}

	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg + "\n"))
}

func random404() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	messages := []string{
		"Impossible. Perhaps the archives are incomplete.",
		"This page is too strong for you, traveler.",
		"This page is in another castle.",
		"One does not simply walk into this page.",
		"It's dangerous to go alone, take this link back home.",
		"\"DID YOU PUT THIS PAGE IN THE GOBLET OF FIRE\", Dumbledore said calmly.",
	}
	return messages[rng.Intn(len(messages))]
}
