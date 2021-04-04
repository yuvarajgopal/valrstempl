package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var flagPath = flag.String("templatepath", "", "Salt template path of the Salt Formula")
var appname = flag.String("appname", "", "Application Name of the Salt Formula")
var apptype = flag.String("apptype", "", "Application Type (engine/service/loader)")

func partialrename(templatepath string, f os.FileInfo, err error) (e error) {

	// check each file if starts with the word "dumb_"
	//fmt.Println(f.Name())
	if strings.HasPrefix(f.Name(), "appname_") {
		base := filepath.Base(templatepath) // get the file's basename
		dir := filepath.Dir(templatepath)

		renameto := filepath.Join(dir, strings.Replace(base, "appname_", *appname, 1))
		os.Rename(templatepath, renameto)
	}
	return
}

func init() {
	flag.Parse()
}

func visit(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	matched, err := filepath.Match("*.*", fi.Name())

	if err != nil {
		panic(err)
	}

	if matched {
		read, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		// fmt.Println(string(read))
		// fmt.Println(path)

		//newContents := strings.Replace(string(read), "<APPNAME>", os.Args[2], -1)
		StrReplace := strings.NewReplacer("<APPNAME>", *appname, "<APPTYPE>", *apptype)
		newContents := StrReplace.Replace(string(read))

		// fmt.Println(newContents)

		err = ioutil.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			panic(err)
		}

	}

	return nil
}

func main() {

	if *flagPath == "" {
		flag.Usage() // if no, prompt usage
		os.Exit(0)   // and exit
	}

	folderinfo, err1 := os.Stat(*flagPath)
	if os.IsNotExist(err1) {
		fmt.Println("")
		log.Fatal("\t\t Appname Folder does not exist. Please clone the template files\n\n")

	}
	log.Println(folderinfo)

	if *appname == "" {
		flag.Usage() // if no, prompt usage
		os.Exit(0)   // and exit
	}

	if *apptype == "" {
		flag.Usage() // if no, prompt usage
		os.Exit(0)   // and exit
	}

	// walk through the files in the given path and perform partialrename()
	// function
	filepath.Walk(*flagPath, partialrename)

	err := filepath.Walk(*flagPath, visit)
	if err != nil {
		panic(err)
	}
	//fmt.Println(*flagPath)
	//fmt.Println(*appname)
	os.Rename(*flagPath, *appname)
}
