package middleware

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
	"application/x-gzip",
	"application/gzip",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type Middleware func(http.Handler) http.Handler

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func UnzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			req, err := http.NewRequest(r.Method, r.RequestURI, reader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, req)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func ZipResponse(next http.Handler) http.Handler {
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

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, r)
	})
}

func GzipMiddleware(next http.Handler) http.Handler {
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
			req, err := http.NewRequest(r.Method, r.RequestURI, reader)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, req)
		} else {
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, r)
		}
	})
}
