# Permutation_web_service
A small web service for printing similar words in the English language.
The web service listen on localhost, port 8080.

This project was created using the following technologies:
 * Backend:
    * Server: golang
    * DB: MySQL


## Application Requirements

* Make sure you have mySQL installed on your computer, otherwise use this link https://dev.mysql.com/doc/mysql-getting-started/en/
* Make sure you have download go compiler. https://go.dev/dl/
* In local.env file set your MySQL's username/password as well as your database name.
```
DBUSER=username
DBPASS=password
DB_NAME=database_name
```

## Build

```bash

# initial http server including database initilization.
go run main.go -DBCln=true

or

# initial http server without database initilization.
go run main.go

```

## REST API

Now can now use the web service with the following requests:

http://localhost:8080//api/v1/similar?word=your_word

```
The result format is a JSON object as follows:
{
    similar:[list,of,words,that,are,similar,to,provided,word]
}

For example:
http://localhost:8000/api/v1/similar?word=apple
{"similar":["appel","pepla"]}
```

http://localhost:8080/api/v1/stats

```
Return general statistics about the program:
Total number of words in the dictionary
Total number of requests (not including "stats" requests)
Average time for request handling in nano seconds (not including "stats" requests)

The output is a JSON object structured as follows:
{
    totalWords:int
    totalRequests:int
    avgProcessingTimeNs:int
}

For example:
http://localhost:8080/api/v1/stats
{"totalWords":351075,"totalRequests":9,"avgProcessingTimeNs":45239}
```

## How does it works?

Two words w_1 and w_2 are considered similar if w_1 is a letter permutation of w_2 (e.g., "stressed" and "desserts").

In order to intintify if 2 words are equal we are using hash map to count the number of times the letters appeared in the first, and compare it to the second.

we are doing initial process to all words in words_clean.txt.

we creates the following table:

* Table name : "words", Columns : word
* Table for more than one similar words we found, the table's will be a key constructed from the words's letters - sorted.

As a result, in runtime we need only to create the key from the given word and find the relevant table to return the similar words to it.

## Dummy client

you can test the scalability of the web service using a HTTP client in a concurrent application we provided here.

```
go run dummy_client/client.go
```


