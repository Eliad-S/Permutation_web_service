package main

import (
	"flag"

	"github.com/Eliad-S/Permutation_web_service/api"
	"github.com/Eliad-S/Permutation_web_service/db"
	"github.com/Eliad-S/Permutation_web_service/statistics"
)

// const keyServerAddr = "serverAddr"

func main() {
	DBClnPtr := flag.Bool("DBCln", false, "usage: {True|Flase}")
	flag.Parse()

	db.ConnectMySql()
	// db.Process_words_from_file()
	if *DBClnPtr {
		statistics.Init()
	}
	api.InitilizeServer()
}
