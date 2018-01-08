package apithread

import (
	"net/http"
	"database/sql"
	"github.com/reo7sp/technopark-db/apiutil"
	"log"
	"github.com/reo7sp/technopark-db/api"
	"github.com/reo7sp/technopark-db/dbutil"
)

func MakeShowThreadHandler(db *sql.DB) func(http.ResponseWriter, *http.Request, map[string]string) {
	f := func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		in, err := showThreadRead(r, ps)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		showThreadAction(w, in, db)
	}
	return f
}

type showThreadInput struct {
	slugOrIdInput
}

type showThreadOutput api.ThreadModel

func showThreadRead(r *http.Request, ps map[string]string) (in showThreadInput, err error) {
	resolveSlugOrIdInput(ps["slug_or_id"], &in.slugOrIdInput)
	return
}

func showThreadAction(w http.ResponseWriter, in showThreadInput, db *sql.DB) {
	var out showThreadOutput

	sqlQuery := "SELECT id, title, author, forumSlug, \"message\", votesCount, slug, createdAt FROM threads"
	if in.HasId {
		sqlQuery += " WHERE id = $1"
	} else {
		sqlQuery += " WHERE slug = $1"
	}

	err := db.QueryRow(sqlQuery, in.Slug).Scan(&out.Id, &out.Title, &out.AuthorNickname, &out.ForumSlug, &out.Message, &out.VotesCount, &out.Slug, &out.CreatedDateStr)

	if err != nil && dbutil.IsErrorAboutNotFound(err) {
		errJson := api.ErrorModel{Message: "Can't find thread"}
		apiutil.WriteJsonObject(w, errJson, 404)
		return
	}
	if err != nil {
		log.Println("error: apithread.showThreadAction: SELECT:", err)
		w.WriteHeader(500)
		return
	}

	apiutil.WriteJsonObject(w, out, 200)
}
