<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-bold">Service List</h1>
      <UButton icon="i-heroicons-plus" @click="handleCreateService">Create Service</UButton>
    </div>

    <UTable :data="services" :columns="columns">
    </UTable>

    <ServiceForm ref="serviceFormRef" @serviceUpdated="fetchServices"/>

    <UModal v-model:open="isModalOpen" :title="modalTitle" :ui="{ footer: 'justify-end' }">
      <template #body>
        <p>{{ modalDescription }}</p>
      </template>

      <template #footer>
        <UButton color="neutral" @click="isModalOpen = false">No</UButton>
        <UButton :color="modalConfirmButtonColor" @click="handleModalConfirm">{{ modalConfirmButtonText }}</UButton>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import {h, ref, onMounted, resolveComponent} from 'vue';
import {clientKKSchedule} from '~/utils/api/client';
import ServiceForm from '~/components/ServiceForm.vue';
import type {PBRegisterService} from '~~/gen/kk_schedule/Base_pb';
import {ServiceList_InputSchema} from '~~/gen/kk_schedule/ServiceList_pb';
import {create} from '@bufbuild/protobuf';
import {ServiceDelete_InputSchema} from '~~/gen/kk_schedule/ServiceDelete_pb';
import {useToast} from '#imports';
import type {TableColumn} from '@nuxt/ui';

const UButton = resolveComponent('UButton')

const services = ref<PBRegisterService[]>([]);
const serviceFormRef = ref<InstanceType<typeof ServiceForm> | null>(null);
const toast = useToast();

const isModalOpen = ref(false);
const modalTitle = ref('');
const modalDescription = ref('');
const modalConfirmButtonText = ref('');
const modalConfirmButtonColor = ref('primary');
let modalConfirmAction: (() => Promise<void>) | null = null;

const columns: TableColumn<PBRegisterService>[] = [
  {accessorKey: 'ServiceName', header: 'Service Name'},
  {accessorKey: 'Target', header: 'Target'},
  {accessorKey: 'AuthToken', header: 'Auth Token'},
  {
    accessorKey: 'actions', header: 'Actions',
    cell: ({row}) => h('div', {class: 'flex flex-wrap gap-2'}, [
      h(UButton, {onClick: () => handleEditService(row.original)}, () => 'Edit'),
      h(UButton, {color: 'error', onClick: () => handleDeleteService(row.original)}, () => 'Delete'),
    ])
  }
];

const fetchServices = async () => {
  try {
    const request = create(ServiceList_InputSchema);
    const response = await clientKKSchedule.serviceList(request);
    services.value = response.ServiceList || [];
  } catch (error) {
    toast.add({title: 'Error fetching service list', description: String(error), color: 'error'});
  }
};

onMounted(async () => {
  await fetchServices();
});

const handleEditService = async (service: PBRegisterService) => {
  serviceFormRef.value?.open(service);
};

const handleCreateService = () => {
  serviceFormRef.value?.open();
};

const handleDeleteService = async (service: PBRegisterService) => {
  modalTitle.value = 'Warning';
  modalDescription.value = `Are you sure you want to delete service "${service.ServiceName}"?`;
  modalConfirmButtonText.value = 'Yes';
  modalConfirmButtonColor.value = 'error';
  modalConfirmAction = async () => {
    try {
      const request = create(ServiceDelete_InputSchema, {ServiceName: service.ServiceName});
      await clientKKSchedule.serviceDelete(request);
      await fetchServices();
      toast.add({title: 'Service deleted successfully', color: 'success'});
    } catch (error) {
      toast.add({title: 'Error deleting service', description: String(error), color: 'error'});
    }
  };
  isModalOpen.value = true;
};

const handleModalConfirm = async () => {
  if (modalConfirmAction) {
    await modalConfirmAction();
  }
  isModalOpen.value = false;
};
</script>