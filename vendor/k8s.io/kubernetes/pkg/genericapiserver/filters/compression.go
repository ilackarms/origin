/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filters

import (
	"net/http"

	"github.com/golang/glog"
	"compress/gzip"
)

type compressionResponseWriter struct{
	wrapped    http.ResponseWriter
	compressor *gzip.Writer
}

func NewCompressionResponseWriter(w http.ResponseWriter) *compressionResponseWriter {
	return &compressionResponseWriter{
		wrapped: w,
		compressor: gzip.NewWriter(w),
	}
}

// compressionResponseWriter implements http.Responsewriter Interface
var _ http.ResponseWriter = &compressionResponseWriter{}

func (w *compressionResponseWriter) Header() http.Header{
	return w.wrapped.Header()
}

// compress data according to compression method
func (w *compressionResponseWriter) Write(p []byte) (int, error) {
	//w.Header().Set("Content-Encoding", "gzip")
	//w.Header().Set("X-Content-Encoding", "gzip")
	//w.Header().Set("X-Fuck-My-Wife", "false")
	//glog.V(0).Infof("WROTE %v bytes to \nWRITER: %+v\n with headers %+v", len(p), w, w.Header())
	defer w.compressor.Close()
	return w.compressor.Write(p)
	//return w.wrapped.Write(p)
}

func (w *compressionResponseWriter) WriteHeader(status int) {
	w.wrapped.WriteHeader(status)
}

// WithCompression wraps an http Handler to check for
// the Accept-Encoding header. If it is set to a
// known compression encoding, WithCompression
// will attempt to compress the response body
// returned by the handler
func WithCompression(handler http.Handler) http.Handler {
	glog.V(0).Infof("WRAPPING HANDLER: %+v", handler)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//if strings.Contains(req.URL.Path, "watch") {
			req.Header.Del("Accept-Encoding")
		//}
		//glog.V(0).Infof("GOT REQUEST: %+v req, COMPRESSING!", req)
		//var compressor io.Writer
		//switch req.Header.Get("Accept-Encoding") {
		//case "flate":
		//	compressor, _ = flate.NewWriter(w, flate.DefaultCompression)
		//case "gzip":
		//	compressor =
		//case "zlib":
		//	compressor = zlib.NewWriter(w)
		//default:
		//	handler.ServeHTTP(w, req)
		//	return
		//}
		//compressionWriter := &compressionResponseWriter{
		//	wrapped:    w,
		//	compressor: gzip.NewWriter(w),
		//}
		//compressionWriter.Header().Set("Content-Encoding", req.Header.Get("Accept-Encoding"))
		handler.ServeHTTP(w, req)
	})
}