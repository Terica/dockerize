// +build !mocker

package docker

import (
	"local/dockerize/progress"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

var cli *Client

func TestMain(t *testing.T) {
	cli = Connect()
}

func TestPull(t *testing.T) {
	fil := filters.NewArgs()
	fil.Add("reference", "busybox:latest")
	// Do we already have a busybox:latest?
	images := cli.ListImages(types.ImageListOptions{All: true, Filters: fil})
	if len(images) > 0 {
		// let's back it up and remove it
		cli.ImageTag("busybox:latest", "busybox:backup")
		cli.ImageRemove("busybox:latest", types.ImageRemoveOptions{})
	}
	// did it remove ok?
	images = cli.ListImages(types.ImageListOptions{All: true, Filters: fil})
	if len(images) > 0 {
		t.Error("There should not be a busybox:latest image at this point")
	}
	// get the replacement image
	cli.Pull("busybox:latest", progress.New("Pulling busybox:latest", 0, 0))
	images = cli.ListImages(types.ImageListOptions{All: true, Filters: fil})
	if len(images) != 1 {
		t.Error("We didn't pull busybox:latest")
	}
	// do it again without progress
	cli.Pull("busybox:latest", nil)
	// Cleanup
	cli.ImageRemove("busybox:latest", types.ImageRemoveOptions{})
	fil = filters.NewArgs()
	fil.Add("reference", "busybox:backup")
	// Did we make a backup?
	images = cli.ListImages(types.ImageListOptions{All: true, Filters: fil})
	if len(images) > 0 {
		// Restore it
		cli.ImageTag("busybox:backup", "busybox:latest")
		cli.ImageRemove("busybox:backup", types.ImageRemoveOptions{})
	}
}
