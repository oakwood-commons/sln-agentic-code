# `shorty` -- Sample Go URL Shortener

Minimal Go standard-library HTTP service used in the guided walkthrough.

## Run Locally

~~~bash
cd examples/go-app
go run ./cmd/shorty -addr :8080
~~~

## Endpoints

### `POST /links` -- Create a Short Link

Request:

~~~bash
curl -s -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/oakwood-commons/scafctl"}'
~~~

Response (`201 Created`):

~~~json
{
  "code": "a1b2c3",
  "url": "https://github.com/oakwood-commons/scafctl"
}
~~~

### `GET /{code}` -- Redirect to Original URL

~~~bash
curl -i http://localhost:8080/a1b2c3
~~~

Returns `302 Found` with a `Location` header pointing to the original `http` or
`https` URL, or `404 Not Found` if the code does not exist.
