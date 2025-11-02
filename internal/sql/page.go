package sql

import (
	"context"
	"log"
	"strconv"
	"time"
)

// title
// id         | integer                        |           | not null | nextval('title_id_seq'::regclass) | plain    |              |
// key        | character varying(255)         |           | not null |                                   | extended |              |
// value      | character varying(255)         |           | not null |                                   | extended |              |
// mod_id     | integer                        |           | not null |                                   | plain    |              |
// language   | character varying(255)         |           | not null |                                   | extended |              |
// model      | character varying(255)         |           | not null |                                   | extended |              |
// created_at | timestamp(0) without time zone |           | not null |                                   | plain    |              |
// updated_at | timestamp(0) without time zone |           | not null |                                   | plain    |              |

// content
// id         | integer                        |           | not null | nextval('content_id_seq'::regclass) | plain    |              |
// key        | character varying(255)         |           | not null |                                     | extended |              |
// value      | text                           |           | not null |                                     | extended |              |
// mod_id     | integer                        |           | not null |                                     | plain    |              |
// language   | character varying(255)         |           | not null |                                     | extended |              |
// model      | character varying(255)         |           | not null |                                     | extended |              |
// created_at | timestamp(0) without time zone |           | not null |                                     | plain    |              |
// updated_at | timestamp(0) without time zone |           | not null |                                     | plain    |              |

type Page struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetPage(pageID string, lang string) (*Page, error) {

	_pageID, err := strconv.Atoi(pageID)
	if err != nil {
		log.Fatalf("%s is not a instance of int", pageID)
	}

	params := []any{}
	params = append(params, _pageID)
	params = append(params, lang)
	query := `select t.value as title, c.value as content, c.created_at from page_title_contents ptc join title t on t.id = ptc.title_id join content c on ptc.content_id = c.id where c.language = $2 and ptc.page_id = $1`
	var page Page

	err = Pgsql.QueryRow(context.Background(), query, params...).Scan(&page.Title, &page.Content, &page.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func GetLaw(pageID string, lang string) (*Page, error) {

	_pageID, err := strconv.Atoi(pageID)
	if err != nil {
		log.Fatalf("%s is not a instance of int", pageID)
	}

	params := []any{}
	params = append(params, _pageID)
	params = append(params, lang)
	query := `select t.value as title, c.value as content, c.created_at from legal_title_contents ptc join title t on t.id = ptc.title_id join content c on ptc.content_id = c.id where c.language = $2 and ptc.legal_id = $1`
	var page Page

	err = Pgsql.QueryRow(context.Background(), query, params...).Scan(&page.Title, &page.Content, &page.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &page, nil
}
