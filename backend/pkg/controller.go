package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Controller struct {
	rootHandler http.Handler
}

type event struct {
	Input string `json:"input"`
}

var once sync.Once

func (c *Controller) RootHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		kdp := path.Join(os.Getenv("KO_DATA_PATH"))
		if !strings.HasSuffix(kdp, "/") {
			kdp = kdp + "/"
		}
		c.rootHandler = http.FileServer(http.Dir(kdp))
	})
	c.rootHandler.ServeHTTP(w, r)
}

func (c *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error occured reading body: %v", err)
		json.NewEncoder(w).Encode("Error reding request")
		return
	}

	defer r.Body.Close()
	e := &event{}
	err = json.Unmarshal(body, e)
	if err != nil {
		log.Printf("error occured unmarshalling body: %v", err)
		json.NewEncoder(w).Encode("Error unmarshalling request")
		return
	}

	if e.Input != "" {
		e.Input = "hello: " + e.Input + ". How are you today?"
	}

	w.Write([]byte(e.Input))
}
