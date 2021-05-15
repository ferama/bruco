package python

import (
	"fmt"
	"log"
)

// Pool ...
type Pool struct {
	// filename(sessionName) -> python instance
	pythonMap map[string]*Python

	workerQueue chan *Python
}

// GetPoolInstance ...
func NewPool(size int) *Pool {
	pool := &Pool{
		pythonMap:   make(map[string]*Python),
		workerQueue: make(chan *Python, size),
	}

	for i := 0; i < size; i++ {
		name := fmt.Sprintf("worker%d", i)
		pool.createPythonInstance(name)

	}
	return pool
}

func (p *Pool) GetWorker() *Python {
	return <-p.workerQueue
}

// createPythonInstance ...
func (p *Pool) createPythonInstance(name string) *Python {
	python := NewPython(name, p.workerQueue)
	p.pythonMap[name] = python

	log.Println("New python instance started: " + name)
	return python
}

// Destroy ...
func (p *Pool) Destroy() {
	for key, python := range p.pythonMap {
		python.Kill()
		python = nil
		delete(p.pythonMap, key)
	}
}
