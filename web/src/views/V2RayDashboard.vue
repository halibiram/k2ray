<template>
  <div class="v2ray-dashboard">
    <!-- Kaldırılmış olan göstergeler -->
    <div class="dashboard-header">
      <h1 class="page-title">V2Ray Sunucu Yönetimi</h1>
      <div class="connection-status">
        <div class="status-indicator" :class="{ 'online': isConnected }"></div>
        <span class="status-text">{{ connectionStatus }}</span>
      </div>
    </div>

    <!-- Server Panels Grid -->
    <div class="servers-grid">
      <!-- V2Ray VMess Server -->
      <div class="server-panel">
        <div class="panel-header">
          <h3 class="panel-title">V2RAY VMESS SUNUCU</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="servers.vmess.enabled" @change="toggleServer('vmess')">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content">
          <div class="server-status">
            <span class="status-label">Durum:</span>
            <span class="status-value" :class="{ 'active': servers.vmess.enabled }">
              {{ servers.vmess.enabled ? 'Çalışıyor' : 'Devre dışı' }}
            </span>
          </div>
          <div v-if="servers.vmess.enabled" class="server-details">
            <div class="detail-item">
              <span class="detail-label">Sunucu adresi:</span>
              <span class="detail-value">{{ servers.vmess.address }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Port:</span>
              <span class="detail-value">{{ servers.vmess.port }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Erişim:</span>
              <div class="access-controls">
                <button class="btn-secondary" @click="showQRCode('vmess')">QR Kod</button>
                <button class="btn-secondary" @click="copyConfig('vmess')">Yapılandırmayı Kopyala</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- V2Ray VLESS Server -->
      <div class="server-panel">
        <div class="panel-header">
          <h3 class="panel-title">V2RAY VLESS SUNUCU</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="servers.vless.enabled" @change="toggleServer('vless')">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content">
          <div class="server-status">
            <span class="status-label">Durum:</span>
            <span class="status-value" :class="{ 'active': servers.vless.enabled }">
              {{ servers.vless.enabled ? 'Çalışıyor' : 'Devre dışı' }}
            </span>
          </div>
          <div v-if="servers.vless.enabled" class="server-details">
            <div class="detail-item">
              <span class="detail-label">Sunucu adresi:</span>
              <span class="detail-value">{{ servers.vless.address }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Port:</span>
              <span class="detail-value">{{ servers.vless.port }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- V2Ray Trojan Server -->
      <div class="server-panel">
        <div class="panel-header">
          <h3 class="panel-title">V2RAY TROJAN SUNUCU</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="servers.trojan.enabled" @change="toggleServer('trojan')">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content">
          <div class="server-status">
            <span class="status-label">Durum:</span>
            <span class="status-value" :class="{ 'active': servers.trojan.enabled }">
              {{ servers.trojan.enabled ? 'Çalışıyor' : 'Devre dışı' }}
            </span>
          </div>
          <div v-if="servers.trojan.enabled" class="server-details">
            <div class="detail-item">
              <span class="detail-label">Sunucu adresi:</span>
              <span class="detail-value">{{ servers.trojan.address }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Port:</span>
              <span class="detail-value">{{ servers.trojan.port }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- V2Ray Shadowsocks Server -->
      <div class="server-panel">
        <div class="panel-header">
          <h3 class="panel-title">SHADOWSOCKS SUNUCU</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="servers.shadowsocks.enabled" @change="toggleServer('shadowsocks')">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content">
          <div class="server-status">
            <span class="status-label">Durum:</span>
            <span class="status-value" :class="{ 'active': servers.shadowsocks.enabled }">
              {{ servers.shadowsocks.enabled ? 'Çalışıyor' : 'Devre dışı' }}
            </span>
          </div>
          <div v-if="servers.shadowsocks.enabled" class="server-details">
            <div class="detail-item">
              <span class="detail-label">Şifreleme:</span>
              <span class="detail-value">{{ servers.shadowsocks.encryption }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Connection Statistics Panel -->
      <div class="server-panel stats-panel">
        <div class="panel-header">
          <h3 class="panel-title">BAĞLANTI İSTATİSTİKLERİ</h3>
        </div>
        <div class="panel-content">
          <div class="stats-grid">
            <div class="stat-item">
              <span class="stat-label">Aktif Bağlantılar:</span>
              <span class="stat-value">{{ stats.activeConnections }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Toplam Trafik:</span>
              <span class="stat-value">{{ stats.totalTraffic }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Çalışma Süresi:</span>
              <span class="stat-value">{{ stats.uptime }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- System Status Panel -->
      <div class="server-panel system-panel">
        <div class="panel-header">
          <h3 class="panel-title">SİSTEM DURUMU</h3>
        </div>
        <div class="panel-content">
          <div class="system-stats">
            <div class="system-item">
              <span class="system-label">CPU Kullanımı:</span>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: systemStats.cpu + '%' }"></div>
              </div>
              <span class="system-value">{{ systemStats.cpu }}%</span>
            </div>
            <div class="system-item">
              <span class="system-label">Bellek Kullanımı:</span>
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: systemStats.memory + '%' }"></div>
              </div>
              <span class="system-value">{{ systemStats.memory }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- QR Code Modal -->
    <div v-if="showQRModal" class="modal-overlay" @click="closeQRModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>QR Kod - {{ currentQRServer.toUpperCase() }}</h3>
          <button class="modal-close" @click="closeQRModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="qr-code-container">
            <canvas ref="qrCanvas" class="qr-code"></canvas>
          </div>
          <div class="config-text">
            <textarea readonly :value="currentConfig" class="config-textarea"></textarea>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, nextTick } from 'vue'

// Connection status
const isConnected = ref(true)
const connectionStatus = computed(() => isConnected.value ? 'Bağlı' : 'Bağlantı Kesildi')

// Server states
const servers = reactive({
  vmess: {
    enabled: false,
    address: 'keenetic.local:1080',
    port: '1080',
    uuid: 'abcd1234-5678-90ef-ghij-klmnopqrstuv',
    alterId: '0'
  },
  vless: {
    enabled: false,
    address: 'keenetic.local:1081',
    port: '1081',
    uuid: 'efgh5678-90ab-cdef-1234-567890abcdef'
  },
  trojan: {
    enabled: false,
    address: 'keenetic.local:1082',
    port: '1082',
    password: 'trojan_password_123'
  },
  shadowsocks: {
    enabled: false,
    address: 'keenetic.local:1083',
    port: '1083',
    encryption: 'aes-256-gcm',
    password: 'ss_password_456'
  }
})

// Statistics
const stats = reactive({
  activeConnections: 0,
  totalTraffic: '0 MB',
  uptime: '0m'
})

const systemStats = reactive({
  cpu: 25,
  memory: 45
})

// QR Code modal
const showQRModal = ref(false)
const currentQRServer = ref('')
const currentConfig = ref('')
const qrCanvas = ref<HTMLCanvasElement | null>(null)

// Server management functions
const toggleServer = async (serverType: string) => {
  console.log(`Toggling ${serverType} server`)
  // Here you would make API call to backend
  updateStats()
}

const showQRCode = (serverType: string) => {
  currentQRServer.value = serverType
  currentConfig.value = generateConfig(serverType)
  showQRModal.value = true
  
  nextTick(() => {
    if (qrCanvas.value) {
      generateQRCode(currentConfig.value)
    }
  })
}

const closeQRModal = () => {
  showQRModal.value = false
}

const copyConfig = (serverType: string) => {
  const config = generateConfig(serverType)
  navigator.clipboard.writeText(config)
  console.log(`${serverType} config copied to clipboard`)
}

const generateConfig = (serverType: string): string => {
  const server = servers[serverType as keyof typeof servers]
  switch (serverType) {
    case 'vmess':
      return `vmess://${btoa(JSON.stringify({
        v: '2',
        ps: 'Keenetic V2Ray VMess',
        add: server.address.split(':')[0],
        port: server.port,
        id: server.uuid,
        aid: server.alterId,
        net: 'tcp',
        type: 'none',
        host: '',
        path: '',
        tls: ''
      }))}`
    case 'vless':
      return `vless://${server.uuid}@${server.address}?type=tcp&security=none#Keenetic+V2Ray+VLESS`
    case 'trojan':
      return `trojan://${server.password}@${server.address}?type=tcp&security=tls#Keenetic+V2Ray+Trojan`
    case 'shadowsocks':
      return `ss://${btoa(`${server.encryption}:${server.password}`)}@${server.address}#Keenetic+Shadowsocks`
    default:
      return ''
  }
}

const generateQRCode = (text: string) => {
  // Simple QR code generation placeholder
  if (qrCanvas.value) {
    const ctx = qrCanvas.value.getContext('2d')
    if (ctx) {
      ctx.fillStyle = '#ffffff'
      ctx.fillRect(0, 0, 200, 200)
      ctx.fillStyle = '#000000'
      ctx.font = '12px Arial'
      ctx.fillText('QR Code', 75, 100)
      ctx.fillText('Placeholder', 65, 120)
    }
  }
}

const updateStats = () => {
  const enabledServers = Object.values(servers).filter(s => s.enabled).length
  stats.activeConnections = enabledServers * 2
  stats.totalTraffic = `${enabledServers * 150} MB`
  stats.uptime = '2h 30m'
}

// Initialize
updateStats()
</script>

<style scoped>
.v2ray-dashboard {
  @apply p-6 space-y-6;
}

.dashboard-header {
  @apply flex items-center justify-between mb-6;
}

.page-title {
  @apply text-2xl font-bold text-cyan-400;
}

.connection-status {
  @apply flex items-center space-x-2;
}

.status-indicator {
  @apply w-3 h-3 rounded-full bg-red-500;
}

.status-indicator.online {
  @apply bg-green-400;
  box-shadow: 0 0 10px rgba(34, 197, 94, 0.5);
}

.status-text {
  @apply text-sm text-slate-300;
}

.servers-grid {
  @apply grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6;
}

.server-panel {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg overflow-hidden;
  backdrop-filter: blur(10px);
}

.panel-header {
  @apply flex items-center justify-between p-4 bg-slate-700/30 border-b border-slate-600;
}

.panel-title {
  @apply text-sm font-semibold text-cyan-400 uppercase tracking-wide;
}

.panel-toggle {
  @apply flex items-center;
}

.toggle-switch {
  @apply relative inline-flex items-center cursor-pointer;
}

.toggle-switch input[type="checkbox"] {
  @apply sr-only;
}

.slider {
  @apply w-11 h-6 bg-slate-600 rounded-full transition-colors duration-200;
  position: relative;
}

.slider::before {
  content: '';
  @apply absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full transition-transform duration-200;
}

.toggle-switch input[type="checkbox"]:checked + .slider {
  @apply bg-cyan-500;
}

.toggle-switch input[type="checkbox"]:checked + .slider::before {
  transform: translateX(1.25rem);
}

.panel-content {
  @apply p-4 space-y-3;
}

.server-status {
  @apply flex items-center justify-between;
}

.status-label {
  @apply text-sm text-slate-400;
}

.status-value {
  @apply text-sm font-medium text-red-400;
}

.status-value.active {
  @apply text-green-400;
}

.server-details {
  @apply space-y-2 pt-2 border-t border-slate-600;
}

.detail-item {
  @apply flex items-center justify-between;
}

.detail-label {
  @apply text-sm text-slate-400;
}

.detail-value {
  @apply text-sm text-slate-200 font-mono;
}

.access-controls {
  @apply flex space-x-2;
}

.btn-secondary {
  @apply px-3 py-1 text-xs text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-all duration-200;
}

.stats-panel, .system-panel {
  @apply md:col-span-2 xl:col-span-1;
}

.stats-grid {
  @apply space-y-3;
}

.stat-item {
  @apply flex items-center justify-between;
}

.stat-label {
  @apply text-sm text-slate-400;
}

.stat-value {
  @apply text-sm font-semibold text-cyan-400;
}

.system-stats {
  @apply space-y-4;
}

.system-item {
  @apply space-y-2;
}

.system-label {
  @apply text-sm text-slate-400;
}

.progress-bar {
  @apply w-full h-2 bg-slate-600 rounded-full overflow-hidden;
}

.progress-fill {
  @apply h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all duration-300;
}

.system-value {
  @apply text-sm text-slate-300;
}

/* Modal Styles */
.modal-overlay {
  @apply fixed inset-0 bg-black/50 flex items-center justify-center z-50;
}

.modal-content {
  @apply bg-slate-800 border border-slate-600 rounded-lg max-w-md w-full mx-4;
}

.modal-header {
  @apply flex items-center justify-between p-4 border-b border-slate-600;
}

.modal-header h3 {
  @apply text-lg font-semibold text-cyan-400;
}

.modal-close {
  @apply text-slate-400 hover:text-white text-xl font-bold;
}

.modal-body {
  @apply p-4 space-y-4;
}

.qr-code-container {
  @apply flex justify-center;
}

.qr-code {
  @apply border border-slate-600 rounded;
  width: 200px;
  height: 200px;
}

.config-text {
  @apply space-y-2;
}

.config-textarea {
  @apply w-full h-24 p-3 bg-slate-700 border border-slate-600 rounded text-sm text-slate-200 font-mono resize-none;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .v2ray-dashboard {
    @apply p-4 space-y-4;
  }
  
  .servers-grid {
    @apply grid-cols-1;
  }
  
  .dashboard-header {
    @apply flex-col items-start space-y-2;
  }
}
</style>