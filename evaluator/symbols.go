package evaluator

import (
	"math"
	"time"

	"github.com/contactkeval/expressioneval/datatype"
)

// The symbol table with the predefined variables
var SymbolTable = map[string]interface{}{
	"DateTime.Today": datatype.DateTime(time.Now()),
	"Math.PI":        datatype.Double(math.Pi),
	"MyInt":          datatype.Int(100),
	"MyDouble":       datatype.Double(400.0),
	"MyString":       datatype.String("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"),
	"JSONString":     datatype.String(`{ "store": {"book": [{ "category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century",  "price": 8.95},{ "category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99},{ "category": "fiction","author": "Herman Melville", "title": "Moby Dick" , "isbn": "0-553-21311-3", "price": 8.99}, { "category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings" , "isbn": "0-395-19395-8" , "price": 22.99 } ] , "bicycle": { "color": "red", "price": 19.95 } } }`),
	"TestArray": datatype.List([]datatype.DataType{
		datatype.Double(1),
		datatype.Double(10),
		datatype.Double(100),
		datatype.Double(1000),
		datatype.Double(10000),
	}),
}
