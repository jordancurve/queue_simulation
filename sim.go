package main

// Example: `go run sim.go --arrival_rate 9 --departure_rate 10 --simulation_duration 200000`

import (
  "flag"
  "fmt"
  "log"
  "math/rand"
  "sort"
  "time"

  "github.com/fschuetz04/simgo"
)

type cust struct {
  id          int
  arrivalTime float64
}

var sizes []int
var times []float64
var queue []cust
var nextCustomerId int

func arrive(proc simgo.Process, rate float64) {
  for {
    proc.Wait(proc.Timeout(rand.ExpFloat64() / rate))
    sizes = append(sizes, len(queue))
    nextCustomerId++
    queue = append(queue, cust{nextCustomerId, proc.Now()})
  }
}

func depart(proc simgo.Process, rate float64) {
  for {
    proc.Wait(proc.Timeout(rand.ExpFloat64() / rate))
    sizes = append(sizes, len(queue))
    if len(queue) == 0 {
      continue
    }
    times = append(times, proc.Now()-queue[0].arrivalTime)
    queue = queue[1:]
  }
}

func main() {
  var arrivalRate, departureRate, simulationDur float64
  
  flag.Float64Var(&arrivalRate, "arrival_rate", 0.0, "arrival rate")
  flag.Float64Var(&departureRate, "departure_rate", 0.0, "departure rate")
  flag.Float64Var(&simulationDur, "simulation_duration", 0.0, "simulation duration")
  
  flag.Parse()
  
  if arrivalRate == 0.0 {
    log.Fatalf("Must specify --arrival_rate")
  }
  if departureRate == 0.0 {
    log.Fatalf("Must specify --departure_rate")
  }
  if simulationDur == 0.0 {
    log.Fatalf("Must specify --simulation_duration")
  }
  if arrivalRate >= departureRate {
    log.Fatalf("Arrival rate must be less than departure rate")
  }

  rand.Seed(time.Now().UnixNano())

  sim := simgo.Simulation{}
  sim.ProcessReflect(arrive, arrivalRate)
  sim.ProcessReflect(depart, departureRate)

  sim.RunUntil(simulationDur)

  sum := 0.0
  for _, s := range sizes {
    sum += float64(s)
  }
  fmt.Printf("average number of customers in system: %.2f (queueing theory prediction: %.2f)\n", sum/float64(len(sizes)), arrivalRate/(departureRate-arrivalRate))

  sum = 0.0
  for _, s := range times {
    sum += s
  }
  fmt.Printf("average time in system: %.2f (queueing theory prediction: %.2f)\n", sum/float64(len(times)), 1/(departureRate-arrivalRate))

  sort.Float64s(times)
  fmt.Printf("99%%ile time in system: %.2f\n", times[int(0.99*float64(len(times)))])
}
