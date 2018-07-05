package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//	"path/filepath"
	"strings"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/contactkeval/expressioneval/evaluator"
	"github.com/contactkeval/expressioneval/tokenizer"
)

func main() {
	expressionEval("StartsWith(\"Mr. John Smith Jr.\", [\"Miss.\", \"Mrs.\", \"Sir\"])")
	//processFile("./ExpressionEngine/zfix.txt")

	/*	files, _ := filepath.Glob("./ExpressionEngine/*.txt")
		for _, fileName := range files {
			fmt.Println("\n--------------------------------------")
			fmt.Println(fileName, "\n--------------------------------------")
			processFile(fileName)
		}
	*/
}

func processFile(fileName string) {

	if file, err := os.Open(fileName); err == nil {

		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			if len(strings.TrimSpace(scanner.Text())) > 0 && scanner.Text()[0:2] != "//" {
				expressionEval(scanner.Text())
			}
		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Fatal(err)
	}
}

func expressionEval(text string) {

	if tokens, err := tokenizer.Tokenize(text); err != nil {
		fmt.Printf("Tokens: %#v , Error: %#v\n", tokens, err)
	} else {
		tokens = tokens.WithoutWhitespace()

		res, err := evaluator.Evaluate(tokens)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			fmt.Printf("result => <%s> %s\n", res.DataType(), datatype.ToPrint(res))
		}
	}

}
