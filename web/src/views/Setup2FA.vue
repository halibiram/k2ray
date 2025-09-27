<template>
  <div class="p-8 space-y-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
      Set Up Two-Factor Authentication
    </h2>

    <!-- Step 1: Start -->
    <div v-if="step === 1">
      <p class="text-gray-600 dark:text-gray-400">
        Strengthen your account security by enabling 2FA.
      </p>
      <button @click="startSetup" class="btn-primary" :disabled="isLoading">
        <span v-if="isLoading">Generating...</span>
        <span v-else>Enable 2FA</span>
      </button>
    </div>

    <!-- Step 2: Scan QR Code -->
    <div v-if="step === 2" class="space-y-4">
      <p>Scan the image below with your authenticator app (e.g., Google Authenticator).</p>
      <div class="flex justify-center p-4 bg-white rounded-md">
        <img :src="qrCodeDataUrl" alt="2FA QR Code" />
      </div>
      <p class="text-sm text-gray-500">
        Or manually enter this code:
        <code class="px-2 py-1 font-mono bg-gray-200 rounded dark:bg-gray-700">{{ secret }}</code>
      </p>
      <form @submit.prevent="verifyAndEnable">
        <label for="verify-code" class="block text-sm font-medium">Verification Code</label>
        <input
          id="verify-code"
          v-model="verificationCode"
          type="text"
          required
          class="w-full max-w-xs mt-1 input"
          placeholder="123456"
        />
        <button type="submit" class="mt-4 btn-primary" :disabled="isVerifying">
          <span v-if="isVerifying">Verifying...</span>
          <span v-else>Verify & Enable</span>
        </button>
      </form>
    </div>

    <!-- Step 3: Recovery Codes -->
    <div v-if="step === 3" class="space-y-4">
      <h3 class="text-xl font-semibold">Save Your Recovery Codes</h3>
      <p>
        If you lose access to your authenticator app, you can use these codes to log in.
        <strong>Store them in a safe place.</strong>
      </p>
      <div class="p-4 space-y-2 bg-gray-100 rounded-md dark:bg-gray-800">
        <div v-for="code in recoveryCodes" :key="code" class="font-mono text-lg">
          {{ code }}
        </div>
      </div>
      <button @click="finishSetup" class="btn-primary">
        I have saved my codes
      </button>
    </div>

    <p v-if="error" class="text-red-600">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const step = ref(1)
const isLoading = ref(false)
const isVerifying = ref(false)
const error = ref<string | null>(null)
const qrCodeDataUrl = ref('')
const secret = ref('')
const verificationCode = ref('')
const recoveryCodes = ref<string[]>([])

const authStore = useAuthStore()
const router = useRouter()

async function startSetup() {
  isLoading.value = true
  error.value = null
  try {
    const response = await authStore.enable2FA()
    qrCodeDataUrl.value = response.qrCode
    secret.value = response.secret
    step.value = 2
  } catch (err) {
    error.value = (err as Error).message
  } finally {
    isLoading.value = false
  }
}

async function verifyAndEnable() {
  isVerifying.value = true
  error.value = null
  try {
    const response = await authStore.verifyAndEnable2FA(verificationCode.value)
    recoveryCodes.value = response.recoveryCodes
    step.value = 3
  } catch (err) {
    error.value = (err as Error).message
  } finally {
    isVerifying.value = false
  }
}

function finishSetup() {
  // Redirect to the main settings page or dashboard
  router.push({ name: 'Dashboard' })
}
</script>

<style scoped>
@reference "../assets/styles/main.css";

.btn-primary {
  @apply px-4 py-2 font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500;
}
.input {
  @apply px-3 py-2 border border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500;
}
</style>