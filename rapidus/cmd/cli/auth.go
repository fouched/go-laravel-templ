package main

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

func doAuth() error {
	// migrations
	dbType := rap.DB.Type
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := rap.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := rap.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth-tables."+dbType+".sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens cascade;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	// run migrations
	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	// copy files
	err = copyFileFromTemplate("templates/data/user.go.txt", rap.RootPath+"/data/user.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", rap.RootPath+"/data/token.go")
	if err != nil {
		exitGracefully(err)
	}

	// copy middleware
	err = copyFileFromTemplate("templates/middleware/auth.go.txt", rap.RootPath+"/middleware/auth.go")
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/middleware/auth-token.go.txt", rap.RootPath+"/middleware/auth-token.go")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("  - users, tokens and remember_tokens migrations created and executed")
	color.Yellow("  - user and token models created")
	color.Yellow("  - auth middleware created")
	color.Yellow("")
	color.Yellow("  - Don't forget to add user and token models in data/models.go, and to add appropriate middleware to your routes!")

	return nil
}
