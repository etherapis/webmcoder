#!/bin/bash
#
# Contains the WebM encoder and stream patcher.
#
# Usage: encode.sh <input> <output>
#
# Optional environment variables:
#   AUDIO_CHANNELS   - Number of audio channels to generate
#   AUDIO_BITRATE    - Bitrate of the output audio stream
#   VIDEO_RESOLUTION - Resolution to re-encode the video with
#   VIDEO_BITRATE    - Bitrate of the output video stream

# Stop execution if any error occurs
set -e

# Generate the ffmpeg option strings based on the env vars
if [ "$AUDIO_CHANNELS" != "" ] && [ "$AUDIO_CHANNELS" != "0" ] ; then
  channels="-ac $AUDIO_CHANNELS"
fi

if [ "$AUDIO_BITRATE" != "" ] && [ "$AUDIO_BITRATE" != "0" ]; then
  audiorate="-b:a $AUDIO_BITRATE"
fi

if [ "$VIDEO_RESOLUTION" != "" ]; then
  resolution="-s $VIDEO_RESOLUTION"
fi

if [ "$VIDEO_BITRATE" != "" ] && [ "$VIDEO_BITRATE" != "0" ]; then
  videorate="-b:v $VIDEO_BITRATE"
fi

# Convert the video stream into WebM
ffmpeg -y -re -i /input/$1 -acodec libvorbis $channels $audiorate -vcodec libvpx $resolution $videorate /output/$2

# Patch up the output file to support media source streaming
mse_webm_remuxer /output/$2 /output/$2.remux
mv -f /output/$2.remux /output/$2
