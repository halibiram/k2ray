<template>
  <div>
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold text-gray-800">Configuration Manager</h1>
      <button
        @click="openCreateModal"
        class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none"
      >
        Create New Config
      </button>
    </div>
    <p class="mt-2 text-gray-600">Manage your V2Ray configurations.</p>

    <div v-if="configStore.isLoading" class="mt-8 text-center">
      <p>Loading configurations...</p>
    </div>
    <div v-else-if="configStore.error" class="mt-8 text-red-500">
      <p>Error: {{ configStore.error }}</p>
    </div>
    <div v-else class="mt-8">
      <ConfigList
        :configs="configStore.configs"
        @delete="handleDelete"
        @set-active="handleSetActive"
        @edit="openEditModal"
      />
    </div>

    <ConfigEditor
      :is-open="isModalOpen"
      :config="editingConfig"
      @close="closeModal"
      @save="handleSave"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useConfigStore, type V2rayConfig, type ConfigPayload } from '../stores/configs'
import ConfigList from '../components/config/ConfigList.vue'
import ConfigEditor from '../components/config/ConfigEditor.vue'

const configStore = useConfigStore()

const isModalOpen = ref(false)
const editingConfig = ref<V2rayConfig | null>(null)

onMounted(() => {
  configStore.fetchConfigs()
})

const handleDelete = async (id: number) => {
  if (window.confirm('Are you sure you want to delete this configuration?')) {
    await configStore.deleteConfig(id)
  }
}

const handleSetActive = async (id: number) => {
  await configStore.setActiveConfig(id)
  // Optionally, add a notification to the user
}

const closeModal = () => {
  isModalOpen.value = false
  editingConfig.value = null
}

const openCreateModal = () => {
  editingConfig.value = null
  isModalOpen.value = true
}

const openEditModal = (config: V2rayConfig) => {
  editingConfig.value = config
  isModalOpen.value = true
}

const handleSave = async (payload: ConfigPayload) => {
  let success = false
  if (editingConfig.value) {
    // Update existing config
    success = await configStore.updateConfig(editingConfig.value.id, payload)
  } else {
    // Create new config
    success = await configStore.createConfig(payload)
  }
  if (success) {
    closeModal()
  }
}
</script>
