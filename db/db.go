package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/Eliad-S/Permutation_web_service/algorithms"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB = nil

// How many rows you want to operate on each batch
var batchSize = 10000

type Word struct {
	word                  string
	permutation_table_key string
}

func ConnectMySql() {
	// Capture connection properties.
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "words_permutation",
	}
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr.Error())
	}
	fmt.Println("Connected!")
}

func Process_words_from_file() {
	a := []string{}
	file, err := os.Open("words_clean.txt")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		a = append(a, scanner.Text())
		// fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// create table for all words in words_clean.txt
	drop_table("words")
	Create_words_table_if_not_exists()
	Append_words_to_db(a)

	Process_words_permutaion(a)
}

func drop_table(table string) {
	if db == nil {
		log.Fatal("db not connected")
	}

	drop_query := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	fmt.Println(drop_query)
	_, err := db.Exec(drop_query)
	if err != nil {
		panic(err)
	}
}

func Add_permotaion_table_to_db(permotaion_table_name string, words []string) {
	if db == nil {
		log.Fatal("db not connected")
	}
	drop_table(permotaion_table_name)
	// Create permotaion table
	create_table_query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (word VARCHAR(255) NOT NULL, PRIMARY KEY (word))", permotaion_table_name)
	fmt.Println(create_table_query)
	_, err := db.Exec(create_table_query)
	if err != nil {
		panic(err)
	}
	// add items to table
	insert_query := fmt.Sprintf("INSERT INTO `%s` (word) VALUES %s", permotaion_table_name, Join(words))
	fmt.Println(insert_query)
	_, err = db.Exec(insert_query)
	if err != nil {
		panic(err)
	}
}

func Create_words_table_if_not_exists() {
	if db == nil {
		log.Fatal("db not connected")
	}
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `words` (word VARCHAR(255) NOT NULL, permutation_table_key VARCHAR(255), PRIMARY KEY (word))")
	if err != nil {
		panic(err)
	}
}

func Process_words_permutaion(words []string) {
	hash_set := make(map[string]bool)
	map_table := make(map[string][]string)
	for _, word := range words {
		_, ok := hash_set[word]
		if !ok {
			key := algorithms.Generate_key(word)
			map_table[key] = append(map_table[key], word)
			hash_set[word] = true
		}
	}
	create_permotaion_tables(map_table)
}

func create_permotaion_tables(map_table map[string][]string) {
	maxGoroutines := 30
	guard := make(chan struct{}, maxGoroutines)

	for table_name, words := range map_table {
		if len(words) > 1 {
			guard <- struct{}{} // would block if guard channel is already filled
			go func(table_name string, words []string) {
				fmt.Println("Key:", table_name, "=>", "words:", words)
				Add_permotaion_table_to_db(table_name, words)
				Update_words_on_db(table_name, words)
				<-guard
			}(table_name, words)
		}
	}
}

func Append_words_to_db(words []string) {
	if db == nil {
		log.Fatal("db not connected")
	}

	// fixed - panic: packet for query is too large. Try adjusting the 'max_allowed_packet' variable on the server
	for i := 0; i < len(words); i += batchSize {
		j := math.Min(float64(i+batchSize), float64(len(words)))
		select_query := fmt.Sprintf("INSERT INTO `words` (word) VALUES %s", Join(words[i:int(j)]))

		// fmt.Println("INSERT INTO words (word) VALUES ", Join(words[i:int(j)]))
		_, err := db.Exec(select_query)
		if err != nil {
			panic(err.Error())
		}
	}
}

func Update_words_on_db(permotaion_table_name string, words []string) {
	if db == nil {
		log.Fatal("db not connected")
	}
	for _, word := range words {
		update_query := fmt.Sprintf("UPDATE words SET permutation_table_key = '%s' WHERE word='%s'", permotaion_table_name, word)
		_, err := db.Exec(update_query)
		if err != nil {
			panic(err.Error())
		}
	}
}

func Join(strs []string) string {
	var b strings.Builder
	first := true
	for _, word := range strs {
		if first {
			fmt.Fprintf(&b, "(\"%s\")", word)
			first = false
		} else {
			fmt.Fprintf(&b, ", (\"%s\")", word)
		}
	}
	fmt.Println(b.String())
	return b.String()
}

func Get_permutation_table_key(word string) (Word, error) {
	var w Word

	if db == nil {
		return w, fmt.Errorf("db not connected")
	}
	select_query := fmt.Sprintf("select * FROM `words`_permutation.words WHERE word='test'")
	fmt.Println(select_query)

	// select, err := db.Query("SELECT INTO words (word) VALUES ('test3')")
	row := db.QueryRow(select_query)

	if err := row.Scan(&w.word, &w.permutation_table_key); err != nil {
		if err == sql.ErrNoRows {
			return w, fmt.Errorf("albumsById %s: no such album", word)
		}
		return w, fmt.Errorf("albumsById %s: %v", word, err)
	}
	return w, nil
}

func Get_similar_words(word string) ([]string, error) {
	similar_words := []string{}
	if db == nil {
		log.Fatal("db not connected")
	}
	key_table := algorithms.Generate_key(word)
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE word != '%s'", key_table, word)
	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get permutation_table by key %s: %v", word, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var next_word string
		if err := rows.Scan(&next_word); err != nil {
			return nil, fmt.Errorf("get permutation_table by key %s: %v", word, err)
		}
		similar_words = append(similar_words, next_word)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get permutation_table by key %s: %v", word, err)
	}
	fmt.Println("similar_words", similar_words)
	return similar_words, nil
}

func Get_total_words() uint32 {
	if db == nil {
		log.Fatal("db not connected")
	}
	var count uint32
	fmt.Println("SELECT COUNT(*) FROM words")
	row := db.QueryRow("SELECT COUNT(*) FROM words")
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
	return count
}

// func Init_stats(table_name string) statistics.Statistics {
// 	if db == nil {
// 		log.Fatal("db not connected")
// 	}

// 	create_stats_table := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`",
// 		"(id AUTOINCREMENT PRIMARY KEY,",
// 		"TotalWords INT,",
// 		"TotalRequests INT,",
// 		"AvgProcessingTimeNs FLOAT)", table_name)
// 	fmt.Println(create_stats_table)
// 	_, err := db.Exec(create_stats_table)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// if

// 	insert_query := fmt.Sprintf("INSERT INTO `%s` (TotalWordsm, TotalRequests, AvgProcessingTimeNs) VALUES (0, 0, 0)", table_name)
// 	fmt.Println(insert_query)
// 	_, err = db.Exec(insert_query)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return statistics.Statistics{0, 0, 0}
// }

// func Import_stat_from_db(table_name string) statistics.Statistics {

// }
