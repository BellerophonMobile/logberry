Logberry <img src="https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/logberry.png" height="64" title="Logberry" alt="Picture of a strawberry" />
========

Logberry is a structured logging package for [Go](http://golang.org/)
services and applications.  It is focused on generating logs, rather
than managing them, and tries to be lightweight while capturing more
semantics and structure than is typical, in both readable and easily
parsed forms.

[![Build Status](https://travis-ci.org/BellerophonMobile/logberry.svg)](https://travis-ci.org/BellerophonMobile/logberry) [![GoDoc](https://godoc.org/github.com/BellerophonMobile/logberry?status.svg)](https://godoc.org/github.com/BellerophonMobile/logberry) 


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

## Installation

Logberry can be installed directly through Go, like most Go libraries:

    go get github.com/BellerophonMobile/logberry


## Minimal Example

At its most minimal, it can be used much like most logging interfaces,
except with structured data:

```go
package main

import (
	"github.com/BellerophonMobile/logberry"
)

func main() {
	
	logberry.Main.Info("Demo is functional")

	logberry.Main.Info("Computed data", logberry.D{"X": 23, "Label": "priority"})

	logberry.Main.Failure("Arbritrary failure")

}
```

In a terminal this produces the output:

![Colored Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/minimal-colors.png)

Note that the output has been spaced by default to work reasonably on
both wide and standard terminals, in the latter case implicitly
placing the identifiers and structured data on the line following the
primary human text.

A simple switch to JSON output produces:

![JSON Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/minimal-json.png)

This is already significant as the structured output promotes better
reporting and easier log extraction of critical data.  Also note that
the error has automatically included the source code location.


## Small Example

Besides structured event data, Logberry has a very basic notion of
program structure, represented by Task objects.  A small example:


```go
package main

import (
	"github.com/BellerophonMobile/logberry"
	// "os"
)


func main() {

	// Uncomment this and "os" import for JSON output
	// logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))

	
	// Report build information; a script generates buildmeta
	logberry.Main.BuildMetadata(buildmetadata)

	// Report that the program is initialized & running
	logberry.Main.Ready()


	// Create some structured application data and log it
	var data = struct {
		DataLabel string
		DataInt    int
	}{"alpha", 9}

	logberry.Main.Info("Reporting some data", data)


	// Create a program component---a long-running, multi-use entity.
	computerlog := logberry.Main.Component("computer")

	
	// Execute a task within that component, which may fail
	task := computerlog.Task("Compute numbers", &data)
	res, err := somecomputation()
	if err != nil {
		task.Error(err)
		return
	}
	task.Success(logberry.D{"Result": res})


	// Generate an application specific event reporting some other data
	var req = struct {
		User string
	}{"tjkopena"}

	computerlog.Event("request", "Received request", req)


	// Run a function under the component
	if e := arbitraryfunc(computerlog); e != nil {
		// Handle the error here
	}


	// The component ends
	computerlog.End()

	// The program shuts down
	logberry.Main.Stopped()

}


func somecomputation() (int, error) {
	return 7, nil
}


func arbitraryfunc(component *logberry.Task) error {

	// Start a long-running task, using Begin() to log start & begin timer
	task := component.Task("Arbitrary computation")

	// Report some intermediate progress
	task.Info("Intermediate progress", logberry.D{"Best": 9})


	// An error has occurred out of nowhere!  Log & return an error
	// noting that this task has failed, data associated with the error,
	// wrapping the underlying cause, and noting this source location
	return task.Failure("Random unrecoverable error",
		logberry.D{"Bounds": "x-axis"})

}
```

Note that the `buildmetadata` object is in a separate file, generated
by a Logberry utility script.

In a terminal this produces the output:

![Colored Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/small-colors.png)

In the JSON output this looks as follows:

![JSON Logberry terminal output.](https://raw.githubusercontent.com/BellerophonMobile/logberry/master/docs/figures/small-json.png)

Of note in this log:

 * The verbose, automatically generated build information, identifying
   all (Git) repositories found in the host project folder.
 * Every event is situated within a uniquely identified Task.  The
   hierarchical relationship between Tasks is also logged.  Together
   these enable the ready decoupling of interleaved events generated
   by parallel goroutines and the reconstruction of a causal chain of
   computation.   
 * Tasks are logged in a systematized fashion that promotes outputting
   all relevant data for both errors and success, without messy and
   duplicative marshaling code.
 * Long running tasks are automatically timed.
 * Common program events such as configuration, start, and errors are
   all identified, as well as application specific event types.
   Associated data is captured and provided in structured form.

## Read More

Links to documentation:

 * [Related Work](https://github.com/BellerophonMobile/logberry/blob/master/docs/related.md) --- Links to and notes on some other logging packages.
 * [Motivations](https://github.com/BellerophonMobile/logberry/blob/master/docs/motivations.md) --- Lengthy discussion on the design rationale behind Logberry.
 * [Concepts](https://github.com/BellerophonMobile/logberry/blob/master/docs/concepts.md) --- Top level concepts and use of Logberry.
 * [GoDocs](https://godoc.org/github.com/BellerophonMobile/logberry) --- Automatically generated API docs.

## Changelog

 * **2015/06/18: Release 2.0!** The API has been conceptually
   simplified, errors structured, and underlying code improved.
 * **2014/10/30: Release 1.0!** Though definitely not mature at all,
   we consider Logberry to be usable.


## License

Logberry is provided under the open source
[MIT license](http://opensource.org/licenses/MIT):

> The MIT License (MIT)
>
> Copyright (c) 2014, 2015 [Bellerophon Mobile](http://bellerophonmobile.com/)
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
