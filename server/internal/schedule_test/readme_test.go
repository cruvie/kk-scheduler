package schedule_test

import (
	"testing"

	"github.com/cruvie/kk-schedule/server/kk_schedule"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestForREADME(t *testing.T) {
	// create a client for kk-schedule
	conn, err := grpc.NewClient("127.0.0.1:8666",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	defer conn.Close() //nolint
	if err != nil {
		t.Fatal(err)
	}
	client := kk_schedule.NewKKScheduleClient(conn)

	myServiceName := "my-service"
	testJob := func() *kk_schedule.PBRegisterJob {
		j := &kk_schedule.PBRegisterJob{}
		j.SetDescription("test job")
		j.SetServiceName(myServiceName)
		j.SetFuncName("Func1")
		return j
	}()
	testService := func() *kk_schedule.PBRegisterService {
		s := &kk_schedule.PBRegisterService{}
		s.SetServiceName(myServiceName)
		s.SetTarget("127.0.0.1:8000")
		return s
	}()
	{
		// put the running service info to kk-schedule
		input := &kk_schedule.ServicePut_Input{}
		input.SetService(testService)
		resp, err := client.ServicePut(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// put a job to kk-schedule with the service name
		input := &kk_schedule.JobPut_Input{}
		input.SetJob(testJob)
		resp, err := client.JobPut(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// set job spec
		input := &kk_schedule.JobSetSpec_Input{}
		input.SetServiceName(testJob.GetServiceName())
		input.SetFuncName(testJob.GetFuncName())
		input.SetSpec("* * * * *")
		resp, err := client.JobSetSpec(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
	{
		// enable job to be triggered with the spec
		input := &kk_schedule.JobEnable_Input{}
		input.SetServiceName(testJob.GetServiceName())
		input.SetFuncName(testJob.GetFuncName())
		resp, err := client.JobEnable(t.Context(), input)
		assert.NoError(t, err)
		t.Log(resp)
	}
}
