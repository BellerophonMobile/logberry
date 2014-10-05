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
	"errors"
	"os"
)

func somecomputation(data interface{}) (int, error) {
	return 7, nil
}

func main() {

	// Uncomment for JSON output
	// logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))

	// Output autogenerated build information and a hello
	logberry.Main.Build(buildmeta)

	logberry.Main.Info("Start program")

	// Construct a new component of our program
	cmplog := logberry.Main.Component("MyComponent")

	// Create some structured application data
	var data = struct {
		MyString string
		MyInt    int
	}{"alpha", 9}

	// Do some activity on that data, which may fail, within the component
	tlog := cmplog.Task("Some computation", &data)
	res, err := somecomputation(data)
	if err != nil {
		tlog.Error(err)
		os.Exit(1)
	}
	tlog.Complete(&logberry.D{"Result": res})

	// Shut down the component
	cmplog.Finalize()

	// An error has occurred out of nowhere!
	logberry.Main.Fatal("Unrecoverable error", errors.New("Arbitrary fault"))

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

![Colored Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/minimal-colors.png)

Uncommenting to switch to JSON output and re-running produces:

![JSON Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/minimal-json.png)

Of note in these logs are:

 * The verbose build information, automatically populated using an
   included script.
 * There are two types of logging contexts: Components and tasks.
   Both have thread-safe unique identifiers, which are included in all
   log entries.
 * The hierarchical relationship between contexts is automatically
   logged.
 * Tasks are logged in a systematized fashion that promotes outputting
   all relevant data for both errors and success, without messy and
   duplicative marshaling code.
 * Log entries have a structured, parseable format, with easily
   extracted data fields.
 * The API is very flexible, from fairly structured usage as above to
   [completely flat
   logging](https://github.com/BellerophonMobile/logberry/blob/master/examples/toplevel/main.go)
   similar to a more typical logging package.

## Read More

Links to documentation:

 * [Related Work](https://github.com/BellerophonMobile/logberry/blob/master/docs/related.md) --- Links to and notes on some other logging packages.
 * [Motivations](https://github.com/BellerophonMobile/logberry/blob/master/docs/motivations.md) --- Lengthy discussion on the design rationale behind Logberry.
 * [Examples](https://github.com/BellerophonMobile/logberry/blob/master/examples/) --- Simple programs demonstrating Logberry usage; notables include:
   *  [flightpath](https://github.com/BellerophonMobile/logberry/blob/master/examples/flightpath/main.go) --- A fanciful but somewhat structured program.
   * [blueberry](https://github.com/BellerophonMobile/logberry/blob/master/examples/blueberry/main.go) --- An example that actually does something.


## Changelog

 * **2014/10/foo: Release 1.0!** Though definitely not mature at all,
   we consider Logberry to be usable.


## License

Logberry is provided under the open source
[MIT license](http://opensource.org/licenses/MIT):

> The MIT License (MIT)
>
> Copyright (c) 2014 [Bellerophon Mobile](http://bellerophonmobile.com/)
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
