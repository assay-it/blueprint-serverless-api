//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/blueprint-serverless-api
//

package suite

import (
	"github.com/assay-it/sdk-go/assay"
	c "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

//
type Book struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

//
type Books []Book

// Value ...
func (seq Books) Value(i int) interface{} { return seq[i] }
func (seq Books) Len() int                { return len(seq) }
func (seq Books) Swap(i, j int)           { seq[i], seq[j] = seq[j], seq[i] }
func (seq Books) Less(i, j int) bool      { return seq[i].ID < seq[j].ID }
func (seq Books) String(i int) string     { return seq[i].ID }

//
//
var sut = assay.Host("")

func lookup(book *Book) assay.Arrow {
	return http.Join(
		ø.GET("%s/books/%s", sut, &book.ID),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&book),
	)
}

func create(book *Book) assay.Arrow {
	return http.Join(
		ø.POST("%s/books", sut),
		ø.ContentJSON(),
		ø.Send(&book),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&book),
	)
}

func remove(id *string) assay.Arrow {
	return http.Join(
		ø.DELETE("%s/books/%s", sut, id),
		ƒ.Code(http.StatusCodeOK),
	)
}

func update(book *Book) assay.Arrow {
	return http.Join(
		ø.PUT("%s/books/%s", sut, &book.ID),
		ø.ContentJSON(),
		ø.Send(&book),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&book),
	)
}

//
func Create() assay.Arrow {
	book := Book{
		ID:    "book:hobbit",
		Title: "There and Back Again",
	}

	return create(&book).
		Then(
			c.Value(&book.ID).String("book:hobbit"),
			c.Value(&book.Title).String("There and Back Again"),
		)
}

//
func Update() assay.Arrow {
	book := Book{
		ID:    "book:hobbit",
		Title: "The Hobbit",
	}

	return update(&book).
		Then(
			c.Value(&book.ID).String("book:hobbit"),
			c.Value(&book.Title).String("The Hobbit"),
		)
}

//
func Lookup() assay.Arrow {
	book := Book{
		ID: "book:hobbit",
	}

	return lookup(&book).
		Then(
			c.Value(&book.ID).String("book:hobbit"),
			c.Value(&book.Title).String("The Hobbit"),
		)
}

//
func Remove() assay.Arrow {
	id := "book:hobbit"

	return remove(&id)
}

//
func Lifecycle() assay.Arrow {
	book := Book{Title: "The Lord of the Rings"}

	return assay.Join(
		//
		create(&book).Then(
			c.Defined(&book.ID),
			c.Value(&book.Title).String("The Lord of the Rings"),
		),

		//
		c.FMap(func() error {
			book.Title = "The Lord of the Flies"
			return nil
		}),
		update(&book).Then(
			c.Value(&book.Title).String("The Lord of the Flies"),
		),

		//
		lookup(&book),
		c.Value(&book.Title).String("The Lord of the Flies"),

		//
		remove(&book.ID),
	)
}
