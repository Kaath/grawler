# Grawler

## Purpose

Not much really.

## Usage

### As a library

```
type Page struct {
	url string
	body []byte
}
```

The `Page` object represents the content of
a web page read by the crawler.

```
	// ...
	type SaveFunc func(p *Page)
	// ...

	// ...
	StartCrawl(urls []string, treatments ...SaveFunc)
	// ...
```

`SaveFunc` is a type used to define any function applying a
treatment on a `Page` object. This is the default interface
to change the behaviour of the crawler.

`StartCrawl` will begin the crawl on all `urls` simultaneously
(via goroutines) and run all`treatments` sequentialy on all `Pages`.

