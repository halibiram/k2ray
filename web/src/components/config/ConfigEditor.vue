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
                  <input v-model="form.name" type="text" id="name" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3" required />
                </div>
                <div>
                  <label for="protocol" class="block text-sm font-medium text-gray-700">Protocol</label>
                  <select v-model="form.protocol" @change="onProtocolChange" id="protocol" class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3" :disabled="isEditing">
                    <option value="vmess">VMess</option>
                    <option value="vless">VLESS</option>
                    <option value="shadowsocks">Shadowsocks</option>
                    <option value="trojan">Trojan</option>
                  </select>
                </div>

                <!-- Dynamic Fields based on Protocol -->
                <div v-if="form.protocol" class="space-y-4">
                  <!-- Common Fields for VMess/VLESS -->
                  <template v-if="form.protocol === 'vmess' || form.protocol === 'vless'">
                    <div>
                      <label :for="`${form.protocol}-address`" class="block text-sm font-medium text-gray-700">Address</label>
                      <input v-model="form.config_data.add" type="text" :id="`${form.protocol}-address`" class="mt-1 block w-full" required />
                    </div>
                    <div>
                      <label :for="`${form.protocol}-port`" class="block text-sm font-medium text-gray-700">Port</label>
                      <input v-model.number="form.config_data.port" type="number" :id="`${form.protocol}-port`" class="mt-1 block w-full" required />
                    </div>
                    <div>
                      <label :for="`${form.protocol}-uuid`" class="block text-sm font-medium text-gray-700">User ID (UUID)</label>
                      <input v-model="form.config_data.id" type="text" :id="`${form.protocol}-uuid`" class="mt-1 block w-full" required />
                    </div>
                  </template>

                  <!-- VMess Specific Fields -->
                  <template v-if="form.protocol === 'vmess'">
                    <div>
                      <label for="vmess-aid" class="block text-sm font-medium text-gray-700">Alter ID</label>
                      <input v-model.number="form.config_data.aid" type="number" id="vmess-aid" class="mt-1 block w-full" required />
                    </div>
                  </template>

                  <!-- VLESS Specific Fields -->
                  <template v-if="form.protocol === 'vless'">
                    <div>
                      <label for="vless-encryption" class="block text-sm font-medium text-gray-700">Encryption</label>
                      <input v-model="form.config_data.encryption" type="text" id="vless-encryption" disabled class="mt-1 block w-full bg-gray-100" />
                    </div>
                     <div>
                      <label for="vless-flow" class="block text-sm font-medium text-gray-700">Flow</label>
                      <input v-model="form.config_data.flow" type="text" id="vless-flow" class="mt-1 block w-full" />
                    </div>
                  </template>

                  <!-- Common Fields for Trojan/Shadowsocks -->
                  <template v-if="form.protocol === 'trojan' || form.protocol === 'shadowsocks'">
                     <div>
                      <label :for="`${form.protocol}-server`" class="block text-sm font-medium text-gray-700">Server Address</label>
                      <input v-model="form.config_data.server" type="text" :id="`${form.protocol}-server`" class="mt-1 block w-full" required />
                    </div>
                    <div>
                      <label :for="`${form.protocol}-port`" class="block text-sm font-medium text-gray-700">Server Port</label>
                      <input v-model.number="form.config_data.server_port" type="number" :id="`${form.protocol}-port`" class="mt-1 block w-full" required />
                    </div>
                     <div>
                      <label :for="`${form.protocol}-password`" class="block text-sm font-medium text-gray-700">Password</label>
                      <input v-model="form.config_data.password" type="password" :id="`${form.protocol}-password`" class="mt-1 block w-full" required />
                    </div>
                  </template>

                   <!-- Trojan Specific Fields -->
                  <template v-if="form.protocol === 'trojan'">
                    <div>
                      <label for="trojan-sni" class="block text-sm font-medium text-gray-700">SNI (Server Name Indication)</label>
                      <input v-model="form.config_data.sni" type="text" id="trojan-sni" class="mt-1 block w-full" />
                    </div>
                  </template>

                  <!-- Shadowsocks Specific Fields -->
                  <template v-if="form.protocol === 'shadowsocks'">
                     <div>
                      <label for="ss-method" class="block text-sm font-medium text-gray-700">Method</label>
                      <input v-model="form.config_data.method" type="text" id="ss-method" class="mt-1 block w-full" required />
                    </div>
                  </template>

                  <!-- Transport Settings -->
                  <div class="border-t pt-4 mt-4" v-if="form.protocol !== 'shadowsocks'">
                    <h4 class="text-md font-medium text-gray-800">Transport Settings</h4>
                    <div class="grid grid-cols-2 gap-4 mt-2">
                       <div>
                        <label for="transport-net" class="block text-sm font-medium text-gray-700">Network</label>
                        <select v-model="form.config_data.net" id="transport-net" class="mt-1 block w-full">
                          <option value="tcp">TCP</option>
                          <option value="kcp">mKCP</option>
                          <option value="ws">WebSocket</option>
                          <option value="h2">HTTP/2</option>
                          <option value="grpc">gRPC</option>
                        </select>
                      </div>
                      <div>
                        <label for="transport-tls" class="block text-sm font-medium text-gray-700">Security</label>
                        <select v-model="form.config_data.tls" id="transport-tls" class="mt-1 block w-full">
                          <option value="none">None</option>
                          <option value="tls">TLS</option>
                        </select>
                      </div>
                    </div>
                    <!-- WebSocket Settings -->
                    <div v-if="form.config_data.net === 'ws'" class="mt-4 space-y-2">
                       <div>
                        <label for="ws-path" class="block text-sm font-medium text-gray-700">WebSocket Path</label>
                        <input v-model="form.config_data.wsSettings.path" id="ws-path" type="text" class="mt-1 block w-full" />
                      </div>
                    </div>
                    <!-- gRPC Settings -->
                     <div v-if="form.config_data.net === 'grpc'" class="mt-4 space-y-2">
                       <div>
                        <label for="grpc-service-name" class="block text-sm font-medium text-gray-700">gRPC Service Name</label>
                        <input v-model="form.config_data.grpcSettings.serviceName" id="grpc-service-name" type="text" class="mt-1 block w-full" />
                      </div>
                    </div>
                  </div>
                </div>

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

const defaultWsSettings = () => ({ path: '/', headers: {} });
const defaultGrpcSettings = () => ({ serviceName: '' });

const defaultTransport = () => ({
  net: 'tcp',
  tls: 'none',
  wsSettings: defaultWsSettings(),
  grpcSettings: defaultGrpcSettings(),
});

const defaultVmessData = () => ({
  v: '2',
  add: '',
  port: 443,
  id: uuidv4(),
  aid: 0,
  type: 'none',
  host: '',
  path: '',
  ...defaultTransport(),
})

const defaultVlessData = () => ({
  id: uuidv4(),
  add: '',
  port: 443,
  encryption: 'none',
  flow: '',
  ...defaultTransport(),
})

const defaultShadowsocksData = () => ({
  server: '',
  server_port: 8388,
  password: '',
  method: 'aes-256-gcm',
  // No transport settings needed for SS typically
})

const defaultTrojanData = () => ({
  server: '',
  server_port: 443,
  password: '',
  sni: '',
  ...defaultTransport(),
})

const form = ref<ConfigPayload>({
  name: '',
  protocol: 'vmess',
  config_data: defaultVmessData(),
})

const onProtocolChange = () => {
  if (form.value.protocol === 'vmess') {
    form.value.config_data = defaultVmessData()
  } else if (form.value.protocol === 'vless') {
    form.value.config_data = defaultVlessData()
  } else if (form.value.protocol === 'shadowsocks') {
    form.value.config_data = defaultShadowsocksData()
  } else if (form.value.protocol === 'trojan') {
    form.value.config_data = defaultTrojanData()
  }
}

// Deep merge utility
const deepMerge = (target, source) => {
  for (const key in source) {
    if (source[key] instanceof Object && key in target) {
      Object.assign(source[key], deepMerge(target[key], source[key]))
    }
  }
  Object.assign(target || {}, source)
  return target
}

watch(() => props.config, (newConfig) => {
  if (newConfig) {
    form.value.name = newConfig.name
    form.value.protocol = newConfig.protocol

    // Get the correct default data structure for the protocol
    let defaultConfigData;
    switch (newConfig.protocol) {
      case 'vmess': defaultConfigData = defaultVmessData(); break;
      case 'vless': defaultConfigData = defaultVlessData(); break;
      case 'trojan': defaultConfigData = defaultTrojanData(); break;
      case 'shadowsocks': defaultConfigData = defaultShadowsocksData(); break;
      default: defaultConfigData = {};
    }

    // Deep merge the loaded config into the default structure
    const loadedConfigData = JSON.parse(newConfig.config_data)
    form.value.config_data = deepMerge(defaultConfigData, loadedConfigData)

  } else {
    // Reset form for creation
    form.value.name = ''
    // Keep the current protocol or default to vmess, then reset config_data
    onProtocolChange()
  }
}, { immediate: true, deep: true })

const closeModal = () => {
  emit('close')
}

const submitForm = () => {
  emit('save', form.value)
}
</script>