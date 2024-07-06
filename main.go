/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"colligendis/cmd/root"
	"colligendis/internal/common/structs"
	_ "embed"
)

//go:embed templates/tex/preamble
var preamble []byte

//go:embed templates/tex/stats.tmpl
var template []byte

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

	root.Execute(tmpls)
}
