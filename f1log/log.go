// F1GopherLib - Copyright (C) 2022 f1gopher
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package f1log

import (
	"fmt"
	"io"
	"log"
)

type F1GopherLibLog struct {
	output *log.Logger
}

func CreateLog() *F1GopherLibLog {
	return &F1GopherLibLog{
		output: nil,
	}
}

func (l *F1GopherLibLog) SetLogOutput(w io.Writer) {
	l.output = log.New(w, "", log.LstdFlags|log.Lmicroseconds)
}

func (l *F1GopherLibLog) Info(msg string) {
	l.msg("INF", msg)
}

func (l *F1GopherLibLog) Infof(format string, a ...any) {
	l.msgf("INF", format, a...)
}

func (l *F1GopherLibLog) Warn(msg string) {
	l.msg("WRN", msg)
}

func (l *F1GopherLibLog) Warnf(format string, a ...any) {
	l.msgf("WRN", format, a...)
}

func (l *F1GopherLibLog) Error(msg string) {
	l.msg("ERR", msg)
}

func (l *F1GopherLibLog) Errorf(format string, a ...any) {
	l.msgf("ERR", format, a...)
}

func (l *F1GopherLibLog) Fatal(msg string) {
	l.msg("FTL", msg)
}

func (l *F1GopherLibLog) Fatalf(format string, a ...any) {
	l.msgf("FTL", format, a...)
}

func (l *F1GopherLibLog) msgf(prefix string, format string, a ...any) {
	if l.output == nil {
		return
	}

	l.msg(prefix, fmt.Sprintf(format, a...))
}

func (l *F1GopherLibLog) msg(prefix string, msg string) {
	if l.output == nil {
		return
	}

	// TODO - work out calling filename like log normally does?
	l.output.Print(prefix + ": " + msg)
}
