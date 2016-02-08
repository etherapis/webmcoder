## Background

*The [Media Source Extension (MSE)](https://w3c.github.io/media-source/) is a
[W3C](https://www.w3.org/) specification that allows JavaScript to send byte
streams to media codecs within web browsers that support HTML5 video. This allows
the implementation of client-side prefetching and buffering code for streaming
media entirely in JavaScript.* ~[Wikipedia](https://en.wikipedia.org/wiki/Media_Source_Extensions)

To formulate it more clearly, MSE is a way to generate a video stream programatically
from JavaScript within a browser. This is an essential feature when the data source
behind a video stream is not a simple flat file that the browser can pull from a web
server (e.g. multiple sources, AJAX queries, etc).

## WebMCoder

Unfortunately, the MSE is experimental at best, each major browser having various
issues with assembling and playing chunked streams of video data. These include
limited video/audio codec support, the need for specialized stream segmentation,
half baked implementations requiring various hacks to work around.

WebMCoder is a tiny tool to solve the painful problem of converting a video file
into a format suitable for streaming through the MSE. It should be able to convert
most input media streams into the open standard WebM format in a suitable way for
streaming via the media source extensions.

## Installation and usage

WebMCoder is a [Go](https://golang.org/dl/) tool using [Docker](https://docs.docker.com/engine/installation/)
images. These must be installed prior to being able to use the encoder.

Installing the Go wrapper can be done via:

```
$ go get github.com/etherapis/webmcoder
```

### Usage

In its crudest form, WebMCoder can be invoked to simply convert an input video
file into a MSE suitable output WebM file:

```
$ webmcoder input.ext output.webm
```

A few video and audio conversion options are also supported:

 * `--achan`: Number of audio channels to generate (0 = same as input)
 * `--arate`: Audio bitrate to encode the output to (0 = same as input)
 * `--vrate`: Video bitrate to encode the output to (0 = same as input)
 * `--vres`: Video resolution (WxH) to encode the output to (empty = same as input)

A full command to recode the [Elephants Dream](https://orange.blender.org/) movie
for streaming via media source extensions with a bit of reprocessing (mono audio,
640 x 360 video resolution, 674 kbit/s video bitrate) would be:

*Note, the Elephants Dream HD is a 815MB download!*

```
$ curl -L -O http://video.blendertestbuilds.de/download.blender.org/ED/ED_HD.avi
$ webmcoder --achan=1 --vres=640x360 --vrate=691200 ./ED_HD.avi ./elephants-dream.webm

... Approximately 10 mins later ...

$ ls -al
-rw-r-----  1 karalabe karalabe 854537054 Feb  7 18:00 ED_HD.avi
-rw-r--r--  1 karalabe karalabe 60962855 Feb  7 18:10 elephants-dream.webm
```

### Embedding

Just to provide a rough sketch on how to embed the above generated media stream
into a webpage, first an HTML5 `video` element is needed.

```html
<video id="videosink" width="640" height="360" controls>
  <source src="" type="video/webm">Your browser does not support the video tag.
</video>
```

After which we need to retrieve our video stream chunks from some arbitrary data
source, and assemble them into a `MediaSource` element. Please note, the below
code is only pseudocode highlighting the stream parts, the rest is up to the user
to implement.

```js
// Create a new media source to fill with the video stream
var player = document.getElementById("videosink");
var source = new MediaSource;

source.addEventListener('sourceopen', function() {
  // Media source ready, create a video buffer and an async queue to fill with data
  var queue  = [];

  var buffer = source.addSourceBuffer('video/webm; codecs="vorbis,vp8"');
  buffer.addEventListener('updateend', function() {
    if (queue.length > 0) {
      buffer.appendBuffer(queue.shift());
    }
  }, false);

  // Start downloading the video stream, inserting either into the queue or the buffer
  while (!complete) {
    var chunk = downloadNextChunk();

    // Queue up the next chunk or insert directly
    if (queue.length > 0 || buffer.updating) {
      queue.push(data);
    } else {
      buffer.appendBuffer(data);
    }
  }
});

// Initialize the video player with the chunked stream source
player.src = URL.createObjectURL(source);
```

## Credits

The conversion tool is based on:

 * [Docker](https://www.docker.com/): containerizing the tool dependencies
 * [FFmpeg](https://www.ffmpeg.org/): converting the data streams to WebM
 * [mse-tools](https://github.com/acolwell/mse-tools): prepping WebM files for streaming

## License

WebMCoder is licensed under the BSD 3-Clause License. For details please consult
the [LICENSE](https://github.com/etherapis/webmcoder/blob/master/LICENSE) file
contained within the repository root.
