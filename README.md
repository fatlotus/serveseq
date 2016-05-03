# ServeSeq: http.ServeMux, but for sequencing

```go
import "github.com/fatlotus/serveseq"
```

[![Circle CI](https://circleci.com/gh/fatlotus/serveseq.svg?style=svg)](https://circleci.com/gh/fatlotus/serveseq)

ServeSeq is like http.ServeMux, except that it runs one http.Handler after
another. It is designed to allow you to do two things:

1. Authentication (pre-filtering) middleware.
2. Overlaying static handlers with dynamic handlers.

In this way, it does most of what a web framework would need to, without adding
a hugely complicated API, or new types.

## Example

Suppose we want to build an app with authentication and a custom error page. We
start by defining a basic login handler. ServeSeq skips the remaining handlers
if we write a response, so if we are logged in, the handler simply falls
through.

```go
// If the user is not logged in, redirect to the login page.
func RequireLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		fmt.Fprintf(w, "[login form]")

	} else if r.FormValue("session") == "" {
		http.Redirect(w, r, "/login", 307)
	}
}
```

Next, we define what happens _after_ the application. In this case, we override
the default 404 handler and provide our own message.

```go
// Capture any unhandled requests.
func CustomNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "not found :(")
}
```

We can glue these two functions before and after a standard `http.ServeMux`. 

```go
func main() {
	// Create a conventional application.
	mux := http.NewServeMux()
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world!")
	})

	// Compose it with some middleware.
	seq := serveseq.New()
	seq.NextFunc(RequireLogin)
	seq.Next(mux)
	seq.NextFunc(CustomNotFound)

  http.ListenAndServe(":9000", seq)
}
```

Voilà! Simple, type-safe composition.

## License

ServeSeq is licensed under the MIT License:

> Copyright (c) 2016 Jeremy Archer
> 
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
> 
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
> 
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
> BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
> ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
