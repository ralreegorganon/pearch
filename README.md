# PEARCH

PilotEdge live data archiver.

## Use it

Create a new database user, and give it the password of ````flylikeaneagle```` when prompted

    createuser -P pearch

Create a new database
    
    createdb -O pearch pearch

export PEARCH_CONNECTION_STRING="user=pearch password=flylikeaneagle dbname=pearch sslmode=disable"


Install goose to run the database migrations
    
    go get bitbucket.org/liamstask/goose/cmd/goose

Run the migrations.
    goose up

Run it locally

    go build
    ./pearch
