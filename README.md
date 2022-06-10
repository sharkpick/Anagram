# Anagram

## About
first anagram reads /usr/share/dict/words and maps the anagrams together. this is typically a fairly fast process, it takes ~1 second on my machine. it only does this once at boot though. and really if you needed to you could dump the resultant dictionary into a sql server or something. it's probably not necessary.

anagram then launches a webserver listening at port 8080 that serves the static directory of files. requests from index.js are sent to the server to query the dictionary, and results are returned in JSON format.

## Usage
install by running `go build` from the main package, then launch with `./anagram`.

access the API through http://localhost:8080/ when the API is running.