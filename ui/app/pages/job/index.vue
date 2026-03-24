<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-bold">Job List</h1>
      <UButton icon="i-heroicons-plus" @click="handleCreateJob">Create Job</UButton>
    </div>

    <UTable :data="jobs" :columns="columns">

    </UTable>

    <JobForm ref="jobFormRef" @jobUpdated="fetchJobs"/>
    <JobSetSpecForm ref="jobSetSpecFormRef" @jobUpdated="fetchJobs"/>

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
import {JobList_InputSchema} from '~~/gen/kk_schedule/JobList_pb';
import type {PBJob} from '~~/gen/kk_schedule/Job_pb';
import {create} from "@bufbuild/protobuf";
import {JobEnable_InputSchema} from "~~/gen/kk_schedule/JobEnable_pb";
import {JobDisable_InputSchema} from "~~/gen/kk_schedule/JobDisable_pb";
import {JobDelete_InputSchema} from "~~/gen/kk_schedule/JobDelete_pb";
import {JobTrigger_InputSchema} from "~~/gen/kk_schedule/JobTrigger_pb";
import JobForm from '~/components/JobForm.vue';
import JobSetSpecForm from '~/components/JobSetSpecForm.vue';
import {useToast} from '#imports';
import type {TableColumn} from '@nuxt/ui';

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const jobs = ref<PBJob[]>([]);
const toast = useToast();

const isModalOpen = ref(false);
const modalTitle = ref('');
const modalDescription = ref('');
const modalConfirmButtonText = ref('');
const modalConfirmButtonColor = ref('primary');
let modalConfirmAction: (() => Promise<void>) | null = null;

const columns: TableColumn<PBJob>[] = [
  {accessorKey: 'EntryID', header: 'Entry ID'},
  {accessorKey: 'Description', header: 'Description'},
  {accessorKey: 'FuncName', header: 'Function Name'},
  {accessorKey: 'Spec', header: 'Spec'},
  {accessorKey: 'ServiceName', header: 'Service Name'},
  {
    accessorKey: 'Next',
    header: 'Next',
    cell: ({row}) =>
        row.original.Next?.seconds ? new Date(Number(row.original.Next.seconds) * 1000).toLocaleString() : ''
  },
  {
    accessorKey: 'Prev',
    header: 'Prev',
    cell: ({row}) =>
        row.original.Prev?.seconds ? new Date(Number(row.original.Prev.seconds) * 1000).toLocaleString() : ''
  },
  {
    accessorKey: 'Enabled',
    header: 'Status',
    cell: ({row}) => {
      const color = row.original.Enabled ? 'success' : 'error';
      const text = row.original.Enabled ? 'Enabled' : 'Disabled';
      return h(UBadge, {color}, () => text);
    }
  },
  {
    accessorKey: 'actions', header: 'Actions',
    cell: ({row}) => h('div', {class: 'flex flex-wrap gap-2 w-[200px]'}, [
      h(UButton, {disabled: row.original.Enabled, onClick: () => handleEnableJob(row.original)}, () => 'Enable'),
      h(UButton, {disabled: !row.original.Enabled, onClick: () => handleDisableJob(row.original)}, () => 'Disable'),
      h(UButton, {onClick: () => handleEditJob(row.original)}, () => 'Edit'),
      h(UButton, {onClick: () => handleSetSpecJob(row.original)}, () => 'Set Spec'),
      h(UButton, {color: 'primary', onClick: () => handleTriggerJob(row.original)}, () => 'Trigger'),
      h(UButton, {color: 'error', onClick: () => handleDeleteJob(row.original)}, () => 'Delete'),
    ])
  }
];

const jobFormRef = ref<InstanceType<typeof JobForm> | null>(null);
const jobSetSpecFormRef = ref<InstanceType<typeof JobSetSpecForm> | null>(null);


const fetchJobs = async () => {
  try {
    const param = create(JobList_InputSchema);
    const out = await clientKKSchedule.jobList(param);
    jobs.value = out.JobList || [];
  } catch (error) {
    toast.add({title: 'Error fetching job list', description: String(error), color: 'error'});
  }
};

onMounted(async () => {
  await fetchJobs();
});

const handleDisableJob = async (job: PBJob) => {
  try {
    const request = create(JobDisable_InputSchema, {serviceName: job.ServiceName, funcName: job.FuncName});
    await clientKKSchedule.jobDisable(request);
    await fetchJobs();
    toast.add({title: 'Job disabled successfully', color: 'success'});
  } catch (error) {
    toast.add({title: 'Error disabling job', description: String(error), color: 'error'});
  }
};

const handleDeleteJob = async (job: PBJob) => {
  modalTitle.value = 'Warning';
  modalDescription.value = `Are you sure you want to delete job "${job.FuncName}" from service "${job.ServiceName}"?`;
  modalConfirmButtonText.value = 'Yes';
  modalConfirmButtonColor.value = 'error';
  modalConfirmAction = async () => {
    try {
      const request = create(JobDelete_InputSchema, {serviceName: job.ServiceName, funcName: job.FuncName});
      await clientKKSchedule.jobDelete(request);
      await fetchJobs();
      toast.add({title: 'Job deleted successfully', color: 'success'});
    } catch (error) {
      toast.add({title: 'Error deleting job', description: String(error), color: 'error'});
    }
  };
  isModalOpen.value = true;
};

const handleEnableJob = async (job: PBJob) => {
  try {
    const request = create(JobEnable_InputSchema, {serviceName: job.ServiceName, funcName: job.FuncName});
    await clientKKSchedule.jobEnable(request);
    await fetchJobs();
    toast.add({title: 'Job enabled successfully', color: 'success'});
  } catch (error) {
    toast.add({title: 'Error enabling job', description: String(error), color: 'error'});
  }
};

const handleEditJob = async (job: PBJob) => {
  jobFormRef.value?.open(job);
};

const handleCreateJob = () => {
  jobFormRef.value?.open();
};

const handleSetSpecJob = (job: PBJob) => {
  jobSetSpecFormRef.value?.open(job);
};


const handleTriggerJob = async (job: PBJob) => {
  try {
    toast.add({title: 'Triggering job...', color: 'info'});
    const request = create(JobTrigger_InputSchema, {serviceName: job.ServiceName, funcName: job.FuncName});
    await clientKKSchedule.jobTrigger(request);
    toast.add({title: `Job "${job.FuncName}" triggered successfully`, color: 'success'});
  } catch (error) {
    toast.add({title: 'Error triggering job', description: String(error), color: 'error'});
  }
};

const handleModalConfirm = async () => {
  if (modalConfirmAction) {
    await modalConfirmAction();
  }
  isModalOpen.value = false;
};
</script>