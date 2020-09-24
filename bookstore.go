//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/blueprint-serverless-api
//

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fogfish/dynamo"
	µ "github.com/fogfish/gouldian"
	"github.com/fogfish/gouldian/header"
	"github.com/fogfish/gouldian/path"
	"github.com/fogfish/guid"
	"github.com/fogfish/iri"
)

// Book is a struct used by api and storage
type Book struct {
	iri.ID
	Title string `dynamodbav:"title,omitempty" json:"title,omitempty"`
}

var title = dynamo.Thing(Book{}).Field("Title")

// Books is a sequence of books
type Books []Book

// Value and other sequence functions supports transformation of
// seneric DynamoDB sequence into typed Books
func (seq Books) Value(i int) interface{} { return seq[i] }
func (seq Books) Len() int                { return len(seq) }
func (seq Books) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq Books) Less(i, j int) bool      { return seq[i].ID.IRI.String() < seq[j].ID.IRI.String() }
func (seq Books) String(i int) string     { return seq[i].ID.IRI.String() }

// Join is a monoid to append generic element into sequence
func (seq *Books) Join(gen dynamo.Gen) (iri.Thing, error) {
	val := Book{}
	if fail := gen.To(&val); fail != nil {
		return nil, fail
	}
	*seq = append(*seq, val)
	return &val, nil
}

// CRUD REST API implementation
type CRUD struct {
	db dynamo.KeyVal
}

// lookupBooks is the endpoint to fetch all books from DynamoDB
func (api *CRUD) lookupBooks() µ.Endpoint {
	return µ.GET(
		µ.Path(path.Is("books")),
		µ.FMap(func() error {
			seq := Books{}
			if _, err := api.db.Match(iri.New("books")).FMap(seq.Join); err != nil {
				return µ.InternalServerError(err)
			}
			return µ.Ok().JSON(seq)
		}),
	)
}

// createBook stores a new book to DynamoDB
func (api *CRUD) createBook() µ.Endpoint {
	var book Book

	return µ.POST(
		µ.Path(path.Is("books")),
		µ.Header(header.ContentJSON()),
		µ.Body(&book),
		µ.FMap(func() error {
			if book.ID.IRI.String() == "" {
				book.ID = iri.New("books:%s", guid.Seq.ID())
			}

			if err := api.db.Put(&book); err != nil {
				return µ.InternalServerError(err)
			}
			return µ.Ok().JSON(book)
		}),
	)
}

// lookupBook by unique id from DynamoDB
func (api *CRUD) lookupBook() µ.Endpoint {
	var (
		id string
	)

	return µ.GET(
		µ.Path(path.Is("books"), path.String(&id)),
		µ.FMap(func() error {
			book := Book{ID: iri.New(id)}
			if err := api.db.Get(&book); err != nil {
				return µ.InternalServerError(err)
			}
			return µ.Ok().JSON(book)
		}),
	)
}

// updateBook changes existing book at DynamoDB
func (api *CRUD) updateBook() µ.Endpoint {
	var (
		id   string
		book Book
	)

	return µ.PUT(
		µ.Path(path.Is("books"), path.String(&id)),
		µ.Header(header.ContentJSON()),
		µ.Body(&book),
		µ.FMap(func() error {
			book.ID = iri.New(id)
			if err := api.db.Update(&book, title.Exists()); err != nil {
				return µ.InternalServerError(err)
			}
			return µ.Ok().JSON(book)
		}),
	)
}

// removeBook from DynamoDB
func (api *CRUD) removeBook() µ.Endpoint {
	var (
		id string
	)

	return µ.DELETE(
		µ.Path(path.Is("books"), path.String(&id)),
		µ.FMap(func() error {
			book := Book{ID: iri.New(id)}
			if err := api.db.Remove(&book); err != nil {
				return µ.InternalServerError(err)
			}
			return µ.Ok().JSON(book)
		}),
	)
}

// spawn lambda function and init its api
func main() {
	api := CRUD{
		db: dynamo.Must(dynamo.New(os.Getenv("CONFIG_DDB"))),
	}

	lambda.Start(
		µ.Serve(
			api.removeBook(),
			api.updateBook(),
			api.createBook(),
			api.lookupBook(),
			api.lookupBooks(),
		),
	)
}
