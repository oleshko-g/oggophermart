// Package query embeds SQL statements into global variables to query databases
package query

import _ "embed"

// InsertString is the SQL statement to insert a user into a db
//
//go:embed insertUser.sql
var InsertUser string

// SelectUserIDByLogin is the SQL state to select a user id by their login
//
//go:embed selectUserIDByLogin.sql
var SelectUserIDByLogin string

// InsertUserOrder is the the SQL statement to insert an order into a db
//
//go:embed insertOrder.sql
var InsertOrder string
