package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func saveFileToGorMD(title string, markdown string, slug string, date string, tags []string) {
	var format string
	format = "---\ndate: %s\nlayout: post\ntitle: %s\npermalink: '%s'\ncategories:\n%s\ntags:\n%s\n---\n\n"
	var filePath string
	filePath = "dist/" + slug + ".md"
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var tagg string
	for _, t := range tags {
		tagg = tagg + "- " + t + "\n"
	}
	f.WriteString(fmt.Sprintf(format, date, title, slug, "- 随记", tagg))
	f.WriteString(markdown)

}
func main() {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select id,title,markdown,slug,published_at from posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var markdown string
		var slug string
		var public_time time.Time
		rows.Scan(&id, &title, &markdown, &slug, &public_time)
		ct := public_time.Unix() / 1000
		if ct > 0 {
			// fmt.Println(time.Unix(public_time.Unix()/1000, 0))
			// fmt.Println(tn)
			row2, err2 := db.Query("select name from tags where id in (select tag_id from posts_tags where post_id=" + fmt.Sprintf("%d", id) + ")")

			if err2 != nil {
				log.Fatal(err2)
			}
			defer row2.Close()
			var tag_name []string
			for row2.Next() {
				var tn string
				row2.Scan(&tn)
				tag_name = append(tag_name, tn)
				// fmt.Println(tn)
			}
			// fmt.Println(tag_name)
			tt := time.Unix(public_time.Unix()/1000, 0)
			saveFileToGorMD(title, markdown, slug, tt.Format("2006-01-02"), tag_name)
		}
	}
	rows.Close()
}
