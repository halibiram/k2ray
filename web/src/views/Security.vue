<template>
  <div class="security-page">
    <div class="page-header">
      <h1 class="page-title">Güvenlik Ayarları</h1>
      <div class="security-status">
        <div class="security-indicator" :class="{ 'secure': overallSecurity >= 80 }"></div>
        <span class="security-text">Güvenlik Seviyesi: {{ securityLevel }}</span>
      </div>
    </div>

    <!-- Security Overview -->
    <div class="security-overview">
      <div class="overview-card threat-detection">
        <div class="card-header">
          <h3 class="card-title">Tehdit Tespiti</h3>
          <div class="threat-status" :class="{ 'active': threatDetection.enabled }">
            {{ threatDetection.enabled ? 'Aktif' : 'Pasif' }}
          </div>
        </div>
        <div class="card-content">
          <div class="threat-stats">
            <div class="stat-item">
              <span class="stat-value">{{ threatDetection.blockedAttacks }}</span>
              <span class="stat-label">Engellenen Saldırı</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ threatDetection.suspiciousIPs }}</span>
              <span class="stat-label">Şüpheli IP</span>
            </div>
          </div>
          <button 
            class="toggle-btn" 
            :class="{ 'enabled': threatDetection.enabled }"
            @click="toggleThreatDetection"
          >
            {{ threatDetection.enabled ? 'Devre Dışı Bırak' : 'Etkinleştir' }}
          </button>
        </div>
      </div>

      <div class="overview-card firewall">
        <div class="card-header">
          <h3 class="card-title">Güvenlik Duvarı</h3>
          <div class="firewall-status" :class="{ 'active': firewall.enabled }">
            {{ firewall.enabled ? 'Aktif' : 'Pasif' }}
          </div>
        </div>
        <div class="card-content">
          <div class="firewall-info">
            <p class="info-text">Aktif Kurallar: {{ firewall.activeRules }}</p>
            <p class="info-text">Engellenen Bağlantı: {{ firewall.blockedConnections }}</p>
          </div>
          <button 
            class="config-btn"
            @click="showFirewallConfig"
          >
            Yapılandır
          </button>
        </div>
      </div>

      <div class="overview-card encryption">
        <div class="card-header">
          <h3 class="card-title">Şifreleme</h3>
          <div class="encryption-status active">Güçlü</div>
        </div>
        <div class="card-content">
          <div class="encryption-info">
            <p class="info-text">TLS 1.3 Aktif</p>
            <p class="info-text">AES-256-GCM</p>
          </div>
          <button class="config-btn" @click="showEncryptionSettings">
            Ayarlar
          </button>
        </div>
      </div>
    </div>

    <!-- Security Settings Panels -->
    <div class="security-panels">
      <!-- Access Control Panel -->
      <div class="security-panel">
        <div class="panel-header">
          <h3 class="panel-title">ERİŞİM KONTROLÜ</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="accessControl.enabled" @change="updateAccessControl">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content" v-if="accessControl.enabled">
          <div class="access-settings">
            <div class="setting-item">
              <label class="setting-label">Maksimum Bağlantı Sayısı:</label>
              <input 
                type="number" 
                v-model="accessControl.maxConnections" 
                class="setting-input"
                min="1"
                max="1000"
              >
            </div>
            <div class="setting-item">
              <label class="setting-label">Oturum Süresi (dakika):</label>
              <input 
                type="number" 
                v-model="accessControl.sessionTimeout" 
                class="setting-input"
                min="5"
                max="1440"
              >
            </div>
            <div class="setting-item">
              <label class="setting-label">IP Whitelist:</label>
              <div class="ip-list">
                <div v-for="(ip, index) in accessControl.whitelist" :key="index" class="ip-item">
                  <span class="ip-address">{{ ip }}</span>
                  <button class="remove-btn" @click="removeFromWhitelist(index)">×</button>
                </div>
                <div class="add-ip">
                  <input 
                    type="text" 
                    v-model="newWhitelistIP" 
                    placeholder="192.168.1.1" 
                    class="ip-input"
                    @keyup.enter="addToWhitelist"
                  >
                  <button class="add-btn" @click="addToWhitelist">Ekle</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Rate Limiting Panel -->
      <div class="security-panel">
        <div class="panel-header">
          <h3 class="panel-title">HIZ SINIRLAMA</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="rateLimiting.enabled" @change="updateRateLimiting">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content" v-if="rateLimiting.enabled">
          <div class="rate-settings">
            <div class="setting-item">
              <label class="setting-label">Dakikada Maksimum İstek:</label>
              <input 
                type="number" 
                v-model="rateLimiting.requestsPerMinute" 
                class="setting-input"
                min="10"
                max="10000"
              >
            </div>
            <div class="setting-item">
              <label class="setting-label">Bant Genişliği Limiti (MB/s):</label>
              <input 
                type="number" 
                v-model="rateLimiting.bandwidthLimit" 
                class="setting-input"
                min="1"
                max="1000"
              >
            </div>
            <div class="setting-item">
              <label class="setting-label">Eşzamanlı Bağlantı Limiti:</label>
              <input 
                type="number" 
                v-model="rateLimiting.concurrentLimit" 
                class="setting-input"
                min="1"
                max="500"
              >
            </div>
          </div>
        </div>
      </div>

      <!-- DDoS Protection Panel -->
      <div class="security-panel">
        <div class="panel-header">
          <h3 class="panel-title">DDoS KORUNMASI</h3>
          <div class="panel-toggle">
            <label class="toggle-switch">
              <input type="checkbox" v-model="ddosProtection.enabled" @change="updateDDoSProtection">
              <span class="slider"></span>
            </label>
          </div>
        </div>
        <div class="panel-content" v-if="ddosProtection.enabled">
          <div class="ddos-settings">
            <div class="protection-level">
              <label class="setting-label">Koruma Seviyesi:</label>
              <select v-model="ddosProtection.level" class="level-select">
                <option value="low">Düşük</option>
                <option value="medium">Orta</option>
                <option value="high">Yüksek</option>
                <option value="paranoid">Paranoyak</option>
              </select>
            </div>
            <div class="ddos-stats">
              <div class="stat-row">
                <span class="stat-label">Son 24 Saat Engellenen:</span>
                <span class="stat-value">{{ ddosProtection.blockedAttacks24h }}</span>
              </div>
              <div class="stat-row">
                <span class="stat-label">Şu Anda Engellenen IP:</span>
                <span class="stat-value">{{ ddosProtection.currentlyBlocked }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security Logs Panel -->
      <div class="security-panel logs-panel">
        <div class="panel-header">
          <h3 class="panel-title">GÜVENLİK GÜNLÜKLERI</h3>
          <button class="refresh-logs-btn" @click="refreshSecurityLogs">Yenile</button>
        </div>
        <div class="panel-content">
          <div class="logs-container">
            <div class="log-filters">
              <select v-model="logFilter" class="log-filter-select">
                <option value="all">Tüm Günlükler</option>
                <option value="threat">Tehdit Tespiti</option>
                <option value="firewall">Güvenlik Duvarı</option>
                <option value="ddos">DDoS Korunması</option>
                <option value="access">Erişim Kontrolü</option>
              </select>
            </div>
            <div class="logs-list">
              <div 
                v-for="log in filteredLogs" 
                :key="log.id" 
                class="log-entry"
                :class="log.severity"
              >
                <div class="log-time">{{ log.timestamp }}</div>
                <div class="log-type">{{ log.type }}</div>
                <div class="log-message">{{ log.message }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'

interface SecurityLog {
  id: string
  timestamp: string
  type: string
  severity: 'info' | 'warning' | 'danger'
  message: string
  category: 'threat' | 'firewall' | 'ddos' | 'access'
}

// Security overview data
const overallSecurity = ref(85)
const securityLevel = computed(() => {
  if (overallSecurity.value >= 90) return 'Çok Güvenli'
  if (overallSecurity.value >= 70) return 'Güvenli'
  if (overallSecurity.value >= 50) return 'Orta'
  return 'Risk'
})

const threatDetection = reactive({
  enabled: true,
  blockedAttacks: 127,
  suspiciousIPs: 15
})

const firewall = reactive({
  enabled: true,
  activeRules: 45,
  blockedConnections: 89
})

// Security settings
const accessControl = reactive({
  enabled: true,
  maxConnections: 100,
  sessionTimeout: 60,
  whitelist: ['192.168.1.100', '10.0.0.50']
})

const rateLimiting = reactive({
  enabled: true,
  requestsPerMinute: 1000,
  bandwidthLimit: 50,
  concurrentLimit: 100
})

const ddosProtection = reactive({
  enabled: true,
  level: 'medium' as 'low' | 'medium' | 'high' | 'paranoid',
  blockedAttacks24h: 34,
  currentlyBlocked: 8
})

// Logs
const securityLogs = ref<SecurityLog[]>([])
const logFilter = ref('all')
const newWhitelistIP = ref('')

const filteredLogs = computed(() => {
  if (logFilter.value === 'all') return securityLogs.value
  return securityLogs.value.filter(log => log.category === logFilter.value)
})

// Methods
const toggleThreatDetection = () => {
  threatDetection.enabled = !threatDetection.enabled
  console.log(`Threat detection ${threatDetection.enabled ? 'enabled' : 'disabled'}`)
}

const showFirewallConfig = () => {
  console.log('Opening firewall configuration')
}

const showEncryptionSettings = () => {
  console.log('Opening encryption settings')
}

const updateAccessControl = () => {
  console.log('Access control updated:', accessControl)
}

const updateRateLimiting = () => {
  console.log('Rate limiting updated:', rateLimiting)
}

const updateDDoSProtection = () => {
  console.log('DDoS protection updated:', ddosProtection)
}

const addToWhitelist = () => {
  if (newWhitelistIP.value && !accessControl.whitelist.includes(newWhitelistIP.value)) {
    accessControl.whitelist.push(newWhitelistIP.value)
    newWhitelistIP.value = ''
  }
}

const removeFromWhitelist = (index: number) => {
  accessControl.whitelist.splice(index, 1)
}

const refreshSecurityLogs = async () => {
  console.log('Refreshing security logs')
  loadSecurityLogs()
}

const loadSecurityLogs = () => {
  // Mock security logs data
  const mockLogs: SecurityLog[] = [
    {
      id: 'log-001',
      timestamp: '2024-10-07 15:30:45',
      type: 'DDoS Korunması',
      severity: 'warning',
      message: 'IP 203.0.113.42 adresinden şüpheli trafik tespit edildi',
      category: 'ddos'
    },
    {
      id: 'log-002', 
      timestamp: '2024-10-07 15:28:12',
      type: 'Güvenlik Duvarı',
      severity: 'info',
      message: 'Yeni güvenlik kuralı eklendi: Block-Suspicious-IPs',
      category: 'firewall'
    },
    {
      id: 'log-003',
      timestamp: '2024-10-07 15:25:33',
      type: 'Tehdit Tespiti',
      severity: 'danger',
      message: 'Brute force saldırı girişimi engellendi',
      category: 'threat'
    },
    {
      id: 'log-004',
      timestamp: '2024-10-07 15:20:18',
      type: 'Erişim Kontrolü',
      severity: 'info',
      message: 'Whitelist IP 192.168.1.100 başarılı bağlantı',
      category: 'access'
    }
  ]
  
  securityLogs.value = mockLogs
}

onMounted(() => {
  loadSecurityLogs()
})
</script>

<style scoped>
.security-page {
  @apply p-6 space-y-6;
}

.page-header {
  @apply flex items-center justify-between;
}

.page-title {
  @apply text-2xl font-bold text-cyan-400;
}

.security-status {
  @apply flex items-center space-x-2;
}

.security-indicator {
  @apply w-3 h-3 rounded-full bg-yellow-500;
}

.security-indicator.secure {
  @apply bg-green-400;
  box-shadow: 0 0 10px rgba(34, 197, 94, 0.5);
}

.security-text {
  @apply text-sm text-slate-300;
}

.security-overview {
  @apply grid grid-cols-1 md:grid-cols-3 gap-6;
}

.overview-card {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg p-4;
}

.card-header {
  @apply flex items-center justify-between mb-3;
}

.card-title {
  @apply text-lg font-semibold text-cyan-400;
}

.threat-status, .firewall-status, .encryption-status {
  @apply px-2 py-1 text-xs font-semibold rounded-full bg-red-500/20 text-red-400;
}

.threat-status.active, .firewall-status.active, .encryption-status.active {
  @apply bg-green-500/20 text-green-400;
}

.card-content {
  @apply space-y-3;
}

.threat-stats {
  @apply flex justify-between;
}

.stat-item {
  @apply text-center;
}

.stat-value {
  @apply block text-lg font-bold text-cyan-400;
}

.stat-label {
  @apply text-xs text-slate-400;
}

.toggle-btn {
  @apply w-full py-2 text-sm font-semibold rounded border transition-colors duration-200;
}

.toggle-btn.enabled {
  @apply bg-red-500/20 text-red-400 border-red-400 hover:bg-red-500/30;
}

.toggle-btn:not(.enabled) {
  @apply bg-green-500/20 text-green-400 border-green-400 hover:bg-green-500/30;
}

.firewall-info, .encryption-info {
  @apply space-y-1;
}

.info-text {
  @apply text-sm text-slate-300;
}

.config-btn {
  @apply w-full py-2 text-sm font-semibold text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.security-panels {
  @apply grid grid-cols-1 lg:grid-cols-2 gap-6;
}

.security-panel {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg overflow-hidden;
}

.logs-panel {
  @apply lg:col-span-2;
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
  @apply p-4;
}

.access-settings, .rate-settings, .ddos-settings {
  @apply space-y-4;
}

.setting-item {
  @apply flex flex-col space-y-2;
}

.setting-label {
  @apply text-sm font-semibold text-slate-300;
}

.setting-input, .level-select {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.ip-list {
  @apply space-y-2;
}

.ip-item {
  @apply flex items-center justify-between p-2 bg-slate-700/50 rounded;
}

.ip-address {
  @apply font-mono text-sm text-slate-200;
}

.remove-btn {
  @apply w-6 h-6 flex items-center justify-center text-red-400 hover:bg-red-500/20 rounded;
}

.add-ip {
  @apply flex space-x-2;
}

.ip-input {
  @apply flex-1 px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.add-btn {
  @apply px-4 py-2 text-sm font-semibold text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.protection-level {
  @apply mb-4;
}

.ddos-stats {
  @apply space-y-2;
}

.stat-row {
  @apply flex justify-between;
}

.refresh-logs-btn {
  @apply px-3 py-1 text-sm text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.logs-container {
  @apply space-y-4;
}

.log-filters {
  @apply flex;
}

.log-filter-select {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded;
}

.logs-list {
  @apply space-y-2 max-h-64 overflow-y-auto;
}

.log-entry {
  @apply p-3 bg-slate-700/30 border-l-4 rounded-r;
}

.log-entry.info {
  @apply border-blue-400;
}

.log-entry.warning {
  @apply border-yellow-400;
}

.log-entry.danger {
  @apply border-red-400;
}

.log-time {
  @apply text-xs text-slate-400 mb-1;
}

.log-type {
  @apply text-sm font-semibold text-cyan-400 mb-1;
}

.log-message {
  @apply text-sm text-slate-300;
}

/* Responsive */
@media (max-width: 768px) {
  .security-page {
    @apply p-4 space-y-4;
  }
  
  .security-overview {
    @apply grid-cols-1;
  }
  
  .security-panels {
    @apply grid-cols-1;
  }
  
  .logs-panel {
    @apply lg:col-span-1;
  }
}
</style>