// ServeSeq is like http.ServeMux, except that it runs one http.Handler after
// another. It is designed to allow you to do two things:
//
//  1. Authentication (pre-filtering) middleware.
//  2. Overlaying static handlers with dynamic handlers.
//
// it does most of what a web framework would need to, without adding
// a hugely complicated API, or new types.
package serveseq
