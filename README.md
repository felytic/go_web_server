# Go web server
Test task

## Instalation
Install additional go packages:

`$ go get github.com/joho/godotenv`

`$ go get github.com/mattn/go-sqlite3`

Clone current project:

`$ git clone https://github.com/felytic/go_web_server.git`

Enter project folder:

`$ cd cd go_web_server`

## Usage

Run server for mocking requests:

`$ go run mocked/main.go`

In separate console run main server:

`$ go run main.go`

Browse http://0.0.0.0:8081/

Select any plan. Back-end will try to reach each of four mocked servers (Apple, Google, PayPal, Stripe). Each server has 15% chance to fail, it means that in about 50% of all requests main server will not collect all data, so  keep refreshing to see the difference between OK and failure.

You can simulate DB error by removing/renaming `db.sqlite3` file.

## Testing
Run:
`$ go test`

__TODO__ Remove app logs from tests
