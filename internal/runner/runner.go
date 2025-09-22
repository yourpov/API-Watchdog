package runner

import (
	"api-watchdog/internal/alerts"
	"context"
	"net/http"
	"time"

	"github.com/yourpov/logrite"
)

type Options struct {
	Webhook     string
	RequestTOms int
}

type Check struct {
	Name      string
	URL       string
	IntervalS int
}

type Runner struct {
	opts   Options
	client *http.Client
	state  map[string]bool
	checks []Check
}

// New creates and returns a Runner
func New(o Options) *Runner {
	to := time.Duration(o.RequestTOms) * time.Millisecond
	if to <= 0 {
		to = 8 * time.Second
	}
	return &Runner{
		opts:   o,
		client: &http.Client{Timeout: to},
		state:  make(map[string]bool),
	}
}

// AddCheck adds a new check to the runner
func (r *Runner) AddCheck(c Check) {
	if c.IntervalS <= 0 {
		logrite.Error("[Config/Checks] <check_interval_seconds> is not set")
	}
	r.checks = append(r.checks, c)
}

// Run starts the checks
func (r *Runner) Run() {
	for _, c := range r.checks {
		cc := c
		go r.loop(cc)
	}
}

// loop runs checks in a for loop
func (r *Runner) loop(c Check) {
	t := time.NewTicker(time.Duration(c.IntervalS) * time.Second)
	defer t.Stop()
	for {
		healthy, why := r.ping(c)
		r.handleTransition(c, healthy, why)
		<-t.C
	}
}

// ping returns status and reason
func (r *Runner) ping(c Check) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), r.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.URL, nil)
	if err != nil {
		return false, "request build failed"
	}
	start := time.Now()
	resp, err := r.client.Do(req)
	lat := time.Since(start).Milliseconds()
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "status " + resp.Status
	}
	logrite.Info("[%s] ok (%dms)", c.Name, lat)
	return true, "ok"
}

// handleTransition handles state changes
func (r *Runner) handleTransition(c Check, healthy bool, why string) {
	prev, seen := r.state[c.Name]
	r.state[c.Name] = healthy

	if !seen && !healthy {
		r.notify(c, "DOWN", why)
		return
	}

	if seen && prev != healthy {
		if healthy {
			r.notify(c, "RECOVERED", "ok")
		} else {
			r.notify(c, "DOWN", why)
		}
	}
}

// notify sends an embed on state change
func (r *Runner) notify(c Check, kind, desc string) {
	logrite.Info("[%s] %s: %s", c.Name, kind, desc)
	if r.opts.Webhook == "" {
		return
	}
	_ = alerts.SendDiscord(r.opts.Webhook, c.Name, kind, desc)
}
