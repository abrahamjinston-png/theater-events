package logger

import "log/slog"

// The filterable dimensions of a log line. One constructor per field is what
// keeps the key spelled the same everywhere: a CloudWatch filter that misses
// half its lines to a typo reports nothing wrong, so the compiler has to.

// ChatID identifies the user. It is the key users are stored under, so it is
// what a per-person query filters on.
func ChatID(id int64) slog.Attr { return slog.Int64("chat_id", id) }

// Theater names the source a line concerns.
func Theater(name string) slog.Attr { return slog.String("theater", name) }
