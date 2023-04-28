package youtube_dl

import (
	"fmt"
	"io"
	"os/exec"
	"github.com/nnn-revo2012/livedl/procs/base"
)

var cmdList = []string{
	"./bin/yt-dlp/yt-dlp",
	"./bin/yt-dlp",
	"./yt-dlp/yt-dlp",
	"./yt-dlp",
	"yt-dlp",
}

func Open(opt... string) (cmd *exec.Cmd, stdout, stderr io.ReadCloser, err error) {
	cmd, _, stdout, stderr, err = base.Open(&cmdList, false, true, true, false, opt)
	if cmd == nil {
		err = fmt.Errorf("yt-dlp not found")
		return
	}
	return
}
