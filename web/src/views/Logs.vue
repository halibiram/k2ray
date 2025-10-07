<template>
  <div class="logs-page">
    <div class="page-header">
      <h1 class="page-title">Sistem Günlükleri</h1>
      <div class="header-actions">
        <button class="btn-refresh" @click="refreshLogs">
          <span class="refresh-icon">↻</span>
          Yenile
        </button>
        <button class="btn-export" @click="exportLogs">
          <span class="export-icon">↓</span>
          Dışa Aktar
        </button>
        <button class="btn-clear" @click="clearLogs">
          <span class="clear-icon">🗑</span>
          Temizle
        </button>
      </div>
    </div>

    <!-- Log Filters -->
    <div class="log-filters-section">
      <div class="filters-grid">
        <div class="filter-item">
          <label class="filter-label">Kategori:</label>
          <select v-model="filters.category" class="filter-select" @change="applyFilters">
            <option value="">Tümü</option>
            <option value="v2ray">V2Ray</option>
            <option value="system">Sistem</option>
            <option value="security">Güvenlik</option>
            <option value="network">Ağ</option>
            <option value="auth">Kimlik Doğrulama</option>
          </select>
        </div>
        
        <div class="filter-item">
          <label class="filter-label">Seviye:</label>
          <select v-model="filters.level" class="filter-select" @change="applyFilters">
            <option value="">Tümü</option>
            <option value="debug">Debug</option>
            <option value="info">Info</option>
            <option value="warning">Warning</option>
            <option value="error">Error</option>
            <option value="critical">Critical</option>
          </select>
        </div>
        
        <div class="filter-item">
          <label class="filter-label">Tarih Aralığı:</label>
          <select v-model="filters.timeRange" class="filter-select" @change="applyFilters">
            <option value="1h">Son 1 Saat</option>
            <option value="6h">Son 6 Saat</option>
            <option value="24h">Son 24 Saat</option>
            <option value="7d">Son 7 Gün</option>
            <option value="30d">Son 30 Gün</option>
          </select>
        </div>
        
        <div class="filter-item">
          <label class="filter-label">Arama:</label>
          <input 
            type="text" 
            v-model="filters.search" 
            placeholder="Günlüklerde ara..."
            class="filter-input"
            @input="applyFilters"
          >
        </div>
      </div>
    </div>

    <!-- Log Statistics -->
    <div class="log-stats">
      <div class="stat-card">
        <div class="stat-icon debug">D</div>
        <div class="stat-content">
          <div class="stat-value">{{ logStats.debug }}</div>
          <div class="stat-label">Debug</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon info">I</div>
        <div class="stat-content">
          <div class="stat-value">{{ logStats.info }}</div>
          <div class="stat-label">Info</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon warning">W</div>
        <div class="stat-content">
          <div class="stat-value">{{ logStats.warning }}</div>
          <div class="stat-label">Warning</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon error">E</div>
        <div class="stat-content">
          <div class="stat-value">{{ logStats.error }}</div>
          <div class="stat-label">Error</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon critical">C</div>
        <div class="stat-content">
          <div class="stat-value">{{ logStats.critical }}</div>
          <div class="stat-label">Critical</div>
        </div>
      </div>
    </div>

    <!-- Log Viewer -->
    <div class="log-viewer">
      <div class="viewer-header">
        <div class="viewer-info">
          <span class="total-logs">Toplam: {{ filteredLogs.length }} kayıt</span>
          <span class="last-update">Son güncelleme: {{ lastUpdate }}</span>
        </div>
        <div class="viewer-options">
          <label class="auto-refresh">
            <input type="checkbox" v-model="autoRefresh" @change="toggleAutoRefresh">
            <span class="checkbox-label">Otomatik Yenileme</span>
          </label>
          <button 
            class="follow-btn" 
            :class="{ 'active': followMode }"
            @click="toggleFollowMode"
          >
            {{ followMode ? 'Takibi Durdur' : 'Takip Et' }}
          </button>
        </div>
      </div>

      <div class="log-container" ref="logContainer">
        <div 
          v-for="log in paginatedLogs" 
          :key="log.id"
          class="log-entry"
          :class="[log.level, log.category]"
          @click="selectLog(log)"
        >
          <div class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</div>
          <div class="log-level" :class="log.level">{{ log.level.toUpperCase() }}</div>
          <div class="log-category">{{ getCategoryDisplay(log.category) }}</div>
          <div class="log-message">{{ log.message }}</div>
          <div class="log-source" v-if="log.source">{{ log.source }}</div>
        </div>
        
        <div v-if="isLoading" class="loading-indicator">
          <div class="spinner"></div>
          <span>Günlükler yükleniyor...</span>
        </div>
        
        <div v-if="filteredLogs.length === 0 && !isLoading" class="no-logs">
          <div class="no-logs-icon">📄</div>
          <div class="no-logs-text">Filtrelere uygun günlük kaydı bulunamadı</div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="filteredLogs.length > pageSize" class="pagination">
        <button 
          class="page-btn" 
          :disabled="currentPage === 1"
          @click="goToPage(currentPage - 1)"
        >
          Önceki
        </button>
        
        <span class="page-info">
          Sayfa {{ currentPage }} / {{ totalPages }}
        </span>
        
        <button 
          class="page-btn" 
          :disabled="currentPage === totalPages"
          @click="goToPage(currentPage + 1)"
        >
          Sonraki
        </button>
      </div>
    </div>

    <!-- Log Detail Modal -->
    <div v-if="selectedLog" class="modal-overlay" @click="closeLogDetail">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Günlük Detayı</h3>
          <button class="modal-close" @click="closeLogDetail">&times;</button>
        </div>
        <div class="modal-body">
          <div class="log-detail">
            <div class="detail-row">
              <span class="detail-label">Zaman:</span>
              <span class="detail-value">{{ selectedLog.timestamp }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Seviye:</span>
              <span class="detail-value" :class="selectedLog.level">{{ selectedLog.level.toUpperCase() }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Kategori:</span>
              <span class="detail-value">{{ getCategoryDisplay(selectedLog.category) }}</span>
            </div>
            <div class="detail-row" v-if="selectedLog.source">
              <span class="detail-label">Kaynak:</span>
              <span class="detail-value">{{ selectedLog.source }}</span>
            </div>
            <div class="detail-row full-width">
              <span class="detail-label">Mesaj:</span>
              <div class="detail-message">{{ selectedLog.message }}</div>
            </div>
            <div class="detail-row full-width" v-if="selectedLog.stackTrace">
              <span class="detail-label">Stack Trace:</span>
              <pre class="detail-stack">{{ selectedLog.stackTrace }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'

interface LogEntry {
  id: string
  timestamp: string
  level: 'debug' | 'info' | 'warning' | 'error' | 'critical'
  category: 'v2ray' | 'system' | 'security' | 'network' | 'auth'
  message: string
  source?: string
  stackTrace?: string
}

// State
const logs = ref<LogEntry[]>([])
const selectedLog = ref<LogEntry | null>(null)
const isLoading = ref(false)
const autoRefresh = ref(false)
const followMode = ref(false)
const lastUpdate = ref('')
const logContainer = ref<HTMLElement | null>(null)

// Filters
const filters = reactive({
  category: '',
  level: '',
  timeRange: '24h',
  search: ''
})

// Pagination
const currentPage = ref(1)
const pageSize = 50

// Auto refresh interval
let refreshInterval: number | null = null

// Computed
const filteredLogs = computed(() => {
  let filtered = logs.value

  if (filters.category) {
    filtered = filtered.filter(log => log.category === filters.category)
  }

  if (filters.level) {
    filtered = filtered.filter(log => log.level === filters.level)
  }

  if (filters.search) {
    const searchLower = filters.search.toLowerCase()
    filtered = filtered.filter(log => 
      log.message.toLowerCase().includes(searchLower) ||
      (log.source && log.source.toLowerCase().includes(searchLower))
    )
  }

  // Apply time range filter
  const now = new Date()
  const timeRangeMs = getTimeRangeMs(filters.timeRange)
  const cutoff = new Date(now.getTime() - timeRangeMs)
  
  filtered = filtered.filter(log => new Date(log.timestamp) >= cutoff)

  return filtered.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
})

const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return filteredLogs.value.slice(start, end)
})

const totalPages = computed(() => Math.ceil(filteredLogs.value.length / pageSize))

const logStats = computed(() => {
  const stats = {
    debug: 0,
    info: 0,
    warning: 0,
    error: 0,
    critical: 0
  }

  filteredLogs.value.forEach(log => {
    stats[log.level]++
  })

  return stats
})

// Methods
const getTimeRangeMs = (range: string): number => {
  const ranges: Record<string, number> = {
    '1h': 60 * 60 * 1000,
    '6h': 6 * 60 * 60 * 1000,
    '24h': 24 * 60 * 60 * 1000,
    '7d': 7 * 24 * 60 * 60 * 1000,
    '30d': 30 * 24 * 60 * 60 * 1000
  }
  return ranges[range] || ranges['24h']
}

const formatTimestamp = (timestamp: string): string => {
  return new Date(timestamp).toLocaleString('tr-TR')
}

const getCategoryDisplay = (category: string): string => {
  const displays: Record<string, string> = {
    'v2ray': 'V2Ray',
    'system': 'Sistem',
    'security': 'Güvenlik', 
    'network': 'Ağ',
    'auth': 'Kimlik Doğrulama'
  }
  return displays[category] || category
}

const loadLogs = async () => {
  isLoading.value = true
  
  // Mock data - replace with actual API call
  const mockLogs: LogEntry[] = [
    {
      id: 'log-001',
      timestamp: '2024-10-07 15:45:32',
      level: 'info',
      category: 'v2ray',
      message: 'V2Ray VMess sunucusu başlatıldı',
      source: 'v2ray-core'
    },
    {
      id: 'log-002',
      timestamp: '2024-10-07 15:43:18',
      level: 'warning',
      category: 'security',
      message: 'Şüpheli bağlantı girişimi tespit edildi: 203.0.113.42',
      source: 'security-monitor'
    },
    {
      id: 'log-003',
      timestamp: '2024-10-07 15:40:55',
      level: 'error',
      category: 'network',
      message: 'Ağ bağlantısı kesildi, yeniden bağlanıyor...',
      source: 'network-manager'
    },
    {
      id: 'log-004',
      timestamp: '2024-10-07 15:38:20',
      level: 'debug',
      category: 'v2ray',
      message: 'İstemci bağlantısı kabul edildi: 192.168.1.100:45382',
      source: 'v2ray-inbound'
    },
    {
      id: 'log-005',
      timestamp: '2024-10-07 15:35:47',
      level: 'critical',
      category: 'system',
      message: 'Bellek kullanımı kritik seviyeye ulaştı: %95',
      source: 'system-monitor',
      stackTrace: 'SystemError: Memory usage critical\n  at monitor.check(monitor.js:45)\n  at Timer.schedule(timer.js:12)'
    },
    {
      id: 'log-006',
      timestamp: '2024-10-07 15:32:15',
      level: 'info',
      category: 'auth',
      message: 'Kullanıcı admin başarıyla giriş yaptı',
      source: 'auth-service'
    },
    {
      id: 'log-007',
      timestamp: '2024-10-07 15:30:03',
      level: 'warning',
      category: 'v2ray',
      message: 'V2Ray VLESS sunucusu yavaş yanıt veriyor',
      source: 'v2ray-outbound'
    },
    {
      id: 'log-008',
      timestamp: '2024-10-07 15:27:41',
      level: 'info',
      category: 'system',
      message: 'Sistem güncellemesi tamamlandı, yeniden başlatma gerekiyor',
      source: 'update-manager'
    }
  ]

  // Simulate loading delay
  await new Promise(resolve => setTimeout(resolve, 500))
  
  logs.value = mockLogs
  lastUpdate.value = new Date().toLocaleString('tr-TR')
  isLoading.value = false
}

const refreshLogs = async () => {
  await loadLogs()
  if (followMode.value) {
    scrollToBottom()
  }
}

const applyFilters = () => {
  currentPage.value = 1
}

const goToPage = (page: number) => {
  currentPage.value = page
}

const selectLog = (log: LogEntry) => {
  selectedLog.value = log
}

const closeLogDetail = () => {
  selectedLog.value = null
}

const toggleAutoRefresh = () => {
  if (autoRefresh.value) {
    refreshInterval = setInterval(refreshLogs, 5000) // 5 seconds
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
}

const toggleFollowMode = () => {
  followMode.value = !followMode.value
  if (followMode.value) {
    scrollToBottom()
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const exportLogs = () => {
  const csvContent = generateCSV(filteredLogs.value)
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', `logs_${new Date().toISOString().split('T')[0]}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const generateCSV = (logs: LogEntry[]): string => {
  const headers = ['Timestamp', 'Level', 'Category', 'Message', 'Source']
  const rows = logs.map(log => [
    log.timestamp,
    log.level,
    log.category,
    log.message.replace(/"/g, '""'),
    log.source || ''
  ])
  
  return [
    headers.join(','),
    ...rows.map(row => row.map(field => `"${field}"`).join(','))
  ].join('\n')
}

const clearLogs = () => {
  if (confirm('Tüm günlükleri silmek istediğinizden emin misiniz?')) {
    logs.value = []
    console.log('Logs cleared')
  }
}

onMounted(() => {
  loadLogs()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.logs-page {
  @apply p-6 space-y-6;
}

.page-header {
  @apply flex items-center justify-between;
}

.page-title {
  @apply text-2xl font-bold text-cyan-400;
}

.header-actions {
  @apply flex space-x-2;
}

.btn-refresh, .btn-export, .btn-clear {
  @apply flex items-center space-x-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-200 rounded-lg transition-colors duration-200;
}

.btn-clear {
  @apply bg-red-600 hover:bg-red-500;
}

.refresh-icon, .export-icon, .clear-icon {
  @apply text-lg;
}

.log-filters-section {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg p-4;
}

.filters-grid {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4;
}

.filter-item {
  @apply flex flex-col space-y-2;
}

.filter-label {
  @apply text-sm font-semibold text-slate-300;
}

.filter-select, .filter-input {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.log-stats {
  @apply flex flex-wrap gap-4;
}

.stat-card {
  @apply flex items-center p-3 bg-slate-800/50 border border-slate-700 rounded-lg;
}

.stat-icon {
  @apply w-8 h-8 flex items-center justify-center rounded-full text-sm font-bold mr-3;
}

.stat-icon.debug {
  @apply bg-gray-500/20 text-gray-400;
}

.stat-icon.info {
  @apply bg-blue-500/20 text-blue-400;
}

.stat-icon.warning {
  @apply bg-yellow-500/20 text-yellow-400;
}

.stat-icon.error {
  @apply bg-red-500/20 text-red-400;
}

.stat-icon.critical {
  @apply bg-red-600/20 text-red-300;
}

.stat-content {
  @apply flex flex-col;
}

.stat-value {
  @apply text-lg font-bold text-cyan-400;
}

.stat-label {
  @apply text-sm text-slate-400;
}

.log-viewer {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg overflow-hidden;
}

.viewer-header {
  @apply flex items-center justify-between p-4 border-b border-slate-600;
}

.viewer-info {
  @apply flex items-center space-x-4 text-sm text-slate-400;
}

.viewer-options {
  @apply flex items-center space-x-4;
}

.auto-refresh {
  @apply flex items-center space-x-2 text-sm text-slate-300 cursor-pointer;
}

.checkbox-label {
  @apply select-none;
}

.follow-btn {
  @apply px-3 py-1 text-sm text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.follow-btn.active {
  @apply bg-cyan-400 text-slate-900;
}

.log-container {
  @apply max-h-96 overflow-y-auto;
}

.log-entry {
  @apply grid grid-cols-12 gap-2 p-3 border-b border-slate-700 hover:bg-slate-700/30 cursor-pointer transition-colors duration-150;
  font-size: 13px;
}

.log-timestamp {
  @apply col-span-2 text-slate-400 font-mono;
}

.log-level {
  @apply col-span-1 font-semibold uppercase;
}

.log-level.debug {
  @apply text-gray-400;
}

.log-level.info {
  @apply text-blue-400;
}

.log-level.warning {
  @apply text-yellow-400;
}

.log-level.error {
  @apply text-red-400;
}

.log-level.critical {
  @apply text-red-300;
}

.log-category {
  @apply col-span-2 text-cyan-400 font-medium;
}

.log-message {
  @apply col-span-6 text-slate-200;
}

.log-source {
  @apply col-span-1 text-slate-500 text-xs;
}

.loading-indicator {
  @apply flex items-center justify-center p-8 space-x-3 text-slate-400;
}

.spinner {
  @apply w-5 h-5 border-2 border-slate-600 border-t-cyan-400 rounded-full animate-spin;
}

.no-logs {
  @apply flex flex-col items-center justify-center p-8 text-slate-500;
}

.no-logs-icon {
  @apply text-4xl mb-2;
}

.pagination {
  @apply flex items-center justify-center p-4 space-x-4;
}

.page-btn {
  @apply px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:bg-slate-800 disabled:text-slate-600 text-slate-200 rounded transition-colors duration-200;
}

.page-info {
  @apply text-sm text-slate-400;
}

/* Modal Styles */
.modal-overlay {
  @apply fixed inset-0 bg-black/50 flex items-center justify-center z-50;
}

.modal-content {
  @apply bg-slate-800 border border-slate-600 rounded-lg max-w-4xl w-full mx-4 max-h-96 overflow-hidden;
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
  @apply p-4 overflow-y-auto max-h-80;
}

.log-detail {
  @apply space-y-3;
}

.detail-row {
  @apply flex items-start;
}

.detail-row.full-width {
  @apply flex-col;
}

.detail-label {
  @apply w-24 flex-shrink-0 text-sm font-semibold text-slate-400 mr-4;
}

.detail-value {
  @apply text-sm text-slate-200;
}

.detail-value.debug {
  @apply text-gray-400;
}

.detail-value.info {
  @apply text-blue-400;
}

.detail-value.warning {
  @apply text-yellow-400;
}

.detail-value.error {
  @apply text-red-400;
}

.detail-value.critical {
  @apply text-red-300;
}

.detail-message {
  @apply text-sm text-slate-200 mt-1 p-2 bg-slate-700/50 rounded;
}

.detail-stack {
  @apply text-xs text-slate-300 mt-1 p-3 bg-slate-900 rounded font-mono overflow-x-auto;
}

/* Responsive */
@media (max-width: 768px) {
  .logs-page {
    @apply p-4 space-y-4;
  }
  
  .filters-grid {
    @apply grid-cols-1;
  }
  
  .log-stats {
    @apply grid grid-cols-2 gap-2;
  }
  
  .log-entry {
    @apply grid-cols-1 gap-1;
  }
  
  .viewer-header {
    @apply flex-col items-start space-y-2;
  }
  
  .viewer-options {
    @apply w-full justify-between;
  }
}
</style>