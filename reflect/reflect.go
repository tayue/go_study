package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"time"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}
type ITableName interface {
	TableName() string
}

func (u User) TableName() string {
	return "Users"
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// Get Data Type for sqlite3 Dialect
func DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// Parse a struct to a Schema instance
func Parse(dest interface{}) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	fmt.Printf("modelType:%+v\n", modelType)
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	fmt.Println("tableName:", tableName)
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		fmt.Printf("Field(%d):%+v\n", i, p)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			fmt.Printf("\nfield:%+v\n", field)
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func main() {
	schema := Parse(User{"tayue", 30})
	fmt.Printf("schema:%+v\n", schema)
}