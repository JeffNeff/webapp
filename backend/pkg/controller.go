package controller

import (
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Controller struct {
	rootHandler http.Handler
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
