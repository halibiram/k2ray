<template>
  <div class="p-4 border dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 shadow-sm">
    <h2 class="text-xl font-semibold mb-4 text-gray-800 dark:text-gray-100">{{ $t('configImport.title') }}</h2>

    <div v-if="error" class="p-3 mb-4 text-red-800 bg-red-100 rounded-lg">
      <p>Error: {{ error }}</p>
    </div>

    <!-- Camera Scanner -->
    <div class="mb-4">
      <h3 class="text-lg font-medium mb-2 text-gray-700 dark:text-gray-200">{{ $t('configImport.scan') }}</h3>
      <qrcode-stream @decode="onDecode" @init="onInit" />
    </div>

    <!-- File Uploader -->
    <div>
      <h3 class="text-lg font-medium mb-2 text-gray-700 dark:text-gray-200">{{ $t('configImport.upload') }}</h3>
      <qrcode-capture @decode="onDecode" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QrcodeStream, QrcodeCapture } from 'vue-qrcode-reader'

const emit = defineEmits(['import'])
const error = ref<string | null>(null)

const onDecode = (decodedString: string) => {
  error.value = null
  try {
    const configData = JSON.parse(decodedString)
    // Basic validation to check if it's a plausible config object
    if (typeof configData === 'object' && configData !== null) {
      emit('import', configData)
    } else {
      throw new Error('QR code does not contain a valid configuration format.')
    }
  } catch (e) {
    error.value = 'Failed to parse QR code. Please ensure it contains a valid JSON configuration.'
    console.error(e)
  }
}

const onInit = async (promise: Promise<any>) => {
  try {
    await promise
  } catch (e: any) {
    if (e.name === 'NotAllowedError') {
      error.value = "Camera access was not allowed. Please grant permission to use the scanner."
    } else if (e.name === 'NotFoundError') {
      error.value = "No camera found on this device."
    } else if (e.name === 'NotSupportedError') {
      error.value = "Secure context (HTTPS) is required for camera access."
    } else {
      error.value = `Camera error: ${e.message}`
    }
  }
}
</script>