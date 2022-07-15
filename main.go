package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/karantan/gofp"
)

const DATA = "gocsv.csv"

type User struct {
	Id        int
	FirstName string
	LastName  string
	Age       int
	State     string
}

func main() {
	users := loadCSV(DATA)
	gofp.ForEach(func(u User) User {
		fmt.Println(u)
		return u
	}, users)

	dumpCSV("gocsv_copy.csv", users)
}

func dumpCSV(filename string, users []User) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		log.Fatalln("error opening a file", err)
		return
	}

	w := csv.NewWriter(f)
	for _, u := range users {
		if err := w.Write(userToRow(u)); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	// Write any buffered data to the file
	w.Flush()
	// Check for errors occurred during the Flush
	if err := w.Error(); err != nil {
		log.Fatalln(err)
	}
}

// loadCSV opens file `filename`, reads the data in it and returns populated slice of
// `User` structs.
// The given `filename` csv file needs to fit the `User` struct.
func loadCSV(filename string) (users []User) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("error opening a file", err)
		return
	}

	r := csv.NewReader(f)
	for {
		// Note: you might want to skip the first line (csv header)
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
			return
		}
		users = append(users, rowToUser(row))
	}
	return
}

// rowToUser transforms slice of strings to `User` struct.
// E.g. []string{"3", "Jim", "Todd", "43", "WV"} -> User{3, "Jim", "Todd", 43, "WV"}
func rowToUser(row []string) User {

	id, _ := strconv.Atoi(row[0])
	age, _ := strconv.Atoi(row[3])

	return User{
		Id:        id,
		FirstName: row[1],
		LastName:  row[2],
		Age:       age,
		State:     row[4],
	}
}

// userToRow transforms `User` struct to a slice of strings. It's the inverse operation
// of the `rowToUser`.
func userToRow(u User) []string {
	return []string{fmt.Sprint(u.Id), u.FirstName, u.LastName, fmt.Sprint(u.Age), u.State}
}
