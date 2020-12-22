package repositories

import (
	"fmt"
	"net/http"
)

func testFirst(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This worked")
}
