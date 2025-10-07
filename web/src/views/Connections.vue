<template>
  <div class="connections-page">
    <div class="page-header">
      <h1 class="page-title">Aktif Bağlantılar</h1>
      <div class="header-actions">
        <button class="btn-refresh" @click="refreshConnections">
          <span class="refresh-icon">↻</span>
          Yenile
        </button>
      </div>
    </div>

    <!-- Connection Statistics -->
    <div class="stats-overview">
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <div class="stat-value">{{ totalConnections }}</div>
          <div class="stat-label">Toplam Bağlantı</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📈</div>
        <div class="stat-content">
          <div class="stat-value">{{ activeConnections }}</div>
          <div class="stat-label">Aktif Bağlantı</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-value">{{ totalBandwidth }}</div>
          <div class="stat-label">Toplam Bant Genişliği</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">⚡</div>
        <div class="stat-content">
          <div class="stat-value">{{ avgLatency }}</div>
          <div class="stat-label">Ortalama Gecikme</div>
        </div>
      </div>
    </div>

    <!-- Active Connections Table -->
    <div class="connections-table-container">
      <div class="table-header">
        <h2 class="table-title">Aktif Bağlantılar</h2>
        <div class="table-filters">
          <select v-model="filterServer" class="filter-select">
            <option value="">Tüm Sunucular</option>
            <option value="vmess">VMess</option>
            <option value="vless">VLess</option>
            <option value="trojan">Trojan</option>
            <option value="shadowsocks">Shadowsocks</option>
          </select>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="connections-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Sunucu Tipi</th>
              <th>Kaynak IP</th>
              <th>Hedef</th>
              <th>Protokol</th>
              <th>Bağlantı Süresi</th>
              <th>Up/Down</th>
              <th>Durum</th>
              <th>İşlemler</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="connection in filteredConnections" :key="connection.id">
              <td class="connection-id">{{ connection.id }}</td>
              <td>
                <span class="server-badge" :class="connection.serverType">
                  {{ connection.serverType.toUpperCase() }}
                </span>
              </td>
              <td class="source-ip">{{ connection.sourceIP }}</td>
              <td class="destination">{{ connection.destination }}</td>
              <td class="protocol">{{ connection.protocol }}</td>
              <td class="duration">{{ connection.duration }}</td>
              <td class="bandwidth">
                <div class="bandwidth-info">
                  <span class="upload">↑{{ connection.upload }}</span>
                  <span class="download">↓{{ connection.download }}</span>
                </div>
              </td>
              <td>
                <span class="status-badge" :class="connection.status">
                  {{ connection.status === 'active' ? 'Aktif' : 'Kapalı' }}
                </span>
              </td>
              <td class="actions">
                <button 
                  class="action-btn disconnect" 
                  @click="disconnectConnection(connection.id)"
                  title="Bağlantıyı Kes"
                >
                  ✕
                </button>
                <button 
                  class="action-btn details" 
                  @click="showConnectionDetails(connection)"
                  title="Detayları Göster"
                >
                  ℹ
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Connection Details Modal -->
    <div v-if="showDetailsModal" class="modal-overlay" @click="closeDetailsModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Bağlantı Detayları - {{ selectedConnection?.id }}</h3>
          <button class="modal-close" @click="closeDetailsModal">&times;</button>
        </div>
        <div class="modal-body" v-if="selectedConnection">
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">Bağlantı ID:</span>
              <span class="detail-value">{{ selectedConnection.id }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Sunucu Tipi:</span>
              <span class="detail-value">{{ selectedConnection.serverType }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Kaynak IP:</span>
              <span class="detail-value">{{ selectedConnection.sourceIP }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Kaynak Port:</span>
              <span class="detail-value">{{ selectedConnection.sourcePort }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Hedef:</span>
              <span class="detail-value">{{ selectedConnection.destination }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Protokol:</span>
              <span class="detail-value">{{ selectedConnection.protocol }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Başlangıç Zamanı:</span>
              <span class="detail-value">{{ selectedConnection.startTime }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Toplam Upload:</span>
              <span class="detail-value">{{ selectedConnection.totalUpload }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Toplam Download:</span>
              <span class="detail-value">{{ selectedConnection.totalDownload }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'

interface Connection {
  id: string
  serverType: 'vmess' | 'vless' | 'trojan' | 'shadowsocks'
  sourceIP: string
  sourcePort: string
  destination: string
  protocol: string
  duration: string
  upload: string
  download: string
  totalUpload: string
  totalDownload: string
  status: 'active' | 'closed'
  startTime: string
}

// Stats
const totalConnections = ref(0)
const activeConnections = ref(0)
const totalBandwidth = ref('0 MB/s')
const avgLatency = ref('0 ms')

// Filters
const filterServer = ref('')

// Connections data
const connections = ref<Connection[]>([])

// Modal
const showDetailsModal = ref(false)
const selectedConnection = ref<Connection | null>(null)

// Computed
const filteredConnections = computed(() => {
  if (!filterServer.value) return connections.value
  return connections.value.filter(conn => conn.serverType === filterServer.value)
})

// Methods
const refreshConnections = async () => {
  console.log('Refreshing connections...')
  await loadConnections()
}

const loadConnections = async () => {
  // Mock data - replace with actual API call
  const mockConnections: Connection[] = [
    {
      id: 'conn-001',
      serverType: 'vmess',
      sourceIP: '192.168.1.100',
      sourcePort: '58392',
      destination: 'google.com:443',
      protocol: 'TCP',
      duration: '00:15:32',
      upload: '2.3 KB/s',
      download: '45.7 KB/s',
      totalUpload: '2.1 MB',
      totalDownload: '42.3 MB',
      status: 'active',
      startTime: '2024-10-07 14:30:15'
    },
    {
      id: 'conn-002',
      serverType: 'vless',
      sourceIP: '192.168.1.105',
      sourcePort: '58393',
      destination: 'youtube.com:443',
      protocol: 'TCP',
      duration: '00:08:47',
      upload: '1.2 KB/s',
      download: '125.4 KB/s',
      totalUpload: '635 KB',
      totalDownload: '65.2 MB',
      status: 'active',
      startTime: '2024-10-07 14:37:28'
    },
    {
      id: 'conn-003',
      serverType: 'trojan',
      sourceIP: '192.168.1.108',
      sourcePort: '58394',
      destination: 'twitter.com:443',
      protocol: 'TCP',
      duration: '00:03:12',
      upload: '0.8 KB/s',
      download: '12.3 KB/s',
      totalUpload: '154 KB',
      totalDownload: '2.4 MB',
      status: 'active',
      startTime: '2024-10-07 14:42:03'
    }
  ]
  
  connections.value = mockConnections
  updateStats()
}

const updateStats = () => {
  totalConnections.value = connections.value.length
  activeConnections.value = connections.value.filter(c => c.status === 'active').length
  totalBandwidth.value = '183.4 KB/s'
  avgLatency.value = '45 ms'
}

const disconnectConnection = async (connectionId: string) => {
  console.log(`Disconnecting connection: ${connectionId}`)
  // API call to disconnect
  const index = connections.value.findIndex(c => c.id === connectionId)
  if (index !== -1) {
    connections.value[index].status = 'closed'
    updateStats()
  }
}

const showConnectionDetails = (connection: Connection) => {
  selectedConnection.value = connection
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedConnection.value = null
}

onMounted(() => {
  loadConnections()
})
</script>

<style scoped>
.connections-page {
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

.btn-refresh {
  @apply flex items-center space-x-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-cyan-400 rounded-lg transition-colors duration-200;
}

.refresh-icon {
  @apply text-lg;
}

.stats-overview {
  @apply grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-4;
}

.stat-card {
  @apply flex items-center p-4 bg-slate-800/50 border border-slate-700 rounded-lg;
}

.stat-icon {
  @apply text-2xl mr-3;
}

.stat-content {
  @apply flex-1;
}

.stat-value {
  @apply text-lg font-bold text-cyan-400;
}

.stat-label {
  @apply text-sm text-slate-400;
}

.connections-table-container {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg overflow-hidden;
}

.table-header {
  @apply flex items-center justify-between p-4 border-b border-slate-600;
}

.table-title {
  @apply text-lg font-semibold text-cyan-400;
}

.table-filters {
  @apply flex space-x-2;
}

.filter-select {
  @apply px-3 py-1 bg-slate-700 border border-slate-600 text-slate-200 rounded text-sm;
}

.table-wrapper {
  @apply overflow-x-auto;
}

.connections-table {
  @apply w-full;
}

.connections-table th {
  @apply px-4 py-3 text-left text-xs font-semibold text-slate-300 uppercase tracking-wide bg-slate-700/50;
}

.connections-table td {
  @apply px-4 py-3 text-sm text-slate-200 border-t border-slate-700;
}

.connection-id {
  @apply font-mono text-xs;
}

.server-badge {
  @apply px-2 py-1 text-xs font-semibold rounded-full;
}

.server-badge.vmess {
  @apply bg-blue-500/20 text-blue-400;
}

.server-badge.vless {
  @apply bg-green-500/20 text-green-400;
}

.server-badge.trojan {
  @apply bg-purple-500/20 text-purple-400;
}

.server-badge.shadowsocks {
  @apply bg-orange-500/20 text-orange-400;
}

.source-ip, .destination {
  @apply font-mono text-xs;
}

.bandwidth-info {
  @apply flex flex-col text-xs;
}

.upload {
  @apply text-red-400;
}

.download {
  @apply text-green-400;
}

.status-badge {
  @apply px-2 py-1 text-xs font-semibold rounded-full;
}

.status-badge.active {
  @apply bg-green-500/20 text-green-400;
}

.status-badge.closed {
  @apply bg-red-500/20 text-red-400;
}

.actions {
  @apply flex space-x-1;
}

.action-btn {
  @apply w-8 h-8 flex items-center justify-center rounded text-xs font-bold transition-colors duration-200;
}

.action-btn.disconnect {
  @apply text-red-400 hover:bg-red-500/20;
}

.action-btn.details {
  @apply text-cyan-400 hover:bg-cyan-500/20;
}

/* Modal Styles */
.modal-overlay {
  @apply fixed inset-0 bg-black/50 flex items-center justify-center z-50;
}

.modal-content {
  @apply bg-slate-800 border border-slate-600 rounded-lg max-w-2xl w-full mx-4;
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
  @apply p-4;
}

.detail-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-4;
}

.detail-item {
  @apply flex flex-col space-y-1;
}

.detail-label {
  @apply text-sm font-semibold text-slate-400;
}

.detail-value {
  @apply text-sm text-slate-200 font-mono;
}

/* Responsive */
@media (max-width: 768px) {
  .connections-page {
    @apply p-4 space-y-4;
  }
  
  .page-header {
    @apply flex-col items-start space-y-2;
  }
  
  .stats-overview {
    @apply grid-cols-2;
  }
}
</style>