# Permutation_web_service
A small web service for printing similar words in the English language.
The web service listen on localhost, port 8000.

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

## Getting Started

First we need to install all the Python and Node.js packages used in this project:

```
cd [project dir]
npm install
pip install -r requirements.txt
```

After that, you need to go to ```node_modules\react-scripts\config\webpackDevServer.config.js``` file and add the following code under ```watchOptions```:

    watchOptions: {
      ignored: [ ignoredFiles(paths.appSrc), paths.appPublic ]
    },


Now you ready to use the web service to get the following:
```
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
http://localhost:8000/api/v1/stats
{"totalWords":351075,"totalRequests":9,"avgProcessingTimeNs":45239}

