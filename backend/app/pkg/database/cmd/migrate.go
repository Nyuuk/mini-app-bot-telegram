package main

import "github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"

func main(){
	database.RunMigration()
}
