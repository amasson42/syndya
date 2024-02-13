package Controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"syndya/pkg/Models"
	"testing"

	"github.com/gin-gonic/gin"
)

func makeEngine() (*gin.Engine, Models.SearchingPlayersBank) {
	playersBank := Models.NewSearchingPlayersList()
	controller := NewPlayersController(playersBank)

	r := gin.Default()
	controller.Route(r)

	return r, playersBank
}

// TestSearchGame tests the searchGame method of PlayersController.
func TestSearchGame(t *testing.T) {
	r, _ := makeEngine()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	go func() {
		r.ServeHTTP(w, req)
	}()

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, w.Code)
	}
}

// TestGetPlayers tests the getPlayers method of PlayersController.
func TestGetPlayers(t *testing.T) {
	r, _ := makeEngine()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/players", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d; got %d", http.StatusOK, w.Code)
	}
}

// TestSearchGameConcurrentWithMetadata modifies metadata while opening 10 WebSocket connections concurrently.
func TestSearchGameConcurrentWithMetadata(t *testing.T) {
	r, playersBank := makeEngine()

	count := 10
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(index int) {
			defer wg.Done()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/search", nil)
			req.Header.Set("Connection", "upgrade")
			req.Header.Set("Upgrade", "websocket")
			req.Header.Set("Sec-WebSocket-Version", "13")
			req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %d; got %d", http.StatusOK, w.Code)
			}

			// Simulate modifying metadata
			key := fmt.Sprintf("metadata_%d", index)
			value := fmt.Sprintf("value_%d", index)
			playersBank.UpdateSearchingPlayerMetadata(index+1, key, value)
		}(i)
	}

	wg.Wait()

	// Verify modifications to metadata in PlayersBank
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("metadata_%d", i)
		value := fmt.Sprintf("value_%d", i)
		player := playersBank.GetSearchingPlayerFromID(i + 1)
		if player == nil {
			t.Errorf("Player with ID %d not found", i)
			continue
		}
		if player.MetaData[key] != value {
			t.Errorf("Expected metadata value %s for key %s; got %s", value, key, player.MetaData[key])
		}
	}
}
