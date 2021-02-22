package bot

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Middleware func allows to create some middlewares to use them in handlers.
type Middleware func(req HandlerRequest, next HandlerFunc) HandlerFunc

// LogRequest is a dev middleware, that just logs the request data.
func LogRequest(req HandlerRequest, next HandlerFunc) HandlerFunc {
	return func(req HandlerRequest) {
		t := time.Now()
		log.Infof("process %d message: %s", req.Update.UpdateID, req.Update.Message.Text)
		next(req)
		log.Infof("update %d handled in: %v", req.Update.UpdateID, time.Since(t))
	}
}
