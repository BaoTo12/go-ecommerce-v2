package query

import (
	"fmt"
	"strings"
)

// Builder provides a fluent API for building SQL queries
type Builder struct {
	operation   string
	table       string
	columns     []string
	values      []interface{}
	where       []string
	whereArgs   []interface{}
	orderBy     []string
	limit       int
	offset      int
	joins       []string
	groupBy     []string
	having      string
	havingArgs  []interface{}
	returning   []string
}

// Select creates a SELECT query builder
func Select(columns ...string) *Builder {
	return &Builder{
		operation: "SELECT",
		columns:   columns,
	}
}

// Insert creates an INSERT query builder
func Insert(table string) *Builder {
	return &Builder{
		operation: "INSERT",
		table:     table,
	}
}

// Update creates an UPDATE query builder
func Update(table string) *Builder {
	return &Builder{
		operation: "UPDATE",
		table:     table,
	}
}

// Delete creates a DELETE query builder
func Delete(table string) *Builder {
	return &Builder{
		operation: "DELETE",
		table:     table,
	}
}

// From sets the table for SELECT/DELETE
func (b *Builder) From(table string) *Builder {
	b.table = table
	return b
}

// Set adds column-value pairs for INSERT/UPDATE
func (b *Builder) Set(column string, value interface{}) *Builder {
	b.columns = append(b.columns, column)
	b.values = append(b.values, value)
	return b
}

// Where adds a WHERE condition
func (b *Builder) Where(condition string, args ...interface{}) *Builder {
	b.where = append(b.where, condition)
	b.whereArgs = append(b.whereArgs, args...)
	return b
}

// AndWhere adds an AND WHERE condition
func (b *Builder) AndWhere(condition string, args ...interface{}) *Builder {
	return b.Where(condition, args...)
}

// OrWhere adds an OR WHERE condition
func (b *Builder) OrWhere(condition string, args ...interface{}) *Builder {
	if len(b.where) > 0 {
		b.where[len(b.where)-1] = "(" + b.where[len(b.where)-1] + ") OR (" + condition + ")"
	} else {
		b.where = append(b.where, condition)
	}
	b.whereArgs = append(b.whereArgs, args...)
	return b
}

// Join adds a JOIN clause
func (b *Builder) Join(join string) *Builder {
	b.joins = append(b.joins, "JOIN "+join)
	return b
}

// LeftJoin adds a LEFT JOIN clause
func (b *Builder) LeftJoin(join string) *Builder {
	b.joins = append(b.joins, "LEFT JOIN "+join)
	return b
}

// OrderBy adds ORDER BY columns
func (b *Builder) OrderBy(columns ...string) *Builder {
	b.orderBy = append(b.orderBy, columns...)
	return b
}

// Limit sets the LIMIT
func (b *Builder) Limit(limit int) *Builder {
	b.limit = limit
	return b
}

// Offset sets the OFFSET
func (b *Builder) Offset(offset int) *Builder {
	b.offset = offset
	return b
}

// GroupBy adds GROUP BY columns
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.groupBy = append(b.groupBy, columns...)
	return b
}

// Having adds a HAVING clause
func (b *Builder) Having(condition string, args ...interface{}) *Builder {
	b.having = condition
	b.havingArgs = args
	return b
}

// Returning adds a RETURNING clause (PostgreSQL)
func (b *Builder) Returning(columns ...string) *Builder {
	b.returning = columns
	return b
}

// Build generates the SQL query and arguments
func (b *Builder) Build() (string, []interface{}) {
	var sql strings.Builder
	var args []interface{}
	argIndex := 1

	switch b.operation {
	case "SELECT":
		sql.WriteString("SELECT ")
		if len(b.columns) == 0 {
			sql.WriteString("*")
		} else {
			sql.WriteString(strings.Join(b.columns, ", "))
		}
		sql.WriteString(" FROM ")
		sql.WriteString(b.table)

	case "INSERT":
		sql.WriteString("INSERT INTO ")
		sql.WriteString(b.table)
		sql.WriteString(" (")
		sql.WriteString(strings.Join(b.columns, ", "))
		sql.WriteString(") VALUES (")
		placeholders := make([]string, len(b.values))
		for i := range b.values {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			argIndex++
		}
		sql.WriteString(strings.Join(placeholders, ", "))
		sql.WriteString(")")
		args = append(args, b.values...)

	case "UPDATE":
		sql.WriteString("UPDATE ")
		sql.WriteString(b.table)
		sql.WriteString(" SET ")
		sets := make([]string, len(b.columns))
		for i, col := range b.columns {
			sets[i] = fmt.Sprintf("%s = $%d", col, argIndex)
			argIndex++
		}
		sql.WriteString(strings.Join(sets, ", "))
		args = append(args, b.values...)

	case "DELETE":
		sql.WriteString("DELETE FROM ")
		sql.WriteString(b.table)
	}

	// Joins
	for _, join := range b.joins {
		sql.WriteString(" ")
		sql.WriteString(join)
	}

	// Where
	if len(b.where) > 0 {
		sql.WriteString(" WHERE ")
		for i, w := range b.where {
			if i > 0 {
				sql.WriteString(" AND ")
			}
			// Replace ? with $N for PostgreSQL
			replaced := w
			for strings.Contains(replaced, "?") {
				replaced = strings.Replace(replaced, "?", fmt.Sprintf("$%d", argIndex), 1)
				argIndex++
			}
			sql.WriteString(replaced)
		}
		args = append(args, b.whereArgs...)
	}

	// Group By
	if len(b.groupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupBy, ", "))
	}

	// Having
	if b.having != "" {
		sql.WriteString(" HAVING ")
		sql.WriteString(b.having)
		args = append(args, b.havingArgs...)
	}

	// Order By
	if len(b.orderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderBy, ", "))
	}

	// Limit/Offset
	if b.limit > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", b.limit))
	}
	if b.offset > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", b.offset))
	}

	// Returning
	if len(b.returning) > 0 {
		sql.WriteString(" RETURNING ")
		sql.WriteString(strings.Join(b.returning, ", "))
	}

	return sql.String(), args
}

// String returns just the SQL string
func (b *Builder) String() string {
	sql, _ := b.Build()
	return sql
}
