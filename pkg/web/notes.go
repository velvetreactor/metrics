package web

import (
	"context"
	"net/http"
)

func (w *Web) CreateNote(rw http.ResponseWriter, req *http.Request) {
	dbConn := w.getDbConn()
	defer dbConn.Close(context.Background())

	//err := req.ParseForm()
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//commandTag, err := dbConn.Exec(
	//	context.Background(),
	//	"INSERT INTO notes (body) VALUES ($1)",
	//	req.Form.Get("body"),
	//)
	//
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//if commandTag.RowsAffected() != 1 {
	//	log.Printf("Expected to affect 1 row, affected: %d", commandTag.RowsAffected())
	//	return
	//}
}
