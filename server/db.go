package main

import (
	"database/sql"
	"strconv"
)

// QueryArgument is an argument for the sql query
type QueryArgument struct {
	Connector string // and, or
	Field     string
	Op        string
	Value     interface{}
}

func Query(cr *sql.DB, arguments []QueryArgument) User {
	var data dbUser
	var args []interface{}

	query := "SELECT id,name,email,total_balance,create_date,last_update FROM users"

	for i := range arguments {
		arg := arguments[i]
		var subQuery string
		if len(arg.Connector) > 0 {
			subQuery = " " + arg.Connector + " " + arg.Field + " " + arg.Op + " $" + strconv.Itoa(i+1)
		} else {
			subQuery = " " + arg.Field + " " + arg.Op + " $" + strconv.Itoa(i+1)
		}

		args = append(args, arg.Value)

		if i == 0 {
			query += " WHERE" + subQuery
		} else {
			query += subQuery
		}
	}

	err := cr.QueryRow(query, args...).Scan(
		&data.ID,
		&data.Name,
		&data.Email,
		&data.TotalBalance,
		&data.CreateDate,
		&data.LastUpdate,
	)

	if err != nil {
		panic(err)
	}

	res := User{Data: &data}

	return res
}
