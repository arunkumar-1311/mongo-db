package main

import (
	dbconnection "github.com/arunkumar-1311/mongo-db/dbconnection"
	"github.com/arunkumar-1311/mongo-db/routers"
)

func main() {
	dbconnection.DBconnetion()
	defer dbconnection.Disconnect()

	routers.Router()
}
