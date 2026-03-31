<template>
  <UModal v-model:open="dialogVisible" :title="isEdit ? 'Edit Service' : 'Create Service'" :ui="{ footer: 'justify-end' }">
    <template #body>
      <UForm :state="form" class="space-y-4">
        <UFormField label="Service Name" name="ServiceName">
          <UTextarea v-model="form.ServiceName" :disabled="isEdit" autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>
        <UFormField label="Target" name="Target">
          <UTextarea v-model="form.Target" autoresize :rows="1" class="w-full"></UTextarea>
        </UFormField>
        <UFormField label="Auth Token" name="AuthToken">
          <UTextarea v-model="form.AuthToken" autoresize :rows="3" class="w-full"></UTextarea>
        </UFormField>
      </UForm>
    </template>
    <template #footer>
      <UButton color="neutral" @click="dialogVisible = false">Cancel</UButton>
      <UButton type="submit" @click="submitForm">Confirm</UButton>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { create } from "@bufbuild/protobuf";
import { clientKKSchedule } from '~/utils/api/client';
import { ServicePut_InputSchema } from '~~/gen/kk_scheduler/Service_pb';
import type { PBRegisterService } from '~~/gen/kk_scheduler/Base_pb';
import { PBRegisterServiceSchema } from '~~/gen/kk_scheduler/Base_pb';
import { useToast } from '#imports';

const dialogVisible = ref(false);
const isEdit = ref(false);
const toast = useToast();
const form = reactive<PBRegisterService>(create(PBRegisterServiceSchema));

const emit = defineEmits(['serviceUpdated']);

const open = (service?: PBRegisterService) => {
  dialogVisible.value = true;
  if (service) {
    isEdit.value = true;
    Object.assign(form, service);
  } else {
    isEdit.value = false;
    // Reset form for new service
    form.ServiceName = '';
    form.Target = '';
    form.AuthToken = '';
  }
};

const submitForm = async () => {
  try {
    const request = create(ServicePut_InputSchema, {
      Service: form,
    });
    await clientKKSchedule.servicePut(request);
    dialogVisible.value = false;
    toast.add({title: 'Service saved successfully!', color: 'success'});
    emit('serviceUpdated');
  } catch (error) {
    toast.add({title: 'Error submitting service form', description: String(error), color: 'error'});
  }
};

defineExpose({ open });
</script>