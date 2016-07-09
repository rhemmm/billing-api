package people

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/py150504/billingps/src/global"
)

// Person : data type people
type Person struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	JoinDate time.Time `json:"-"`
	Status   int       `json:"-"`
}

func (p *Person) save() error {
	p.JoinDate = time.Now()
	p.Status = 1
	resultInsert, errInsert := queryPerson.insert.Exec(p.Name, p.Phone, p.JoinDate, p.Status)
	if errInsert != nil {
		log.Printf(errInsert.Error())
		return nil
	}
	lastID, errResult := resultInsert.LastInsertId()
	if errResult != nil {
		log.Printf(errResult.Error())
		return nil
	}
	p.ID = lastID

	return nil
}

func (p *Person) load() error {
	errSelect := queryPerson.selectPerson.QueryRow(p.ID).Scan(
		&p.ID,
		&p.Name,
		&p.Phone,
		&p.JoinDate)

	if errSelect != nil {
		global.LogError.Printf(errSelect.Error())
		return errSelect
	}

	return nil
}

func (p *Person) delete() error {
	resultDelete, errDelete := queryPerson.delete.Exec(p.ID)
	if errDelete != nil {
		global.LogError.Printf(errDelete.Error())
		return errDelete
	}
	affectedRow, errResult := resultDelete.RowsAffected()
	if errResult != nil {
		global.LogError.Printf(errResult.Error())
		return errResult
	}
	if affectedRow == 0 {
		global.LogError.Printf(fmt.Sprintf("%d", affectedRow))
		return nil
	}
	return nil
}

func (p *Person) update() error {
	resultUpdate, errUpdate := queryPerson.update.Exec(p.Name, p.Phone, p.ID)
	if errUpdate != nil {
		global.LogError.Printf(errUpdate.Error())
		return errUpdate
	}
	affectedRow, errResult := resultUpdate.RowsAffected()
	if errResult != nil {
		global.LogError.Printf(errResult.Error())
		return errResult
	}
	if affectedRow == 0 {
		global.LogError.Printf(fmt.Sprintf("%d", affectedRow))
		return nil
	}
	return nil
}

func getPerson(id int64) *Person {
	person := new(Person)
	person.ID = id
	person.load()

	return person
}

func getPeople() []*Person {
	people := []*Person{}
	rows, errSelect := queryPerson.selectPeople.Query()
	defer rows.Close()
	if errSelect != nil {
		log.Printf(errSelect.Error())
	}
	for rows.Next() {
		person := new(Person)
		errScan := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Phone,
			&person.JoinDate)
		if errScan != nil {
			log.Printf(errScan.Error())
		}
		people = append(people, person)
	}
	return people
}

// MapPerson : map person
func MapPerson(p *Person, detail bool) interface{} {
	var attributes interface{}
	if detail {
		attributes = map[string]interface{}{
			"name":      p.Name,
			"phone":     p.Phone,
			"join_date": p.JoinDate.Format("02 January 2006, 15:04"),
		}
	}
	person := map[string]interface{}{
		"id":         strconv.FormatInt(p.ID, 10),
		"type":       "person",
		"attributes": attributes,
	}
	return person
}
