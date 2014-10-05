# Related Work

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
