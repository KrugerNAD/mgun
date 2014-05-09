package work

import (
	"math"
	"time"
)

type TaskReport struct {
	Hostname string
	Port int
	Concurrency int
	TotalRequests int
	CompleteRequests int
	FailedRequests int
	Availability float64
	RequestsPerSecond float64
	TotalTransferred int64
	TotalTime float64
}

func NewStageReport(hit *Hit) *StageReport {
	return new(StageReport).create(hit)
}

type StageReport struct {
	count int
	startTime time.Time
	endTime time.Time
	MinTime float64
	MaxTime float64
	AvgTime float64
	CompleteRequests int
	FailedRequests int
	RequestsPerSecond float64
	TotalTransferred int64
	TotalTime float64
}

func (this *StageReport) create(hit *Hit) *StageReport {
	timeRequest := this.getDiffSeconds(hit)
	this.MinTime = timeRequest
	this.MaxTime = timeRequest
	this.TotalTime = timeRequest
	this.updateCount()
	this.updateTotalTransferred(hit)
	this.checkResponseStatusCode(hit)
	this.startTime = hit.StartTime
	this.endTime = hit.EndTime
	return this
}

func (this *StageReport) getDiffSeconds(hit *Hit) float64 {
	return hit.EndTime.Sub(hit.StartTime).Seconds()
}

func (this *StageReport) checkResponseStatusCode(hit *Hit) {
	shot := hit.Shot
	if hit.Request != nil && hit.Response != nil {
		statusCode := hit.Response.StatusCode
		if this.inArray(statusCode, shot.FailedStatusCodes) {
			this.FailedRequests++
		} else if this.inArray(statusCode, shot.GetSuccessStatusCodes()) {
			this.CompleteRequests++
		} else {
			this.FailedRequests++
		}
	} else {
		this.FailedRequests++
	}
}

func (this *StageReport) inArray(a int, array []int) bool {
	for _, b := range array {
		if a == b {
			return true
		}
	}
	return false
}

func (this *StageReport) updateCount() {
	this.count++
}

func (this *StageReport) updateTotalTransferred(hit *Hit) {
	if hit.Response != nil {
		this.TotalTransferred += hit.Response.ContentLength
	}
}

func (this *StageReport) Update(hit *Hit) *StageReport {
	timeRequest := this.getDiffSeconds(hit)
	this.MinTime = math.Min(this.MinTime, timeRequest)
	this.MaxTime = math.Max(this.MaxTime, timeRequest)
	this.TotalTime += timeRequest
	this.updateCount()
	this.updateTotalTransferred(hit)
	this.checkResponseStatusCode(hit)
	this.endTime = hit.EndTime
	return this
}

func (this *StageReport) GetAvgTime() float64 {
	return (this.MinTime + this.MaxTime) / 2
}

func (this *StageReport) GetAvailability() float64 {
	return float64(this.CompleteRequests) * 100 / float64(this.count)
}

func (this *StageReport) GetRequestPerSeconds() float64 {
	if this.endTime.Equal(this.startTime) {
		return float64(this.count)
	} else {
		return float64(this.count) / this.endTime.Sub(this.startTime).Seconds()
	}
}
