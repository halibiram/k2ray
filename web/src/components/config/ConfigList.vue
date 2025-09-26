<template>
  <div class="overflow-x-auto bg-white dark:bg-gray-700 rounded-lg shadow">
    <table class="min-w-full">
      <thead class="bg-gray-50 dark:bg-gray-600">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">{{ $t('configList.name') }}</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">{{ $t('configList.protocol') }}</th>
          <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">{{ $t('configList.actions') }}</th>
        </tr>
      </thead>
      <tbody class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-600">
        <tr v-if="configs.length === 0">
          <td colspan="3" class="px-6 py-4 text-center text-gray-500 dark:text-gray-400">{{ $t('configList.noConfigs') }}</td>
        </tr>
        <tr v-for="config in configs" :key="config.id">
          <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">{{ config.name }}</td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">{{ config.protocol }}</td>
          <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
            <button @click="$emit('setActive', config.id)" class="text-blue-600 hover:text-blue-400">{{ $t('configList.setActive') }}</button>
            <button @click="$emit('edit', config)" class="text-indigo-600 hover:text-indigo-400">{{ $t('configList.edit') }}</button>
            <button @click="$emit('showQr', config)" class="text-green-600 hover:text-green-400">{{ $t('configList.qrCode') }}</button>
            <button @click="$emit('delete', config.id)" class="text-red-600 hover:text-red-400">{{ $t('configList.delete') }}</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { V2rayConfig } from '../../stores/configs'

defineProps<{
  configs: V2rayConfig[]
}>()

defineEmits(['edit', 'delete', 'setActive', 'showQr'])
</script>
