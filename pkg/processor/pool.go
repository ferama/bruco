package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type EventCallback func(event *Response)

// Pool ...
type Pool struct {
	pythonMap   map[string]*Python
	wrapperPath string
	lambdaPath  string

	availableWorkers chan *Python
}

// GetPoolInstance ...
func NewPool(cfg *ProcessorConf) *Pool {
	data, _ := pythonWrapper.ReadFile("wrapper/wrapper.py")
	file, err := ioutil.TempFile(os.TempDir(), "python-wrapper-")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = file.Write(data); err != nil {
		log.Fatal("[PROCESSOR]  failed to write to temporary file", err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("[PROCESSOR] unpacked wrapper to %s", file.Name())

	pool := &Pool{
		pythonMap:        make(map[string]*Python),
		availableWorkers: make(chan *Python, cfg.Workers),
		wrapperPath:      file.Name(),
		lambdaPath:       cfg.LambdaPath,
	}
	for i := 0; i < cfg.Workers; i++ {
		name := fmt.Sprintf("worker%d", i)
		pool.createPythonInstance(name)

	}
	log.Printf("[PROCESSOR] allocated %d workers", cfg.Workers)
	return pool
}

func (p *Pool) HandleEventAsync(data []byte, callback EventCallback) error {
	python := <-p.availableWorkers
	err := python.handleEvent(data)
	go func() {
		response := <-python.eventResponse
		if callback != nil {
			callback(&response)
		}
	}()
	return err
}

func (p *Pool) HandleEvent(data []byte) (*Response, error) {
	python := <-p.availableWorkers
	err := python.handleEvent(data)
	if err != nil {
		return nil, err
	}
	response := <-python.eventResponse
	return &response, nil
}

// createPythonInstance ...
func (p *Pool) createPythonInstance(name string) *Python {
	python := NewPython(name, p.availableWorkers, p.wrapperPath, p.lambdaPath)
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
	os.Remove(p.wrapperPath)
}
