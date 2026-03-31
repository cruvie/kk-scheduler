<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-bold">Task Execution Records</h1>
      <UButton icon="i-heroicons-arrow-path" @click="fetchExecutions">Refresh</UButton>
    </div>

    <UTable :data="executions" :columns="columns">

    </UTable>

    <USlideover v-model:open="isLogPanelOpen" title="Execution Log" side="right" :ui="{ content: 'w-2/3 max-w-none' }">
      <template #body>
        <pre class="whitespace-pre-wrap text-sm bg-gray-100 dark:bg-gray-800 p-3 rounded h-full overflow-auto">{{
            selectedLog
          }}</pre>
      </template>
    </USlideover>

    <UModal v-model:open="isDeleteModalOpen" title="Warning" :ui="{ footer: 'justify-end' }">
      <template #body>
        <p>Are you sure you want to delete this task execution record?</p>
      </template>

      <template #footer>
        <UButton color="neutral" @click="isDeleteModalOpen = false">No</UButton>
        <UButton color="error" @click="handleConfirmDelete">Yes</UButton>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import {h, ref, onMounted, resolveComponent} from 'vue';
import {clientKKSchedule} from '~/utils/api/client';
import {
  TaskExecutionList_InputSchema,
  TaskExecutionDelete_InputSchema,
  TaskExecutionGet_InputSchema
} from '~~/gen/kk_scheduler/TaskExecution_pb';
import type {PBTaskExecution} from '~~/gen/kk_scheduler/TaskExecution_pb';
import {TaskExecutionStatus} from '~~/gen/kk_scheduler/TaskExecution_pb';
import {create} from "@bufbuild/protobuf";
import {useToast} from '#imports';
import type {TableColumn} from '@nuxt/ui';

const UButton = resolveComponent('UButton')
const UBadge = resolveComponent('UBadge')

const executions = ref<PBTaskExecution[]>([]);
const toast = useToast();

const isLogPanelOpen = ref(false);
const selectedLog = ref('');
const isDeleteModalOpen = ref(false);
let deleteTargetId: string | null = null;

const statusBadge = (status: TaskExecutionStatus) => {
  const map: Record<TaskExecutionStatus, { color: string; text: string }> = {
    [TaskExecutionStatus.UNSPECIFIED]: {color: 'neutral', text: 'Unknown'},
    [TaskExecutionStatus.Init]: {color: 'warning', text: 'Init'},
    [TaskExecutionStatus.RUNNING]: {color: 'info', text: 'Running'},
    [TaskExecutionStatus.COMPLETED]: {color: 'success', text: 'Completed'},
    [TaskExecutionStatus.FAILED]: {color: 'error', text: 'Failed'},
  };
  const {color, text} = map[status] || map[TaskExecutionStatus.UNSPECIFIED];
  return h(UBadge, {color}, () => text);
};

const formatTimestamp = (ts?: { seconds: bigint }) => {
  if (!ts?.seconds) return '';
  return new Date(Number(ts.seconds) * 1000).toLocaleString();
};

const columns: TableColumn<PBTaskExecution>[] = [
  {accessorKey: 'Id', header: 'Id'},
  {accessorKey: 'JobId', header: 'Job ID'},
  {
    accessorKey: 'Status',
    header: 'Status',
    cell: ({row}) => statusBadge(row.original.Status),
  },
  {
    accessorKey: 'StartedAt',
    header: 'Started At',
    cell: ({row}) => formatTimestamp(row.original.StartedAt),
  },
  {
    accessorKey: 'FinishedAt',
    header: 'Finished At',
    cell: ({row}) => formatTimestamp(row.original.FinishedAt),
  },
  {
    accessorKey: 'actions', header: 'Actions',
    cell: ({row}) => h('div', {class: 'flex flex-wrap gap-2 w-[200px]'}, [
      h(UButton, {onClick: () => handleViewLog(row.original)}, () => 'View Log'),
      h(UButton, {color: 'error', onClick: () => handleDeleteExecution(row.original)}, () => 'Delete'),
    ]),
  },
];

const fetchExecutions = async () => {
  try {
    const param = create(TaskExecutionList_InputSchema);
    const out = await clientKKSchedule.taskExecutionList(param);
    executions.value = out.TaskExecutionList || [];
  } catch (error) {
    toast.add({title: 'Error fetching task execution list', description: String(error), color: 'error'});
  }
};

onMounted(async () => {
  await fetchExecutions();
});

const handleViewLog = async (execution: PBTaskExecution) => {
  try {
    const request = create(TaskExecutionGet_InputSchema, {Id: execution.Id});
    const out = await clientKKSchedule.taskExecutionGet(request);
    selectedLog.value = out.TaskExecution?.Log || '(empty)';
    isLogPanelOpen.value = true;
  } catch (error) {
    toast.add({title: 'Error fetching task execution detail', description: String(error), color: 'error'});
  }
};

const handleDeleteExecution = (execution: PBTaskExecution) => {
  deleteTargetId = execution.Id;
  isDeleteModalOpen.value = true;
};

const handleConfirmDelete = async () => {
  if (!deleteTargetId) return;
  try {
    const request = create(TaskExecutionDelete_InputSchema, {Id: deleteTargetId});
    await clientKKSchedule.taskExecutionDelete(request);
    await fetchExecutions();
    toast.add({title: 'Task execution deleted successfully', color: 'success'});
  } catch (error) {
    toast.add({title: 'Error deleting task execution', description: String(error), color: 'error'});
  }
  isDeleteModalOpen.value = false;
  deleteTargetId = null;
};
</script>
