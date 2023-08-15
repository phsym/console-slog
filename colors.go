package console

type color string

const (
	reset color = "\x1b[0m"
	bold  color = "\x1b[1m"

	colorTimestamp  color = "\x1b[90m"
	colorSource     color = bold + "\x1b[90m"
	colorErrorValue color = "\x1b[91m"
	colorMessage    color = "\x1b[97m"
	colorAttrKey    color = "\x1b[36m"
	colorAttrValue  color = "\x1b[90m"
	colorLevelError color = bold + colorErrorValue
	colorLevelWarn  color = "\x1b[93m"
	colorLevelInfo  color = "\x1b[92m"
	colorLevelDebug color = "\x1b[95m"
)
