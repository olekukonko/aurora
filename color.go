package aurora

import "github.com/fatih/color"

// Add these new types and constants near the top of the file
type ColorOption func(*color.Color)

// Update the Value struct to support multiple attributes
type Value struct {
	value string
	attrs []color.Attribute
}

// Add color combination support
func (v Value) Colorize(attrs ...color.Attribute) Value {
	return Value{v.value, append(v.attrs, attrs...)}
}

// Update String() method to handle multiple attributes
func (v Value) String() string {
	if len(v.attrs) == 0 {
		return v.value
	}
	c := color.New(v.attrs...)
	return c.Sprint(v.value)
}

// Color constructors (foreground colors)
func Black(s string) Value   { return Value{s, []color.Attribute{color.FgBlack}} }
func Red(s string) Value     { return Value{s, []color.Attribute{color.FgRed}} }
func Green(s string) Value   { return Value{s, []color.Attribute{color.FgGreen}} }
func Yellow(s string) Value  { return Value{s, []color.Attribute{color.FgYellow}} }
func Blue(s string) Value    { return Value{s, []color.Attribute{color.FgBlue}} }
func Magenta(s string) Value { return Value{s, []color.Attribute{color.FgMagenta}} }
func Cyan(s string) Value    { return Value{s, []color.Attribute{color.FgCyan}} }
func White(s string) Value   { return Value{s, []color.Attribute{color.FgWhite}} }

// Bright foreground colors
func BrightBlack(s string) Value   { return Value{s, []color.Attribute{color.FgHiBlack}} }
func BrightRed(s string) Value     { return Value{s, []color.Attribute{color.FgHiRed}} }
func BrightGreen(s string) Value   { return Value{s, []color.Attribute{color.FgHiGreen}} }
func BrightYellow(s string) Value  { return Value{s, []color.Attribute{color.FgHiYellow}} }
func BrightBlue(s string) Value    { return Value{s, []color.Attribute{color.FgHiBlue}} }
func BrightMagenta(s string) Value { return Value{s, []color.Attribute{color.FgHiMagenta}} }
func BrightCyan(s string) Value    { return Value{s, []color.Attribute{color.FgHiCyan}} }
func BrightWhite(s string) Value   { return Value{s, []color.Attribute{color.FgHiWhite}} }

// Background colors
func BgBlack(s string) Value   { return Value{s, []color.Attribute{color.BgBlack}} }
func BgRed(s string) Value     { return Value{s, []color.Attribute{color.BgRed}} }
func BgGreen(s string) Value   { return Value{s, []color.Attribute{color.BgGreen}} }
func BgYellow(s string) Value  { return Value{s, []color.Attribute{color.BgYellow}} }
func BgBlue(s string) Value    { return Value{s, []color.Attribute{color.BgBlue}} }
func BgMagenta(s string) Value { return Value{s, []color.Attribute{color.BgMagenta}} }
func BgCyan(s string) Value    { return Value{s, []color.Attribute{color.BgCyan}} }
func BgWhite(s string) Value   { return Value{s, []color.Attribute{color.BgWhite}} }

// Bright background colors
func BgBrightBlack(s string) Value   { return Value{s, []color.Attribute{color.BgHiBlack}} }
func BgBrightRed(s string) Value     { return Value{s, []color.Attribute{color.BgHiRed}} }
func BgBrightGreen(s string) Value   { return Value{s, []color.Attribute{color.BgHiGreen}} }
func BgBrightYellow(s string) Value  { return Value{s, []color.Attribute{color.BgHiYellow}} }
func BgBrightBlue(s string) Value    { return Value{s, []color.Attribute{color.BgHiBlue}} }
func BgBrightMagenta(s string) Value { return Value{s, []color.Attribute{color.BgHiMagenta}} }
func BgBrightCyan(s string) Value    { return Value{s, []color.Attribute{color.BgHiCyan}} }
func BgBrightWhite(s string) Value   { return Value{s, []color.Attribute{color.BgHiWhite}} }

// Text styles
func Bold(s string) Value      { return Value{s, []color.Attribute{color.Bold}} }
func Faint(s string) Value     { return Value{s, []color.Attribute{color.Faint}} }
func Italic(s string) Value    { return Value{s, []color.Attribute{color.Italic}} }
func Underline(s string) Value { return Value{s, []color.Attribute{color.Underline}} }
func Blink(s string) Value     { return Value{s, []color.Attribute{color.BlinkSlow}} }
func BlinkFast(s string) Value { return Value{s, []color.Attribute{color.BlinkRapid}} }
func Reverse(s string) Value   { return Value{s, []color.Attribute{color.ReverseVideo}} }
func Conceal(s string) Value   { return Value{s, []color.Attribute{color.Concealed}} }
func Strike(s string) Value    { return Value{s, []color.Attribute{color.CrossedOut}} }

// Chainable color methods
func (v Value) Black() Value           { return v.Colorize(color.FgBlack) }
func (v Value) Red() Value             { return v.Colorize(color.FgRed) }
func (v Value) Green() Value           { return v.Colorize(color.FgGreen) }
func (v Value) Yellow() Value          { return v.Colorize(color.FgYellow) }
func (v Value) Blue() Value            { return v.Colorize(color.FgBlue) }
func (v Value) Magenta() Value         { return v.Colorize(color.FgMagenta) }
func (v Value) Cyan() Value            { return v.Colorize(color.FgCyan) }
func (v Value) White() Value           { return v.Colorize(color.FgWhite) }
func (v Value) BrightBlack() Value     { return v.Colorize(color.FgHiBlack) }
func (v Value) BrightRed() Value       { return v.Colorize(color.FgHiRed) }
func (v Value) BrightGreen() Value     { return v.Colorize(color.FgHiGreen) }
func (v Value) BrightYellow() Value    { return v.Colorize(color.FgHiYellow) }
func (v Value) BrightBlue() Value      { return v.Colorize(color.FgHiBlue) }
func (v Value) BrightMagenta() Value   { return v.Colorize(color.FgHiMagenta) }
func (v Value) BrightCyan() Value      { return v.Colorize(color.FgHiCyan) }
func (v Value) BrightWhite() Value     { return v.Colorize(color.FgHiWhite) }
func (v Value) BgBlack() Value         { return v.Colorize(color.BgBlack) }
func (v Value) BgRed() Value           { return v.Colorize(color.BgRed) }
func (v Value) BgGreen() Value         { return v.Colorize(color.BgGreen) }
func (v Value) BgYellow() Value        { return v.Colorize(color.BgYellow) }
func (v Value) BgBlue() Value          { return v.Colorize(color.BgBlue) }
func (v Value) BgMagenta() Value       { return v.Colorize(color.BgMagenta) }
func (v Value) BgCyan() Value          { return v.Colorize(color.BgCyan) }
func (v Value) BgWhite() Value         { return v.Colorize(color.BgWhite) }
func (v Value) BgBrightBlack() Value   { return v.Colorize(color.BgHiBlack) }
func (v Value) BgBrightRed() Value     { return v.Colorize(color.BgHiRed) }
func (v Value) BgBrightGreen() Value   { return v.Colorize(color.BgHiGreen) }
func (v Value) BgBrightYellow() Value  { return v.Colorize(color.BgHiYellow) }
func (v Value) BgBrightBlue() Value    { return v.Colorize(color.BgHiBlue) }
func (v Value) BgBrightMagenta() Value { return v.Colorize(color.BgHiMagenta) }
func (v Value) BgBrightCyan() Value    { return v.Colorize(color.BgHiCyan) }
func (v Value) BgBrightWhite() Value   { return v.Colorize(color.BgHiWhite) }
func (v Value) Bold() Value            { return v.Colorize(color.Bold) }
func (v Value) Faint() Value           { return v.Colorize(color.Faint) }
func (v Value) Italic() Value          { return v.Colorize(color.Italic) }
func (v Value) Underline() Value       { return v.Colorize(color.Underline) }
func (v Value) Blink() Value           { return v.Colorize(color.BlinkSlow) }
func (v Value) BlinkFast() Value       { return v.Colorize(color.BlinkRapid) }
func (v Value) Reverse() Value         { return v.Colorize(color.ReverseVideo) }
func (v Value) Conceal() Value         { return v.Colorize(color.Concealed) }
func (v Value) Strike() Value          { return v.Colorize(color.CrossedOut) }
