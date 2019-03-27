package threads

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

// Node is an immutable snapshot of the story, like a commit
// its Children are the story versions based on it
// its content can be rendered using the String method
type Node struct {
	ID    int64
	Title string
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "commonwriter.db")
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
	stmt := `CREATE TABLE IF NOT EXISTS stories
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		parent_id INTEGER);
		
		REPLACE INTO stories (id, title, content, parent_id)
		VALUES (1, "", "", null);`

	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatalf("opening database: %v", err)
	}
}

// Content returns the whole story
// TODO change signature to HTML to allow some markup
func (n *Node) Content() string {
	stmt := "SELECT content FROM stories WHERE id = $1"
	row := db.QueryRow(stmt, n.ID)
	var content string
	err := row.Scan(&content)
	if err != nil {
		log.Panicf("on .Content: %v", err)
	}
	return content
}

// Children returns the list of Nodes based on n
// TODO stop leaking hella memory
func (n *Node) Children() []*Node {
	stmt := "SELECT id,title FROM stories WHERE parent_id = $1"
	rows, err := db.Query(stmt, n.ID)
	if err != nil {
		log.Panicf("on .Children query: %v", err)
	}
	var ret []*Node
	for rows.Next() {
		child := new(Node)
		err = rows.Scan(&child.ID, &child.Title)
		if err != nil {
			log.Panicf("on .Children scan: %v", err)
		}
		ret = append(ret, child)
	}
	if err = rows.Err(); err != nil {
		log.Panicf("on .Children scan: %v", err)
	}
	return ret
}

// Append makes a node n get a child with content, author and title.
// It then returns the new node's ID.
func (n *Node) Append(content, author, title string) int64 {
	stmt := "INSERT INTO stories VALUES (null, $1, $2, $3);"
	res, err := db.Exec(stmt, title, content, n.ID)
	if err != nil {
		log.Panicf("on .Append: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Panicf("on .Append: %v", err)
	}
	return id
}

// Get a Node with the corresponding id
func Get(id int64) (*Node, error) {
	stmt := "SELECT title FROM stories WHERE id = $1"
	row := db.QueryRow(stmt, id)
	n := new(Node)
	err := row.Scan(&n.Title)
	if err != nil {
		log.Panicf("on .Content: %v", err)
	}
	n.ID = id
	return n, nil // TODO real error reporting
}

// Roots returns all first drafts
func Roots() []*Node {
	var n Node
	n.ID = 1
	return n.Children()
}
