package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"golang.org/x/term"
)

const (
	defaultNumLayers   = 8
	defaultInitialLen  = 15
	defaultAngleDeg      = 40.0
	defaultLeafLen       = 4
	defaultWaitTime      = 0.0
	defaultBranchChars   = "~;:= "
	defaultLeafChars     = "&%#@ "
	defaultWindowWidth   = 80
	defaultWindowHeight  = 25
	version              = "1.0.0"
	desc                 = "GoBonsai procedurally generates ASCII art trees in your terminal."
)

var rng = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(os.Getpid())))

type Options struct {
	NumLayers   int
	InitialLen  int
	AngleMean   float64
	LeafLen     int
	Instant     bool
	WaitTime    float64
	BranchChars string
	LeafChars   string
	Type        int
	UserSetType bool
	FixedWindow bool
	WindowWidth int
	WindowHeight int
}

func defaultWindowSize() (int, int) {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			return min(w, defaultWindowWidth), min(h, defaultWindowHeight)
		}
	}
	return defaultWindowWidth, defaultWindowHeight
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func defaultOptions() *Options {
	width, height := defaultWindowSize()

	return &Options{
		NumLayers:   defaultNumLayers,
		InitialLen:  defaultInitialLen,
		AngleMean:   defaultAngleDeg * math.Pi / 180,
		LeafLen:     defaultLeafLen,
		WaitTime:    defaultWaitTime,
		BranchChars: defaultBranchChars,
		LeafChars:   defaultLeafChars,
		Type:        rng.IntN(4),
		WindowWidth: width,
		WindowHeight: height,
	}
}

const noValue = "\x00"

type arg struct {
	name  string
	value string
}

func parseArgs(args []string) []arg {
	var parsed []arg

	for i, x := range args {
		if len(x) == 0 || x[0] != '-' {
			continue
		}

		value := argValue(args, i)

		if len(x) > 2 && x[1] != '-' {
			for _, c := range x[1:] {
				parsed = append(parsed, arg{"-" + string(c), value})
			}
			continue
		}

		parsed = append(parsed, arg{x, value})
	}

	return parsed
}

func argValue(args []string, inx int) string {
	if inx+1 >= len(args) {
		return noValue
	}
	if v := args[inx+1]; len(v) > 0 && v[0] != '-' {
		return v
	}
	return noValue
}

var shortOptions = map[string]string{
	"-h": "--help",
	"-i": "--instant",
	"-c": "--branch-chars",
	"-C": "--leaf-chars",
	"-w": "--wait",
	"-x": "--width",
	"-y": "--height",
	"-t": "--type",
	"-s": "--seed",
	"-S": "--start-len",
	"-L": "--leaf-len",
	"-l": "--layers",
	"-a": "--angle",
	"-f": "--fixed-window",
}

func getOptions(rawArgs []string) *Options {
	options := defaultOptions()

	for _, a := range parseArgs(rawArgs) {
		name := a.name
		if len(name) > 1 && name[1] != '-' {
			long, ok := shortOptions[name]
			if !ok {
				showInvalid(name)
			}
			name = long
		}
		setOption(options, name, a.value)
	}

	return options
}

func setOption(options *Options, name, value string) {
	switch name {
	case "--layers":
		options.NumLayers = parseInt(name, value)
	case "--start-len":
		options.InitialLen = parseInt(name, value)
	case "--angle":
		options.AngleMean = float64(parseInt(name, value)) * math.Pi / 180
	case "--leaf-len":
		options.LeafLen = parseInt(name, value)
	case "--instant":
		options.Instant = true
	case "--wait":
		options.WaitTime = parseFloat(name, value)
	case "--branch-chars":
		requireValue(name, value)
		options.BranchChars = stripQuotes(value)
	case "--leaf-chars":
		requireValue(name, value)
		options.LeafChars = stripQuotes(value)
	case "--type":
		options.Type = parseInt(name, value)
		options.UserSetType = true
	case "--width":
		options.WindowWidth = parseInt(name, value)
	case "--height":
		options.WindowHeight = parseInt(name, value)
	case "--help":
		showHelp()
	case "--version":
		showVersion()
	case "--seed":
		setSeed(options, parseInt(name, value))
	case "--fixed-window":
		options.FixedWindow = true
	default:
		showInvalid(name)
	}
}

func parseInt(name, value string) int {
	requireValue(name, value)
	v, err := strconv.Atoi(stripQuotes(value))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid value for %s: %s\n", name, value)
		os.Exit(1)
	}
	return v
}

func parseFloat(name, value string) float64 {
	requireValue(name, value)
	v, err := strconv.ParseFloat(stripQuotes(value), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid value for %s: %s\n", name, value)
		os.Exit(1)
	}
	return v
}

func requireValue(name, value string) {
	if value == noValue {
		fmt.Fprintf(os.Stderr, "Missing value for %s\n", name)
		os.Exit(1)
	}
}

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if q := s[0]; (q == '\'' || q == '"') && s[len(s)-1] == q {
		return s[1 : len(s)-1]
	}
	return s
}

func setSeed(options *Options, seed int) {
	rng = rand.New(rand.NewPCG(uint64(int64(seed)), 0))

	if !options.UserSetType {
		options.Type = rng.IntN(4)
	}
}

func showHelp() {
	fmt.Printf("USEAGE gobonsai [OPTION]...\n\n")
	fmt.Println(desc)
	fmt.Println(optionDescs())
	os.Exit(0)
}

func optionDescs() string {
	return fmt.Sprintf(`
OPTIONS:
    -h, --help            display help
        --version         display version

    -s, --seed            seed for the random number generator

    -i, --instant         instant mode: display finished tree immediately
    -w, --wait            time delay between drawing characters when not in instant mode [default %g]

    -c, --branch-chars    string of chars randomly chosen for branches [default "%s"]
    -C, --leaf-chars      string of chars randomly chosen for leaves [default "%s"]

    -x, --width           maximum width of the tree [default %d]
    -y, --height          maximum height of the tree [default %d]

    -t, --type            tree type: integer between 0 and 3 inclusive [default random]
    -S, --start-len       length of the root branch [default %d]
    -L, --leaf-len        length of each leaf [default %d]
    -l, --layers          number of branch layers: more => more branches [default %d]
    -a, --angle           mean angle of branches to their parent, in degrees; more => more arched trees [default %g]

    -f, --fixed-window    do not allow window height to increase when tree grows off screen
`,
		defaultWaitTime, defaultBranchChars, defaultLeafChars,
		defaultWindowWidth, defaultWindowHeight,
		defaultInitialLen, defaultLeafLen, defaultNumLayers, defaultAngleDeg)
}

func showVersion() {
	fmt.Printf("GoBonsai version %s\n", version)
	os.Exit(0)
}

func showInvalid(name string) {
	fmt.Fprintf(os.Stderr, "Invalid option: %s. Use gobonsai --help for useage.\n", name)
	os.Exit(1)
}