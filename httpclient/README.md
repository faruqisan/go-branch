# httpclient

`httpclient` package aim to provide easy to use http client with builtin circuit breaker.
It is still in development and subject to API changes, but very encouraged to be used
by internal `affiliate` code.

When API changes happen, it is the responsibility of the package maintainer to modify
the code that using this package.

ALL `affiliate` code should use `httpclient` to do HTTP request call.
If this package doesn't meet your requirement, we should improve it.

## Using httpclient.Client

Your package should use single client per service you call.
So, you should use different client for `topads` and `tome`.
This is because different target service should use different circuit breaker.

Creates client with default options.
 
```go
client := httpclient.NewClient()
```

Creates client with the given timeout
```go
client := httpclient.NewClient(
    httpclient.WithHTTPTimeout ( 10 * time.Second),
    httpclient.WithMaxConcurrentRequest(10),
    )
```

### Using Do & DoJSON method of 

Do HTTP request, given the `request` object

```go
resp, err = client.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
....
```

Do HTTP request, and parse the JSON response
```go
var data myStruct
resp, err := client.DoJSON(req, &data)
if err != nil {
    return err
}
// no need to close the response body
```

### Using GET convenient helper

GET request
```go
headers := map[string][]string{
        "K1": {"v1"},
        "K2": {"v2"},
}

// execute
resp, err := client.Get(context.Background(), url, headers)
....
defer resp.Body.Close()
```

Do GET request, and parse the JSON response
```go
headers := map[string][]string{
        "K1": {"v1"},
        "K2": {"v2"},
}
data data myStruct
// execute
resp, err := client.GetJSON(context.Background(), url, headers, &data)
....
// no need to close the response body
```


## Using global default client

It is NOT recommended to use this global default client.

Please only use this API in these conditions:
- to call endpoints outside of the tokopedia 
- you don't care about circuit breaker

Do HTTP request, given the `request` object

```go
resp, err = httpclient.Do(req)
if err != nil {
    return err
}
defer resp.Body.Close()
....
```

Do HTTP request, and parse the JSON response
```go
var data myStruct
resp, err := httpclient.DoJSON(req, &data)
if err != nil {
    return err
}
// no need to close the response body
// do something with the data
```

