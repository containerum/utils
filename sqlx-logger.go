package utils

import (
	"github.com/jmoiron/sqlx"
	"io"
	"database/sql"
	"fmt"
)

type SQLXQueryLogger struct {
	q sqlx.Queryer
	w io.Writer
}

func NewSQLXQueryLogger(queryer sqlx.Queryer, logWriter io.Writer) *SQLXQueryLogger {
	return &SQLXQueryLogger{
		q: queryer,
		w: logWriter,
	}
}

func (q *SQLXQueryLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Fprintln(q.w, append([]interface{}{query}, args...)...)
	return q.q.Query(query, args...)
}

func (q *SQLXQueryLogger) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	fmt.Fprintln(q.w, append([]interface{}{query}, args...)...)
	return q.q.Queryx(query, args...)
}

func (q *SQLXQueryLogger) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	fmt.Fprintln(q.w, append([]interface{}{query}, args...)...)
	return q.q.QueryRowx(query, args...)
}

type SQLXExecLogger struct {
	e sqlx.Execer
	w io.Writer
}

func NewSQLXExecLogger(execer sqlx.Execer, logWriter io.Writer) *SQLXExecLogger {
	return &SQLXExecLogger{
		e: execer,
		w: logWriter,
	}
}

func (e *SQLXExecLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	fmt.Fprintln(e.w, append([]interface{}{query}, args...)...)
	return e.e.Exec(query, args...)
}


