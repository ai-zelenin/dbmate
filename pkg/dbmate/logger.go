package dbmate

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
)

var Log = logger{
	w: os.Stdout,
}

type logger struct {
	w io.Writer
}

func (l *logger) Fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func (l *logger) Fprintln(w io.Writer, args ...interface{}) {
	_, _ = fmt.Fprintln(w, args...)
}

func (l *logger) Fprint(w io.Writer, args ...interface{}) {
	_, _ = fmt.Fprint(w, args...)
}
func (l *logger) FPrintColor(w io.Writer, a color.Attribute, format string, args ...interface{}) {
	_, _ = color.New(a).Fprintln(w, fmt.Sprintf(format, args...))
}
func (l *logger) PrintLnColor(a color.Attribute, format string, args ...interface{}) {
	_, _ = color.New(a).Fprintln(l.w, fmt.Sprintf(format, args...))
}
func (l *logger) PrintColor(a color.Attribute, format string, args ...interface{}) {
	_, _ = color.New(a).Fprint(l.w, fmt.Sprintf(format, args...))
}
