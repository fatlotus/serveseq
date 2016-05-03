// ServeSeq is like http.ServeMux, except that it runs one http.Handler after
// another. It is designed to allow you to do two things:
//
//  1. Authentication (pre-filtering) middleware.
//  2. Overlaying static handlers with dynamic handlers.
//
// In this way, it encapsulates most of the requirements of a web framework
// without adding a hugely complicated API.
package serveseq
