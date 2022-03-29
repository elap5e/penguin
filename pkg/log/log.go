// Copyright 2022 Elapse and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"log"

	"github.com/mattn/go-colorable"
)

var (
	logger log.Logger

	Panicf  = logger.Panicf
	Panicln = logger.Panicln

	Fatalf  = logger.Fatalf
	Fatalln = logger.Fatalln

	Print   = logger.Print
	Printf  = logger.Printf
	Println = logger.Println
)

func init() {
	logger.SetFlags(log.Ltime | log.Lmicroseconds)
	logger.SetOutput(colorable.NewColorableStdout())
}

func Trace(format string, v ...any) {
	// logger.Println("[TRACE]", fmt.Sprintf(format, v...))
}

func Debug(format string, v ...any) {
	logger.Println("\x1b[37m[DEBUG]\x1b[0m", fmt.Sprintf(format, v...))
}

func Chat(format string, v ...any) {
	logger.Println("\x1b[37;1m[CHAT]", fmt.Sprintf(format, v...)+"\x1b[0m")
}

func Info(format string, v ...any) {
	logger.Println("\x1b[36m[INFO]", fmt.Sprintf(format, v...)+"\x1b[0m")
}

func Warn(format string, v ...any) {
	logger.Println("\x1b[33m[WARN]", fmt.Sprintf(format, v...)+"\x1b[0m")
}

func Error(format string, v ...any) {
	logger.Println("\x1b[31m[ERROR]", fmt.Sprintf(format, v...)+"\x1b[0m")
}

func Fatal(format string, v ...any) {
	logger.Println("\x1b[35m[FATAL]\x1b[0m", fmt.Sprintf(format, v...))
}

func Panic(format string, v ...any) {
	logger.Println("\x1b[34m[PANIC]\x1b[0m", fmt.Sprintf(format, v...))
}
