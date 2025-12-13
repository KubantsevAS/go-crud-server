package stat

import (
	"demo/go-server/configs"
	"demo/go-server/pkg/middleware"
	"fmt"
	"net/http"
	"time"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

type LinkResponse struct {
	NewLink string
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const (
			day   = "day"
			month = "month"
		)

		from, err := parseDateParam("from", req)
		if err != nil {
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}

		to, err := parseDateParam("to", req)
		if err != nil {
			http.Error(w, "Invalid to", http.StatusBadRequest)
			return
		}

		by := req.URL.Query().Get("by")
		if by != day && by != month && by != "" {
			http.Error(w, "Invalid by", http.StatusBadRequest)
			return
		}

		fmt.Println(from)
		fmt.Println(to)
		fmt.Println(by)
	}
}

func parseDateParam(dateParam string, req *http.Request) (time.Time, error) {
	var date time.Time
	paramStr := req.URL.Query().Get(dateParam)

	if paramStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", paramStr)
		if err != nil {
			return date, err
		}
	}

	return date, nil
}
