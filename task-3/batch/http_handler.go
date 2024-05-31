package batch

import (
	"net/http"
)

func (b *Service) httpHandler(w http.ResponseWriter, r *http.Request) {
	//try to load data from field data
	data := ""
	switch r.Method {
	case http.MethodGet:
		data = r.URL.Query().Get("data")
	case http.MethodPost:
		data = r.FormValue("data")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println(data)
	err := b.Push(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}
