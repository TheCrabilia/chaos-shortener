package chaos

import (
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type FailureType int

const (
	FailureTypeNone FailureType = iota
	FailureTypeLatency
	FailureTypeError
	FailureTypeConnDrop
	FailureTypeOutage
)

type Injector struct {
	latencyRate  atomic.Int32
	errorRate    atomic.Int32
	connDropRate atomic.Int32
	outageRate   atomic.Int32

	maxLatency time.Duration
}

func NewInjector() *Injector {
	i := &Injector{
		maxLatency: time.Second * 2,
	}

	return i
}

func (i *Injector) SetLatencyRate(rate float64) {
	i.latencyRate.Store(int32(rate * 100))
}

func (i *Injector) GetLatencyRate() float64 {
	return float64(i.latencyRate.Load()) / 100
}

func (i *Injector) SetErrorRate(rate float64) {
	i.errorRate.Store(int32(rate * 100))
}

func (i *Injector) GetErrorRate() float64 {
	return float64(i.errorRate.Load()) / 100
}

func (i *Injector) SetConnDropRate(rate float64) {
	i.connDropRate.Store(int32(rate * 100))
}

func (i *Injector) GetConnDropRate() float64 {
	return float64(i.connDropRate.Load()) / 100
}

func (i *Injector) SetOutageRate(rate float64) {
	i.outageRate.Store(int32(rate * 100))
}

func (i *Injector) GetOutageRate() float64 {
	return float64(i.outageRate.Load()) / 100
}

func (i *Injector) ShouldInject() (bool, FailureType) {
	probabilities := []float64{
		i.GetLatencyRate(),
		i.GetErrorRate(),
		i.GetConnDropRate(),
		i.GetOutageRate(),
		1 - i.failureRate(),
	}

	failureTypes := []FailureType{
		FailureTypeLatency,
		FailureTypeError,
		FailureTypeConnDrop,
		FailureTypeOutage,
		FailureTypeNone,
	}

	failure := selectWeightedRandom(failureTypes, probabilities)

	if failure == FailureTypeNone {
		return false, failure
	}
	return true, failure
}

func (i *Injector) failureRate() float64 {
	return i.GetLatencyRate() + i.GetErrorRate() + i.GetConnDropRate() + i.GetOutageRate()
}

func (i *Injector) InjectLatency() {
	l := time.Duration(rand.Int63n(int64(i.maxLatency)))
	time.Sleep(l)
}

func (i *Injector) InjectError(w http.ResponseWriter) {
	http.Error(w, "chaos error", http.StatusInternalServerError)
}

func (i *Injector) InjectConnDrop(w http.ResponseWriter) {
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		slog.Error("failed to hijack connection", "error", err)
	}

	conn.(*net.TCPConn).SetLinger(0)
	conn.Close()
}

func (i *Injector) InjectOutage() {
	time.Sleep(time.Second * 30)
}
