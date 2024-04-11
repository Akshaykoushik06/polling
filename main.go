package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	username := ""
	password := ""
	host := "akshaysqlserver.mysql.database.azure.com"
	port := 3306
	database := "testdb"

	// build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	// Open the connection
	_db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	DB = _db
}

func createAzVM(serverID int) {
	fmt.Println("creating VM")
	_, err := DB.Exec("UPDATE p_servers SET status = 'TODO' WHERE server_id = ?;", serverID)
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	_, err = DB.Exec("UPDATE p_servers SET status = 'IN PROGRESS' WHERE server_id = ?;", serverID)
	if err != nil {
		panic(err)
	}
	fmt.Println("VM creation in progress")
	time.Sleep(5 * time.Second)

	_, err = DB.Exec("UPDATE p_servers SET status = 'DONE' WHERE server_id = ?;", serverID)
	if err != nil {
		panic(err)
	}
	fmt.Println("VM creation completed")
}

func main() {
	ge := gin.Default()

	ge.POST("/servers", func(ctx *gin.Context) {
		go createAzVM(1)
		ctx.JSON(200, map[string]interface{}{"submitted": "ok"})
	})

	ge.GET("/short/status/:server_id", func(ctx *gin.Context) {
		serverID := ctx.Param("server_id")

		var status string
		row := DB.QueryRow("SELECT status from p_servers WHERE server_id = ?;", serverID)
		if row.Err() != nil {
			panic(row.Err())
		}
		row.Scan(&status)

		ctx.JSON(200, map[string]interface{}{"status": status})
	})

	ge.GET("/long/status/:server_id", func(ctx *gin.Context) {
		serverID := ctx.Param("server_id")
		currentStatus := ctx.Query("status")

		var status string
		for {
			row := DB.QueryRow("SELECT status FROM p_servers WHERE server_id = ?;", serverID)
			if row.Err() != nil {
				panic(row.Err())
			}

			row.Scan(&status)

			if currentStatus != status {
				break
			}

			time.Sleep(1 * time.Second)
		}

		ctx.JSON(200, map[string]interface{}{"status": status})
	})

	ge.Run(":9000")
}
