// Package display implements a display server (displayd) similar to the full
// featured X.org or Wayland on modern linux distributions, just very tiny and
// for a text-based LCD display.
//
// In contrast to a normal lcd driver which can output characters at a specified
// position on the lcd display, displayd introduces the metaphor of "lines"
// (which have a position and a fixed width text with optional scrolling) and
// "windows" which group together a named set of lines. There can be more lines
// in a window then the lcd has physically, since windows can be "moved" and
// "truncated" to implement menus and similar constructs on the client-side.
// Clients can have many windows which they can swap easily by "switching"
// to a certain active window.
//
// For writing on the actual LCD display it relies on a driver programm
// that reads lines from stdin with this simple protocol:
//
//    <line>[,<offset> <Text...>
//
//  Examples:
//
//    0 Text at line 0 (first line)
//    1,5 Text at line 1 and offset 5
//
// It is expected that the driver manages to not re-render unchanged areas.
// Displayd will then print the LCD matrix to the driver periodically.
// Why no event based loop? because the usual usecase of several scrolling
// lines will actually update more often than a polling based approach.
//
// By it's architecture it also enables the simultaneous write to the display
// in an ordered fashion and makes it easy for the programmer to use normal
// utf8 encoding while silently subsituting it with lcd compatible codepoints.
//
// displayd is controlled by a simple line based text protocol and supports
// currently the following commands:
//
// switch <win>               -- Switch to <win>, creating it if needed.
// line <win> <pos> <text>    -- Write <text> in, possibly new, line <pos> of <win>.
// move <win> <off>           -- Move <win> down by <off> (may be negative)
// truncate <win> <max>       -- Truncates <win> to <max> lines.
// render                     -- Outputs the current active window to the socket.
// close                      -- Terminates the connection.
// quit                       -- Terminates displayd.
// scroll <win> <pos> <delay> -- Make line <pos> of <win> scrolled with speed <delay>
//                               (default: 0 -> disabled)
//
package display
