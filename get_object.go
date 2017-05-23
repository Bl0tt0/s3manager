package s3manager

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// GetObjectHandler downloads an object to the client.
func GetObjectHandler(s3 S3) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucketName := vars["bucketName"]
		objectName := vars["objectName"]

		object, err := s3.GetObject(bucketName, objectName)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set(headerContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", objectName))
		w.Header().Set(HeaderContentType, contentTypeOctetStream)

		_, err = io.Copy(w, object)
		if err != nil {
			code := http.StatusInternalServerError
			if err.Error() == ErrBucketDoesNotExist || err.Error() == ErrKeyDoesNotExist {
				code = http.StatusNotFound
			}
			handleHTTPError(w, code, err)
			return
		}
	})
}