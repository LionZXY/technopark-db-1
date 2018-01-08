package apithread

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/dbutil"
)

func MakeEditThreadHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := editThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		editThreadAction(w, in, db)
	}
	return f
}

type editThreadInput struct {
	slugOrIdInput

	Title   *string `json:"title"`
	Message *string `json:"message"`
}

type editThreadOutput api.ThreadModel

func editThreadRead(r *http.Request, ps map[string]string) (in editThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	err = apiutil.ReadJsonObject(r, &in)
	return
}

func editThreadAction(w http.ResponseWriter, in editThreadInput, db *sql.DB) {
	var out editThreadOutput

	sqlQuery := "UPDATE threads SET title = COALESCE($1, title), \"message\" = COALESCE($2, \"message\")"
	sqlFields := []interface{}{in.Title, in.Message, nil}
	if in.HasId {
		sqlQuery += " WHERE id = $3"
		sqlFields[2] = in.Id
	} else {
		sqlQuery += " WHERE slug = $3"
		sqlFields[2] = in.Slug
	}
	sqlQuery += " RETURNING author, createdAt, forumSlug, id, \"message\", slug, title"

	err := db.QueryRow(sqlQuery, sqlFields...).Scan(&out.AuthorNickname, &out.CreatedDateStr, &out.ForumSlug, &out.Id, &out.Message, &out.Slug, &out.Title)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.editThreadAction: UPDATE:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
