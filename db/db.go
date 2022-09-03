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

const DB_NAME = "words_permutation"

// How many rows you want to operate on each batch
var batchSize = 10000

type Word struct {
	word                  string
	permutation_table_key string
}

func ConnectMySql() (*sql.DB, error) {
	// Capture connection properties.
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
		return nil, err
	}

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: os.Getenv("DB_NAME"),
	}
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	fmt.Println("Connected!")

	err = Select_database(os.Getenv("DB_NAME"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Process_words_from_file(file_name string) error {
	words := []string{}
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatalf("Can't open file %s. Err: %s", file_name, err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		words = append(words, scanner.Text())
		// fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error occured while scanning rows in 'Process_words_from_file'. Err: %s", err)
		return err
	}

	// create table for all words in words_clean.txt
	if err = drop_table("words"); err != nil {
		return err
	}

	Create_words_table_if_not_exists()
	fmt.Println(words)
	Append_words_to_db(words)

	Process_words_permutaion(words)

	return nil
}

func drop_table(table string) error {
	if db == nil {
		return fmt.Errorf("db not connected")
	}

	drop_query := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	fmt.Println(drop_query)
	_, err := db.Exec(drop_query)
	if err != nil {
		log.Fatalf("Error occured during drop_table '%s.%s', Err: %s", DB_NAME, table, err)
		return err
	}

	return nil
}

func Add_permotaion_table_to_db(permotaion_table_name string, words []string) error {
	if db == nil {
		return fmt.Errorf("db not connected")
	}
	drop_table(permotaion_table_name)
	// Create permotaion table
	create_table_query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (word VARCHAR(255) NOT NULL, PRIMARY KEY (word))", permotaion_table_name)
	fmt.Println(create_table_query)
	_, err := db.Exec(create_table_query)
	if err != nil {
		return fmt.Errorf("Error : In Add_permotaion_table_to_db. Err: %s", err)
	}
	// add items to table
	insert_query := fmt.Sprintf("INSERT INTO `%s` (word) VALUES %s", permotaion_table_name, Join(words))
	fmt.Println(insert_query)
	_, err = db.Exec(insert_query)
	if err != nil {
		return fmt.Errorf("Error : In Add_permotaion_table_to_db. Err: %s", err)
	}
	return nil
}

func Create_words_table_if_not_exists() {
	if db == nil {
		log.Fatal("db not connected")
	}
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `words` (word VARCHAR(255) NOT NULL, permutation_table_key VARCHAR(255), PRIMARY KEY (word))")
	if err != nil {
		log.Fatalf("Error occured in 'Create_words_table_if_not_exists'. Err: %s", err)
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
	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)

	// var wg sync.WaitGroup

	for table_name, words := range map_table {
		if len(words) > 1 {
			// wg.Add(1)
			guard <- struct{}{} // would block if guard channel is already filled
			go func(table_name string, words []string) {
				fmt.Println("Key:", table_name, "=>", "words:", words)
				if err := Add_permotaion_table_to_db(table_name, words); err != nil {
					log.Fatalf("Error occured in 'Add_permotaion_table_to_db'. Err: %s", err)
					return
				}
				// defer wg.Done()
				Update_words_on_db(table_name, words)
				<-guard
			}(table_name, words)
		}
	}
	// wg.Wait()
}

func Append_words_to_db(words []string) {
	if db == nil {
		log.Fatal("db not connected")
	}

	// fixed - panic: packet for query is too large. Try adjusting the 'max_allowed_packet' variable on the server
	for i := 0; i < len(words); i += batchSize {
		j := math.Min(float64(i+batchSize), float64(len(words)))
		select_query := fmt.Sprintf("INSERT INTO `words` (word) VALUES %s", Join(words[i:int(j)]))

		//fmt.Println("INSERT INTO words (word) VALUES ", Join(words[i:int(j)]))
		_, err := db.Exec(select_query)
		if err != nil {
			log.Fatalf("Error occured in 'Append_words_to_db'. Err: %s", err)
			return
		}
	}
}

func Update_words_on_db(permotaion_table_name string, words []string) {
	if db == nil {
		log.Fatal("db not connected")
		return
	}
	for _, word := range words {
		update_query := fmt.Sprintf("UPDATE words SET permutation_table_key = '%s' WHERE word='%s'", permotaion_table_name, word)
		_, err := db.Exec(update_query)
		if err != nil {
			log.Fatalf("Error occured in 'Update_words_on_db'. Err: %s", err)
			return
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
	// fmt.Println(b.String())
	return b.String()
}

func Get_permutation_table_key(word string) (Word, error) {
	var w Word

	if db == nil {
		return w, fmt.Errorf("db not connected")
	}
	select_query := fmt.Sprintf("select * FROM `words`_permutation.words WHERE word='%s'", word)
	fmt.Println(select_query)

	// select, err := db.Query("SELECT INTO words (word) VALUES ('test3')")
	row := db.QueryRow(select_query)

	if err := row.Scan(&w.word, &w.permutation_table_key); err != nil {
		if err == sql.ErrNoRows {
			return w, fmt.Errorf("Word %s: no such word", word)
		}
		return w, fmt.Errorf("Word %s: %v", word, err)
	}
	return w, nil
}

func Get_similar_words(word string) ([]string, error) {
	similar_words := []string{}
	if db == nil {
		log.Fatal("db not connected")
		return similar_words, fmt.Errorf("db not connected")
	}
	key_table := algorithms.Generate_key(word)
	query := fmt.Sprintf("select * from `%s`", key_table)
	fmt.Println(query)
	table, table_check := db.Query(query)

	switch {
	case table_check == sql.ErrNoRows:
		fmt.Printf("Table for word '%s' doesn't exist, no similar words, Err: %v", word, table_check)
	case table_check != nil:
		fmt.Printf("Err: %v", table_check)
		return similar_words, nil
	default:

	}
	defer table.Close()

	query = fmt.Sprintf("SELECT * FROM `%s` WHERE word != '%s'", key_table, word)
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

	fmt.Println("similar_words", similar_words)
	return similar_words, nil
}

func Get_total_words() (uint32, error) {
	if db == nil {
		log.Fatal("db not connected")
		return 0, fmt.Errorf("db not connected")
	}
	var count uint32
	fmt.Println("SELECT COUNT(*) FROM words")
	row := db.QueryRow("SELECT COUNT(*) FROM words")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	fmt.Println(count)
	return count, nil
}

func Creat_database_if_not_exists(database_name string) error {
	if db == nil {
		return fmt.Errorf("db not connected")
	}

	drop_query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", database_name)
	fmt.Println(drop_query)
	_, err := db.Exec(drop_query)
	if err != nil {
		log.Fatalf("Error occured creating database '%s' Err: %s", database_name, err)
		return err
	}

	return nil
}

func Select_database(database_name string) error {
	if db == nil {
		return fmt.Errorf("db not connected")
	}

	select_db := fmt.Sprintf("use `%s`", database_name)
	fmt.Println(select_db)
	_, err := db.Exec(select_db)
	if err != nil {
		log.Fatalf("Error occured select database '%s', Err: %s", database_name, err)
		return err
	}

	return nil
}
