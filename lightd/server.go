package lightd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/studentkittens/eulenfunk/util"
)

func max(r uint8, g uint8, b uint8) uint8 {
	max := r
	if max < g {
		max = g
	}
	if max < b {
		max = b
	}
	return max
}

type effectQueue struct {
	StdInPipe io.Writer
	Blocked   chan bool
}

func (q *effectQueue) Push(e effect) {
	c := e.ComposeEffect()

	for color := range c {
		colorValue := fmt.Sprintf("%d %d %d\n", color.R, color.G, color.B)
		if _, err := q.StdInPipe.Write([]byte(colorValue)); err != nil {
			log.Printf("Failed to write to driver: %v", err)
		}
	}
}

func newEffectQueue(driverBinary string) (*effectQueue, error) {
	cmd := exec.Command(driverBinary, "cat")
	stdinpipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	blocked := make(chan bool, 1)
	blocked <- false

	return &effectQueue{
		Blocked:   blocked,
		StdInPipe: stdinpipe,
	}, cmd.Start()
}

type rgbColor struct {
	R, G, B uint8
}

type effect interface {
	ComposeEffect() chan rgbColor
}

// Common properties
type properties struct {
	Delay  time.Duration
	Color  rgbColor
	Repeat int
}

////////////////////////
// INDIVIDUAL EFFECTS //
////////////////////////

// Fade to single color and back to black
type fadeEffect struct {
	properties
}

// Shortly flash a single color
type flashEffect struct {
	properties
}

// Blend between two different colors
type blendEffect struct {
	StartColor rgbColor
	EndColor   rgbColor
	Duration   time.Duration
}

// Warm, orange fire effect for nostalgic reasons.
type fireEffect struct {
	properties
}

/////////////////////
// COMPOSE METHODS //
/////////////////////

func (color *rgbColor) ComposeEffect() chan rgbColor {
	c := make(chan rgbColor, 1)
	c <- rgbColor{color.R, color.G, color.B}
	close(c)
	return c
}

func (effect *flashEffect) ComposeEffect() chan rgbColor {
	c := make(chan rgbColor, 1)
	keepLooping := false
	if effect.Repeat < 0 {
		keepLooping = true
	}

	go func() {
		for {
			if effect.Repeat <= 0 && !keepLooping {
				break
			}

			c <- effect.Color
			time.Sleep(effect.Delay)
			c <- rgbColor{0, 0, 0}
			time.Sleep(effect.Delay)
			effect.Repeat--
		}
		close(c)
	}()
	return c
}

func (effect *fadeEffect) ComposeEffect() chan rgbColor {
	c := make(chan rgbColor, 1)

	keepLooping := false
	if effect.Repeat < 0 {
		keepLooping = true
	}

	max := max(effect.Color.R, effect.Color.B, effect.Color.G)
	go func() {
		for {

			if effect.Repeat <= 0 && !keepLooping {
				break
			}

			r := int(math.Floor(float64(effect.Color.R) / float64(max) * 100.0))
			g := int(math.Floor(float64(effect.Color.G) / float64(max) * 100.0))
			b := int(math.Floor(float64(effect.Color.B) / float64(max) * 100.0))

			for i := 0; i < int(max); i++ {
				c <- rgbColor{uint8((i * r) / 100), uint8((i * g) / 100), uint8((i * b) / 100)}
				time.Sleep(effect.Delay)
			}

			for i := int(max - 1); i >= 0; i-- {
				c <- rgbColor{uint8((i * r) / 100), uint8((i * g) / 100), uint8((i * b) / 100)}
				time.Sleep(effect.Delay)
			}
			effect.Repeat--
		}
		close(c)
	}()
	return c
}

func (effect *blendEffect) ComposeEffect() chan rgbColor {
	c := make(chan rgbColor, 1)
	go func() {
		// How much colors should be generated during the effect?
		N := 20 * effect.Duration.Seconds()

		sr := float64(effect.StartColor.R)
		sg := float64(effect.StartColor.G)
		sb := float64(effect.StartColor.B)

		for i := 0; i < int(N); i++ {
			sr += (float64(effect.EndColor.R) - float64(effect.StartColor.R)) / N
			sg += (float64(effect.EndColor.G) - float64(effect.StartColor.G)) / N
			sb += (float64(effect.EndColor.B) - float64(effect.StartColor.B)) / N

			c <- rgbColor{uint8(sr), uint8(sg), uint8(sb)}
			time.Sleep(time.Duration(1/N*1000) * time.Millisecond)
		}

		close(c)
	}()

	return c
}

func (effect *fireEffect) ComposeEffect() chan rgbColor {
	fn := func(t, n, jitter int, fac float64) uint8 {
		j := float64(rand.Int()%(jitter<<1) - jitter)
		f := float64(t - n>>1)
		return uint8(fac * (float64(-255.0/262144.0)*f*f + 255 + j))
	}

	c := make(chan rgbColor, 1)
	go func() {
		defer close(c)

		for t := 0; t < effect.Repeat; t++ {
			c <- rgbColor{
				fn(t, effect.Repeat, 50, 1.00),
				fn(t, effect.Repeat, 70, 0.10),
				fn(t, effect.Repeat, 80, 0.01),
			}
			time.Sleep(effect.Delay)
		}
	}()

	return c
}

////////////// EFFECT SPEC PARSING //////////////////

var (
	regexColorr     = regexp.MustCompile(`(c|color|)\{(\d{1,3}),(\d{1,3}),(\d{1,3})\}`)
	regexProperties = regexp.MustCompile(`.*?\{(.*)\|(.*)\|(.*)\}`)
)

func parseColor(s string) (*rgbColor, error) {
	matches := regexColorr.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Bad color: %s", s)
	}

	triple := []uint8{}
	for _, str := range matches[2:] {
		c, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("Bad color value: %s", str)
		}

		if c < 0 || c > 255 {
			return nil, fmt.Errorf("Color value must be in 0-255")
		}

		triple = append(triple, uint8(c))
	}

	return &rgbColor{triple[0], triple[1], triple[2]}, nil
}

func parseBlendEffect(s string) (*blendEffect, error) {
	// Same regex as for properties (by chance)
	matches := regexProperties.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Bad blend effect: %s", s)
	}

	srcColor, err := parseColor(matches[1])
	if err != nil {
		return nil, fmt.Errorf("Bad source color: `%s`: %v", matches[1], err)
	}

	dstColor, err := parseColor(matches[2])
	if err != nil {
		return nil, fmt.Errorf("Bad destination color: `%s`: %v", matches[2], err)
	}

	duration, err := time.ParseDuration(matches[3])
	if err != nil {
		return nil, fmt.Errorf("Bad duration: `%s`: %v", matches[3], err)
	}

	return &blendEffect{*srcColor, *dstColor, duration}, nil
}

func parseProperties(s string) (*properties, error) {
	matches := regexProperties.FindStringSubmatch(s)
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("Bad effect properties: %s", s)
	}

	duration, err := time.ParseDuration(matches[1])
	if err != nil {
		return nil, fmt.Errorf("Bad duration: `%s`: %v", matches[1], err)
	}

	color, err := parseColor(matches[2])
	if err != nil {
		return nil, fmt.Errorf("Bad color: `%s`: %v", matches[2], err)
	}

	repeatCnt, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("Bad repeat count: `%s`: %v", matches[3], err)
	}

	return &properties{duration, *color, repeatCnt}, nil
}

func parseFadeEffect(s string) (*fadeEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "fade"))
	if err != nil {
		return nil, err
	}

	return &fadeEffect{*props}, nil
}

func parseFlashEffect(s string) (*flashEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "flash"))
	if err != nil {
		return nil, err
	}

	return &flashEffect{*props}, nil
}

func parseFireEffect(s string) (*fireEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "fire"))
	if err != nil {
		return nil, err
	}

	return &fireEffect{*props}, nil
}

func parseEffect(s string) (effect, error) {
	sepIdx := strings.Index(s, "{")
	name, rest := s[:sepIdx], s[sepIdx:]

	switch name {
	case "", "c", "color":
		return parseColor(rest)
	case "fade":
		return parseFadeEffect(rest)
	case "flash":
		return parseFlashEffect(rest)
	case "fire":
		return parseFireEffect(rest)
	case "blend":
		return parseBlendEffect(rest)
	default:
		return nil, fmt.Errorf("Bad effect name: `%s`", name)
	}
}

//////////// SERVER MAIN //////////////

func lock(queue *effectQueue) {
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-queue.Blocked:
		break
	case <-timer.C:
		break
	}
}

func unlock(queue *effectQueue) {
	timer := time.NewTimer(5 * time.Second)
	select {
	case queue.Blocked <- false:
		break
	case <-timer.C:
		break
	}
}

func handleRequest(conn io.ReadWriteCloser, queue *effectQueue) {
	defer util.Closer(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Bad line:
		if len(line) == 0 {
			continue
		}

		switch line {
		case "!lock":
			lock(queue)

			if _, err := conn.Write([]byte("OK\n")); err != nil {
				log.Printf("Failed to answer lock response: %v", err)
			}

			continue
		case "!unlock":
			unlock(queue)

			if _, err := conn.Write([]byte("OK\n")); err != nil {
				log.Printf("Failed to answer unlock response: %v", err)
			}

			continue
		case "!close":
			break
		}

		effect, err := parseEffect(line)
		if err != nil {
			log.Printf("Unable to process effect: %v", err)
			continue
		}

		lock(queue)
		queue.Push(effect)
		unlock(queue)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("line scanning failed: %v", err)
	}
}

//////////////////////
// PUBLIC INTERFACE //
//////////////////////

// Config offers some adjustment screws to the user of lightd
type Config struct {
	// Host where lightd will run
	Host string
	// Port wich lightd will listen on
	Port int
	// DriverBinary is the name of the binary lightd will output rgb triples on.
	DriverBinary string
}

func cancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// Run starts lightd with the options specified in `cfg`,
// cancelling services when `ctx` is cancelled.
func Run(cfg *Config, ctx context.Context) error {
	queue, err := newEffectQueue(cfg.DriverBinary)
	if err != nil {
		log.Printf("Unable to hook up to lightd: %v", err)
		return err
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error listening: %v", err.Error())
		return err
	}

	defer util.Closer(lsn)
	log.Println("Listening on " + addr)

	for !cancelled(ctx) {
		if tcpLsn, ok := lsn.(*net.TCPListener); ok {
			if err := tcpLsn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
				log.Printf("Setting deadline failed: %v", err)
				return err
			}
		}

		conn, err := lsn.Accept()
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			continue
		}

		if err != nil {
			log.Printf("Error accepting: %v", err.Error())
			return err
		}

		go handleRequest(conn, queue)
	}

	return nil
}
