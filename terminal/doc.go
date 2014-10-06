/*

Package terminal provides platform specific terminal detection
functions.  It is part of the Logberry logging framework, used by some
of the output drivers for tasks such as determining whether or not to
use terminal colors.  Its interface consists entirely of a single
function, IsTerminal(fd), which returns true if the file descriptor fd
is a terminal.

*/
package terminal
