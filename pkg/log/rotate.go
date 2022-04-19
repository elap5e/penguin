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
	"os"
	"path"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

type RotateWriter struct {
	mu   sync.Mutex
	base string
	time time.Time
	file *os.File
}

func New(base string) *RotateWriter {
	return &RotateWriter{base: base, time: time.Now()}
}

func (w *RotateWriter) Rotate() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil || time.Now().Day() != w.time.Day() {
		if w.file != nil {
			err = w.file.Close()
			w.file = nil
			if err != nil {
				return
			}
		}
		name := path.Join(w.base, time.Now().Format("2006-01-02")+".log")
		w.file, err = os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	return
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	if err := w.Rotate(); err != nil {
		return 0, err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return colorable.NewNonColorable(w.file).Write(output)
}
