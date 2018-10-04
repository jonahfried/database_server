package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if parts[1] == "api" {
		dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			dbUser, dbPassword, dbName,
		)
		db, err := sql.Open("postgres", dbinfo)
		if err != nil {
			log.Fatal(err)
		}
		if parts[2] == "table" && parts[3] == "militaryequipment" {
			table := getAllFromMilitaryEquipment(db)
			// fmt.Println(table)
			json := "["
			for i := 0; i < table.rowCount; i++ {
				json += fmt.Sprintf("{\"id\":%d, \"name\":\"%s\", \"classification\":\"%s\", \"manID\":%d},", table.ids[i], table.names[i], table.classifications[i], table.manIDs[i])
			}
			json = json[:len(json)-1] + "]"
			// fmt.Println(json)
			fmt.Fprint(w, json)
		} else if parts[2] == "table" && parts[3] == "manufacturers" {
			table := getAllFromManufacturers(db)
			// fmt.Println(table)
			json := "["
			for i := 0; i < table.rowCount; i++ {
				json += fmt.Sprintf("{\"id\":%d, \"name\":\"%s\"},", table.ids[i], table.names[i])
			}
			json = json[:len(json)-1] + "]"
			// fmt.Println(json)
			fmt.Fprint(w, json)
		} else if parts[2] == "update" {
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				log.Fatal(err)
			}
			newName := parts[4]
			newClassification := parts[5]
			updateInfo(db, int64(id), newName, newClassification)
			fmt.Fprint(w, "true")
		} else if parts[2] == "delete" {
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				log.Fatal(err)
			}
			deleteByID(db, int64(id))
			fmt.Fprint(w, "true")
		} else if parts[2] == "deletemanufacturer" {
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				log.Fatal(err)
			}
			deleteManufacturerByID(db, int64(id))
			fmt.Fprint(w, "true")
		} else if parts[2] == "add" {
			name := parts[3]
			classification := parts[4]
			manufacturer := parts[5]
			insert(db, name, classification, manufacturer)
		} else if parts[2] == "addmanufacturer" {
			name := parts[3]
			insertManufacturer(db, name)
		} else if parts[2] == "updatemanufacturers" {
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				log.Fatal(err)
			}
			name := parts[4]
			updateManufacturer(db, int64(id), name)
		} else if parts[2] == "post" {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Print(err)
			}
			fileText := string(bytes)
			fileLines := strings.Split(fileText, "\n")
			fileLines = fileLines[4 : len(fileLines)-2]

			readCsvData(db, fileLines)
		}
	} else if parts[1] == "militaryEquipment" {
		http.ServeFile(w, r, "militaryEquipment.html")
	} else if parts[1] == "manufacturers" {
		http.ServeFile(w, r, "manufacturers.html")
	} else if parts[1] == "loaddata" {
		http.ServeFile(w, r, "loaddata.html")
	} else {
		http.ServeFile(w, r, "index.html")
	}
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName,
	)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	runsuite(db)

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", requestHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
