# go-logtag
Colored tags before log messages

![image](https://user-images.githubusercontent.com/25147494/233410171-d3cf37b2-422f-4f8d-aace-cbd1b091120c.png)

## Usage

- if you have not done so yet, add the following code (with correct username) to your `.gitconfig` file
```
[url "https://YOUR_GITHUB_USERNAME@github.com/"]
    insteadOf = https://github.com/
```

- Add the repo to the GOPRIVATE go environment variable `go env -w GOPRIVATE=github.com/roboticeyes/go-logtag` (Note: need to look how to append)
- Import the package in code same way as other modules
- Define your tags and colormap and call `ConfigureLogger`:
```go
// define tags
const (
	TagSystem     = "SYSTEM"
	TagConfig     = "Config"
	TagRepository = "Repository"
	TagHttp       = "HTTP"
)

//define colormap
tags := map[string]logtag.LogColor{
	TagSystem:     logtag.Magenta,
	TagConfig:     logtag.BrightRed,
	TagRepository: logtag.Green,
	TagHttp:       logtag.BrightBlue,
}
 
 logtag.ConfigureLogger(tags)
```
- Start logging with colored tags:
```go
logtag.Printf(TagSystem, "This is a system message")
logtag.Printf(TagConfig, "This is a config message")
logtag.Printf(TagRepository, "This is a repository message")
logtag.Println(TagHttp, "This is a HTTP message")
```

Also supports log functions `Info`, `Infof`, `Warn`, `Warnf`, `Error`, `Errorf`, `Fatal`, `Fatalf` for easy migration from other logging libraries

## Log levels

It is possible to define a minimum log level to limit how much is logged. E.g if you want only errors and fatals:
```go
lt.SetMinimumLogLevel(lt.LevelError)
```

## Gin middleware function

![image](https://user-images.githubusercontent.com/25147494/233414767-20375971-7baa-4d5c-9321-d52d63d3279c.png)

Use as gin logging middleware:
```go
engine *gin.Engine = gin.New()
engine.Use(logtag_gin.GinLogTag(TagHttp))
```

## GRPC interceptors

There are GRPC logging interceptors for unary and streaming calls for both client and server side.
Here's an example for a uniary server interceptor:
```go
s := grpc.NewServer(grpc.UnaryInterceptor(logtag_grpc.GrpcLogTagServerInterceptor(lt.TagGrpc)))
```
