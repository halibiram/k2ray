<template>
  <div>
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold text-gray-800 dark:text-gray-100">{{ $t('configManager.title') }}</h1>
      <button
        @click="openCreateModal"
        class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none"
      >
        {{ $t('configManager.createNew') }}
      </button>
    </div>
    <p class="mt-2 text-gray-600 dark:text-gray-300">{{ $t('configManager.description') }}</p>

    <div v-if="configStore.isLoading" class="mt-8 text-center">
      <p class="dark:text-gray-300">{{ $t('configManager.loading') }}</p>
    </div>
    <div v-else-if="configStore.error" class="mt-8 text-red-500">
      <p>{{ $t('configManager.error', { error: configStore.error }) }}</p>
    </div>
    <div v-else class="mt-8">
      <ConfigList
        :configs="configStore.configs"
        @delete="handleDelete"
        @set-active="handleSetActive"
        @edit="openEditModal"
        @show-qr="handleShowQr"
      />
    </div>

    <div class="mt-8">
      <ConfigImport @import="handleImport" />
    </div>

    <ConfigEditor
      :is-open="isModalOpen"
      :config="editingConfig"
      @close="closeModal"
      @save="handleSave"
    />

    <QrCodeModal
      :is-open="isQrModalOpen"
      :qr-value="qrCodeValue"
      @close="closeQrModal"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useConfigStore, type V2rayConfig, type ConfigPayload } from '../stores/configs'
import ConfigList from '../components/config/ConfigList.vue'
import ConfigEditor from '../components/config/ConfigEditor.vue'
import QrCodeModal from '../components/config/QrCodeModal.vue'
import ConfigImport from '../components/config/ConfigImport.vue'

const configStore = useConfigStore()

const isModalOpen = ref(false)
const editingConfig = ref<V2rayConfig | null>(null)

const isQrModalOpen = ref(false)
const qrCodeValue = ref('')

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

const handleShowQr = (config: V2rayConfig) => {
  // Assuming the QR code should contain the main config object, stringified.
  // This can be adjusted to a specific format like VMess URI if needed.
  qrCodeValue.value = JSON.stringify(config.data)
  isQrModalOpen.value = true
}

const closeQrModal = () => {
  isQrModalOpen.value = false
  qrCodeValue.value = ''
}

const handleSave = async (payload: ConfigPayload) => {
  let success = false
  if (editingConfig.value && editingConfig.value.id) {
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

const handleImport = (configData: any) => {
  // Create a new config object from the imported data
  const newConfig = {
    id: 0, // No ID yet, as it's a new config
    name: 'New Imported Config', // Default name
    data: configData,
    // Add other fields from V2rayConfig with default values if needed
  }

  // Set this as the editing config and open the modal
  editingConfig.value = newConfig as V2rayConfig
  isModalOpen.value = true
}
</script>
