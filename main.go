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

const (
	inDir   = `C:\Majestic Software\KHS\data`
	outFile = "dbfdump.json"
)

func main() {
	// skiplist, dbf dir, output dir args
	out := flag.String("out", inDir, "output directory")
	dir := flag.String("dir", inDir, "directory containing DBF files")
	skip := flag.String("skip", "dictfinal,outlines,carriers,smtp,letcfg,oclm", "DBF files to ignore")
	ui := flag.Bool("ui", true, "show GUI")
	flag.Parse()

	if *ui {
		showUI(*out, *dir, *skip)
	} else {
		err := processDBF(*dir, *out, skipFiles(*skip))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func skipFiles(skip string) map[string]bool {
	items := strings.Split(skip, ",")
	skipMap := make(map[string]bool, len(items))
	for _, item := range items {
		skipMap[strings.ToLower(item)] = true
	}

	return skipMap
}

func processDBF(dbfDir, outputDir string, skip map[string]bool) error {
	files, err := ioutil.ReadDir(dbfDir)
	if err != nil {
		return err
	}

	outPath := filepath.Join(outputDir, outFile)
	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("unable to open output file %s: %v", outPath, err)
	}
	defer out.Close()

	jsonOut := make(map[string]interface{})

	dbfCount := 0
	for _, file := range files {
		// skip non-dbf files
		if strings.ToLower(filepath.Ext(file.Name())) != ".dbf" {
			continue
		}
		// skip files we were told to ignore
		withoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		if skip[strings.ToLower(withoutExt)] {
			continue
		}
		dbfPath := filepath.Join(dbfDir, file.Name())
		db, err := dbf.OpenFile(dbfPath, new(dbf.Win1250Decoder))
		if err != nil {
			fmt.Printf("error opening %s: %v", file.Name(), err)
			continue
		}

		dbfCount++
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

	if dbfCount == 0 {
		return fmt.Errorf("no dbf files found in %s", dbfDir)
	}

	je := json.NewEncoder(out)
	err = je.Encode(jsonOut)
	if err != nil {
		return fmt.Errorf("unable to write output: %v", err)
	}
	return nil
}
