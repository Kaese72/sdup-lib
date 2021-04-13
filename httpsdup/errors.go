package httpsdup

import (
	"net/http"

	"github.com/Kaese72/sdup-lib/sduptemplates"
)

//HTTPStatusCode crudely translates error into http status code
func HTTPStatusCode(err error) int {
	switch err {
	case sduptemplates.NoSuchAttribute, sduptemplates.NoSuchCapability, sduptemplates.NoSuchDevice:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
