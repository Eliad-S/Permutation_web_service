package main

import (
	"flag"
	"log"

	"github.com/Eliad-S/Permutation_web_service/api"
	"github.com/Eliad-S/Permutation_web_service/db"
	"github.com/Eliad-S/Permutation_web_service/statistics"
)

// const keyServerAddr = "serverAddr"

func main() {
	DBClnPtr := flag.Bool("DBCln", false, "usage: {True|Flase}")
	flag.Parse()

	sql_db, err := db.ConnectMySql()
	if err != nil {
		log.Fatal("Error ConnectMySql" + err.Error())
	}
	defer sql_db.Close()

	if *DBClnPtr {
		err = db.Process_words_from_file("words_clean.txt")
		if err != nil {
			log.Fatal("Error in Process_words_from_file" + err.Error())

		}
	}

	statistics.Init()
	api.InitilizeServer()
}
