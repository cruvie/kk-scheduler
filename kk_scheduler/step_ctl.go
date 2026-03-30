package kk_scheduler

// StepCtl provides control for step execution
type StepCtl struct {
	addLog func(message string)
}

// Log adds a log message to the execution record
func (c *StepCtl) Log(message string) {
	c.addLog(message)
}
