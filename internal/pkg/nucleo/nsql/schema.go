package nsql

import (
	"fmt"
	"strings"
)

type Schema struct {
	TableName string
	Columns   []string
	// Private
	autoIncrement bool
	primaryKey    string
	alias         string
	// Generated
	allColumns         string
	insertColumns      string
	insertNamedColumns string
	updateNamedColumns string
}

func (s *Schema) SelectAllColumns() string {
	return s.allColumns
}

func (s *Schema) InsertColumns() string {
	return s.insertColumns
}

func (s *Schema) InsertNamedColumns() string {
	return s.insertNamedColumns
}

func (s *Schema) UpdateNamedColumns() string {
	return s.updateNamedColumns
}

func (s *Schema) SelectJoinColumns() string {
	return CreateSelectJoinColumns(s.alias, s.Columns)
}

func (s *Schema) SelectAllExplicitColumns() string {
	return CreateSelectExplicitColumns(s.TableName, s.Columns)
}

func (s *Schema) JoinAlias() string {
	return CreateTableAlias(s.TableName, s.alias)
}

func (s *Schema) TableAlias() string {
	return fmt.Sprintf(`"%s"`, s.alias)
}

func (s *Schema) TableQuery() string {
	return fmt.Sprintf(`"%s"`, s.TableName)
}

func (s *Schema) FindById() string {
	q := fmt.Sprintf(`SELECT %s FROM "%s" WHERE "%s" = $1`, s.SelectAllColumns(), s.TableName, s.primaryKey)
	return q
}

func NewSchema(tableName string, columns []string, args ...SchemaOption) *Schema {
	// Evaluate options
	o := evaluateSchemaOptions(args)

	// Init schema
	s := Schema{
		TableName:     tableName,
		Columns:       columns,
		autoIncrement: o.autoIncrement,
		primaryKey:    o.primaryKey,
		alias:         o.alias,
	}

	// Create all column query
	s.allColumns = CreateColumns(s.Columns)

	// Create insert query
	var insertColumns []string
	if o.autoIncrement {
		// Remove PK from columns list
		insertColumns = make([]string, len(s.Columns)-1)
		i := 0
		for _, col := range s.Columns {
			if s.primaryKey == col {
				continue
			}
			insertColumns[i] = col
			i += 1
		}
	} else {
		insertColumns = s.Columns
	}

	s.insertNamedColumns = CreateNamedColumns(insertColumns)
	s.insertColumns = CreateColumns(insertColumns)
	s.updateNamedColumns = CreateSetNamedColumns(insertColumns)

	return &s
}

// --------------
// Schema Options
// --------------

type schemaOption struct {
	alias         string
	primaryKey    string
	autoIncrement bool
}

var defaultSchemaOption = &schemaOption{
	alias:         "",
	primaryKey:    "id",
	autoIncrement: true,
}

type SchemaOption = func(*schemaOption)

func evaluateSchemaOptions(opts []SchemaOption) *schemaOption {
	optCopy := &schemaOption{}
	*optCopy = *defaultSchemaOption
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

func WithAlias(str string) SchemaOption {
	return func(o *schemaOption) {
		o.alias = str
	}
}

func WithAutoIncrement(ai bool) SchemaOption {
	return func(o *schemaOption) {
		o.autoIncrement = ai
	}
}

func WithPrimaryKey(pk string) SchemaOption {
	return func(o *schemaOption) {
		o.primaryKey = pk
	}
}

// ---------
// Utilities
// ---------

func CreateColumns(cols []string) string {
	result := make([]string, len(cols))
	for i, v := range cols {
		result[i] = fmt.Sprintf(`"%s"`, v)
	}
	return strings.Join(result, ", ")
}

func CreateSelectJoinColumns(alias string, columns []string) string {
	// Add alias
	result := make([]string, len(columns))
	for i, col := range columns {
		col = strings.TrimSpace(col)
		result[i] = fmt.Sprintf(`"%s"."%s" AS "%s.%s"`, alias, col, alias, col)
	}
	return strings.Join(result, ", ")
}

func CreateSelectExplicitColumns(alias string, columns []string) string {
	// Add alias
	result := make([]string, len(columns))
	for i, col := range columns {
		col = strings.TrimSpace(col)
		result[i] = fmt.Sprintf(`"%s"."%s"`, alias, col)
	}
	return strings.Join(result, ", ")
}

func CreateNamedColumns(columns []string) string {
	// Add alias
	result := make([]string, len(columns))
	for i, col := range columns {
		col = strings.TrimSpace(col)
		result[i] = fmt.Sprintf(`:%s`, col)
	}
	return strings.Join(result, ", ")
}

func CreateTableAlias(tableName string, alias string) string {
	return fmt.Sprintf(`"%s" as %s`, tableName, alias)
}

func CreateSetNamedColumns(columns []string) string {
	// Add alias
	result := make([]string, len(columns))
	for i, col := range columns {
		col = strings.TrimSpace(col)
		result[i] = fmt.Sprintf(`"%s" = :%s`, col, col)
	}
	return strings.Join(result, ", ")
}
