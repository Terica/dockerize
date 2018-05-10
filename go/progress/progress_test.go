package progress

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	pb := New("testing", 0, 0)
	if pb.Summary != "testing" {
		t.Error("Summary isn't 'testing'")
	}
	if pb.Position != 0 {
		t.Error("Position isn't 0")
	}
	if pb.Total != 0 {
		t.Error("Total isn't 0")
	}
}

func TestOutput(t *testing.T) {
	t.Run("spinner", func(*testing.T) {
		pb := New("spinner", 0, 0)
		out := pb.Output()
		if out != "\rspinner:  / " {
			t.Errorf("Result was '%s'", out)
		}
		out = pb.Output()
		if out != "\rspinner:  - " {
			t.Errorf("Result was '%s'", out)
		}
		pb.Output()
		out = pb.Output()
		if out != "\rspinner:  | " {
			t.Errorf("Result was '%s'", out)
		}
	})
	t.Run("perc", func(*testing.T) {
		pb := New("perc", 0, 100)
		out := pb.Output()
		if out != "\rperc:   0%" {
			t.Errorf("Result was '%s'", out)
		}
		pb.SetPosition(50)
		out = pb.Output()
		if out != "\rperc:  50%" {
			t.Errorf("Result was '%s'", out)
		}
	})
}

func TestAdd(t *testing.T) {
	pb := New("add", 0, 0)
	sub := pb.OnlyAdd("part1", 50, 100)
	if sub.Parent != pb {
		t.Error("pb should be sub.Parent")
	}
	result := pb.OnlyAdd("part1", 1, 1000)
	if result != nil {
		t.Error("Added a second part1 somehow")
	}
	if pb.Position != 50 || pb.Total != 100 {
		t.Error("pb did not get updated to value of part1")
	}
	pb.Add("part2", 50, 100)
	if pb.Position != 100 || pb.Total != 200 {
		t.Error("pb is not the aggregate of part1 and part2")
	}
	pb.Add("part2", 100, 100)
	if pb.Position != 150 || pb.Total != 200 {
		t.Error("pb is not the aggregate of part1 and part2")
	}
	sub.SetPosition(75)
	sub.SetTotal(75)
	if pb.Position != 175 || pb.Total != 175 {
		t.Error("updating the total of part1 did not update the parent")
	}
	sub.Add("moprogress", 100, 100)
	if pb.Position != 200 || pb.Total != 200 {
		t.Error("updating the total of moprogress did not replace part1 and update pb")
	}
}

func TestDisplay(t *testing.T) {
	bufReadWriter := &bytes.Buffer{}
	tmpBuf := make([]byte, 80)
	pb := New("display", 0, 0)
	// bufReadWriter is an io.Writer, so we can use it as a destination
	// This means we can test the actual output functions.
	pb.Destination = bufReadWriter
	if !pb.LastDisplay.IsZero() {
		t.Error("LastDisplay isn't 0")
	}
	t.Run("output", func(*testing.T) {
		pb.Display()
		if pb.Destination != bufReadWriter {
			t.Errorf("Somehow the destination got switched")
		}
		// tmpBuf is not sized to the amount read, so just compare the slice of bytes read.
		if n, _ := io.ReadFull(bufReadWriter, tmpBuf); string(tmpBuf[:n]) != "\rdisplay:  / " {
			t.Errorf("Buffer doesn't match %+v", tmpBuf[:n])
		}
		if pb.LastDisplay.IsZero() {
			t.Error("LastDisplay didn't get updated")
		}
	})
	// The next test should run quickly enough that we trigger the throttle
	old := pb.LastDisplay
	t.Run("nooutput", func(*testing.T) {
		pb.Display()
		if pb.LastDisplay != old {
			t.Error("It seems to have taken a second between running Display()")
		}
		// Nothing new should have been written to this stream.
		if n, _ := io.ReadFull(bufReadWriter, tmpBuf); n != 0 {
			t.Error("Display outputted early")
		}
	})
	t.Run("cleanup", func(*testing.T) {
		pb.Done()
		if n, _ := io.ReadFull(bufReadWriter, tmpBuf); string(tmpBuf[:n]) != "\r\x1b[0K" {
			t.Errorf("Buffer doesn't match %+v", tmpBuf[:n])
		}
	})
}
