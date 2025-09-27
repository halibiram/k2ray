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

    <ConfirmModal
      :is-open="isConfirmModalOpen"
      :title="t('configManager.confirmDelete.title')"
      :message="t('configManager.confirmDelete.message')"
      :confirm-button-text="t('common.delete')"
      :cancel-button-text="t('common.cancel')"
      @confirm="confirmDelete"
      @close="closeConfirmModal"
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
import ConfirmModal from '../components/common/ConfirmModal.vue'
import { useI18n } from 'vue-i18n'

const configStore = useConfigStore()
const { t } = useI18n()

const isModalOpen = ref(false)
const editingConfig = ref<V2rayConfig | null>(null)

const isQrModalOpen = ref(false)
const qrCodeValue = ref('')

const isConfirmModalOpen = ref(false)
const configToDeleteId = ref<number | null>(null)

onMounted(() => {
  configStore.fetchConfigs()
})

const handleDelete = (id: number) => {
  configToDeleteId.value = id
  isConfirmModalOpen.value = true
}

const confirmDelete = async () => {
  if (configToDeleteId.value !== null) {
    await configStore.deleteConfig(configToDeleteId.value)
  }
  closeConfirmModal()
}

const closeConfirmModal = () => {
  isConfirmModalOpen.value = false
  configToDeleteId.value = null
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
  // The QR code value should be the raw config data string.
  // This can be adjusted later to a specific format like VMess URI.
  qrCodeValue.value = config.config_data;
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
  // Create a new config object from the imported data.
  // We provide a full V2rayConfig structure to satisfy the type.
  // The user will edit details in the modal.
  const newConfig: V2rayConfig = {
    id: 0, // Temporary ID for a new item
    user_id: 0, // Assuming a default/placeholder
    name: 'New Imported Config',
    protocol: 'vmess', // Default protocol, user can change in editor
    config_data: JSON.stringify(configData),
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };

  editingConfig.value = newConfig;
  isModalOpen.value = true;
}
</script>
