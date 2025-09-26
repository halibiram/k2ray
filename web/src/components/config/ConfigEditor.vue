<template>
  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="closeModal" class="relative z-10">
      <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0" enter-to="opacity-100" leave="duration-200 ease-in" leave-from="opacity-100" leave-to="opacity-0">
        <div class="fixed inset-0 bg-black bg-opacity-25" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4 text-center">
          <TransitionChild as="template" enter="duration-300 ease-out" enter-from="opacity-0 scale-95" enter-to="opacity-100 scale-100" leave="duration-200 ease-in" leave-from="opacity-100 scale-100" leave-to="opacity-0 scale-95">
            <DialogPanel class="w-full max-w-lg transform overflow-hidden rounded-2xl bg-white p-6 text-left align-middle shadow-xl transition-all">
              <DialogTitle as="h3" class="text-lg font-medium leading-6 text-gray-900">
                {{ isEditing ? 'Edit' : 'Create' }} Configuration
              </DialogTitle>
              <form @submit.prevent="submitForm" class="mt-4 space-y-4">
                <div>
                  <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
                  <input v-model="form.name" type="text" id="name" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                </div>
                <div>
                  <label for="protocol" class="block text-sm font-medium text-gray-700">Protocol</label>
                  <select v-model="form.protocol" @change="onProtocolChange" id="protocol" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" :disabled="isEditing">
                    <option value="vmess">VMess</option>
                    <option value="vless">VLESS</option>
                  </select>
                </div>

                <!-- VMess Fields -->
                <template v-if="form.protocol === 'vmess'">
                  <div>
                    <label for="address" class="block text-sm font-medium text-gray-700">Address</label>
                    <input v-model="form.config_data.add" type="text" id="address" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                  </div>
                <div>
                  <label for="port" class="block text-sm font-medium text-gray-700">Port</label>
                  <input v-model.number="form.config_data.port" type="number" id="port" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                </div>
                <div>
                  <label for="uuid" class="block text-sm font-medium text-gray-700">User ID (UUID)</label>
                  <input v-model="form.config_data.id" type="text" id="uuid" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                </div>
                 <div>
                  <label for="alterId" class="block text-sm font-medium text-gray-700">Alter ID</label>
                  <input v-model.number="form.config_data.aid" type="number" id="alterId" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                </div>
                </template>

                <!-- VLESS Fields -->
                <template v-if="form.protocol === 'vless'">
                   <div>
                    <label for="vless-address" class="block text-sm font-medium text-gray-700">Address</label>
                    <input v-model="form.config_data.add" type="text" id="vless-address" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                  </div>
                  <div>
                    <label for="vless-port" class="block text-sm font-medium text-gray-700">Port</label>
                    <input v-model.number="form.config_data.port" type="number" id="vless-port" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                  </div>
                  <div>
                    <label for="vless-uuid" class="block text-sm font-medium text-gray-700">User ID (UUID)</label>
                    <input v-model="form.config_data.id" type="text" id="vless-uuid" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" required />
                  </div>
                   <div>
                    <label for="vless-encryption" class="block text-sm font-medium text-gray-700">Encryption</label>
                    <input v-model="form.config_data.encryption" type="text" id="vless-encryption" value="none" disabled class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 bg-gray-100 focus:outline-none sm:text-sm" />
                  </div>
                </template>

                <div class="mt-6 flex justify-end space-x-2">
                  <button type="button" @click="closeModal" class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none">
                    Cancel
                  </button>
                  <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none">
                    Save
                  </button>
                </div>
              </form>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { TransitionRoot, TransitionChild, Dialog, DialogPanel, DialogTitle } from '@headlessui/vue'
import { v4 as uuidv4 } from 'uuid'
import type { V2rayConfig, ConfigPayload } from '../../stores/configs'

const props = defineProps<{
  isOpen: boolean
  config: V2rayConfig | null
}>()

const emit = defineEmits(['close', 'save'])

const isEditing = computed(() => !!props.config)

const form = ref<ConfigPayload>({
  name: '',
  protocol: 'vmess',
  config_data: {
    v: '2',
    add: '',
    port: 443,
    id: '', // UUID
    aid: 0,
    net: 'tcp',
    type: 'none',
    host: '',
    path: '',
    tls: '',
  },
})

watch(() => props.config, (newConfig) => {
  if (newConfig) {
    form.value.name = newConfig.name
    form.value.protocol = newConfig.protocol
    form.value.config_data = JSON.parse(newConfig.config_data)
  } else {
    // Reset form for creation, generating a new UUID for the config's ID
    form.value = {
      name: '',
      protocol: 'vmess',
      config_data: { v: '2', add: '', port: 443, id: uuidv4(), aid: 0, net: 'tcp', type: 'none', host: '', path: '', tls: '' },
    }
  }
})

const closeModal = () => {
  emit('close')
}

const submitForm = () => {
  emit('save', form.value)
}

const onProtocolChange = () => {
  if (form.value.protocol === 'vmess') {
    form.value.config_data = { v: '2', add: '', port: 443, id: uuidv4(), aid: 0, net: 'tcp', type: 'none', host: '', path: '', tls: '' }
  } else if (form.value.protocol === 'vless') {
    form.value.config_data = { id: uuidv4(), add: '', port: 443, encryption: 'none', flow: '' }
  }
}
</script>
