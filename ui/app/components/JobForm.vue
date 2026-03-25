<template>
  <UModal v-model:open="dialogVisible" :title="isEdit ? 'Edit Job' : 'Create Job'">
    <template #body>
      <UForm :state="form" class="space-y-4">
        <UFormField label="Description" name="Description">
          <UTextarea v-model="form.Description" autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>

        <UFormField label="Function Name" name="FuncName">
          <UTextarea v-model="form.FuncName" :disabled="isEdit" autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>

        <UFormField label="Service Name" name="ServiceName">
          <USelect v-model="form.ServiceName" :items="serviceOptions" class="w-full" placeholder="Select a service"></USelect>
        </UFormField>
      </UForm>
    </template>
    <template #footer>
      <div class="flex justify-end gap-2">
        <UButton color="neutral" @click="dialogVisible = false">Cancel</UButton>
        <UButton type="submit" @click="handleSave">Confirm</UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { clientKKSchedule } from '~/utils/api/client';
import { JobPut_InputSchema } from '~~/gen/kk_scheduler/JobPut_pb';
import {type PBJob} from '~~/gen/kk_scheduler/Job_pb';
import {create} from "@bufbuild/protobuf";
import { PBRegisterJobSchema, type PBRegisterJob } from '~~/gen/kk_scheduler/Base_pb';
import { ServiceList_InputSchema } from '~~/gen/kk_scheduler/ServiceList_pb';
import { useToast } from '#imports';

const dialogVisible = ref(false);
const isEdit = ref(false);
const toast = useToast();
const serviceOptions = ref<{label: string, value: string}[]>([]);

const form = reactive<PBRegisterJob>(create(PBRegisterJobSchema, {
  Description: '',
  FuncName: '',
  ServiceName: '',
}));

const emit = defineEmits(['jobUpdated']);

const open = async (job?: PBJob) => {
  dialogVisible.value = true;

  // Fetch service list for dropdown
  try {
    const response = await clientKKSchedule.serviceList(create(ServiceList_InputSchema, {}));
    serviceOptions.value = response.ServiceList.map(service => ({
      label: service.ServiceName,
      value: service.ServiceName,
    }));
  } catch (error) {
    toast.add({title: 'Error fetching service list', description: String(error), color: 'error'});
  }

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