# Go REST

Project template for Go Lang REST API microservice (code boilerplate). Repository contains the most useful and repetitive
code base to start writing REST API microservice almost in pure Go Lang. Areas covered:

- Single YAML configuration file for microservice
- HTTP web server
- Single routes entrypoint
- REST API errors handling
- Middlewares chain
- HTTP request/response logging middleware
- HTTP Header (Api-Key) authorization middleware
- Logging to file with rotation
- Pagination utilities
- Filtering & sorting utilities

Required Go Lang version: **1.22**

Additional dependencies:

```
# UUID handling
github.com/google/uuid v1.6.0

# Handling logs as files with rotate
gopkg.in/natefinch/lumberjack.v2 v2.2.1

# YAML encode/decode
gopkg.in/yaml.v3 v3.0.1
```

## Getting started

Repository contains REST API microservice template ready to use out of the box. Application can be compiled and ran with
usage of **Makefile** just type command below in terminal directly in project root directory.

```
make
```

You can test if web server is listening with usage of curl below.

```
curl --request GET \
--url GET-http://localhost:8080/hello \
--header 'Content-Type: application/json'
```

API should respond with HTTP status 200 and following response:

```json
"Hello World"
```

Project structure looks like:

- **cmd** - application main entrypoints
- **configs** - application configuration files (for example YAML)
- **logs** - application log files
- **internal/app/web** - web server definition with middlewares and routes
- **internal/app/handlers** - business and errors handlers
- **internal/app/settings** - settings file handling (default values and reading from file)
- **internal/app/types** - entities, dto, pagination, etc

## How to use?

Repository should be treated as a template, so you need to modify it's contents before usage for yourself.

Recommendations:

1. Clone repository.
2. Copy/paste sources to new directory.
3. Modify **go.mod** to match your project name.
4. Refactor imports in sources.
5. Rename **./cmd/app** directory to match your project name.
6. Change variable **app** value in **Makefile** to match your project name.
7. You are ready to Go.

### YAML configuration file

Microservice reads it's configuration from single YAML configuration file. If file not exist or file do not contain
certain settings, default configuration will be used. Configuration file should be available for microservice in one of
two locations (order has impact on priority - first configuration file will be loaded):

1. ./app.yml
2. ./configs/app.yml

```yaml
# Default settings

# Logs configuration
log:
  file-enabled: false # Indicates if logs should be saved to file
  max-size: 10 # Max file size in MB before rotate
  max-age: 30 # Max age in days before rotate

# Web server configuration
server:
  host: "0.0.0.0" # Web server host
  port: "8080" # Web server port

# Authorization configuration
authorization:
  enabled: false # Indicates if authorization via header is enabled
  header: "Api-Key" # Name of the header in which key will be provided
  key: "" # Value of correct authorization header key
```

You can add your own settings by modifying **internal/app/settings/settings.go** just write it as in example - add new
structs and fill the defaults.

### Routes definition

Microservice template has single point of routes definition in **internal/app/web/server.go**:

```go
func (s *server) routes(rtr *http.ServeMux) {
	rtr.HandleFunc("GET /hello", handlers.Hello)
}
```

Here you should add your own routes definitions or do whatever you want (need).

### Errors handling

Standard error visible by microservice client have following structure:

```json
{
    "timestamp": "2024-06-28T16:47:23Z",
    "code": "err.validation",
    "message": "Validation Failed",
    "details": [
        {
            "field": "y",
            "code": "val.division_by_zero",
            "message": "Division by zero",
            "value": "0.0",
            "expected": "y != 0"
        }
    ]
}
```

Errors handling takes place in **internal/app/handlers/errors.go**:

```go
func HandleError(err error, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	var apiErr *types.ApiError

	if errors.As(err, &apiErr) {
		slog.Error("Handling expected API error", "err", err.Error())

		switch {
		case apiErr.Status == 401:
			dto := types.NewApiErrorDto("auth.unauthorized", "Unauthorized")
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		case apiErr.Status == 400:
			dto := types.NewApiErrorDto("err.validation", "Validation Failed", apiErr.Details...)
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		case apiErr.Status == 406:
			dto := types.NewApiErrorDto("err.not-acceptable", "Not Acceptable")
			res.WriteHeader(apiErr.Status)
			res.Write(dtoToJson(dto))
			return
		}
	}

	slog.Error("Handling unexpected API error", "err", err.Error())
	errDto := types.NewApiErrorDto("err.internal", "Internal Server Error")
	res.WriteHeader(http.StatusInternalServerError)
	res.Write(dtoToJson(errDto))
}
```

As you can see, right now there are only 4 HTTP codes handled:

1. 400 - validation errors
2. 401 - unauthorized errors
3. 406 - invalid incoming payload (structure) errors
4. 500 - any other application errors

How it works? In **internal/app/types/errors.go** you will find all definitions to understand structures. Here you also
have custom error implemented. Anywhere in the application you can throw the error as follows:

```go
if nam == "" {
  return types.NewApiError(400, "Validation failed", types.NewFieldError("name", "val.not-found", "Template not found"))
}
```

First argument is the most important because it will be used in above switch statement to decide with which HTTP code
error response should be returned. When you return that error somewhere in your application, you should pass it to handler:

```go
if err := cmd.NewDeleteTemplateCmd(h.templatesStore, h.templatesParamsStore).Execute(name); err != nil {
  HandleError(err, res)
  return
}
```

In case you need more errors to be handled by your application just modify the handler.

**NOTICE**: any other errors like HTTP 404 are handled by Go Lang itself.

### Middlewares

Microservice template has implemented middlewares chain mechanism. There is single entry point for middlewares to be
declared in **internal/app/web/server.go**:

```go
func (s *server) middlewares() middleware {
	return middlewaresChain(
		logRequestAndResponseMiddleware,
		authorizationMiddleware(s.settings),
	)
}
```

**NOTICE**: middlewares order is important, first declared, first executed on incoming HTTP request.

Middlewares themself are declared in **internal/app/web/middlewares.go** also with mechanism to create middlewares chain:

```go
func middlewaresChain(mid ...middleware) middleware {
	return func(nxt http.Handler) http.HandlerFunc {
		for i := len(mid) - 1; i >= 0; i-- {
			nxt = mid[i](nxt)
		}

		return nxt.ServeHTTP
	}
}
```

Modify it as you wish.

### Logging middleware

Microservice template has implemented middleware that is responsible for logging incoming HTTP requests and responses.
Logs with request and response have unique UUID assigned, so you can track whole web server operation with single UUID
filtering. Middleware is also logging HTTP request/response body so have it in mind that when your microservice should
handle something more than "application/json" it is not the best solution to log whole request/response body. It should
be modified, you can at least achieve this as simple as that (`logRequestAndResponseMiddleware` in **internal/app/web/middlewares.go**):

```go
uid := uuid.New().String()
var body []byte
reqContentType := req.Header.Get("Content-Type")

if strings.Contains(reqContentType, "multipart/form-data") {
  body = []byte(req.FormValue("payload"))
} else {
  body, _ = io.ReadAll(req.Body)
}

req.Body = io.NopCloser(bytes.NewBuffer(body))
slog.Info("HTTP Request", "uid", uid, "method", req.Method, "path", req.RequestURI, "body", body)
responseWriter := newLogResponseWriter(res)
nxt.ServeHTTP(responseWriter, req)
resContentType := responseWriter.Header().Get("Content-Type")

if resContentType == "application/json" {
  slog.Info("HTTP Response", "uid", uid, "status", responseWriter.status, "body", responseWriter.body.String())
} else {
 slog.Info("HTTP Response", "uid", uid, "status", responseWriter.status)
}
```

### Authorization middleware

Microservice template has implemented simple authorization mechanism based on `key` that should be provided in HTTP request
header. Header name and authorization key are configurable in **settings** described above. If authorization is enabled,
lack of authorization header or invalid key will result with HTTP 401 error response.

### Logging to file

Microservice template could log to two targets:

1. Standard output (console)
2. File (with rotation)

Logging to file is configurable in **settings**. If disabled, logs will be visible only in standard output (console).
If enabled, logs will be visible in standard output and file (with rotation). Rotation is also configurable via **settings**.

Default log files location is **./logs**.

### Pagination

Microservice template has utilities to handle GET HTTP requests with pagination. Core of this mechanism is written in
**internal/app/types/page.go**. By default, page size is 25. To use pagination, first of all, read all rows quantity:

```go
quantity, err := c.sentMessagesStore.CountByFilter(filter)

if err != nil {
  return types.EmptyPageDto(), err
}

page, err := strconv.Atoi(qry.Get("page"))

if err != nil {
  page = types.DefaultPage
}

size, err := strconv.Atoi(qry.Get("size"))

if err != nil {
 size = types.DefaultPageSize
}

pagination := types.NewPagination(page, size, quantity)

if !pagination.IsValid() {
  return types.EmptyPageDto(), nil
}
```

As you can see there is a lot of operations to do. Your microservice should take arguments from URL parameters as
`page` and `size`. When something in pagination math is incorrect, you can use `EmptyPageDto()`. If everything is correct -
you have all rows count and rows themself, you can use following, to return paginated results (just implement shared
interface `PageContent` and prepare your queries for example database queries to use pagination - limit and offset):

```go
messages, err := c.sentMessagesStore.GetByFilter(filter, pagination)

if err != nil {
  return types.EmptyPageDto(), err
}

content := []types.PageContent{}

for _, msg := range messages {
  content = append(content, types.NewSentMessageDto(...))
}

return types.NewPageDto(pagination, content), nil
```

### Filtering & sorting

Microservice template has implemented mechanisms that allows you to filter and sort result on HTTP GET requests that
returns paginated results. Core of this mechanism is written in **internal/app/types/filter.go**. First of all, you
should rewrite URL parameters to filter (this is as simple as `map[string]string`, to use in for example database queries):

```go
sortable := []string{"templateName", "to", "status", "sentAt"}
filter := types.NewFilter("", "sentAt desc", qry.Get("sort"), sortable)
filter.Params["template"] = strings.TrimSpace(qry.Get("template"))
filter.Params["to"] = strings.TrimSpace(qry.Get("to"))
filter.Params["status"] = strings.TrimSpace(qry.Get("status"))
filter.Params["isRetry"] = strings.TrimSpace(qry.Get("isRetry"))
```

What is more important, you can declare sort fields as you wish. In `sortable` above there is declaration on which fields
results can be sorted (any other value is not allowed). Nextly you creates `NewFilter`. Here first argument is `prefix`.
That prefix will be used to create sort query - `t.templateName asc, t.to asc`. Empty prefix will result in `templateName asc, to asc`.
There is placeholder for single prefix "by design" (from my experience sorting on body `root` is the best solution to
handle by default in REST API's). Second argument is `defaultSort` (what to use if sort is not specified by end user).

Microservice handles sort query in format below (as HTTP request URL parameter `sort`):

```
sort=templateName:asc,to:desc
```

As you can see it is available to sort on multiple fields at once. At the end, you should use parsed by `Filter` sort
in your queries for example in database queries `order by $`.
