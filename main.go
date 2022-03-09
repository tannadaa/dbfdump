package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	dbf "github.com/tannadaa/go-foxpro-dbf"
)

func main() {
	// skiplist, dbf dir, output file args
	out := flag.String("out", "dbfdump.json", "output file")
	dir := flag.String("dir", "", "directory containing DBF files")
	skip := flag.String("skip", "", "DBF files to ignore")
	flag.Parse()

	if *out == "" {
		fmt.Println("output file is required")
		os.Exit(1)
	}
	if *dir == "" {
		fmt.Println("dbf directory is required")
		os.Exit(1)
	}

	err := processDBF(*dir, *out, skipFiles(*skip))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func skipFiles(skip string) map[string]struct{} {
	items := strings.Split(skip, ",")
	skipMap := make(map[string]struct{}, len(items))
	for _, item := range items {
		skipMap[strings.ToLower(item)] = struct{}{}
	}

	return skipMap
}

func processDBF(dbfDir, outputFile string, skip map[string]struct{}) error {
	files, err := ioutil.ReadDir(dbfDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no DBF files found")
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("unable to open output file %s: %v", outputFile, err)
	}
	defer out.Close()

	jsonOut := make(map[string]interface{})

	for _, file := range files {
		// skip non-dbf files
		if strings.ToLower(filepath.Ext(file.Name())) != ".dbf" {
			continue
		}
		// skip files we were told to ignore
		withoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		if _, ok := skip[strings.ToLower(withoutExt)]; ok {
			continue
		}
		dbfPath := filepath.Join(dbfDir, file.Name())
		db, err := dbf.OpenFile(dbfPath, new(dbf.Win1250Decoder))
		if err != nil {
			fmt.Printf("error opening %s: %v", file.Name(), err)
			continue
		}

		topKey := strings.ToLower(strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name())))
		jsRecords := make([]map[string]interface{}, 0)

		var i uint32
		for i = 0; i < db.NumRecords(); i++ {
			m, err := db.RecordToMap(i)
			if err != nil {
				fmt.Printf("error converting record to map: %v", err)
				continue
			}
			jsRecord := make(map[string]interface{}, len(m))
			for k, v := range m {
				key := strings.ToLower(k)
				if str, ok := v.(string); ok {
					// get rid of padding
					jsRecord[key] = strings.TrimSpace(str)
				} else {
					jsRecord[key] = v
				}
			}

			jsRecords = append(jsRecords, jsRecord)
		}
		jsonOut[topKey] = jsRecords
		db.Close()
	}

	je := json.NewEncoder(out)
	err = je.Encode(jsonOut)
	if err != nil {
		return fmt.Errorf("unable to write output: %v", err)
	}
	return nil
}
