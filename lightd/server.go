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
)

// TODO: Cleanup, this is a bit older code
//       and the effect syntax is awful.

type EffectQueue struct {
	StdInPipe io.Writer
	Blocked   chan bool
}

func (q *EffectQueue) Push(e Effect) {
	c := e.ComposeEffect()

	for color := range c {
		colorValue := fmt.Sprintf("%d %d %d\n", color.R, color.G, color.B)
		q.StdInPipe.Write([]byte(colorValue))
	}

	q.Blocked <- false
}

func NewEffectQueue(driverBinary string) (*EffectQueue, error) {
	cmd := exec.Command(driverBinary, "cat")
	stdinpipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	blocked := make(chan bool, 1)
	blocked <- false

	return &EffectQueue{
		Blocked:   blocked,
		StdInPipe: stdinpipe,
	}, cmd.Start()
}

type SimpleColor struct {
	R, G, B uint8
}

type Effect interface {
	ComposeEffect() chan SimpleColor
}

type Properties struct {
	Delay  time.Duration
	Color  SimpleColor
	Repeat int
}

type FadeEffect struct {
	Properties
}

type FlashEffect struct {
	Properties
}

func (color *SimpleColor) ComposeEffect() chan SimpleColor {
	c := make(chan SimpleColor, 1)
	c <- SimpleColor{color.R, color.G, color.B}
	close(c)
	return c
}

func (effect *FlashEffect) ComposeEffect() chan SimpleColor {
	c := make(chan SimpleColor, 1)
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
			c <- SimpleColor{0, 0, 0}
			time.Sleep(effect.Delay)
			effect.Repeat--
		}
		close(c)
	}()
	return c
}

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

func (effect *FadeEffect) ComposeEffect() chan SimpleColor {
	c := make(chan SimpleColor, 1)

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

			for i := 0; i < int(max); i += 1 {
				c <- SimpleColor{uint8((i * r) / 100), uint8((i * g) / 100), uint8((i * b) / 100)}
				time.Sleep(effect.Delay)
			}

			for i := int(max - 1); i >= 0; i -= 1 {
				c <- SimpleColor{uint8((i * r) / 100), uint8((i * g) / 100), uint8((i * b) / 100)}
				time.Sleep(effect.Delay)
			}
			effect.Repeat--
		}
		close(c)
	}()
	return c
}

type BlendEffect struct {
	StartColor SimpleColor
	EndColor   SimpleColor
	Duration   time.Duration
}

func (effect *BlendEffect) ComposeEffect() chan SimpleColor {
	c := make(chan SimpleColor, 1)
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

			c <- SimpleColor{uint8(sr), uint8(sg), uint8(sb)}
			time.Sleep(time.Duration(1/N*1000) * time.Millisecond)
		}

		close(c)
	}()

	return c
}

type FireEffect struct {
	Properties
}

func (effect *FireEffect) ComposeEffect() chan SimpleColor {
	fn := func(t, n, jitter int, fac float64) uint8 {
		j := float64(rand.Int()%(jitter<<1) - jitter)
		f := float64(t - n>>1)
		return uint8(fac * (float64(-255.0/262144.0)*f*f + 255 + j))
	}

	c := make(chan SimpleColor, 1)
	go func() {
		defer close(c)

		for t := 0; t < effect.Repeat; t++ {
			c <- SimpleColor{
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
	REGEX_COLOR      = regexp.MustCompile("(c|color|)\\{(\\d{1,3}),(\\d{1,3}),(\\d{1,3})\\}")
	REGEX_PROPERTIES = regexp.MustCompile(".*?\\{(.*)\\|(.*)\\|(.*)\\}")
)

func parseColor(s string) (*SimpleColor, error) {
	matches := REGEX_COLOR.FindStringSubmatch(s)
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

	return &SimpleColor{triple[0], triple[1], triple[2]}, nil
}

func parseBlendEffect(s string) (*BlendEffect, error) {
	// Same regex as for properties (by chance)
	matches := REGEX_PROPERTIES.FindStringSubmatch(s)
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

	return &BlendEffect{*srcColor, *dstColor, duration}, nil
}

func parseProperties(s string) (*Properties, error) {
	matches := REGEX_PROPERTIES.FindStringSubmatch(s)
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

	return &Properties{duration, *color, repeatCnt}, nil
}

func parseFadeEffect(s string) (*FadeEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "fade"))
	if err != nil {
		return nil, err
	}

	return &FadeEffect{*props}, nil
}

func parseFlashEffect(s string) (*FlashEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "flash"))
	if err != nil {
		return nil, err
	}

	return &FlashEffect{*props}, nil
}

func parseFireEffect(s string) (*FireEffect, error) {
	props, err := parseProperties(strings.TrimPrefix(s, "fire"))
	if err != nil {
		return nil, err
	}

	return &FireEffect{*props}, nil
}

func parseEffect(s string) (Effect, error) {
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

func handleRequest(conn net.Conn, queue *EffectQueue) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Bad line:
		if len(line) == 0 {
			continue
		}

		switch line {
		case "!lock":
			log.Printf("Locking queue...")
			<-queue.Blocked

			log.Printf("Lock acquired.")
			if _, err := conn.Write([]byte("OK\n")); err != nil {
				log.Printf("Failed to answer lock response: %v", err)
			}

			return
		case "!unlock":
			log.Printf("Unlocking queue...")
			queue.Blocked <- false

			if _, err := conn.Write([]byte("OK\n")); err != nil {
				log.Printf("Failed to answer unlock response: %v", err)
			}
			return
		}

		effect, err := parseEffect(line)
		if err != nil {
			log.Printf("Unable to process effect: %v", err)
			return
		}

		<-queue.Blocked
		queue.Push(effect)
	}
}

type Config struct {
	Host         string
	Port         int
	DriverBinary string
}

func Run(cfg *Config) error {
	queue, err := NewEffectQueue(cfg.DriverBinary)
	if err != nil {
		log.Printf("Unable to hook up to lightd: %v", err)
		return err
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error listening:", err.Error())
		return err
	}

	defer lsn.Close()
	log.Println("Listening on " + addr)

	for {
		conn, err := lsn.Accept()
		if err != nil {
			log.Printf("Error accepting: ", err.Error())
			return err
		}

		go handleRequest(conn, queue)
	}
}
