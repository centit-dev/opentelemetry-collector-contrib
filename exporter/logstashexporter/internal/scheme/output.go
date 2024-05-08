package scheme // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logstashexporter/internal/scheme

import "go.opentelemetry.io/collector/pdata/plog"

type LogOutput struct {
	Attributes         map[string]any `json:"attributes"`
	Timestamp          int64          `json:"timestamp"`
	TraceID            string         `json:"traceID"`
	SpanID             string         `json:"spanID"`
	SeverityNumber     int16          `json:"severityNumber"`
	SeverityText       string         `json:"severityText"`
	Body               string         `json:"body"`
	ResourceAttributes map[string]any `json:"resource_attributes"`
}

type LogsOutput struct {
	Logs []*LogOutput `json:"logs"`
}

func NewHttpOutputFromPlogs(ld plog.Logs) *LogsOutput {
	var result LogsOutput
	result.Logs = make([]*LogOutput, 0)

	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rl := rls.At(i)

		ills := rl.ScopeLogs()
		res := rl.Resource()
		for j := 0; j < ills.Len(); j++ {
			ils := ills.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				lr.Attributes().AsRaw()

				lsLog := &LogOutput{
					Attributes:         lr.Attributes().AsRaw(),
					ResourceAttributes: res.Attributes().AsRaw(),
					TraceID:            lr.TraceID().String(),
					SpanID:             lr.TraceID().String(),
					SeverityNumber:     int16(lr.SeverityNumber()),
					SeverityText:       lr.SeverityText(),
					Body:               lr.Body().AsString(),
					Timestamp:          lr.Timestamp().AsTime().Unix(),
				}

				result.Logs = append(result.Logs, lsLog)
			}
		}
	}

	return &result
}
