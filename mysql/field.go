package mysql

import (
	//"bytes"
	//"fmt"
	//"math"
	//"os"
	//"reflect"
	//"strconv"
	//"time"
)

// DataType is ...
type DataType string

const (
    // Bool is ...
	Bool   DataType = "bool"
    // Int is ...
	Int    DataType = "int"
    // Uint is ...
	Uint   DataType = "uint"
    // Float is ...
	Float  DataType = "float"
    // String is ...
	String DataType = "string"
    // Time is ...
	Time   DataType = "time"
    // Bytes is ...
	Bytes  DataType = "bytes"
)

// Field is ...
type Field struct {
	Name            string
	DBName          string
	Table           string
	DataType        DataType
	PrimaryKey      bool
	AutoIncrement   bool
	DefaultValue    string
	NotNull         bool
	Unique          bool
	Comment         string
	Size            int
	Scale           int
}


// Field is ...
//type Field struct {
	//Catalog  string
	//Db       string
	//Table    string
	//OrgTable string
	//Name     string
	//OrgName  string
	//DispLen  uint32
	////  Charset  uint16
	//Flags uint16
	//Type  byte
	//Scale byte
//}

