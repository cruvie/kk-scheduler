package kk_scheduler

func (x *PBRegisterService) Check() error {
	if x.GetTarget() == "" {
		return ErrTargetEmpty
	}
	if x.GetServiceName() == "" {
		return ErrServiceNameEmpty
	}
	return nil
}

func (x *PBRegisterJob) Check() error {
	if x.GetFuncName() == "" {
		return ErrFuncNameEmpty
	}
	return nil
}
