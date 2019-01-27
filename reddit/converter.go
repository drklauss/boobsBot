package reddit

import (
	"regexp"

	"github.com/leesper/holmes"
)

var reMedia *regexp.Regexp

// Converter converts SubredditResponse into slice of Elements
type Converter struct {
	items []Element
}

// Element is a ready for sending or storing in db structure
type Element struct {
	URL     string
	Caption string
}

// NewConverter returns converter instanse
func NewConverter() (*Converter, error) {
	reMedia = regexp.MustCompile(`\.(png|jpg|jpeg|gif|mp4)`)
	c := new(Converter)

	return c, nil
}

// Run starts convert process
func (c *Converter) Run(sr *SubRedditResponse) []*Element {
	results := make(chan *Element, 25)
	data := make(chan Data, 25)
	for i := 1; i <= 3; i++ {
		go c.worker(i, data, results)
	}

	go func() {
		for _, oneEl := range sr.Data.Children {
			data <- oneEl.Data
		}
		close(data)
	}()
	var elements []*Element
	for k := 1; k <= len(sr.Data.Children); k++ {
		el := <-results
		if el != nil {
			holmes.Debugf("got converted element: %+v", el)
			elements = append(elements, el)
		}
	}
	close(results)

	return elements
}

// worker processed
func (c *Converter) worker(i int, data chan Data, results chan *Element) {
	for oneEl := range data {
		holmes.Debugf("worker %d processing %v", i, oneEl)
		switch oneEl.Domain {
		case "giant.gfycat.com":
			fallthrough
		case "gfycat.com":
			el, err := c.processingGfycat(oneEl)
			if err != nil {
				holmes.Warnf("could not process element: %v", err)
			}
			results <- el
		case "i.imgur.com":
			fallthrough
		case "imgur.com":
			el, err := c.processingImgur(oneEl)
			if err != nil {
				holmes.Warnf("could not process element: %v", err)
			}

			results <- el
		case "i.redd.it":
			results <- &Element{URL: oneEl.URL, Caption: oneEl.Title}
		default:
			holmes.Warnf("unknown domain name: %s", oneEl.Domain)
			results <- nil
		}
	}
}
