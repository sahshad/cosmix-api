package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	colorReset  = "\x1b[0m"
	colorBlue   = "\x1b[34m"
	colorYellow = "\x1b[33m"
	colorRed    = "\x1b[31m"
	colorPurple = "\x1b[35m"
	colorCyan   = "\x1b[36m"
	colorGreen  = "\x1b[32m"
	colorGray   = "\x1b[90m"
)

// fieldOrder controls the display order of known fields; unknown fields are dropped.
var fieldOrder = []string{
	"service",
	"method",
	"routing_key",
	"full_method",
	"request_id",
	"correlation_id",
	"duration",
	"error",
	"panic",
}

var fieldLabels = map[string]string{
	"service":        "Service",
	"method":         "Method",
	"routing_key":    "Routing Key",
	"full_method":    "Full Method",
	"request_id":     "Request ID",
	"correlation_id": "Correlation",
	"duration":       "Duration",
	"error":          "Error",
	"panic":          "Panic",
}

var fieldColors = map[string]string{
	"method":         colorCyan,
	"routing_key":    colorCyan,
	"duration":       colorGreen,
	"request_id":     colorGray,
	"correlation_id": colorGray,
}

var prettyBufferPool = buffer.NewPool()

type prettyEncoder struct {
	*zapcore.MapObjectEncoder
}

func newPrettyEncoder(zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return &prettyEncoder{MapObjectEncoder: zapcore.NewMapObjectEncoder()}, nil
}

func (enc *prettyEncoder) Clone() zapcore.Encoder {
	clone := zapcore.NewMapObjectEncoder()

	for k, v := range enc.Fields {
		clone.Fields[k] = v
	}

	return &prettyEncoder{MapObjectEncoder: clone}
}

func (enc *prettyEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := zapcore.NewMapObjectEncoder()

	for k, v := range enc.Fields {
		line.Fields[k] = v
	}

	for _, f := range fields {
		f.AddTo(line)
	}

	buf := prettyBufferPool.Get()

	levelColor, levelText := levelDisplay(entry.Level)

	buf.AppendString(entry.Time.Format("15:04:05.000"))
	buf.AppendString("  ")
	buf.AppendString(levelColor)
	buf.AppendString(fmt.Sprintf("%-5s", levelText))
	buf.AppendString(colorReset)
	buf.AppendString(" [")
	buf.AppendString(colorPurple)
	buf.AppendString(loggerDisplayName(entry.LoggerName))
	buf.AppendString(colorReset)
	buf.AppendString("] ")
	buf.AppendString(entry.Message)
	buf.AppendString("\n")

	rows := buildRows(entry, line.Fields)

	colWidth := 0
	for _, r := range rows {
		if len(r.label) > colWidth {
			colWidth = len(r.label)
		}
	}
	colWidth += 2

	for _, r := range rows {
		buf.AppendString("    ")
		buf.AppendString(r.label)
		buf.AppendString(strings.Repeat(" ", colWidth-len(r.label)))
		buf.AppendString(": ")

		if r.color != "" {
			buf.AppendString(r.color)
			buf.AppendString(r.value)
			buf.AppendString(colorReset)
		} else {
			buf.AppendString(r.value)
		}

		buf.AppendString("\n")
	}

	buf.AppendString("\n")

	return buf, nil
}

type row struct {
	label string
	value string
	color string
}

func buildRows(entry zapcore.Entry, values map[string]interface{}) []row {
	rows := make([]row, 0, len(fieldOrder)+1)

	for _, key := range fieldOrder {
		v, ok := values[key]
		if !ok {
			continue
		}

		rows = append(rows, row{
			label: fieldLabels[key],
			value: formatValue(v),
			color: fieldColors[key],
		})
	}

	if entry.Level >= zapcore.ErrorLevel && entry.Caller.Defined {
		rows = append(rows, row{
			label: "Source",
			value: entry.Caller.TrimmedPath(),
			color: colorGray,
		})
	}

	return rows
}

func levelDisplay(lvl zapcore.Level) (color, text string) {
	switch lvl {
	case zapcore.DebugLevel:
		return colorGray, "DEBUG"
	case zapcore.InfoLevel:
		return colorBlue, "INFO"
	case zapcore.WarnLevel:
		return colorYellow, "WARN"
	default:
		return colorRed, strings.ToUpper(lvl.String())
	}
}

func loggerDisplayName(name string) string {
	if name == "" {
		return "app"
	}

	return name
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case time.Duration:
		return val.String()
	case error:
		return val.Error()
	case fmt.Stringer:
		return val.String()
	case string:
		return val
	default:
		return fmt.Sprint(val)
	}
}
