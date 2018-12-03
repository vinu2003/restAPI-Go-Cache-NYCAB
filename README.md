# restAPI-Go-Cache-NYCAB
A basic Newyork Cab trips retrieving API written in GO Programming Language
you can perform,
    - Retrive number of trips made by particular medallion(cob) given a particular date.
    - flush the server cache.
    
The API receives one or more medallions and return how many trips each medallion has made.
Considering that the query creates a heavy load on the database, the results must be cached.
The API allows user to ask for fresh data, ignoring the cache.

Directory Structure:
--------------------
```
dataRep/
    |-- rediscache
        |--cache.go                        - Contains redis client methods - GET, SET, FLUSHDB
        |--cache_test.go                   - Contains units tests for resdis client methods
    |--repo
        |-- dbstore.go                     - Methods interacting with the database
    |-- error.go                           - Implements functions to manipulate errors
    |-- README.md
    |-- ny_cab_data_cab_trip_data_full.sql - sql statements to populate database.
    |-- main.go                            - Entry point of the API, defined router endpoints and also handlers
    |-- middleware.go                      - middleware methods like cache wrapper.
 ```

SETUP
-----
Local Machine I have used is macbookPro - Version 10.13.6

GoLang installment setup:
go version - go version go1.11.1 darwin/amd64

***ASSUMPTION*** : Go is requested to be the prmary programming language so I beleive the go setup must be available.
IDE used : GoLand
Set GOPATH to your src directory where the source files are placed.

Install MYSQL server 8.0 as database.

Install REDIS 5.0 as frontend for database as cache layer.

## Library for unit test
$ go get -u "github.com/stretchr/testify"

## Libraries to handle network routing
$ go get -u "github.com/gorilla/mux"
$ go get -u "github.com/gorilla/handlers"
$ go get -u "github.com/go-sql-driver/mysql"

## for redis
$ got get -u "gopkg.in/redis.v5"


REDIS setup
-----------
setup the server passwd by editing /usr/local/etc/redis.conf -> requirepass United123

restart the redis service
$ brew redis restart service

install redis-cli if required - to manually set/get KEY-VALUE and test.
$ brew install redis-cli
set the client password.

API details
-----------

Implemented GET method with one or more medallion and pickup date as query string and cached the details. 
The option is provided for the user to skipCache or not true/false - can be gievn as query string.

if skipCache not provided default is false. ie DO NOT SKIP THE CACHE READ. MAKE FIRST ATTEMPT TO CHECK IF THE PAGE for GIEVN KEY is CACHED. IF SO GET THE DATA FROM CACHE BEFORE CALLING HANDLER.
IF the cache is not found then handler read from DB.

if skipCache is true - always read from DB.

However, when if skipCache is true(explicitly) or if the page is first time fetched by user(fresh fetch) or page not found in cache so server reads from DB - for all these cases update cache for that particular key.
DEFAULT CACHE time: 15s in main.go.

Keep it simple and user friendly.
Given medallion and dates provided or their combinations not found update the output as "Data not found."

SET and GET methods for redis cache are implemented in cache.go

Next, implemented PUT method - deletes all keys of the currently selected DB.

Unit tests:
----------
for handler routines - main_test.go
for redis methods - cache_test.go


Improvement todo:
----------------
Improve the error handling mechanism - for now it wraps the error from external dependencies passed on to handler and verified.
This can be even more simplefied and logs as well.


