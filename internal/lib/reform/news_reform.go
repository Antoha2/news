// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package reform

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type reformNewsTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *reformNewsTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("news").
func (v *reformNewsTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *reformNewsTableType) Columns() []string {
	return []string{
		"news_id",
		"title",
		"content",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *reformNewsTableType) NewStruct() reform.Struct {
	return new(ReformNews)
}

// NewRecord makes a new record for that table.
func (v *reformNewsTableType) NewRecord() reform.Record {
	return new(ReformNews)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *reformNewsTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// ReformNewsTable represents news view or table in SQL database.
var ReformNewsTable = &reformNewsTableType{
	s: parse.StructInfo{
		Type:    "ReformNews",
		SQLName: "news",
		Fields: []parse.FieldInfo{
			{Name: "Id", Type: "int", Column: "news_id"},
			{Name: "Title", Type: "string", Column: "title"},
			{Name: "Content", Type: "string", Column: "content"},
		},
		PKFieldIndex: 0,
	},
	z: new(ReformNews).Values(),
}

// String returns a string representation of this struct or record.
func (s ReformNews) String() string {
	res := make([]string, 3)
	res[0] = "Id: " + reform.Inspect(s.Id, true)
	res[1] = "Title: " + reform.Inspect(s.Title, true)
	res[2] = "Content: " + reform.Inspect(s.Content, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *ReformNews) Values() []interface{} {
	return []interface{}{
		s.Id,
		s.Title,
		s.Content,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *ReformNews) Pointers() []interface{} {
	return []interface{}{
		&s.Id,
		&s.Title,
		&s.Content,
	}
}

// View returns View object for that struct.
func (s *ReformNews) View() reform.View {
	return ReformNewsTable
}

// Table returns Table object for that record.
func (s *ReformNews) Table() reform.Table {
	return ReformNewsTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *ReformNews) PKValue() interface{} {
	return s.Id
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *ReformNews) PKPointer() interface{} {
	return &s.Id
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *ReformNews) HasPK() bool {
	return s.Id != ReformNewsTable.z[ReformNewsTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.Id = pk.
func (s *ReformNews) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = ReformNewsTable
	_ reform.Struct = (*ReformNews)(nil)
	_ reform.Table  = ReformNewsTable
	_ reform.Record = (*ReformNews)(nil)
	_ fmt.Stringer  = (*ReformNews)(nil)
)

func init() {
	parse.AssertUpToDate(&ReformNewsTable.s, new(ReformNews))
}