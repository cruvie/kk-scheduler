package schedule_test

import (
	"testing"

	"github.com/cruvie/kk-schedule/server/kk_schedule"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func getClient(t *testing.T) kk_schedule.KKScheduleClient {
	var err error
	conn, err = grpc.NewClient("127.0.0.1:8666",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return kk_schedule.NewKKScheduleClient(conn)
}

func down() {
	conn.Close() // nolint
}

const testAuthToken = "sdgoisdglodshlghlshlghdlskg"

var testJob = func() *kk_schedule.PBRegisterJob {
	j := &kk_schedule.PBRegisterJob{}
	j.SetDescription("test job")
	j.SetServiceName("test-service")
	j.SetFuncName("test-func")
	return j
}()

var testService = func() *kk_schedule.PBRegisterService {
	s := &kk_schedule.PBRegisterService{}
	s.SetServiceName("test-service")
	s.SetTarget("127.0.0.1:8000")
	s.SetAuthToken(testAuthToken)
	return s
}()

func TestJobList(t *testing.T) {
	defer down()
	input := &kk_schedule.JobList_Input{}
	jobs, err := getClient(t).JobList(t.Context(), input)
	assert.NoError(t, err)
	for _, job := range jobs.GetJobList() {
		t.Log(job)
		t.Log(job.GetPrev().AsTime(), job.GetNext().AsTime())
	}
}

func TestJobGet(t *testing.T) {
	defer down()
	input := &kk_schedule.JobGet_Input{}
	input.SetServiceName(testJob.GetServiceName())
	input.SetFuncName(testJob.GetFuncName())
	job, err := getClient(t).JobGet(t.Context(), input)
	assert.NoError(t, err)
	t.Log(job.GetJob())
	t.Log(job.GetJob().GetPrev().AsTime(), job.GetJob().GetNext().AsTime())
}

func TestJobSetSpec(t *testing.T) {
	defer down()
	input := &kk_schedule.JobSetSpec_Input{}
	input.SetServiceName(testJob.GetServiceName())
	input.SetFuncName(testJob.GetFuncName())
	input.SetSpec("* * * * *")
	resp, err := getClient(t).JobSetSpec(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobEnable(t *testing.T) {
	defer down()
	input := &kk_schedule.JobEnable_Input{}
	input.SetServiceName(testJob.GetServiceName())
	input.SetFuncName(testJob.GetFuncName())
	resp, err := getClient(t).JobEnable(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobDisable(t *testing.T) {
	defer down()
	input := &kk_schedule.JobDisable_Input{}
	input.SetServiceName(testJob.GetServiceName())
	input.SetFuncName(testJob.GetFuncName())
	resp, err := getClient(t).JobDisable(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestJobPut(t *testing.T) {
	defer down()
	input := &kk_schedule.JobPut_Input{}
	input.SetJob(testJob)
	resp, err := getClient(t).JobPut(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceList(t *testing.T) {
	defer down()
	input := &kk_schedule.ServiceList_Input{}
	resp, err := getClient(t).ServiceList(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServicePut(t *testing.T) {
	defer down()
	input := &kk_schedule.ServicePut_Input{}
	input.SetService(testService)
	resp, err := getClient(t).ServicePut(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceGet(t *testing.T) {
	defer down()
	input := &kk_schedule.ServiceGet_Input{}
	input.SetServiceName(testService.GetServiceName())
	resp, err := getClient(t).ServiceGet(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}

func TestServiceDelete(t *testing.T) {
	defer down()
	input := &kk_schedule.ServiceDelete_Input{}
	input.SetServiceName(testService.GetServiceName())
	resp, err := getClient(t).ServiceDelete(t.Context(), input)
	assert.NoError(t, err)
	t.Log(resp)
}
