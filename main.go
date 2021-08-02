package main

import (
	"log"
	_ "github.com/powerpuffpenguin/webpc/assets/document/statik"
	_ "github.com/powerpuffpenguin/webpc/assets/en-US/statik"
	_ "github.com/powerpuffpenguin/webpc/assets/public/statik"
	_ "github.com/powerpuffpenguin/webpc/assets/zh-Hans/statik"
	_ "github.com/powerpuffpenguin/webpc/assets/zh-Hant/statik"
	"github.com/powerpuffpenguin/webpc/cmd"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
