package grawler

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Page struct {
	url string
	body []byte
}

func Save(p *Page) {
	u, err := url.ParseRequestURI(p.url)
	if err != nil {
		return
	}

	str := REPOSITORY_PATH + strings.Replace(u.Host, ".", "/", -1) + u.Path
	err = os.MkdirAll(str, 0755)
	if err != nil {
		logger.Panic(err)
	}

	var name string
	if u.Path == "/" {
		u.Path = ""
		name = "default.html"
	} else {
		tmp := strings.Split(u.Path, "/")
		name = tmp[len(tmp) - 1]
	}

	name += fmt.Sprintf("-%s", time.Now().Format("2-Jan-2006-15:04:05.000"))
	f, err := os.Create(str + "/" + name)
	if err != nil {
		logger.Panic(err)
	}

	logger.Printf("Created dir: %s/%s\n", str, name)

	f.Write(p.body)
	f.Close()
}
