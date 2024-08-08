# restclient [![Go Report Card](https://goreportcard.com/badge/github.com/raeperd/restclient)](https://goreportcard.com/report/github.com/raeperd/restclient) [![Coverage Status](https://coveralls.io/repos/github/raeperd/restclient/badge.svg?branch=main)](https://coveralls.io/github/raeperd/restclient?branch=main)  

## _REST client for Gophers._

**The problem**: Go's net/http is powerful and versatile, but using it correctly for client requests can be extremely verbose.

**The solution**: The restclient.Builder type is a convenient way to build, send, and handle HTTP requests. Builder has a fluent API with methods returning a pointer to the same struct, which allows for declaratively describing a request by method chaining.

restclient also comes with tools for building custom http transports, include a request recorder and replayer for testing.

## Features

- Simplifies HTTP client usage compared to net/http
- Can't forget to close response body
- Checks status codes by default
- Supports context.Context
- JSON serialization and deserialization helpers
- Easily manipulate URLs and query parameters
- Request recording and replaying for tests
- Customizable transports and validators that are compatible with the standard library and third party libraries
- No third party dependencies
- Good test coverage

## Examples
### Simple GET into a string

<table>
<thead>
<tr>
<th><strong>code with net/http</strong></th>
<th><strong>code with restclient</strong></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
req, err := http.NewRequestWithContext(ctx,
	http.MethodGet, "http://example.com", nil)
if err != nil {
	// ...
}
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
b, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
s := string(b)
```
</td>
<td>

```go
var s string
err := restclient.
	URL("http://example.com").
	ToString(&s).
	Fetch(ctx)
```

</td>
</tr>
<tr><td>11+ lines</td><td>5 lines</td></tr>
</tbody>
</table>


### POST a raw body

<table>
<thead>
<tr>
<th><strong>code with net/http</strong></th>
<th><strong>code with restclient</strong></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
body := bytes.NewReader(([]byte(`hello, world`))
req, err := http.NewRequestWithContext(ctx, http.MethodPost,
	"https://postman-echo.com/post", body)
if err != nil {
	// ...
}
req.Header.Set("Content-Type", "text/plain")
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
_, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
```

</td>
<td>

```go
err := restclient.
	URL("https://postman-echo.com/post").
	BodyBytes([]byte(`hello, world`)).
	ContentType("text/plain").
	Fetch(ctx)
```

</td>
</tr>
<tr><td>12+ lines</td><td>5 lines</td></tr></tbody></table>

### GET a JSON object

<table>
<thead>
<tr>
<th><strong>code with net/http</strong></th>
<th><strong>code with restclient</strong></th>
</tr>
</thead>
<tbody>
<tr>
<td>

```go
var post placeholder
u, err := url.Parse("https://jsonplaceholder.typicode.com")
if err != nil {
	// ...
}
u.Path = fmt.Sprintf("/posts/%d", 1)
req, err := http.NewRequestWithContext(ctx,
	http.MethodGet, u.String(), nil)
if err != nil {
	// ...
}
res, err := http.DefaultClient.Do(req)
if err != nil {
	// ...
}
defer res.Body.Close()
b, err := io.ReadAll(res.Body)
if err != nil {
	// ...
}
err := json.Unmarshal(b, &post)
if err != nil {
	// ...
}
```
</td><td>

```go
var post placeholder
err := restclient.
	URL("https://jsonplaceholder.typicode.com").
	Pathf("/posts/%d", 1).
	ToJSON(&post).
	Fetch(ctx)
```

</td>
</tr>
<tr><td>18+ lines</td><td>7 lines</td></tr></tbody></table>

### POST a JSON object and parse the response

```go
var res placeholder
req := placeholder{
	Title:  "foo",
	Body:   "baz",
	UserID: 1,
}
err := restclient.
	URL("/posts").
	Host("jsonplaceholder.typicode.com").
	BodyJSON(&req).
	ToJSON(&res).
	Fetch(ctx)
// net/http equivalent left as an exercise for the reader
```

### Set custom headers for a request

```go
// Set headers
var headers postman
err := restclient.
	URL("https://postman-echo.com/get").
	UserAgent("bond/james-bond").
	ContentType("secret").
	Header("martini", "shaken").
	Fetch(ctx)
```

### Easily manipulate URLs and query parameters

```go
u, err := restclient.
	URL("https://prod.example.com/get?a=1&b=2").
	Hostf("%s.example.com", "dev1").
	Param("b", "3").
	ParamInt("c", 4).
	URL()
if err != nil { /* ... */ }
fmt.Println(u.String()) // https://dev1.example.com/get?a=1&b=3&c=4
```

### Record and replay responses

```go
// record a request to the file system
var s1, s2 string
err := restclient.URL("http://example.com").
	Transport(restclient.Record(nil, "somedir")).
	ToString(&s1).
	Fetch(ctx)
check(err)

// now replay the request in tests
err = restclient.URL("http://example.com").
	Transport(restclient.Replay("somedir")).
	ToString(&s2).
	Fetch(ctx)
check(err)
assert(s1 == s2) // true
```

## FAQs

[See wiki](https://github.com/earthboundkid/requests/wiki) for more details.

### Why not just use the standard library HTTP client?

Brad Fitzpatrick, long time maintainer of the net/http package, [wrote an extensive list of problems with the standard library HTTP client](https://github.com/bradfitz/exp-httpclient/blob/master/problems.md). His four main points (ignoring issues that can't be resolved by a wrapper around the standard library) are:

> - Too easy to not call Response.Body.Close.
> - Too easy to not check return status codes
> - Context support is oddly bolted on
> - Proper usage is too many lines of boilerplate

restclient solves these issues by always closing the response body, checking status codes by default, always requiring a `context.Context`, and simplifying the boilerplate with a descriptive UI based on fluent method chaining.

### Why restclient and not some other helper library?

There are two major flaws in other libraries as I see it. One is that in other libraries support for `context.Context` tends to be bolted on if it exists at all. Two, many hide the underlying `http.Client` in such a way that it is difficult or impossible to replace or mock out. Beyond that, I believe that none have achieved the same core simplicity that the restclient library has.

### How do I just get some JSON?

```go
var data SomeDataType
err := restclient.
	URL("https://example.com/my-json").
	ToJSON(&data).
	Fetch(ctx)
```

### How do I post JSON and read the response JSON?

```go
body := MyRequestType{}
var resp MyResponseType
err := restclient.
	URL("https://example.com/my-json").
	BodyJSON(&body).
	ToJSON(&resp).
	Fetch(ctx)
```

### How do I just save a file to disk?

It depends on exactly what you need in terms of file atomicity and buffering, but this will work for most cases:

```go
err := restclient.
	URL("http://example.com").
	ToFile("myfile.txt").
	Fetch(ctx)
```

For more advanced use case, use `ToWriter`.

### How do I save a response to a string?

```go
var s string
err := restclient.
	URL("http://example.com").
	ToString(&s).
	Fetch(ctx)
```

### How do I validate the response status?

By default, if no other validators are added to a builder, restclient will check that the response is in the 2XX range. If you add another validator, you can add `builder.CheckStatus(200)` or `builder.AddValidator(restclient.DefaultValidator)` to the validation stack.

To disable all response validation, run `builder.AddValidator(nil)`.

## Contributing

Please [create a discussion](https://github.com/raeperd/restclient/discussions) before submitting a pull request for a new feature.
