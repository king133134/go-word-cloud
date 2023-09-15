# Word Cloud Generator

This is a word cloud generator written in Go.

## Usage

### Install Dependencies

Make sure you have Go language environment installed and execute the following commands to install the required dependencies:
![example](https://github.com/king133134/go-word-cloud/blob/master/test/out.png)
```shell
go get -u github.com/king133134/go-word-cloud
```

## Example
```go
package main

import (
    "github.com/king133134/go-word-cloud"
    "golang.org/x/image/font/opentype"
    "log"
    "math/rand"
    "os"
)

func main() {
    custom()
}

func randSize() {
    cloud := wordcloud.NewWordCloud(800, 450)
    words := []string{"love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子"}

    fontBytes, _ := os.ReadFile("./fonts/SmileySans.ttf")
    f, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatalf("err: %v", err)
    }
    cloud.SetRotate(120, 5)
    cloud.SetFont(f)
    err = cloud.Render(wordcloud.NewWords(words), "test/out.png")
    if err != nil {
        log.Fatalf("err: %v", err)
    }
}

func custom() {
    cloud := wordcloud.NewWordCloud(800, 450)
    words := []string{"love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子", "love", "movie", "animation", "music", "分子", "电子", "松子", "离子", "绝绝子", "孙子", "孟子",
        "孔子", "老子", "原子", "质子", "中子", "上帝粒子"}

    w := make([]*wordcloud.Word, len(words))
    for i, word := range words {
        w[i] = wordcloud.NewWord(word, rand.Intn(45)+15, wordcloud.Color(wordcloud.NewRandColor()))
    }

    fontBytes, _ := os.ReadFile("./fonts/SmileySans.ttf")
    f, err := opentype.Parse(fontBytes)
    if err != nil {
        log.Fatalf("err: %v", err)
    }
    cloud.SetFont(f)
    cloud.SetRotate(120, 5)
    err = cloud.Render(w, "test/out.png")
    if err != nil {
        log.Fatalf("err: %v", err)
    }
}
```

## LICENSE

[MIT](https://github.com/king133134/leetCodeTests/blob/master/LICENSE)