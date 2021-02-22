package reddit

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

var reMedia *regexp.Regexp

// Converter converts SubredditResponse into slice of Elements.
type Converter struct {
	items []Element
}

// Element is a ready for sending or storing in db structure.
type Element struct {
	URL     string
	Caption string
}

// NewConverter returns converter instance.
func NewConverter() *Converter {
	reMedia = regexp.MustCompile(`\.(png|jpg|jpeg|gif|mp4)`)

	return new(Converter)
}

// Run starts convert process.
func (c *Converter) Run(sr *SubRedditResponse) []*Element {
	workersCount := 3
	results := make(chan *Element, workersCount)
	data := make(chan Data, workersCount)
	for i := 1; i <= workersCount; i++ { // run workers
		go c.worker(i, data, results)
	}

	go func() {
		// put request data to channel
		for _, oneEl := range sr.Data.Children {
			data <- oneEl.Data
		}

	}()
	var elements []*Element
	for k := 1; k <= len(sr.Data.Children); k++ {
		el := <-results // wait all results
		if el != nil {
			log.Debugf("got converted element: %+v", el)
			elements = append(elements, el)
		}
	}
	close(data)    // close request channel
	close(results) // close result channel
	log.Debug("converting process finished")

	return elements
}

// worker processes convert.
func (c *Converter) worker(i int, data chan Data, results chan *Element) {
	for oneEl := range data {
		log.Debugf("worker %d processing %v", i, oneEl)
		switch oneEl.Domain {
		case "giant.gfycat.com":
			fallthrough
		case "gfycat.com":
			el, err := c.processingGfycat(oneEl)
			if err != nil {
				log.Warnf("gfycat.com - could not process element: %v", err)
			}
			results <- el
		case "redgifs.com":
			el, err := c.processingRedgifs(oneEl)
			if err != nil {
				log.Warnf("redgifs.com - could not process element: %v", err)
			}
			results <- el
		case "i.imgur.com":
			fallthrough
		case "imgur.com":
			el, err := c.processingImgur(oneEl)
			if err != nil {
				log.Warnf("imgur.com - could not process element: %v", err)
			}
			results <- el
		case "i.redd.it":
			results <- &Element{URL: oneEl.URL, Caption: oneEl.Title}
		default:
			log.Warnf("unknown domain name: '%s'", oneEl.Domain)
			results <- nil
		}
	}
	log.Debug("converting work finished")
}
