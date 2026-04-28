package repository

import (
	"strconv"
	"strings"
)

type Scanner interface {
	Scan(dest ...any) error
}

type PatchBuilder struct {
	args  []any
	argID int
	sb    *strings.Builder
}

func NewPatchBuilder() *PatchBuilder {
	return &PatchBuilder{args: []any{}, argID: 1, sb: new(strings.Builder)}
}

func (c *PatchBuilder) Head(tableName string) {
	c.sb.WriteString("UPDATE ")
	c.sb.WriteString(tableName)
	c.sb.WriteString(" SET ")
}

func (c *PatchBuilder) Add(name string, value any) {
	if c.argID > 1 {
		c.sb.WriteString(", ")
	}

	c.sb.WriteString(name)
	c.sb.WriteString(" = $")
	c.sb.WriteString(strconv.Itoa(c.argID))

	c.argID++
	c.args = append(c.args, value)
}

// Use once at the end
func (c *PatchBuilder) Where(k string, v any) {
	c.sb.WriteString(" WHERE ")
	c.sb.WriteString(k)
	c.sb.WriteString(" = $")
	c.sb.WriteString(strconv.Itoa(c.argID))
	c.args = append(c.args, v)
}

func (c *PatchBuilder) String() string {
	return c.sb.String()
}

func (c *PatchBuilder) Len() int {
	return len(c.args)
}

func (c *PatchBuilder) Args() []any {
	return c.args
}
