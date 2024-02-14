package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nlsra/message-system/pkg/config"
)

func main() {
	config.Run()
}
