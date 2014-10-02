Logberry <img src="https://raw.githubusercontent.com/BellerophonMobile/logberry/master/logberry.png" height="64" title="Logberry" alt="Picture of a strawberry" />
========

Logberry is a structured logging package for [Go](http://golang.org/)
services and applications.  It is focused on generating logs, rather
than managing them, and tries to be lightweight while also taking a
more semantic approach than is typical.

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
	"errors"
	"github.com/BellerophonMobile/logberry"
)

func main() {

	logberry.AddOutput(logberry.NewStdOutput())

  logberry.Build(buildmeta)
  logberry.CommandLine()

  logberry.Info("The bananas have gotten loose!",
		&logberry.D{"Status": "insane"})

  logberry.Error("Could not launch", errors.New("Failure in hyperdrive!"))

	logberry.Warning("Power drain");

  logberry.Info("Continuing on", &logberry.D{"Code": "Red"})

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


## Related Work

There are a number of logging libraries available for Go, including:

 * [log](http://golang.org/pkg/log/): The built-in simple logging
   package.  Has basic levels and just outputs to an IO stream, making
   it easy to swap in different endpoints.

 * [Logrus](https://github.com/sirupsen/logrus): A great library that
   features structured event reports made of key/value pairs.  Also
   includes a set of hooks for targeted special handling, e.g.,
   throwing exceptions to [Airbrake](https://airbrake.io/).

 * [loggo](https://github.com/juju/loggo): Focused around hierarchical
   component labeling for events, where each tier can be configured
   for level, output format, and target independently.

 * [FactorLog](https://github.com/kdar/factorlog): Fairly standard
   logging package, just slightly heavier than the built-in logging
   but more customizable and with several neat features like
   outputting source code snippets.

 * [glog](https://github.com/golang/glog): Fairly standard logging
   package, with some different API syntax stylings.  Primary point
   from many other libraries is that it exposes and encourages using
   the verbosity levels within conditionals to avoid parameter
   expansion.  Unfortunately, without macros as in C/C++ there doesn't
   seem to be any cleaner way to do this that's not excessively
   convoluted, i.e., through Go's file-level build flag conditionals.

 * [log4go](http://code.google.com/p/log4go/): Similar to log4j,
   focused on output management (rotation, etc.), highly customizable
   syntax, and runtime XML configuration.

Logberry shares or borrows ideas and code from several of these.
Notable mentions include structured event reports as well as terminal
detection and coloring from
[Logrus](https://github.com/sirupsen/logrus), and source code location
determination from [FactorLog](https://github.com/kdar/factorlog)
(which seems to have gotten it from
[glog](https://github.com/golang/glog)).


## Motivation

This library was directly motivated by experiences developing a
reasonably complex component used within software created by a large
team spread across multiple corporations and several years of
development.  The most important observations include:

 * Logs should include unambiguous indication of the executing code.
   Build timestamps are useful, but the source control commit ID is
   critical.  This includes indicating whether or not it is a modified
   repository.  Ideally this also captures the build environment, at
   least the machine/user.

 * Concurrent execution paths need to be cleanly differentiated within
   log reports.  Essentially all loggers protect concurrent event
   output.  But in a typical unstructured log API it's all too easy to
   have an error or informational statement that doesn't properly
   indicate which of multiple simultaneous tasks generated the report.

 * Basic timing should be available for easy monitoring of larger
   tasks.  Even simple measures aid top level profiling, identifying
   major issues and confirming non-issues.  Having such tools
   available upfront and built into logging makes it easy to utilize
   early in development and without any effort.  Such measurement
   should however be optional, potentially at runtime, as there is
   typically a performance cost.

 * Logs should be mechanically accessible and manipulable.  Particular
   events, execution strands, task identifiers, call stacks, times,
   and other data should all be easily located and parsed.  Ideally it
   should be straightforward to extract inputs against a particular
   sub-component for re-running in detailed debugging, or to write
   simple scripts and tools to summarize or visualize logs.

 * Like all libraries, logging should be simple to use and quick to
   get started.  However, The logging needs of a small, one-off, fun
   project are not the same as a large, multi-aspect component built
   by a changing pool of developers.  Some additional complexity and
   effort is well worthwhile if it reduces the overall work.  In
   particular, any time gained via trivial logging APIs can easily be
   immensely dwarfed by time spent in processing and deciphering
   generated logs.

 * The distinction between verbosity levels should be fairly clear,
   ideally with fine grained control.  Simple ordinal verbosity levels
   achieve neither, giving no clear guidance to the developer as to
   which to use, and no mechanism to hone in on a particular strand of
   examination in debugging.

 * A frequently encountered problem working with other developers is
   receiving partial logs or mere snippets, excluding critical
   information.  It's therefore advantageous to make events as atomic
   as possible, capturing at once all the information needed to
   decipher what happened.  Strongly counterbalancing this though is
   the verbosity entailed.  Beyond human readability problems,
   somewhat manageable with tools, there are environments and
   deployments where log size is an important concern, e.g.  e.g.,
   live reporting or even post-run harvesting over a network with
   limited bandwidth or frequent reporting.

 * Services and applications do have some differences when it comes to
   logging.  For example, [Logrus](https://github.com/sirupsen/logrus)
   doesn't do any log management, such as rotation, because that
   "should not be a feature of the application-level logger."  That's
   largely true for both services and applications, but a long-lived
   service does need to take management into account to some extent.
   An example includes incorporating some mechanism to ensure output
   streams don't break when a daily log file is rotated out.

More obvious or smaller points include:

 * Components, modules, or libraries need logging controls
   encapsulated to their scope---they can't all be tucked under the
   same logcat tag!  At the same time, logging needs to be manageable
   from the outer software using that component.  This typically means
   setting verbosity and output targets.

 * The library itself cannot rely on its own command line or settings
   file configuration, though it can include tools for doing so.
   There's too much variability if the component is being used in
   applications and services, desktops and mobile devices, and so on.

 * Logs should handle arbitrarily sized event output.  Somebody might
   very well push an 18KB SQL statement through your logger, and that
   might actually be a useful and convenient thing for them to do.

 * Though not a core feature, it can be useful for the library to
   provide tools for easily connecting logging to a live view, such as
   a web page over HTTP, in addition to standard file output and such.

 * On any modern platform with reasonable resources, you essentially
   always want to log at least some things.  In many settings you may
   also want to move from production to debug levels of output without
   recompiling.  The focus many C/C++ libraries place(d) on entirely
   compiling out verbose log statements is not critical.  That said,
   report parameters that are expensive to generate should still be
   able to be skipped entirely.


## API

This section outlines the Logberry API.

### Output Drivers

Logberry's output is generated by output drivers.
