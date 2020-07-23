package db

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Insertable map[string]json.RawMessage

/*
All statements checks are evaluated here before going to the API class
*/
func Insert(entity string, insertable Insertable) (err error) {
	conn := NextPoolCon()
	defer conn.Close()

	cols := ""
	binds := ""
	colID := 0
	vals := make([]interface{}, 0)
	for colName, val := range insertable {
		va := string(val)
		if va[:1] == "\"" {
			va = va[1 : len(va)-1]
			vals = append(vals, va)
		} else {
			vals = append(vals, val)
		}
		colID++
		if cols == "" {
			cols = colName
		} else {
			cols += ", " + colName
		}
		if binds == "" {
			binds = "$" + strconv.Itoa(colID)
		} else {
			binds += ", $" + strconv.Itoa(colID)
		}
	}

	sql := fmt.Sprintf("insert into %s (%s) values(%s)", entity, cols, binds)
	if _, err = conn.Exec(sql, vals...); err != nil {
		logrus.Errorln(err)
		return
	}
	return

}

func Update(entity, field, id string, insertable Insertable) (err error) {
	conn := NextPoolCon()
	defer conn.Close()

	cols := ""
	colID := 0
	vals := make([]interface{}, 0)
	for colName, val := range insertable {
		colID++
		if cols == "" {
			cols = colName + "=$" + strconv.Itoa(colID)
		} else {
			cols += ", " + colName + "=$" + strconv.Itoa(colID)
		}

		va := string(val)
		if va[:1] == "\"" {
			va = va[1 : len(va)-1]
			vals = append(vals, va)
		} else {
			vals = append(vals, val)
		}

	}

	colID++
	sql := fmt.Sprintf("update %s set %s where %s=%s", entity, cols, field, "$"+strconv.Itoa(colID))
	vals = append(vals, id)
	if _, err = conn.Exec(sql, vals...); err != nil {
		logrus.Errorln(err)
		return
	}
	return

}

func Delete(entity, field, id string) (err error) {
	conn := NextPoolCon()
	defer conn.Close()

	sql := fmt.Sprintf("delete from %s where %s=$1", entity, field)
	if _, err = conn.Exec(sql, id); err != nil {
		logrus.Errorln(err)
		return
	}
	return

}
