using Go = import "/go.capnp";
@0xef5a6e7f51007b47;
$Go.package("books");
$Go.import("foo/books");

struct Book {
	title @0 :Text;
	# Title of the book.

	pageCount @1 :Int32;
	# Number of pages in the book.
}