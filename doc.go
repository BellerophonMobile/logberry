/*
Package logberry implements a structured logging framework.  It is
focused on generating logs, rather than managing them, and tries to be
lightweight while capturing more semantics and structure than is
typical, in readable and easily parsed forms.

There are five central concepts:

  Root            - Logging's outgoing interface and output controller.
  Context         - A conceptual structure in which logging is produced.
      Component   - An object, class, or cluster of related functionality.
      Task        - A specific execution sequence toward some objective.
   D               - Specific event or context data to be logged.

More documentation is available from the repository and README:
  https://github.com/BellerophonMobile/logberry

*/
package logberry
