# go-logtag
Colored tags before log messages

![image](https://user-images.githubusercontent.com/25147494/233410171-d3cf37b2-422f-4f8d-aace-cbd1b091120c.png)

## Usage

- Copy the logtag directory into your go project
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

Also supports log functions `Error`, `Errorf`, `Warn`, `Warnf`, `Info`, `Infof`, `Fatal`, `Fatalf` for easy migration from other logging libraries

## Gin middleware function

![image](https://user-images.githubusercontent.com/25147494/233414767-20375971-7baa-4d5c-9321-d52d63d3279c.png)

Use as gin logging middleware:
```go
engine *gin.Engine = gin.New()
engine.Use(logtag.GinLogTag(TagHttp))
```
