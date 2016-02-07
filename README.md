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
into a format suitable for streaming through the MSE. It should be ble to convert
most input media streams into the open standard WebM format in a suitable way for
streaming via the media source extensions.

## Credits

The conversion tool is based on:

 * [Docker](https://www.docker.com/): containerizing the tool dependencies
 * [FFmpeg](https://www.ffmpeg.org/): converting the data streams to WebM
 * [mse-tools](https://github.com/acolwell/mse-tools): prepping WebM files for streaming

## License

WebMCoder is licensed under the BSD 3-Clause License. For details please consult
the [LICENSE](https://github.com/etherapis/webmcoder/blob/master/LICENSE) file
contained within the repository root.
