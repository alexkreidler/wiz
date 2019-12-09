# Guide

Lets get our hands dirty building some cool stuff.

First, you'll need the `wiz` CLI. See our [installation instructions](#installing)

Start off by initializing a package. It'll ask you for some basic information about it. Most of the time, the defaults are fine. (If you come from a Node.js background, this should be very similar).
<!-- We strongly encourage sharing ML packages and code, so this tutorial includes using `git`. -->

```
$ mkdir wiz_example
$ cd wiz_example
$ wiz init
package name: (wiz_example)
type: (program)
version: (1.0.0)
description: A simple example using NLP + ImageNet GANs
training mode: (wiz.yml)
running mode: (run.py) interactive.py
dependencies: bagowords, word2vec, googlenet, fcn
git repo:
keywords: example, nlp, gans, imagenet
author: Oz Wizard
license: (ISC)
Done!
```

Alright, time to string together some packages. Open up wiz.yml and add the following:
```
runtime: python
wiz.linear:
  nlp:
    - bagowords
    - word2vec
  gan:
    init:
      package: googlenet
      layers: end-1
    end: fcn
  output:
    type: tmpfile
    format: png
train:
```
Finally, lets setup the input/deployment method. In `interactive.py`, add:
```
import wiz
import Image
import time

def output(f):
  image = Image.open('File.jpg')
  image.show()
  time.sleep(10)
  image.close()

# Output setup
wiz.outputfn(output)

print("Input an image name, we will generate it. Type quit to exit.")
while True:
  input = read()
  if input == "quit":
    break
  wiz.input(input)
```



## Installing
We love [Docker](https://docker.com) here at Wiz Project. If you use it too, you can just run `docker run -it wizproject/wiz` and you'll get a nice shell with the latest `wiz` already installed.

Some people are wary of using Docker to do ML/DL because it doesn't have great support for GPUs. The Moby Project will split up docker into components over the next few months, and we will be adding a GPU manager component (similary to the nvidia-docker plugin, but with better UX)

Remember, you can mirror a directory on your computer inside the container using a volume. If you run `docker run -it -v $(pwd):/wiz wizproject/wiz` then your current directory will be shared between both and you can use your favorite editor to write the code.

#### Build from source

Make sure you have the `go` tool, and that `$GOPATH/bin` is in your `$PATH`

```
git clone https://github.com/tim15/wiz
cd wiz
go install ./cmd/wiz
```

Then you're all set to run `wiz`

<!--
### MacOS
```
brew install wiz
```
### Debian/Ubuntu
```
$ sudo add-apt-repository download.wizproject.ml
$ sudo apt-get update
$ sudo apt-get install wiz
``` -->
