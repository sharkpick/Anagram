package main

import (
	"dictionary"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	Anagrams = dictionary.New()

	Mux = http.DefaultServeMux
)

func SetupHandlers() {
	Mux.Handle("/", http.FileServer(http.Dir("./static")))
	Mux.HandleFunc("/anagrams", HandleAnagrams)
}

func HandleAnagrams(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Println("/anagrams", r.RemoteAddr, r.Method)
	if err := r.ParseMultipartForm(2048); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	word := r.FormValue("word")
	partial, err := strconv.ParseBool(r.FormValue("partial"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Println("/anagrams request for word", word+". partial matches?:", partial)
	var results *dictionary.Anagrams
	if partial {
		results = Anagrams.GetPartialAnagrams(word)
	} else {
		results = Anagrams.GetStraightAnagrams(word)
	}
	w.Write(results.JSON())
	log.Println("finished responding in", time.Since(start))
}

func main() {
	SetupHandlers()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Kill, os.Interrupt)
	go func() {
		if err := http.ListenAndServe(":8080", Mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("ListenAndServe returned unrecognized error", err)
		}
	}()
	<-done
}
