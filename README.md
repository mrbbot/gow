# Go Watcher

Restarts your app when source files change.

## Docker Example

If we have a web app in a file called `main.go`,

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = r.Run()
}
```

...and the `Dockerfile` shown below running with a volume mounted to `/app`,

```Dockerfile
FROM golang:1.11-alpine
RUN apk --update add git curl bash build-base
WORKDIR /app
RUN go get github.com/mrbbot/gow
CMD ["gow", "main.go"]
EXPOSE 8080
```

...anytime we change a file, the server will restart.