# Concepts

Logberry has four top level concepts/objects:

 * `Root`: A consolidation point between Tasks and OutputDrivers.
 * `OutputDriver`: Serializers for publishing events.
 * `Task`: A component, function, or logic path that generates events.
 * `D`: Data to be published with an event.

Also important are two less fundamental but included concepts/objects:

 * `BuildMetadata`: A simple representation of the build environment.
 * `Error`: A generic structured error report.


## Root

All logging is coordinated by a Root.  Tasks under a Root generate
events, which are given to OutputDrivers to actually publish.
Typically you would only have the default Root in a program.

Generated events are thown into a channel, which may or may not be
buffered.  A separate goroutine continually processes events from the
channel in the background, dispatching them to registered outputs.
The host program needs to call `Root::Stop()` on a BackgroundRoot to
ensure that all generated events are output.  Otherwise it is possible
for the background goroutine to not activate before a short running
program terminates, leaving events in the buffer.


## OutputDriver

Events are actually recorded, aggregated, or otherwise processed by
OutputDrivers.  Applications may implement its interface and provide
their own, but there are two included in Logberry:

 * TextOutput: Human readable output, colorized if outputting
   to a terminal.
 * JSONOutput: Machine readable JSON formatted output.

TextOutputs include a program label on each line.  By default
`logberry.Std` has a registered TextOutput with a program label
derived from the executing process' filename that writes to stdout.

OutputDrivers are registered using `Root::AddOutputDriver()` or
`Root::SetOutputDriver()`.  E.g., to switch the default to JSON
formatting:

```go
  logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))
```

Roots may have multiple OutputDrivers, all of which receive each event
for that Root.  A program may also utilize multiple Roots at once.  A
single OutputDriver instance should not be registered to more than one
Root simultaneously unless its specific documentation notes otherwise.


## Tasks

Log events are generated via Task objects.  These represent a
particular component, function, or related block of logic, ranging
anywhere in scope from an entire program to a single library call.

For example, by default `logberry.Main` is a Task under the
`logberry.Std` Root.  Using it, programs can output events much like
any other flat logging interface, except with structured data, e.g.:

```go
	logberry.Main.Info("Computed data", logberry.D{"X": 23, "Label": "priority"})

	logberry.Main.Failure("Arbritrary failure")
```

At the opposite scope, Tasks can also be used to properly and easily
log fine grained tasks, e.g., calling a function such as opening a
file or performing a specific computation:

```go
	task := computerlog.Task("Compute numbers", &data)
	res, err := somecomputation()
	if err != nil {
		task.Error(err)
		return
	}
	task.Success(logberry.D{"Result": res})
```

This snippet creates a new Task to log a specific function call within
a larger component represented by the `computerlog` Task.  The
function task is given a specific human-oriented activity label
("Compute numbers").  The task as a whole is then associated with some
data (`data`).  After the computation is performed, the success or
failure of the task is reported.  Either report will include the
associated data, eliminating the marshaling redundancy or suboptimal
reporting of more typical logging.  In addition, a successful outcome
reports additional data particular to that outcome.

### Hierarchy

Tasks are created using the `Task` or `Component` functions of either
Roots or Tasks:

```go
	// Create a program component---a long-running, multi-use entity.
	computerlog := logberry.Main.Component("computer")

	
	// Execute a task within that component, which may fail
	task := computerlog.Task("Compute numbers", &data)
```

They thus have a hierachical structure originating in a Root.  This
structure may be reported by the OutputDrivers, as it is by the
built-in drivers, to enable easily reconstructing program execution
structure even across interleaved goroutines.  Each Task has a numeric
identifier unique to that program instance, and the identifiers for
both a Task and its parent are included in the standard outputs.


### Components

Both the Task and Component creation functions return a Task.  The
only difference is one of human semantics.

All Tasks have a component tag included in event reports to indicate
of which functional area the task is part.  E.g., the default for
`logberry.Main` is 'main' while a sub-task might be tagged
'websocket', 'mapper', or any other program specific label.  Tasks
also have a human-oriented activity text, e.g., 'Save configuration'
or 'Connect to database'.

By default Tasks inherit the component tag of their parent and are
given a text label specifying some focused activity.  Tasks created
using the Component creator though are assigned the given component
label, presumably different from that of their parent Task.  Their
activity text is also generated to identify that component.
Termination of the component may then be logged using
`Task::Finalized()` or the error reports as needed.

Note, however, that these component Tasks are just regular Task
objects that apply a few conventions when created.


### Data

Tasks have data associated with them, captured by a D object as
described below.  This data may be aggregated into the object over
time and is reported with all its generated events, alongside any data
given specific to each event.  For example, a task for accessing an
HTTP endpoint might start with only the resource known and associated
with the task.  After the user is authenticated, their identifier
might be added to the task.  Each of these will be included in
subsequent log events.  The task might then terminate on success by
additionally reporting the number of bytes transmitted.

Data to be associated with a Task may be passed to creation functions.
The `Task::AddData(data...)` function may also be used to assert data
as the Task continues.

Event specific data may be included in all of the reporting functions
outlined in the following.  This data does not aggregate into the Task
for output in subsequent calls.


### Reporting

The core logging function is `Task::Event(event, msg, data)`.  This
generates from the Task an event of type indicated by the `event`
string, the human-oriented short message string `msg`, and the
arbitrary variadic event-specific `data` to be reported in a
structured fashion.

A simple example is:

```go
	// Generate an application specific event reporting some other data
	var req = struct {
		User string
	}{"tjkopena"}

	computerlog.Event("request", "Received request", req)
```

Built on top of this are a variety of common event triggers:

 * Informational: Report generic data, human messages, or warnings.
   * `Info`
   * `Warning`
 * State: Mark the end of initialization or pause in processing.
   * `Ready`
   * `Stopped`
 * Lifetime: Report the termination state of components or tasks.
   * `Success`
   * `Error`
   * `Failure`
   * `Finalized`

There are also functions to report program or module initialization:

   * `BuildMetadataEvent`
   * `BuildSignatureEvent`
   * `ConfigurationEvent`
   * `CommandLineEvent`
   * `EnvironmentEvent`
   * `ProcessEvent`

Details on all these functions may be found in the [API
GoDocs](https://godoc.org/github.com/BellerophonMobile/logberry).

## D

Data associated with Tasks and events are captured by D, a simple Go
map from string keys to values of arbitrary type.  Most of the
Logberry functions have a variadic `...interface{}` parameter to pass
data to the call, which depending on the function is either associated
to a Task or reported along with a specific event, e.g.:

```go
// Task creates a new sub-task.  Parameter activity should be a short
// natural language description of the work that the Task represents,
// without any terminating punctuation.  Any data given here will be
// associated with the Task and reported with all its events.
func (x *Task) Task(activity string, data ...interface{}) *Task

// Error generates an error log event reporting an unrecoverable fault
// in an activity or component.  An error is returned wrapping the
// original error with a message reporting that the Task's activity
// has failed.  Continuing to use the Task is discouraged.  The
// variadic data parameter is aggregated as a D and embedded in the
// generated error.  It and the data permanently associated with the
// Task is reported with the event.  The reported source code position
// of the generated task error is adjusted to be the event invocation.
func (x *Task) Error(err error, data ...interface{}) error {	
```

D is just a standard Go map, with its real functionality in several
functions for copying instances from arbitrarily typed and variadic
parameters.  In particular, it will automatically incorporate exposed
fields of structs and keys of maps as keys of the constructed D.  In
this way it is very easy to incorporate objects into a log event and
have the data reported in a structured fashion.  E.g., the following
log event will have `DataLabel` and `DataInt` fields:

```go
	// Create some structured application data and log it
	var data = struct {
		DataLabel string
		DataInt    int
	}{"alpha", 9}

	logberry.Main.Info("Reporting some data", data)
```

Several functions to produce properly delimited but human-readable
`key=value` text printouts of D objects are also provided.  However,
note that as they're simply Go maps, they're trivial to throw into
JSON marshaling or other serialization.


## BuildMetadata

Logberry also has built in reporting functions and simple
representation of a program's build environment, along with script
tools to automatically construct those representations.

First of these is `BuildSignatureEvent()`, which takes and reports a
string as some arbitrary stamp of the binary's build profile.  The
script `util/build-signature.sh` will automatically generate such from
basic host device parameters and the Git repository assumed to be the
working directory, e.g.:

```sh
joe@scully ../github.com/BellerophonMobile/logberry (git)-[master] % ./util/build-signature.sh
logberry master 8aff1c9174c6b23309bb64d094419b90a2687a5d* scully joe 2015-06-19 10:51:08-04:00
```

This string is then easy to pass in to be compiled with a program via
Go's linker flags, e.g.:

```sh
go install -ldflags "-X main.buildsignature '`./util/build-signature.sh
```

More expressively, `BuildMetadataEvent()` takes and reports a more
in-depth, structured representation of the build environment.  The
program `util/build-metadata.go` constructs a BuildMetadata object,
and is intended to be executed using `go run` or `go generate` to
create a Go source code file to be included into the application,
e.g.:

```sh
joe@scully ~/chimerakb/code/workspace (git)-[master] % go run src/github.com/BellerophonMobile/logberry/util/build-metadata.go

/**
 * This file generated automatically.  Do not modify.
 * Generated from workspace: .
 */

package main

import "github.com/BellerophonMobile/logberry"

var buildmetadata = &logberry.BuildMetadata{
  Host:     "scully",
  User:     "joe",
  Date:     "2015-06-19T10:56:27-04:00",

  Repositories: []logberry.RepositoryMetadata {

    logberry.RepositoryMetadata{
      Repository: "workspace",
      Branch:     "master",
      Commit:     "c8cda1e7eeab3486691a207865d51c3f0782d3d8",
      Dirty:      false,
      Path:       ".",
    },

    logberry.RepositoryMetadata{
      Repository: "core",
      Branch:     "master",
      Commit:     "6847dda91a7f940780f446f17a7c6c48f2d8dd10",
      Dirty:      false,
      Path:       "src/chimerakb.com/pkg/core",
    },

    logberry.RepositoryMetadata{
      Repository: "public",
      Branch:     "master",
      Commit:     "cbc639561e957a0c13882478b7d47775c21d93ce",
      Dirty:      true,
      Path:       "src/chimerakb.com/pkg/public",
    },

    logberry.RepositoryMetadata{
      Repository: "logberry",
      Branch:     "master",
      Commit:     "8aff1c9174c6b23309bb64d094419b90a2687a5d",
      Dirty:      true,
      Path:       "src/github.com/BellerophonMobile/logberry",
    },

  },
}
joe@scully ~/chimerakb/code/workspace (git)-[master] % 
```

That example shows the intended use case: A large project made up of
several elements organized in their own repositories is organized and
compiled inside a larger workspace repository.  The script records a
timestamp as well as the user and host for the build, and scans for
all git repositories of and under the build directory.  This is logged
using `Task:BuildMetadata()` to more or less unambiguously identify
the exact build composition of a binary, a critical piece of data for
complex and especially multi-team software.

Due to some of the specifics of `go generate`, local scripts can be a
little tricky to invoke with that tool.  One way to do that is to set
an environment variable identifying the root workspace directory and
execute as such:

```go
//go:generate go run $WORKSPACE/src/github.com/BellerophonMobile/logberry/util/build-metadata.go -workspace=$WORKSPACE -out=build
```

This generates a file `build.go` in the invoking file's directory
containing a BuildMetadata object capturing the entire workspace.
Note that the `-out` flag takes the given target as the prefix of a Go
source file and adds ".go" to the filename.  This is done as `go run`
otherwise interprets arguments ending in ".go" as files to execute.
To specify a different target without this suffix being added, use the
`-target` option.  The package in which the generated object is placed
may be set using the `-package` flag; the default is `main`.

## Error

Finally, Logberry includes a simple generic Error object.  It includes
a human-oriented short message, structured data captured as a D, the
source code location from which the event originates, and optionally
an underlying preceding error that caused this higher level fault.
The linked error may be a generic error, not necessarily a Logberry
Error.  In this way the errors can capture a stacktrace of failures,
each reporting structured data and source code location.  There is
also a code string and a global convenience function to take an error
and check if it has any of a set of code strings.  This is useful as a
quick way to check the outcome of client code for a specific type of
error, e.g., connection closed.


## Summary

At the core, Logberry provides very flexible logging for reporting
events with structured data.  Above that though is built a lightweight
and simple but useful interface for capturing program structure and
event semantics.  Utilizing the task instance hierarchy, event types,
and the structured data output, the lifecycle of a program, its
components, and its activities may be readily captured in a human
readable or machine parseable trace.
