package neural

import (
	"sync"
	"math"
)

type weight struct {
	channel chan *message
	weight  float64
	value 	float64
	messageReceived bool
}

func makeWeight(ch chan *message, seed float64) *weight {
	weight := new(weight)
	weight.channel = ch
	weight.weight = seed
	return weight
}

type neuron struct {
	weights []*weight
	receivedValues []float64
	editMutex sync.RWMutex
	sendChannel  chan *message
}

func (n *neuron) listen(weight *weight, notifier chan bool) {
	notifier <- true
	for {
		select {
		case msg := <- weight.channel:
			if msg.command == _DISPOSE {
				go n.send(msg)
				return
			} else {
				n.editMutex.Lock()
				weight.value = msg.value
				weight.messageReceived = true
				n.editMutex.Unlock()
				n.checkDone()
			}
		}
	}
}

func (n *neuron) send(msg *message) {
	n.sendChannel <- msg
}

func (n *neuron) checkDone() {
	defer n.editMutex.Unlock()

	n.editMutex.Lock()

	for _, weight := range n.weights {
		if !weight.messageReceived {
			return
		}
	}

	result := n.calculate()

	message := makeMessage(_CALC, result, -1)

	for _, weight := range n.weights {
		weight.messageReceived = false
	}

	go n.send(message)
}

func (n *neuron) calculate() float64 {
	var sum float64
	sum += 1.0
	
	for _, weight := range n.weights {
		sum += weight.value * weight.weight
	}

	result := 1.0 / (1.0 + math.Exp(-sum))

	return result
}

func makeNeuron(output chan *message, weights ...*weight) *neuron {
	neuron := new(neuron)
	neuron.weights = weights
	neuron.receivedValues = make([]float64, len(weights))
	notifiers := make([]chan bool, len(weights))
	neuron.sendChannel = output
	for i, weight := range weights {
		notifier := make(chan bool, 1)
		go neuron.listen(weight, notifier)
		notifiers[i] = notifier
	}
	for _, notifier := range notifiers {
		<- notifier
	}
	return neuron
}