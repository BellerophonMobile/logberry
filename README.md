Logberry <img src="https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/logberry.png" height="64" title="Logberry" alt="Picture of a strawberry" />
========

Logberry is a structured logging package for [Go](http://golang.org/)
services and applications.  It is focused on generating logs, rather
than managing them, and tries to be lightweight while capturing more
semantics and structure than is typical, in both readable and easily
parsed forms.


## Introduction

Most log output libraries fall into one of two camps:

 * Super simple, with a global API that's really easy to use but has
   no structure beyond a component label and message level or type;
 * More complex, but focused on extensive formatting controls and/or
   various output serializers, pipes, aggregators, and managers.

Logberry is a bit different, and places more focus on *what* you're
logging, rather than *how*.  At the core, its log events are based
around key/value pairs rather than arbitrary strings, much like
[Logrus](https://github.com/sirupsen/logrus).  On top of that is a
very light, optional structure for capturing execution stacks,
delineating concurrent output, basic task timing, and other generic
semantics that encourage better, more useful event log structure.

Supporting all that are some very simple concrete output options, much
like many other libraries.  In fact, those tools can be easily dropped
in at this layer.  Although it stands alone just fine, there's a good
chance Logberry is complementary to, rather than competing with, your
preferred log output engine.


## Minimal Example

Logberry can be installed directly through Go:

    go get github.com/BellerophonMobile/logberry

A very simple use is then:

```go
package main

import (
	"github.com/BellerophonMobile/logberry"
)

var target = struct {
	Destination string
	Priority int
}{
	"Tau Ceti",
	1,
}

func main() {

	logberry.Main.Build(buildmeta)

	logberry.Main.CommandLine()

	logberry.Main.Info("Spooling engines", &logberry.D{"Block": "C", "Power": 2})
	logberry.Main.Info("Calculating flightpath", target)
	logberry.Main.Info("Ignition")

	logberry.Main.Warning("Power drain");

	e := logberry.Main.Failure("Failure in hyperdrive!")

	logberry.Main.Fatal("Aborting flight", e)

}
```

The buildmeta object is an instance of `logberry.BuildMetadata`.  It
can be managed by hand or generated using the utility script in
`logberry/util/build-stmt-go.sh`, i.e., as part of the build process:

```sh
  $ ./util/build-stmt-go.sh > examples/minimal/build.go
  $ go build github.com/BellerophonMobile/logberry/examples/minimal
  $ ./minimal
```

In a terminal this produces the output:

![Colored logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/minimal-colors.png)

## Read More

 * [Related Work](https://github.com/BellerophonMobile/logberry/blob/master/docs/related.md)
 * [Motivations](https://github.com/BellerophonMobile/logberry/blob/master/docs/motivations.md)

## License

Logberry is provided under the open source
[MIT license](http://opensource.org/licenses/MIT):

> The MIT License (MIT)
>
> Copyright (c) 2014 Bellerophon Mobile
> 
>
> Permission is hereby granted, free of charge, to any person
> obtaining a copy of this software and associated documentation files
> (the "Software"), to deal in the Software without restriction,
> including without limitation the rights to use, copy, modify, merge,
> publish, distribute, sublicense, and/or sell copies of the Software,
> and to permit persons to whom the Software is furnished to do so,
> subject to the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
> BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
> ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
