package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/cznic/ql"
)

var (
	flagDbName = flag.String("dbname", "", "Database (file) name")

	db  *ql.DB
	err error
)

func init() {
	flag.Parse()
}

func main() {

	if "" != *flagDbName {
		db, err = ql.OpenFile(*flagDbName, &ql.Options{})
	} else {
		db, err = ql.OpenMem()
	}
	if nil != err {
		log.Fatalln(err.Error())
	}

	files, err := filepath.Glob("*.ql")
	if nil != err {
		log.Fatalln(err.Error())
	}

	ctx := ql.NewRWCtx()
	for _, file := range files {
		fmt.Println("Reading", file)
		q, err := ioutil.ReadFile(file)
		if nil != err {
			log.Panicln(err.Error())
		}
		if _, _, err = db.Run(ctx, string(q)); nil != err {
			log.Panicln(err.Error())
		}
	}

	info, err := db.Info()
	fmt.Printf("%#v Error=%s\n", info, err)
}
