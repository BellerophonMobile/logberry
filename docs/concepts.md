# Concepts

Logberry has four top level concepts/objects:

 * `Root`: An interface between Tasks and OutputDrivers.
 * `OutputDriver`: Serializer for publishing events.
 * `Task`: A component, function, or logic that generates events.
 * `D`: Data to be published with an event.

Also important are two less fundamental but included concepts/objects:

 * `BuildMetadata`: A simple representation of the build environment.
 * `Error`: A generic structured error report.


## Root

All logging is coordinated by a Root.  Tasks under a Root generate
events, which are given to OutputDrivers to actually publish.

There is a default Root logberry.Std coupled by default to text output
on stdout.  However, applications may make and utilize their own Roots.

In particular, there are two kinds of Roots:

 * ImmediateRoot: Simply dispatches each event to registered outputs.
 * BackgroundRoot: Throws each event into a channel, which may or may
   be buffered.  A separate goroutine then continually processes
   events from the channel, dispatching them to registered outputs.

Both serialize incoming events such that only one is report at a time,
in order of first processing.  From the user's perspective, the big
difference is that the host program needs to call Stop() on the Root
to ensure that all generated events are output.  Otherwise it's
possible for the background goroutine to not activate before the
program terminates, leaving events still in the buffer.  However, that
buffering can be useful for long lived programs using OutputDrivers
which may take some time, e.g., publishing to a remote log service.

There is a default Root logberry.Std, which is an ImmediateRoot so
that it intuitively outputs all events without any additional calls.
However, programs need not make any use of this Root, instead
generating Tasks under custom Roots as described below.


## OutputDriver

Events are actually recorded, aggregated, or otherwise processed by
OutputDrivers.  There are two built-in:

 * TextOutput: Arguably human readable output, colorized if outputting
   to a terminal.
 * JSONOutput: Machine readable JSON formatted output.

TextOutputs include a program label on each line.  By default
logberry.Std has a registered TextOutput with a program label derived
from the executing process' filename.

OutputDrivers are registered using `Root::AddOutputDriver()` or
`Root::SetOutputDriver()`.  E.g., to switch the default to JSON
formatting:

```go
  logberry.Std.SetOutputDriver(logberry.NewJSONOutput(os.Stdout))
```

Roots may have multiple OutputDrivers, all of which receive each event
for that Root, and a program may of course utilize multiple Roots.

## Tasks

Log events are generated via Task objects.  These represent a
particular component, function, or related block of logic, ranging
anywhere in scope from an entire program to a single library call.

For example, by default logberry.Main is a Task under the logberry.Std
Root.  Using it, programs can output log events much like any other
flat logging interface, e.g.:

```go
	logberry.Main.Info("Demo is functional")

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


### Components

Tasks are created using the `Task` or `Component` functions of either
Roots or Tasks.  They thus have a hierachical structure originating in
a Root.  This structure may be reported by the OutputDrivers, as it is
for the built-in drivers, to enable easily reconstructing program
execution structure even across interleaved goroutines.  Each Task has
a numeric identifier unique to that program instance, and the
identifiers for both a Task and its parent are included in the output
of the standard drivers.

The Task and Component creation functions both return a Task.  The
only difference is one of human semantics.  All Tasks have a component
tag to be included in event reports to indicate which component area
the task is part of.  E.g., the default for logberry.Main is 'main'
while a sub-task might be tagged 'websocket', 'mapper', or any other
program specific label.  They also have a human-oriented activity
text, e.g., 'Save configuration' or 'Connect to database'.

By default Tasks inherit the component tag of their parent and are
given a text label specifying some focused activity.  They also do not
log their instantiation, though the `Task::Begin()` function may be
used to do so.  Tasks created using the Component creator though are
given a new component label.  Their activity text is also generated to
identify that component, and their instantiation logged.  Termination
of the component may then be logged using `Task::End()` in addition to
the error reports.

Note, however, that these component Tasks are just regular Task
objects that apply a few conventions when created.  Component tags and
activity texts may also be manually set or changed for all Tasks.


### Data

Tasks have data associated with them, captured by a D object as
described below.  This data may be aggregated into the object over
time and is reported with all log events, alongside any data given
specific to that event.  For example, a task for accessing an HTTP
endpoint might start with only the resource known.  After the user is
authenticated, their identifier might be added to the task.  Each of
these will be included in subsequent log events.  Finally, the task
might terminate by additionally reporting the number of bytes
transmitted.

The core mechanism for this is the `Task::AddData(key, value)`
function.  `Task::Aggregatedata`` does similarly, following the
behavior of D objects as described below.  A number of commonly used
data functions are also incorporated, to standardize when helpful on
some generic keys, e.g.:

 * `Calculation`
 * `File`
 * `Resource`
 * `Service`
 * `User`
 * `URL`
 * `Bytes`
 * `ID`
 * `Endpoint`

Associated data may also be passed in to the Task creation functions.
Event specific data may be included in all of the reporting functions
outlined in the following.  This does not aggregate into the Task for
output in subsequent calls.


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

Built on top of this basic function a variety of common event
reporting functions.  Details on these and a full up-to-date list may
be found in the [API godocs](https://godoc.org/github.com/BellerophonMobile/logberry).

 * Configuration: Report on program or module initialization.
   * `BuildMetadata`
   * `BuildSignature`
   * `Configuration`
   * `CommandLine`
   * `Environment`
   * `Process`
 * Informational: Report generic data, human messages, or warnings.
   * `Info`
   * `Warning`
 * State: Mark the end of initialization or pause in processing.
   * `Ready`
   * `Stopped`
 * Lifetime: Report the start of components or long-running tasks, and
   denote their termination state
   * `Begin`
   * `End`
   * `Success`
   * `Error`
   * `WrapError`
   * `Fatal`
   * `Failure`
   * `Die`

### Utilities

Task may additionally be muted (and unmuted).  In this state they will
produce no log events.  This is useful when using a task simply to
collect other tasks, or to generate and throw but not directly report
a structured error message.

Tasks may also be timed, and their duration reported by the
terminating lifetime events.  To start the clock, invoke
`Task::Time()`.  To additionally report the instationation, use
`Task::Begin()`. Components automatically do the latter.  The current
duration may be fetched using `Task::Clock()`.

## D

## BuildMetadata

## Error
