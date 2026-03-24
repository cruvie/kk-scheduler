<template>
  <UModal v-model:open="dialogVisible" title="Set Job Spec">
    <template #body>
      <UForm :state="form">
        <UFormField label="Service Name" name="ServiceName">
          <UInput v-model="form.ServiceName" disabled/>
        </UFormField>

        <UFormField label="Function Name" name="FuncName">
          <UInput v-model="form.FuncName" disabled/>
        </UFormField>

        <UFormField label="Spec" name="Spec">
          <UTextarea v-model="form.Spec" :rows="5"/>
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
import {ref, reactive} from 'vue';
import {clientKKSchedule} from '~/utils/api/client';
import {JobSetSpec_InputSchema} from '~~/gen/kk_schedule/JobSetSpec_pb';
import {create} from "@bufbuild/protobuf";
import type {PBJob} from '~~/gen/kk_schedule/Job_pb';
import { useToast } from '#imports';

const dialogVisible = ref(false);
const toast = useToast();

const form = reactive({
  ServiceName: '',
  FuncName: '',
  Spec: '',
});

const emit = defineEmits(['jobUpdated']);

const open = (job: PBJob) => {
  dialogVisible.value = true;
  form.ServiceName = job.ServiceName;
  form.FuncName = job.FuncName;
  form.Spec = job.Spec;
};

const handleSave = async () => {
  try {
    const request = create(JobSetSpec_InputSchema, {
      serviceName: form.ServiceName,
      funcName: form.FuncName,
      spec: form.Spec,
    });
    await clientKKSchedule.jobSetSpec(request);
    toast.add({title: 'Job spec updated successfully!', color: 'success'});
    dialogVisible.value = false;
    emit('jobUpdated');
  } catch (error) {
    toast.add({title: 'Error setting job spec', description: String(error), color: 'error'});
  }
};

defineExpose({
  open,
});
</script>