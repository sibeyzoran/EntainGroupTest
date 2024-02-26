# Entain Racing Back End Test Project

## Target User Audience
A user of the wagering business who wants to get updates for racing and sporting events.

## Overview
The racing a sporting service utilizes a front end API to serve as an entry point to collecting information about racing and sporting events. As such there are two projects: 
1. The API front end
1. The gRPC server back end

The front end has two endpoints which users can access. They are:
1. /races
1. /sports

Each endpoint has a single GET and POST methods. These being:
1. GET: /races/{id} (where ID equals the ID of a single race)
1. POST: /races/list-races

### Unique to racing

Race objects returned are made up of the following variables:
1. int64 id
1. int64 meeting_id
1. string name
1. int64 number
1. bool visible
1. datetime advertised_start_time
1. string status

Racing POST requests can implement a filter which is made up by:
1. int64 []meeting_ids  - an array of numbers
1. bool visible_only - true to see only visible races, by default is set to false and will return both visible and hidden races
1. string orderBy - allows users to orderBy any variable in a race e.g. advertised_start_time (by default will orderBy this), name or, ID
1. string sort - allows users to sort by ascending or descending order by entering "asc" or "desc"

### Unique to sports

Sport objects returned are made up of the following variables:
1. int64 id
1. string name
1. datetime advertised_start_time
1. string sport
1. string current_score

Sport POST requests can implement a filter which is made up by:
1. int64 []ids - an array of numbers
1. string sport - the name of a sport. Currently these are limited to: Basketball, AFL, Soccer, Hockey and, Rugby League.
1. string orderBy - allows users to orderBy any variable in a race e.g. advertised_start_time (by default will orderBy this), sport or, ID
1. string sort - allows users to sort by ascending or descending order by entering "asc" or "desc"

## How to use

### Setup
1. Clone the code in the repository
1. Install Go (latest).

```bash
brew install go
```

... or [see here](https://golang.org/doc/install).

1. Install `protoc`

```
brew install protobuf
```

... or [see here](https://grpc.io/docs/protoc-installation/).

1. In a terminal window, start our racing service...

```bash
cd ./racing

go build && ./racing
➜ INFO[0000] gRPC server listening on: localhost:9000
```

1. In another terminal window, start our api service...

```bash
cd ./api

go build && ./api
➜ INFO[0000] API server listening on: localhost:8000
```

Now that both the API and the gRPC server are running and listening on their respective ports we can begin sending HTTP requests to the API.

### Using the GET method
There are multiple ways to send HTTP requests to an endpoint. Here I will provide examples using curl - a unix based cmdlet.
Racing example:
1. Open a terminal and type
```bash
curl -X "GET" "http://localhost:8000/v1/races/83" -H 'Content-Type: application/json'
```

You should receive a JSON response that looks similar to:

```JSON
{
    "race": {
        "id": "83",
        "meetingId": "8",
        "name": "Wisconsin bats",
        "number": "10",
        "visible": true,
        "advertisedStartTime": "2021-03-01T18:49:21Z",
        "status": "CLOSED"
}
}
```
Sports example:
1. Open a terminal and type
```bash
curl -X "GET" "http://localhost:8000/v1/sports/13" -H 'Content-Type: application/json'
```

You should receive a JSON response that looks similar to:

```JSON
{
    "sports": {
        "id": "13",
        "name": "Wisconsin bats vs Alabama Ants",
        "sport": "Basketball"
        "advertisedStartTime": "2021-03-01T18:49:21Z",
        "currentScore": "120-103"
}
}
```


### Using the POST method
There are multiple ways to send HTTP requests to an endpoint. Here I will provide examples using curl - a unix base cmdlet. The POST method allows users to create a filter to filter the list to only the results they want. They can narrow the list down by providing an array of meeting ID's as well as only returning races that are visible. The sports endpoint also allows for filtering via ID's and the type of sport.

Racing example:
1. Open a terminal and type:
```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
    -H 'Content-Type: application/json' \
    -d $'{
        "filter":{
            "visibleOnly":true,
            "meetingIds": [5],
            "orderBy": "name",
            "sort": "desc"
        }
}'
```
You should receive a JSON response similar to:
```JSON
{
    "races": {
        "id": "83",
        "meetingId": "8",
        "name": "Wisconsin bats",
        "number": "10",
        "visible": true,
        "advertisedStartTime": "2021-03-01T18:49:21Z",
        "status": "CLOSED"
    }
    {
        "id": "82",
        "meetingId": "1",
        "name": "Alabama ants",
        "number": "123",
        "visible": false,
        "advertisedStartTime": "2021-03-021T18:49:21Z",
        "status": "OPEN"
    }

}
```

Sports Example:
1. Open a terminal and type:
```bash
curl -X "POST" "http://localhost:8000/v1/list-sports" \
    -H 'Content-Type: application/json' \
    -d $'{
        "filter":{
            "sport":"basketball",
            "orderBy": "name",
            "sort": "desc"
        }
}'
```
You should receive a JSON response similar to:
```JSON
{
    "sports": [
        { "id": "33", "name": "Alabama druids VS Oklahoma chimeras", "advertisedStartTime": "2024-02-26T10:02:56Z", "sport": "basketball", "currentScore": "0-0" },
        { "id": "32", "name": "Alaska sons VS Colorado dwarves", "advertisedStartTime": "2024-02-25T04:29:13Z", "sport": "basketball", "currentScore": "0-0" },
        { "id": "35", "name": "Georgia ants VS Indiana witches", "advertisedStartTime": "2024-02-23T15:28:36Z", "sport": "basketball", "currentScore": "107-3" },
        { "id": "9", "name": "Georgia elephants VS Arkansas gnomes", "advertisedStartTime": "2024-02-25T04:42:52Z", "sport": "basketball", "currentScore": "0-0" }
    ]
}
```

## Future implementations:
The major outstanding deficit in these projects are the lack of unit tests. Some tests that will need to be written but haven't yet are as follows:

1. Creating mock responses to each endpoints HTTP requests
1. Creating unit tests that craft HTTP requests to test the GET and POST requests of both the /races and /sports endpoints and then compare them to the expected result.
1. The complete separation of the sports project from racing. This includes creating a new project called sports and extracting the sports code and logic from being entwined with racing to its own completely separate instance.
