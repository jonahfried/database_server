package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"

	_ "github.com/lib/pq"
)

type militaryEquipment struct {
	ids             []int64
	classifications []string
	names           []string
	manIDs          []int64
	rowCount        int
}

func (me militaryEquipment) displayEquipment() {
	for i := 0; i < me.rowCount; i++ {
		fmt.Println(me.ids[i], me.classifications[i], me.names[i], me.manIDs[i])
	}
}

func insert(db *sql.DB, name, classification, manufacturer string) {
	manID, err := getManufacturerByName(db, manufacturer)
	if err == sql.ErrNoRows {
		err := db.QueryRow("INSERT INTO manufacturers (id, name) VALUES(DEFAULT, $1) RETURNING id;", manufacturer).Scan(&manID)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO militaryEquipment (id, name, classification, manufacturerID) VALUES(DEFAULT, $1, $2, $3);", name, classification, manID)
	if err != nil {
		log.Fatal(err)
	}
}

func insertManufacturer(db *sql.DB, name string) {
	_, err := db.Exec("INSERT INTO manufacturers (id, name) VALUES(DEFAULT, $1);", name)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllFromMilitaryEquipment(db *sql.DB) (results militaryEquipment) {
	results.ids = make([]int64, 0)
	results.classifications = make([]string, 0)
	results.names = make([]string, 0)
	results.manIDs = make([]int64, 0)
	var newID int64
	var newClassification string
	var newName string
	var manID int64

	rows, err := db.Query("SELECT * FROM militaryEquipment ORDER BY id;")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		results.rowCount++
		if err := rows.Scan(&newID, &newName, &newClassification, &manID); err != nil {
			log.Fatal(err)
		}
		results.ids = append(results.ids, newID)
		results.classifications = append(results.classifications, newClassification)
		results.names = append(results.names, newName)
		results.manIDs = append(results.manIDs, manID)
	}
	rows.Close()
	return results
}

func getByID(db *sql.DB, id int64) (newID int64, newType, newName string) {

	err := db.QueryRow("SELECT * FROM militaryEquipment WHERE equipmentId = $1;", id).Scan(&newID, &newType, &newName)
	if err == sql.ErrNoRows {
		fmt.Println("There are no rows with ID", id)
	} else if err != nil {
		log.Fatal(err)
	}
	return newID, newType, newName
}

func getManufacturerByID(db *sql.DB, id int64) (name string) {

	err := db.QueryRow("SELECT name FROM Manufacturers WHERE id = $1;", id).Scan(&name)
	if err == sql.ErrNoRows {
		fmt.Println("There are no rows with ID", id)
	} else if err != nil {
		log.Fatal(err)
	}
	return name
}

func getManufacturerByName(db *sql.DB, name string) (id int64, err error) {

	err = db.QueryRow("SELECT id FROM Manufacturers WHERE name = $1;", name).Scan(&id)

	return id, err
}

func getByType(db *sql.DB, equipmentType string) (results militaryEquipment) {
	results.ids = make([]int64, 0)
	results.classifications = make([]string, 0)
	results.names = make([]string, 0)
	results.manIDs = make([]int64, 0)
	var newID int64
	var newType string
	var newName string
	var manID int64

	rows, err := db.Query("SELECT * FROM militaryEquipment WHERE type = $1;", equipmentType)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		results.rowCount++
		if err := rows.Scan(&newID, &newType, &newName, &manID); err != nil {
			log.Fatal(err)
		}
		results.ids = append(results.ids, newID)
		results.classifications = append(results.classifications, newType)
		results.names = append(results.names, newName)
		results.manIDs = append(results.manIDs, manID)
	}
	rows.Close()
	return results
}

func updateInfo(db *sql.DB, id int64, name, classification string) {
	_, err := db.Exec("UPDATE militaryEquipment SET name=$1, classification=$2 WHERE id = $3;", name, classification, id)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteByID(db *sql.DB, id int64) {
	_, err := db.Exec("DELETE FROM militaryEquipment WHERE id = $1;", id)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM warEquipmentPairs WHERE equipmentID = $1;", id)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
)

func clearDatabases(db *sql.DB) {
	_, err := db.Exec("DROP TABLE militaryEquipment, manufacturers, wars, warEquipmentPairs")
	if err != nil {
		log.Print(err)
	}
	_, err = db.Exec("CREATE TABLE militaryEquipment (id serial, name varchar(20), classification varchar(20), manufacturerID varchar(20));")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE Manufacturers (id serial, name varchar(20));")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE wars (id serial, name varchar(20), deaths integer);")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE warEquipmentPairs (warID integer, equipmentID integer);")
	if err != nil {
		log.Fatal(err)
	}
}

func getEquipmentByMan(db *sql.DB, manufacturer string) (results militaryEquipment) {
	manID, err := getManufacturerByName(db, manufacturer)
	if err == sql.ErrNoRows {
		fmt.Printf("No Manufacturer by that name")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	results.ids = make([]int64, 0)
	results.classifications = make([]string, 0)
	results.names = make([]string, 0)
	results.manIDs = make([]int64, 0)
	var newID int64
	var newType string
	var newName string
	var newManID int64

	rows, err := db.Query("SELECT * FROM militaryEquipment WHERE manufacturerID = $1;", manID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		results.rowCount++
		if err := rows.Scan(&newID, &newType, &newName, &newManID); err != nil {
			log.Fatal(err)
		}
		results.ids = append(results.ids, newID)
		results.classifications = append(results.classifications, newType)
		results.names = append(results.names, newName)
		results.manIDs = append(results.manIDs, newManID)
	}
	rows.Close()
	return results
}

func deleteManufacturer(db *sql.DB, name string) {
	var manID int64
	err := db.QueryRow("DELETE FROM manufacturer WHERE name = $1 RETURNING id", name).Scan(&manID)
	if err != nil {
		log.Fatal(err)
	}
	equipmentIdsToDelete := make([]int64, 0)
	var newID int64
	rows, err := db.Query("SELECT id FROM militaryEquipment WHERE manufacturerID = $1", manID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err = rows.Scan(&newID); err != nil {
			log.Fatal(err)
		}
		equipmentIdsToDelete = append(equipmentIdsToDelete, newID)
	}
	_, err = db.Exec("DELETE FROM warEquipmentPairs WHERE equipmentID = any($1)", pq.Array(equipmentIdsToDelete))
	if err != nil {
		log.Fatal(err)
	}
}

func getEquipmentIDFromName(db *sql.DB, name string) (id int64) {
	err := db.QueryRow("SELECT id FROM militaryEquipment WHERE name = $1", name).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func deleteEquipment(db *sql.DB, name string) {
	_, err := db.Exec("DELETE FROM militaryEquipment WHERE name = $1", name)
	if err != nil {
		log.Fatal(err)
	}
	id := getEquipmentIDFromName(db, name)
	_, err = db.Exec("DELETE FROM warEquipmentPairs WHERE equipmentID = $1", id)
	if err != nil {
		log.Fatal(err)
	}
}

func insertWar(db *sql.DB, name string, deaths int64) {
	_, err := db.Exec("INSERT INTO wars (id, name, deaths) VALUES(DEFAULT, $1, $2);", name, deaths)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteWar(db *sql.DB, name string) {
	var warID = getIDFromWarName(db, name)
	_, err := db.Exec("DELETE FROM wars WHERE name = $1", name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM warEquipmentPairs WHERE warID = $1;", warID)
	if err != nil {
		log.Fatal(err)
	}
}

func addWarEquipmentPair(db *sql.DB, warName string, equipmentName string) {
	warID := getIDFromWarName(db, warName)
	var equipmentID int64
	err := db.QueryRow("SELECT id FROM militaryEquipment WHERE name = $1;", equipmentName).Scan(&equipmentID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO warEquipmentPairs (warID, equipmentID) VALUES($1, $2);", warID, equipmentID)
	if err != nil {
		log.Fatal(err)
	}
}

func getIDFromWarName(db *sql.DB, name string) (warID int64) {
	err := db.QueryRow("SELECT id FROM wars WHERE name = $1;", name).Scan(&warID)
	if err != nil {
		log.Fatal(err)
	}
	return warID
}

func getEquipmentByWar(db *sql.DB, warName string) (results militaryEquipment) {
	warID := getIDFromWarName(db, warName)

	equipmentIDsToPull := make([]int64, 0)
	var newToPull int64
	rows, err := db.Query("SELECT equipmentID FROM warEquipmentPairs WHERE warID = $1;", warID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&newToPull); err != nil {
			log.Fatal(err)
		}
		equipmentIDsToPull = append(equipmentIDsToPull, newToPull)
	}

	results.ids = make([]int64, 0)
	results.classifications = make([]string, 0)
	results.names = make([]string, 0)
	results.manIDs = make([]int64, 0)
	var newID int64
	var newType string
	var newName string
	var newManID int64

	rows, err = db.Query("SELECT * FROM militaryEquipment WHERE id = any($1);", pq.Array(equipmentIDsToPull))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		results.rowCount++
		if err := rows.Scan(&newID, &newType, &newName, &newManID); err != nil {
			log.Fatal(err)
		}
		results.ids = append(results.ids, newID)
		results.classifications = append(results.classifications, newType)
		results.names = append(results.names, newName)
		results.manIDs = append(results.manIDs, newManID)
	}
	rows.Close()
	return results
}

func runsuite(db *sql.DB) {
	clearDatabases(db)
	insertManufacturer(db, "Fokker")
	insert(db, "Fokker E IV", "plane", "Fokker")
	insert(db, "big tank", "tank", "TankCo")
	insert(db, "small tank", "tank", "TankCo")
	insert(db, "BigBoy", "bomber", "bombsRus")
	insert(db, "Gotha G V", "bomber", "bombsRus")
	insertWar(db, "bad war", 100)
	addWarEquipmentPair(db, "bad war", "small tank")
	// table := getEquipmentByWar(db, "bad war")
	// // deleteById(db, 10)
	// // table := getByID(db, 3)
	// table.displayEquipment()
}
