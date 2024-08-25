package router

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

func getLocationsRouteHandler(w http.ResponseWriter, r *http.Request, redisClient RedisClientInterface) {
	sheetId := r.URL.Query().Get("sheetId")
	if sheetId == "" {
		http.Error(w, "Missing sheetId parameter", http.StatusBadRequest)
		return
	}

	// log.Println("Calling Redis Get method with key:", sheetId)
	locationsJson, err := redisClient.Get(r.Context(), sheetId).Bytes()
	if err == redis.Nil {
		http.Error(w, "No locations found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(locationsJson))
}
