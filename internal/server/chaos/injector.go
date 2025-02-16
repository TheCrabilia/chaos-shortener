package chaos

import (
	"fmt"
	"math/rand"
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
	delay := time.Duration(rand.Int63n(int64(i.maxLatency)))

	timer := time.NewTimer(delay)
	defer timer.Stop()

	<-timer.C
}

func (i *Injector) InjectError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "injected error"}`)
}

func (i *Injector) InjectConnDrop(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "injected connection drop"}`)
}

func (i *Injector) InjectOutage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "injected service unavailable"}`)
}
