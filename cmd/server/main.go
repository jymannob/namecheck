package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/jymannob/namecheck"
	_ "github.com/jymannob/namecheck/facebook"
	_ "github.com/jymannob/namecheck/github"
	_ "github.com/jymannob/namecheck/instagram"
	_ "github.com/jymannob/namecheck/twitter"
)

var count uint64
var m = make(map[string]uint64)
var mu sync.Mutex

type Result struct {
	platform  string
	valid     bool
	available bool
	err       error
}

func main() {
	// catch interupt SIGINT,SIGKILL
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	go func() {
		<-signals
		fmt.Println("[signals] Interupt !")
		// write detail in json file
		writeDetailInfile()
		os.Exit(0)
	}()
	//Load detail from json file
	loadDetailFromfile()

	http.Handle("/check", http.HandlerFunc(handle))
	http.Handle("/count", http.HandlerFunc(handleCount))
	http.Handle("/details", http.HandlerFunc(handleDetails))
	startBrowser("http://127.0.0.1:8080/check?username=jymannob")
	startBrowser("http://127.0.0.1:8080/details")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	entity := struct {
		Count uint64 `json:"count"`
	}{
		Count: atomic.LoadUint64(&count),
	}
	dec := json.NewEncoder(w)
	if err := dec.Encode(entity); err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
}

func handleDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	dec := json.NewEncoder(w)
	mu.Lock()
	{
		if err := dec.Encode(m); err != nil {
			http.Error(w, "{}", http.StatusInternalServerError)
			return
		}
	}
	mu.Unlock()
}

func handle(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&count, 1)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	username := r.URL.Query().Get("username")
	mu.Lock()
	{
		m[username]++
	}
	mu.Unlock()
	if len(username) == 0 {
		http.Error(w, "missing 'username' query parameter", http.StatusInternalServerError)
		return
	}
	ch := make(chan Result)
	var wg sync.WaitGroup
	checkers := namecheck.Checkers()
	wg.Add(len(checkers))
	for _, c := range checkers {
		go check(c, username, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	type jsonResult struct {
		Platform  string `json:"platform"`
		Valid     string `json:"valid"`
		Available string `json:"available"`
	}
	jsonResults := make([]jsonResult, 0, len(checkers))
	for result := range ch {
		res := jsonResult{
			Platform:  result.platform,
			Valid:     fmt.Sprintf("%t", result.valid),
			Available: availabilityStatus(result),
		}
		jsonResults = append(jsonResults, res)
	}
	entity := struct {
		Username string       `json:"username"`
		Results  []jsonResult `json:"results"`
	}{
		Username: username,
		Results:  jsonResults,
	}
	dec := json.NewEncoder(w)
	if err := dec.Encode(entity); err != nil {
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
}

func check(
	c namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	ch chan<- Result) {
	defer wg.Done()
	res := Result{
		platform: c.String(),
	}
	valid := c.IsValid(username)
	if !valid {
		ch <- res
		return
	}
	res.valid = true
	avail, err := c.IsAvailable(username)
	if err != nil {
		res.err = err
		ch <- res
		return
	}
	if !avail {
		ch <- res
		return
	}
	res.available = true
	ch <- res
}

func availabilityStatus(res Result) string {
	if res.err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%t", res.available)
}

func writeDetailInfile() {
	f, err := os.OpenFile("detail.json", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Impossible d'écrire dans le fichier : detail.json")
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	mu.Lock()
	{
		if err := enc.Encode(m); err != nil {
			log.Println("Impossible de convertir en Json")
			return
		}
	}
	mu.Unlock()
}

func loadDetailFromfile() {
	dat, err := ioutil.ReadFile("detail.json")
	if err != nil {
		fmt.Println("Impossible de lire le fichier : detail.json")
		return
	}
	fmt.Println(string(dat))

	mu.Lock()
	err = json.Unmarshal(dat, &m)
	mu.Unlock()

	if err != nil {
		fmt.Println("Impossible de convertir le Json")
		return
	}
}

// startBrowser tries to open the URL in a browser, and returns
// whether it succeed.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
