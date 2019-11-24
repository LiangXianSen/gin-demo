package web

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LiangXianSen/gin-demo/config"
)

var srv *Server

var (
	fixturePath string
)

func TestMain(m *testing.M) {
	// Change to project root directory
	if err := os.Chdir("../../"); err != nil {
		log.Fatal(err)
	}

	// testing fixtures path
	cwd, _ := os.Getwd()
	fixturePath = filepath.Join(cwd, "fixtures")

	// Load config from toml file
	conf, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	srv = NewServer(conf)

	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	must := assert.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	srv.Router.ServeHTTP(w, r)

	must.Equal(w.Code, http.StatusOK)
}

func BenchmarkHealth(b *testing.B) {
	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/health", nil)
		srv.Router.ServeHTTP(w, r)
	}
}
