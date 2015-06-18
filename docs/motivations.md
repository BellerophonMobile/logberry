# Motivations

Logberry was directly motivated by experiences developing a reasonably
complex component used within software created by a large team spread
across multiple corporations and several years of development.  The
most important observations include:

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

Logberry certainly doesn't address all these points, and of course is
imperfect in many ways, but these are the goals, design rationale, and
ideas we're working to incorporate.
