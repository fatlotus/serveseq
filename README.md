# ServeSeq: http.ServeMux, but for sequencing

Rubyists and Pythonistas love to compose web applications out of layers of
intermediate software, or middleware, but left out in the cold when switching
to Go, whose thread model makes thread local variables somewhat ugly.

I propose a different model. Essentially, it replaces

```go
func Allowed(w http.ResponseWriter, r *http.Request, role string) bool {
    if r.FormValue("role") != role {
        http.Error(w, r, "forbidden", 403)
        return false
    }
    return true
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    if !auth.Allowed(w, r, "administrator") {
        return
    }
    
    fmt.Fprintf(w, "hello world")
}

func main() {
    http.ListenAndServe(":8080", http.HandleFunc(MyHandler))
}
```

with

```go
func Allowed(string role) {
    return http.HandleFunc(func (w http.ResponseWriter, r *http.Request) {
        if r.FormValue("role") != role {
            http.Error(w, r, "forbidden", 403)
        }
    })
}

func MyHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello world")
}

func main() {
    seq := serveseq.New()
    seq.Append(Allowed("administrator"))
    seq.AppendFunc(MyHandler)

    http.ListenAndServe(":8080", seq)
}
```
