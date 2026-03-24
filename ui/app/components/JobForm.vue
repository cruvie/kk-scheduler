<template>
  <UModal v-model:open="dialogVisible" :title="isEdit ? 'Edit Job' : 'Create Job'">
    <template #body>
      <UForm :state="form" >
        <UFormField label="Description" name="Description">
          <UInput v-model="form.Description" />
        </UFormField>

        <UFormField label="Function Name" name="FuncName">
          <UInput v-model="form.FuncName" :disabled="isEdit" />
        </UFormField>

        <UFormField label="Service Name" name="ServiceName">
          <UInput v-model="form.ServiceName" />
        </UFormField>
      </UForm>
    </template>
    <template #footer>
      <UButton color="neutral" @click="dialogVisible = false">Cancel</UButton>
      <UButton type="submit" @click="handleSave">Confirm</UButton>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { clientKKSchedule } from '~/utils/api/client';
import { JobPut_InputSchema } from '~~/gen/kk_schedule/JobPut_pb';
import {type PBJob} from '~~/gen/kk_schedule/Job_pb';
import {create} from "@bufbuild/protobuf";
import { PBRegisterJobSchema, type PBRegisterJob } from '~~/gen/kk_schedule/Base_pb';
import { useToast } from '#imports';

const dialogVisible = ref(false);
const isEdit = ref(false);
const toast = useToast();

const form = reactive<PBRegisterJob>(create(PBRegisterJobSchema, {
  Description: '',
  FuncName: '',
  ServiceName: '',
}));

const emit = defineEmits(['jobUpdated']);

const open = (job?: PBJob) => {
  dialogVisible.value = true;

  if (job) {
    isEdit.value = true;
    form.Description = job.Description;
    form.FuncName = job.FuncName;
    form.ServiceName = job.ServiceName;
  } else {
    isEdit.value = false;
    Object.assign(form, create(PBRegisterJobSchema, {
      Description: '',
      FuncName: '',
      ServiceName: '',
    }));
  }
};

const handleSave = async () => {
  try {
    if (isEdit.value) {
      const putRequest = create(JobPut_InputSchema, {
        Job: create(PBRegisterJobSchema, {
          Description: form.Description,
          ServiceName: form.ServiceName,
          FuncName: form.FuncName,
        }),
      });
      await clientKKSchedule.jobPut(putRequest);
      toast.add({title: 'Job updated successfully!', color: 'success'});
    } else {
      const request = create(JobPut_InputSchema, {
        Job: create(PBRegisterJobSchema, {
          Description: form.Description,
          ServiceName: form.ServiceName,
          FuncName: form.FuncName,
        }),
      });
      await clientKKSchedule.jobPut(request);
      toast.add({title: 'Job created successfully!', color: 'success'});
    }
    dialogVisible.value = false;
    emit('jobUpdated');
  } catch (error) {
    toast.add({title: 'Error saving job', description: String(error), color: 'error'});
  }
};

defineExpose({
  open,
});
</script>