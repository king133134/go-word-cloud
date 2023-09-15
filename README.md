# Word Cloud Generator

This is a word cloud generator written in Go.

## Usage

### Install Dependencies

Make sure you have Go language environment installed and execute the following commands to install the required dependencies:
![example](https://github.com/king133134/go-word-cloud/blob/master/test/out.png)
```shell
go get -u github.com/king133134/wordcloud
```

## Example
```go
package main

import (
    "github.com/king133134/wordcloud"
    "golang.org/x/image/font/opentype"
    "image/color"
    "log"
    "os"
)

func main() {
    // Generate word cloud with random size
    randSize()
}

func randSize() {
    cloud := wordcloud.NewWordCloud(400, 400)
    words := []string{"分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子"}

    fontBytes, _ := os.ReadFile("./fonts/SmileySans.ttf")
    f, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatalf("err: %v", err)
    }
    cloud.SetFont(f)
    err = cloud.Render(wordcloud.NewWords(words), "out.png")
    if err != nil {
        log.Fatalf("err: %v", err)
    }
}

func custom() {
    cloud := wordcloud.NewWordCloud(400, 400)
    words := []string{"分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "i", "told", "some", "custom", "-ナルト-"}

    w := make([]*wordcloud.Word, len(words))
    for i, word := range words {
        w[i] = wordcloud.NewWord(word, i*2+12, wordcloud.Color(color.RGBA{255, 0, 255, 255}))
    }

    fontBytes, _ := os.ReadFile("./fonts/SmileySans.ttf")
    f, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatalf("err: %v", err)
    }
    cloud.SetFont(f)
    err = cloud.Render(w, "out.png")
    if err != nil {
        log.Fatalf("err: %v", err)
    }
}
```

## LICENSE

[MIT](https://github.com/king133134/leetCodeTests/blob/master/LICENSE)