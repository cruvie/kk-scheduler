<template>
  <UModal v-model:open="dialogVisible" title="Set Job Spec">
    <template #body>
      <UForm :state="form" class="space-y-4">
        <UFormField label="Service Name" name="ServiceName">
          <UTextarea v-model="form.ServiceName" disabled autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>

        <UFormField label="Function Name" name="FuncName">
          <UTextarea v-model="form.FuncName" disabled autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>

        <UFormField label="Spec" name="Spec">
          <UTextarea v-model="form.Spec" :rows="5" autoresize class="w-full"></UTextarea>
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
import {ref, reactive} from 'vue';
import {clientKKSchedule} from '~/utils/api/client';
import {JobSetSpec_InputSchema} from '~~/gen/kk_scheduler/JobSetSpec_pb';
import {create} from "@bufbuild/protobuf";
import type {PBJob} from '~~/gen/kk_scheduler/Job_pb';
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