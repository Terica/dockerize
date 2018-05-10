package progress

import (
	"fmt"
	"io"
	"os"
	"time"
)

var spinner = []rune{
	'|',
	'/',
	'-',
	'\\',
}

type progress interface {
	New() Progress
	Add(string, Progress)
	SetPosition(int)
	SetTotal(int)
	Update(string, Progress)
	Output()
}

// Progress is the structure which stores the state of a progress.
type Progress struct {
	Destination io.Writer
	LastDisplay time.Time
	Summary     string
	Position    int64
	Total       int64
	Parent      *Progress
	Part        map[string]*Progress
}

// New creates a new Progress structure
func New(summary string, position int64, total int64) *Progress {
	return &Progress{Destination: os.Stderr, Summary: summary, Position: position, Total: total, Part: map[string]*Progress{}}
}

// OnlyAdd only creates a part if it doesn't already exist, otherwise it's a noop
// Used for placeholder statuses which you don't know the value of
func (p *Progress) OnlyAdd(index string, position int64, total int64) *Progress {
	if _, ok := p.Part[index]; ok {
		return nil
	}
	return p.Add(index, position, total)
}

// Add adds a section of progress - this makes the parent's value dynamic rather than static
// because it calls Update()
// It lets you track the overall percentage of other identifiable parts
func (p *Progress) Add(index string, position int64, total int64) *Progress {
	var (
		val *Progress
		ok  bool
	)
	if val, ok = p.Part[index]; ok {
		val.Position = position
		val.Total = total
	} else {
		val = New(index, position, total)
		p.Part[index] = val
		val.Parent = p
	}
	p.Update()
	return val
}

// SetPosition sets the position on a progress, but this will be overwritten by any Update() on the same object.
func (p *Progress) SetPosition(position int64) {
	p.Position = position
	if p.Parent != nil {
		p.Parent.Update()
	}
}

// SetTotal sets the total on a progress, but this will be overwritten by any Update() on the same object.
func (p *Progress) SetTotal(total int64) {
	p.Total = total
	if p.Parent != nil {
		p.Parent.Update()
	}
}

// Update populates parent objects values from its children
func (p *Progress) Update() {
	p.Position = 0
	p.Total = 0
	for _, v := range p.Part {
		p.Position = p.Position + v.Position
		p.Total = p.Total + v.Total
	}
	if p.Parent != nil {
		p.Parent.Update()
	}
}

// Output generates the output of the progress's current state.
// It flicks the spinner on one position.
func (p *Progress) Output() string {
	// Use a Spinner if we'd be division by zero
	if p.Total == 0 {
		p.Position = (p.Position + 1) % int64(len(spinner))
		return fmt.Sprintf("\r%s:  %c ", p.Summary, spinner[p.Position])
	}
	return fmt.Sprintf("\r%s: %3d%%", p.Summary, p.Position/(p.Total/100))
}

// Display displays the current progress
func (p *Progress) Display() {
	t := time.Now()
	// Humans can deal with waiting a second for an update.
	if t.Sub(p.LastDisplay) > time.Second {
		// Because relatively this is quite expensive.
		fmt.Fprintf(p.Destination, p.Output())
		p.LastDisplay = t
	}
}

// Done cleans up the display by clearing to the end of the line.
func (p *Progress) Done() {
	fmt.Fprintf(p.Destination, "\r%c[0K", 27)
}
