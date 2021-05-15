package python

import (
	"fmt"
)

// Pool ...
type Pool struct {
	// filename(sessionName) -> python instance
	pythonMap map[string]*Python

	availableWorkers chan *Python
}

// GetPoolInstance ...
func NewPool(size int) *Pool {
	pool := &Pool{
		pythonMap:        make(map[string]*Python),
		availableWorkers: make(chan *Python, size),
	}

	for i := 0; i < size; i++ {
		name := fmt.Sprintf("worker%d", i)
		pool.createPythonInstance(name)

	}
	return pool
}

func (p *Pool) HandleEvent(data []byte) error {
	python := <-p.availableWorkers
	return python.handleEvent(data)
}

// createPythonInstance ...
func (p *Pool) createPythonInstance(name string) *Python {
	python := NewPython(name, p.availableWorkers)
	p.pythonMap[name] = python

	// log.Println("New python instance started: " + name)
	return python
}

// Destroy ...
func (p *Pool) Destroy() {
	for key, python := range p.pythonMap {
		python.kill()
		python = nil
		delete(p.pythonMap, key)
	}
}
