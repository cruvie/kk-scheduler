package task_executor

// StepCtl provides global control for all steps
// All steps share one StepCtl instance
type StepCtl struct {
	addLog func(message string)
	stop   func()
}

// AddLog adds a log message to the execution record
func (c *StepCtl) AddLog(message string) {
	if c.addLog != nil {
		c.addLog(message)
	}
}

// Stop signals the task to stop after current step completes
func (c *StepCtl) Stop() {
	if c.stop != nil {
		c.stop()
	}
}
