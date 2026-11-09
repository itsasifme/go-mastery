# net/http

> 📝 **Planned.** No code yet. This is what I plan to add here.

Small HTTP servers and clients. Big REST API projects don't belong in this repo.

## Examples to add

- `01-server`: routes like `GET /items/{id}` with `http.ServeMux`
- `02-client`: a client with a timeout
- `03-httptest`: test a handler without a real network

## Things to remember

- **The default HTTP client has no timeout.** Always set one.
- **Always close `resp.Body`.**

## Go version notes

- Routes with a method and `{id}` came in Go 1.22.
