# winss
[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://pkg.go.dev/github.com/songlim327/winss) [![License MIT](https://img.shields.io/badge/license-MIT-blue)](https://img.shields.io/badge/license-MIT-blue) [![Go Report Card](https://goreportcard.com/badge/github.com/songlim327/winss)](https://goreportcard.com/report/github.com/songlim327/winss) ![Visitors](https://api.visitorbadge.io/api/visitors?path=github.com%2Fsonglim327%2Fwinss&label=visitors&countColor=%23dce775&style=flat)

* Go library for capturing screenshot on windows

## Install
```
go get github.com/songlim327/winss
```

## Example
```go
    package main

    import (
        "log"

        "github.com/songlim327/winss"
    )

    func main() {
        handler, err := winss.FindWindow("winss example.txt - Notepad")
        if err != nil {
            log.Fatal(err)
        }

        // Pass the handler as argument
        img := winss.Screenshot(handler)

        // Screenshot() will return *image.RGBA, pass it together with filename as arguments to SaveImage()
        err := winss.SaveImage(img, "./winss.png")
        return err
    }
```
## Sample Result 

![screenshot](https://ik.imagekit.io/songlim/winss.png?updatedAt=1697769749995)

## License
MIT License