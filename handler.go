package gelf

import (
	"os"

	"github.com/inconshreveable/log15"
)

// Handler sends logs to Graylog in GELF.
type gelfHandler struct {
	gelfWriter *Writer
	host       string
}

// GelfHandler returns a handler that writes GELF messages to a service at gelfAddr. It is already wrapped
// in log15's CallerFileHandler and SyncHandler helpers. Its error is non-nil if there
// is a problem creating the GELF writer or determining our hostname.
// address is in teh format host:port.
//
//     log.GelfHandler("myhost:12201")
//
func GelfHandler(address string) (log15.Handler, error) {
	w, err := NewWriter(address)
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return log15.CallerFileHandler(log15.LazyHandler(log15.SyncHandler(gelfHandler{
		gelfWriter: w,
		host:       host,
	}))), nil
}

// Log forwards a log message to the specified receiver.
func (h gelfHandler) Log(r *log15.Record) error {

	// extract gelf-specific messages
	short, full := ShortAndFull(r.Msg)
	ctx := CtxToMap(r.Ctx)
	callerFile, callerLine := Caller(ctx)
	delete(ctx, "_caller")

	m := &Message{
		Version:  "1.1",
		Host:     h.host,
		Short:    short,
		Full:     full,
		TimeUnix: float64(r.Time.UnixNano()/1000000) / 1000., // seconds with millis from record
		//TimeUnix: float64(r.Time.UnixNano())/1e9 ,		// full timestamp
		Level: log15LevelsToSyslog[r.Lvl],
		File:  callerFile,
		Line:  callerLine,
		Extra: ctx,
	}

	return h.gelfWriter.WriteMessage(m)
}

// source: http://www.cisco.com/c/en/us/td/docs/security/asa/syslog-guide/syslogs/logsevp.html
var log15LevelsToSyslog = map[log15.Lvl]int32{
	log15.LvlCrit:  2,
	log15.LvlError: 3,
	log15.LvlWarn:  4,
	log15.LvlInfo:  6,
	log15.LvlDebug: 7,
}

var Must muster

func must(h log15.Handler, err error) log15.Handler {
	if err != nil {
		panic(err)
	}
	return h
}

type muster struct{}

// GelfHandler provides a panicking version
func (m muster) GelfHandler(address string) log15.Handler {
	return must(GelfHandler(address))
}
