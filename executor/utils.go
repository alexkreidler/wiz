package executor

import "fmt"

// Data access utilities
func checkProcessorExists(p ProcessorExecutor, processorID string) error {
	_, ok := p.base[processorID]
	if !ok {
		return fmt.Errorf("processor %s not found", processorID)
	}
	return nil
}

func chuckRunExists(p ProcessorExecutor, id string, runID string) error {
	if err := checkProcessorExists(p, id); err != nil {
		return err
	}
	_, ok := p.runMap[id][runID]
	if !ok {
		return fmt.Errorf("run %s not found", runID)
	}
	return nil
}

//getRun returns a pointer because the run, processor, and data need to be updated
func getRun(p ProcessorExecutor, id string, runID string) (*runProcessor, error) {
	err := chuckRunExists(p, id, runID)
	if err != nil {
		return nil, err
	}

	return p.runMap[id][runID], nil
}
