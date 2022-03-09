package nsql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type DatabaseContext struct {
	conn *sqlx.DB
	ctx  context.Context
}

// Prepare prepare sql statements or exit app if fails or error
func (s *DatabaseContext) Prepare(query string) *sqlx.Stmt {
	stmt, err := s.conn.PreparexContext(s.ctx, query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%w)", query, err))
	}
	return stmt
}

// PrepareFmt prepare sql statements from string format or exit app if fails or error
func (s *DatabaseContext) PrepareFmt(queryFmt string, args ...interface{}) *sqlx.Stmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.Prepare(query)
}

// PrepareNamedFmt prepare sql statements from string format with named bindvars or exit app if fails or error
func (s *DatabaseContext) PrepareNamedFmt(queryFmt string, args ...interface{}) *sqlx.NamedStmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.PrepareNamed(query)
}

// PrepareNamed prepare sql statements with named bindvars or exit app if fails or error
func (s *DatabaseContext) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := s.conn.PrepareNamedContext(s.ctx, query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%w)", query, err))
	}
	return stmt
}

// PrepareTemplate prepare sql statements from a string template format or exit app if fails or error
func (s *DatabaseContext) PrepareTemplate(q string, values map[string]string) *sqlx.Stmt {
	for a, v := range values {
		q = strings.ReplaceAll(q, ":"+a, v)
	}
	return s.Prepare(q)
}

// PrepareFmtRebind prepare sql statements from string format and rebind variable or exit app if fails or error
func (s *DatabaseContext) PrepareFmtRebind(queryFmt string, args ...interface{}) *sqlx.Stmt {
	query := fmt.Sprintf(queryFmt, args...)
	query = s.conn.Rebind(query)
	return s.Prepare(query)
}

// PrepareNamedFmtRebind prepare sql statements from string format with named bindvars or exit app if fails or error
func (s *DatabaseContext) PrepareNamedFmtRebind(queryFmt string, args ...interface{}) *sqlx.NamedStmt {
	query := fmt.Sprintf(queryFmt, args...)
	query = s.conn.Rebind(query)
	return s.PrepareNamed(query)
}
