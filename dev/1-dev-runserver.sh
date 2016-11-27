#this script runs the server

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

source $DIR/dev/dev.env
go run $DIR/packetsender/cmd/server.go
