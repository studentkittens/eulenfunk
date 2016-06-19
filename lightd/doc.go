// Package lightd implements a daemon that is able to process effects
// send to it and send them to an external, simpler driver program.
// It also manages the locking of the light resource.
//
// The external driver program should fulfill the following requirements
// For eulenfunk, the driver program is called radio-led.
//
// lightd can be controlled by a simple line based network protocol,
// which currently supports the following commands:
//
// !lock      -- Try to acquire lock or block until available.
// !unlock    -- Give back lock.
// !close     -- Close the connection.
// <effect>   -- Lines starting without ! are parsed as effect spec.
//
// <effect> can be one of the following:
//
//   {<r>,<g>,<b>}
//   blend{<src-color>|<dst-color>|<duration>}
//   flash{<duration>|<color>|<repeat>}
//   fire{<duration>|<color>|<repeat>}
//   fade{<duration>|<color>|<repeat>}
//
// where <*-color> can be:
//
//   {<r>,<g>,<b>}
//
// and where <duration> is something time.ParseDuration() understands.
// <repeat> is a simple integer.
//
// Examples:
//
//   {255,0,255}                     -- The world needs more solid pink.
//   fire{1ms|{255,255,255}|0}       -- Warm fire effect.
//   blend{{255,0,0}|{0,255,0}|2s}   -- Blend from red to green.
//
package lightd
