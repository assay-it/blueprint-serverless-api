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

var sut = assay.Host("")

func create(book *Book) assay.Arrow {
	return http.Join(
		ø.POST("%s/books", sut),
		ø.ContentJSON(),
		ø.Send(book),
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

func TestX() assay.Arrow {
	book := Book{Title: "The Lord of the Rings"}

	isLordOfTheRings := assay.Join(
		c.Defined(&book.ID),
		c.Value(&book.Title).String("The Lord of the Rings"),
	)

	return assay.Join(
		create(&book),
		isLordOfTheRings,

		remove(&book.ID),
	)
}

/*
func (sut *SUT) lookup(books Books) assay.Arrow {
	return http.Join(
		ø.GET("%s/books", sut.URL),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&books),
	)
}





func (sut *SUT) update(book *Book) assay.Arrow {
	return http.Join(
		ø.PUT("%s/books/%s", sut.URL, book.ID),
		ø.ContentJSON(),
		ø.Send(book),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Recv(&book),
	)
}
*/
