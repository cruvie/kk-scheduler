package kk_scheduler

// StepCtl provides control for step execution
type StepCtl struct {
	addLog func(message string)
}

// AddLog adds a log message to the execution record
func (c *StepCtl) AddLog(message string) {
	if c.addLog != nil {
		c.addLog(message)
	}
}
