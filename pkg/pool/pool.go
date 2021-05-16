package pool

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

// Pool ...
type Pool struct {
	// filename(sessionName) -> python instance
	pythonMap   map[string]*Python
	wrapperPath string
	lambdaPath  string

	availableWorkers chan *Python
}

// GetPoolInstance ...
func NewPool(size int, lambdaPath string) *Pool {
	data, _ := pythonWrapper.ReadFile("wrapper/wrapper.py")
	file, err := ioutil.TempFile(os.TempDir(), "python-wrapper-")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = file.Write(data); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("unpacked wrapper to %s", file.Name())

	pool := &Pool{
		pythonMap:        make(map[string]*Python),
		availableWorkers: make(chan *Python, size),
		wrapperPath:      file.Name(),
		lambdaPath:       lambdaPath,
	}

	for i := 0; i < size; i++ {
		name := fmt.Sprintf("worker%d", i)
		pool.createPythonInstance(name)

	}
	return pool
}

func (p *Pool) HandleEventAsync(data []byte) error {
	python := <-p.availableWorkers
	return python.handleEvent(data)
}

func (p *Pool) HandleEvent(data []byte) (string, error) {
	python := <-p.availableWorkers
	err := python.handleEvent(data)
	response := <-python.eventResponse
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}
	return response.Response, nil
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
