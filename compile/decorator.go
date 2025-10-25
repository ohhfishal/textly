package compile

const (
	Reset = "\033[0m"

	Red    = "\033[0;31m"
	Black  = "\033[0;30m"
	Green  = "\033[0;32m"
	Yellow = "\033[0;33m"
	Blue   = "\033[0;34m"
	Purple = "\033[0;35m"
	Cyan   = "\033[0;36m"
	White  = "\033[0;37m"

	UnderlinedRed    = "\033[4;31m"
	UnderlinedBlack  = "\033[4;30m"
	UnderlinedGreen  = "\033[4;32m"
	UnderlinedYellow = "\033[4;33m"
	UnderlinedBlue   = "\033[4;34m"
	UnderlinedPurple = "\033[4;35m"
	UnderlinedCyan   = "\033[4;36m"
	UnderlinedWhite  = "\033[4;37m"

	BackgroundRed    = "\033[0;41m"
	BackgroundBlack  = "\033[0;40m"
	BackgroundGreen  = "\033[0;42m"
	BackgroundYellow = "\033[0;43m"
	BackgroundBlue   = "\033[0;44m"
	BackgroundPurple = "\033[0;45m"
	BackgroundCyan   = "\033[0;46m"
	BackgroundWhite  = "\033[0;47m"

	IntenseRed    = "\033[0;91m"
	IntenseBlack  = "\033[0;90m"
	IntenseGreen  = "\033[0;92m"
	IntenseYellow = "\033[0;93m"
	IntenseBlue   = "\033[0;94m"
	IntensePurple = "\033[0;95m"
	IntenseCyan   = "\033[0;96m"
	IntenseWhite  = "\033[0;97m"
)

var colorMap = map[string]string{
	"red":    Red,
	"blue":   Blue,
	"black":  Black,
	"green":  Green,
	"yellow": Yellow,
	"purple": Purple,
	"cyan":   Cyan,
	"white":  White,
	"grey":   IntenseBlack,

	"_red":    UnderlinedRed,
	"_blue":   UnderlinedBlue,
	"_black":  UnderlinedBlack,
	"_green":  UnderlinedGreen,
	"_yellow": UnderlinedYellow,
	"_purple": UnderlinedPurple,
	"_cyan":   UnderlinedCyan,
	"_white":  UnderlinedWhite,

	"+red":    BackgroundRed,
	"+blue":   BackgroundBlue,
	"+black":  BackgroundBlack,
	"+green":  BackgroundGreen,
	"+yellow": BackgroundYellow,
	"+purple": BackgroundPurple,
	"+cyan":   BackgroundCyan,
	"+white":  BackgroundWhite,

	"!red":    IntenseRed,
	"!blue":   IntenseBlue,
	"!black":  IntenseBlack,
	"!green":  IntenseGreen,
	"!yellow": IntenseYellow,
	"!purple": IntensePurple,
	"!cyan":   IntenseCyan,
	"!white":  IntenseWhite,
}
