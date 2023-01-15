package service

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

var allowedContentTypes = [...]string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// проверяем, имеет ли смысл сжимать контент
		// если контент не входит в перечень allowedContentTypes,
		// передаём управление дальше без изменений
		contentType := r.Header.Get("Content-Type")
		contentTypeOK := false
		for _, t := range allowedContentTypes {
			if strings.Contains(contentType, t) {
				contentTypeOK = true
				break
			}
		}
		if !contentTypeOK {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gzw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				log.Fatalf("Error when creating gzipWriter: %s", err)
				return
			}
			return
		}
		defer func(gzw *gzip.Writer) {
			err := gzw.Close()
			if err != nil {
				log.Fatalf("Error when closing gzipWriter: %s", err)
			}
		}(gzw)

		w.Header().Set("Content-Encoding", "gzip")

		var reader io.Reader

		if strings.Contains(r.Header.Get(`Content-Encoding`), "gzip") {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gzr
			defer func(gzr *gzip.Reader) {
				err := gzr.Close()
				if err != nil {
					log.Fatalf("Error when closing gzipReader: %s", err)
				}
			}(gzr)
		} else {
			reader = r.Body
		}

		req, err := http.NewRequest(r.Method, r.RequestURI, reader)

		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, req)
	})
}
