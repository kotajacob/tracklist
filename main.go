package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type track struct {
	start string
	end   string
	name  string
}

func parseTracks(r io.Reader) ([]track, error) {
	var tracks []track
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		columns := strings.SplitN(text, " ", 2)
		if len(columns) < 2 {
			return nil, fmt.Errorf(
				"error not enough fields in tracklist: %v",
				text,
			)
		}

		var t track
		t.start = columns[0]
		t.name = strings.Trim(columns[1], "- \t")

		// Set last track's ending to current track's start.
		trackCount := len(tracks)
		if trackCount > 0 {
			tracks[trackCount-1].end = t.start
		}

		tracks = append(tracks, t)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tracks, nil
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("usage: tracklist list.txt album.opus")
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tracks, err := parseTracks(file)
	if err != nil {
		log.Fatalln(err)
	}

	albumName := flag.Arg(1)
	ext := filepath.Ext(albumName)
	dir := strings.TrimSuffix(albumName, ext)
	err = os.Mkdir(dir, 0777)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatalf("failed creating output directory: %v\n", err)
	}

	var cmds []*exec.Cmd
	for i, track := range tracks {
		args := []string{
			"-i", albumName,
			"-acodec", "copy",
			"-ss", track.start,
		}

		if track.end != "" {
			args = append(args, "-to", track.end)
		}
		args = append(args, filepath.Join(
			dir,
			strconv.Itoa(i)+" - "+track.name,
		)+ext)
		cmd := exec.Command(
			"ffmpeg",
			args...,
		)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		fmt.Printf("%#v\n", cmd)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
