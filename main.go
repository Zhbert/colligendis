/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"colligendis/cmd/root"
	"colligendis/internal/common/structs"
	"colligendis/internal/db_service"
	_ "embed"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

//go:embed templates/tex/preamble
var preamble []byte

//go:embed templates/tex/stats.tmpl
var template []byte

//go:embed convert_csv.sh
var convertScript []byte

func main() {

	var tmpls []structs.TemplateStruct

	var preambleStruct structs.TemplateStruct
	preambleStruct.Name = "preamble"
	preambleStruct.Data = preamble
	tmpls = append(tmpls, preambleStruct)

	var statsStruct structs.TemplateStruct
	statsStruct.Name = "stats.tmpl"
	statsStruct.Data = template
	tmpls = append(tmpls, statsStruct)

	var convertScriptStruct structs.TemplateStruct
	convertScriptStruct.Name = "convert_csv.sh"
	convertScriptStruct.Data = convertScript
	tmpls = append(tmpls, convertScriptStruct)

	db, err := gorm.Open(sqlite.Open("colligendis.db"),
		&gorm.Config{Logger: logger.Default.LogMode(db_service.GetLogger())})
	if err != nil {
		log.Fatal("Error opening db!\n")
	}

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	root.Execute(tmpls, db)
}
