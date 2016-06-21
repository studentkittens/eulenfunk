// Package automount implements a small mount daemon that is controllable via a
// simple line based textprotocol. The protocol supports these commands currently:
//
// mount <device> <label>   # Mount <device> (e.g. /dev/sda1) to <music_dir>/<label>
//                          # Scan this device and add music files to mpd under
//                          # a playlist named <label>
// unmount <device>         # Unmount the device again.
// close                    # Close the connection early.
// quit                     # Quit automountd.
package automount
