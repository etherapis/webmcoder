// WebMCoder - Media stream extension encoder
// Copyright (c) 2016 EtherAPIs Authors. All rights reserved.
//
// Released under the BSD license.

// Wrapper around the WebMCoder docker container.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// WebM encoding docker containers
var dockerDist = "etherapis/webmcoder"

// Command line arguments to fine tune the encoding
var (
	audioChannels   = flag.Int("achan", 0, "Number of audio channels to generate (0 = same as input)")
	audioBitrate    = flag.Int("arate", 0, "Audio bitrate to encode the output to (0 = same as input)")
	videoResolution = flag.String("vres", "", "Video resolution (WxH) to encode the output to (empty = same as input)")
	videoBitrate    = flag.Int("vrate", 0, "Video bitrate to encode the output to (0 = same as input)")

	dockerImage = flag.String("image", "", "Use custom docker image instead of official distribution")
)

// EncodeFlags is a simple collection of flags to fine tune the WebM encoding.
type EncodeFlags struct {
	AudioChannels   int    // Number of audio channels to generate
	AudioBitrate    int    // Audio bitrate to encode the output to
	VideoResolution string // Video resolution (WxH) to encode the output to
	VideoBitrate    int    // Video bitrate to encode the output to
}

func main() {
	// Retrieve the CLI flags and ensure everything's set
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Printf("Usage: %s [<options>] <input> <output>.webm\n", os.Args[0])
		os.Exit(-1)
	}
	if !strings.HasSuffix(flag.Args()[1], ".webm") {
		fmt.Printf("Output file must be WebM format (.webm)\n")
		os.Exit(-1)
	}
	// Ensure docker is available
	if err := checkDocker(); err != nil {
		log.Fatalf("Failed to check docker installation: %v.", err)
	}
	// Select the image to use, either official or custom
	image := dockerDist
	if *dockerImage != "" {
		image = *dockerImage
	}
	// Check that all required images are available
	found, err := checkDockerImage(image)
	switch {
	case err != nil:
		log.Fatalf("Failed to check docker image availability: %v.", err)
	case !found:
		fmt.Println("not found!")
		if err := pullDockerImage(image); err != nil {
			log.Fatalf("Failed to pull docker image from the registry: %v.", err)
		}
	default:
		fmt.Println("found.")
	}
	// Assemble the encoding docker command and run the encoder
	config := &EncodeFlags{
		AudioChannels:   *audioChannels,
		AudioBitrate:    *audioBitrate,
		VideoResolution: *videoResolution,
		VideoBitrate:    *videoBitrate,
	}
	if err := encode(image, config, flag.Args()[0], flag.Args()[1]); err != nil {
		log.Fatalf("Failed to cross compile package: %v.", err)
	}
}

// checkDocker checks whether a docker installation can be found and is functional.
func checkDocker() error {
	fmt.Println("Checking docker installation...")
	if err := run(exec.Command("docker", "version")); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

// checkDockerImage checks whether a required docker image is available locally.
func checkDockerImage(image string) (bool, error) {
	fmt.Printf("Checking for required docker image %s... ", image)
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte(image)), nil
}

// pullDockerImage pulls an image from the docker registry.
func pullDockerImage(image string) error {
	fmt.Printf("Pulling %s from docker registry...\n", image)
	return run(exec.Command("docker", "pull", image))
}

// encode processes the input file and generates an output WebM encoded video
// file supporting media source extension streaming.
func encode(image string, config *EncodeFlags, input, output string) error {
	// Assemble and run the cross compilation command
	fmt.Printf("WebMCoding %s -> %s...\n", input, output)

	inputPath, err := filepath.Abs(input)
	if err != nil {
		return err
	}
	outputPath, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	args := []string{
		"run", "--rm",
		"-v", filepath.Dir(inputPath) + ":/input:ro",
		"-v", filepath.Dir(outputPath) + ":/output:rw",
		"-e", fmt.Sprintf("AUDIO_CHANNELS=%d", config.AudioChannels),
		"-e", fmt.Sprintf("AUDIO_BITRATE=%d", config.AudioBitrate),
		"-e", fmt.Sprintf("VIDEO_RESOLUTION=%s", config.VideoResolution),
		"-e", fmt.Sprintf("VIDEO_BITRATE=%d", config.VideoBitrate),
	}
	args = append(args, []string{image, filepath.Base(inputPath), filepath.Base(outputPath)}...)
	return run(exec.Command("docker", args...))
}

// run executes a command synchronously, redirecting its output to stdout.
func run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
