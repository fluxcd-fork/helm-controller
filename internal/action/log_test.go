/*
Copyright 2022 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr"
)

func TestLogBuffer_Log(t *testing.T) {
	nowTS = stubNowTS

	tests := []struct {
		name      string
		size      int
		fill      []string
		wantCount int
		want      string
	}{
		{name: "log", size: 2, fill: []string{"a", "b", "c"}, wantCount: 3, want: fmt.Sprintf("%[1]s b\n%[1]s c", stubNowTS().Format(time.RFC3339Nano))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var count int
			l := NewLogBuffer(func(format string, v ...interface{}) {
				count++
			}, tt.size)
			for _, v := range tt.fill {
				l.Log("%s", v)
			}
			if count != tt.wantCount {
				t.Errorf("Inner Log() called %v times, want %v", count, tt.wantCount)
			}
			if got := l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogBuffer_Len(t *testing.T) {
	tests := []struct {
		name string
		size int
		fill []string
		want int
	}{
		{name: "empty buffer", fill: []string{}, want: 0},
		{name: "filled buffer", size: 2, fill: []string{"a", "b"}, want: 2},
		{name: "half full buffer", size: 4, fill: []string{"a", "b"}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogBuffer(NewDebugLog(logr.Discard()), tt.size)
			for _, v := range tt.fill {
				l.Log("%s", v)
			}
			if got := l.Len(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogBuffer_Reset(t *testing.T) {
	bufferSize := 10
	l := NewLogBuffer(NewDebugLog(logr.Discard()), bufferSize)

	if got := l.buffer.Len(); got != bufferSize {
		t.Errorf("Len() = %v, want %v", got, bufferSize)
	}

	for _, v := range []string{"a", "b", "c"} {
		l.Log("%s", v)
	}

	if got := l.String(); got == "" {
		t.Errorf("String() = empty")
	}

	l.Reset()

	if got := l.buffer.Len(); got != bufferSize {
		t.Errorf("Len() = %v after Reset(), want %v", got, bufferSize)
	}
	if got := l.String(); got != "" {
		t.Errorf("String() != empty after Reset()")
	}
}

func TestLogBuffer_String(t *testing.T) {
	nowTS = stubNowTS

	tests := []struct {
		name string
		size int
		fill []string
		want string
	}{
		{name: "empty buffer", fill: []string{}, want: ""},
		{name: "filled buffer", size: 2, fill: []string{"a", "b", "c"}, want: fmt.Sprintf("%[1]s b\n%[1]s c", stubNowTS().Format(time.RFC3339Nano))},
		{name: "duplicate buffer items", fill: []string{"b", "b"}, want: fmt.Sprintf("%[1]s b\n%[1]s b", stubNowTS().Format(time.RFC3339Nano))},
		{name: "duplicate buffer items", fill: []string{"b", "b", "b"}, want: fmt.Sprintf("%[1]s b\n%[1]s b (1 duplicate line omitted)", stubNowTS().Format(time.RFC3339Nano))},
		{name: "duplicate buffer items", fill: []string{"b", "b", "b", "b"}, want: fmt.Sprintf("%[1]s b\n%[1]s b (2 duplicate lines omitted)", stubNowTS().Format(time.RFC3339Nano))},
		{name: "duplicate buffer items", fill: []string{"a", "b", "b", "b", "c", "c"}, want: fmt.Sprintf("%[1]s a\n%[1]s b\n%[1]s b (1 duplicate line omitted)\n%[1]s c\n%[1]s c", stubNowTS().Format(time.RFC3339Nano))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLogBuffer(NewDebugLog(logr.Discard()), tt.size)
			for _, v := range tt.fill {
				l.Log("%s", v)
			}
			if got := l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// stubNowTS returns a fixed time for testing purposes.
func stubNowTS() time.Time {
	return time.Date(2016, 2, 18, 12, 24, 5, 12345600, time.UTC)
}
