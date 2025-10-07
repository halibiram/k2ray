<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">Sistem Ayarları</h1>
      <div class="header-actions">
        <button class="btn-backup" @click="createBackup">
          <span class="backup-icon">💾</span>
          Yedekle
        </button>
        <button class="btn-restore" @click="showRestoreModal = true">
          <span class="restore-icon">📂</span>
          Geri Yükle
        </button>
      </div>
    </div>

    <!-- Settings Navigation -->
    <div class="settings-nav">
      <button 
        v-for="tab in settingsTabs" 
        :key="tab.id"
        class="nav-tab" 
        :class="{ 'active': activeTab === tab.id }"
        @click="activeTab = tab.id"
      >
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-text">{{ tab.name }}</span>
      </button>
    </div>

    <!-- Settings Content -->
    <div class="settings-content">
      <!-- General Settings -->
      <div v-show="activeTab === 'general'" class="settings-section">
        <h2 class="section-title">Genel Ayarlar</h2>
        
        <div class="settings-grid">
          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Sistem Bilgileri</h3>
            </div>
            <div class="setting-content">
              <div class="info-grid">
                <div class="info-item">
                  <span class="info-label">Model:</span>
                  <span class="info-value">{{ systemInfo.model }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Firmware:</span>
                  <span class="info-value">{{ systemInfo.firmware }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">K2Ray Sürümü:</span>
                  <span class="info-value">{{ systemInfo.k2rayVersion }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Çalışma Süresi:</span>
                  <span class="info-value">{{ systemInfo.uptime }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Ağ Ayarları</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">Web Arayüz Portu:</label>
                <input 
                  type="number" 
                  v-model="generalSettings.webPort" 
                  class="setting-input"
                  min="1000"
                  max="65535"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">API Portu:</label>
                <input 
                  type="number" 
                  v-model="generalSettings.apiPort" 
                  class="setting-input"
                  min="1000"
                  max="65535"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">HTTPS Kullan:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="generalSettings.httpsEnabled">
                  <span class="slider"></span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- V2Ray Settings -->
      <div v-show="activeTab === 'v2ray'" class="settings-section">
        <h2 class="section-title">V2Ray Ayarları</h2>
        
        <div class="settings-grid">
          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Genel V2Ray Ayarları</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">V2Ray Core Yolu:</label>
                <input 
                  type="text" 
                  v-model="v2raySettings.corePath" 
                  class="setting-input"
                  placeholder="/opt/bin/v2ray"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Yapılandırma Dizini:</label>
                <input 
                  type="text" 
                  v-model="v2raySettings.configDir" 
                  class="setting-input"
                  placeholder="/opt/etc/v2ray"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Log Seviyesi:</label>
                <select v-model="v2raySettings.logLevel" class="setting-select">
                  <option value="none">None</option>
                  <option value="error">Error</option>
                  <option value="warning">Warning</option>
                  <option value="info">Info</option>
                  <option value="debug">Debug</option>
                </select>
              </div>
              <div class="setting-group">
                <label class="setting-label">Otomatik Başlat:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="v2raySettings.autoStart">
                  <span class="slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">İnbound Ayarları</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">Varsayılan VMess Port:</label>
                <input 
                  type="number" 
                  v-model="v2raySettings.vmessPort" 
                  class="setting-input"
                  min="1000"
                  max="65535"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Varsayılan VLess Port:</label>
                <input 
                  type="number" 
                  v-model="v2raySettings.vlessPort" 
                  class="setting-input"
                  min="1000"
                  max="65535"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Varsayılan Trojan Port:</label>
                <input 
                  type="number" 
                  v-model="v2raySettings.trojanPort" 
                  class="setting-input"
                  min="1000"
                  max="65535"
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security Settings -->
      <div v-show="activeTab === 'security'" class="settings-section">
        <h2 class="section-title">Güvenlik Ayarları</h2>
        
        <div class="settings-grid">
          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Erişim Kontrolü</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">Başarısız Giriş Limiti:</label>
                <input 
                  type="number" 
                  v-model="securitySettings.loginAttemptLimit" 
                  class="setting-input"
                  min="1"
                  max="10"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Kilitleme Süresi (dakika):</label>
                <input 
                  type="number" 
                  v-model="securitySettings.lockoutDuration" 
                  class="setting-input"
                  min="1"
                  max="60"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Oturum Süresi (saat):</label>
                <input 
                  type="number" 
                  v-model="securitySettings.sessionTimeout" 
                  class="setting-input"
                  min="1"
                  max="24"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">IP Kısıtlaması:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="securitySettings.ipRestrictionEnabled">
                  <span class="slider"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">SSL/TLS Sertifikaları</h3>
            </div>
            <div class="setting-content">
              <div class="certificate-info">
                <div class="cert-status" :class="{ 'valid': certificateInfo.valid }">
                  <span class="status-indicator"></span>
                  <span class="status-text">
                    {{ certificateInfo.valid ? 'Sertifika Geçerli' : 'Sertifika Geçersiz' }}
                  </span>
                </div>
                <div class="cert-details">
                  <div class="cert-item">
                    <span class="cert-label">Son Kullanma:</span>
                    <span class="cert-value">{{ certificateInfo.expiryDate }}</span>
                  </div>
                  <div class="cert-item">
                    <span class="cert-label">Konu:</span>
                    <span class="cert-value">{{ certificateInfo.subject }}</span>
                  </div>
                </div>
              </div>
              <button class="btn-upload-cert">Sertifika Yükle</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Monitoring Settings -->
      <div v-show="activeTab === 'monitoring'" class="settings-section">
        <h2 class="section-title">İzleme Ayarları</h2>
        
        <div class="settings-grid">
          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Log Ayarları</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">Log Seviyesi:</label>
                <select v-model="monitoringSettings.logLevel" class="setting-select">
                  <option value="debug">Debug</option>
                  <option value="info">Info</option>
                  <option value="warning">Warning</option>
                  <option value="error">Error</option>
                </select>
              </div>
              <div class="setting-group">
                <label class="setting-label">Log Saklama Süresi (gün):</label>
                <input 
                  type="number" 
                  v-model="monitoringSettings.logRetentionDays" 
                  class="setting-input"
                  min="1"
                  max="90"
                >
              </div>
              <div class="setting-group">
                <label class="setting-label">Maksimum Log Boyutu (MB):</label>
                <input 
                  type="number" 
                  v-model="monitoringSettings.maxLogSize" 
                  class="setting-input"
                  min="1"
                  max="1000"
                >
              </div>
            </div>
          </div>

          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Performans İzleme</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">CPU İzleme:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="monitoringSettings.cpuMonitoring">
                  <span class="slider"></span>
                </label>
              </div>
              <div class="setting-group">
                <label class="setting-label">Bellek İzleme:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="monitoringSettings.memoryMonitoring">
                  <span class="slider"></span>
                </label>
              </div>
              <div class="setting-group">
                <label class="setting-label">Ağ İzleme:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="monitoringSettings.networkMonitoring">
                  <span class="slider"></span>
                </label>
              </div>
              <div class="setting-group">
                <label class="setting-label">İzleme Aralığı (saniye):</label>
                <input 
                  type="number" 
                  v-model="monitoringSettings.monitoringInterval" 
                  class="setting-input"
                  min="1"
                  max="300"
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Backup & Restore -->
      <div v-show="activeTab === 'backup'" class="settings-section">
        <h2 class="section-title">Yedekleme ve Geri Yükleme</h2>
        
        <div class="settings-grid">
          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Otomatik Yedekleme</h3>
            </div>
            <div class="setting-content">
              <div class="setting-group">
                <label class="setting-label">Otomatik Yedekleme:</label>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="backupSettings.autoBackupEnabled">
                  <span class="slider"></span>
                </label>
              </div>
              <div v-if="backupSettings.autoBackupEnabled" class="setting-group">
                <label class="setting-label">Yedekleme Sıklığı:</label>
                <select v-model="backupSettings.backupFrequency" class="setting-select">
                  <option value="daily">Günlük</option>
                  <option value="weekly">Haftalık</option>
                  <option value="monthly">Aylık</option>
                </select>
              </div>
              <div class="setting-group">
                <label class="setting-label">Maksimum Yedek Sayısı:</label>
                <input 
                  type="number" 
                  v-model="backupSettings.maxBackups" 
                  class="setting-input"
                  min="1"
                  max="10"
                >
              </div>
            </div>
          </div>

          <div class="setting-card">
            <div class="setting-header">
              <h3 class="setting-title">Mevcut Yedekler</h3>
            </div>
            <div class="setting-content">
              <div class="backup-list">
                <div 
                  v-for="backup in backupList" 
                  :key="backup.id"
                  class="backup-item"
                >
                  <div class="backup-info">
                    <div class="backup-name">{{ backup.name }}</div>
                    <div class="backup-date">{{ formatDate(backup.date) }}</div>
                    <div class="backup-size">{{ backup.size }}</div>
                  </div>
                  <div class="backup-actions">
                    <button class="btn-restore-backup" @click="restoreBackup(backup.id)">
                      Geri Yükle
                    </button>
                    <button class="btn-delete-backup" @click="deleteBackup(backup.id)">
                      Sil
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Save Button -->
    <div class="settings-footer">
      <button class="btn-save-settings" @click="saveSettings">
        <span class="save-icon">💾</span>
        Ayarları Kaydet
      </button>
    </div>

    <!-- Restore Modal -->
    <div v-if="showRestoreModal" class="modal-overlay" @click="closeRestoreModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Yedek Geri Yükleme</h3>
          <button class="modal-close" @click="closeRestoreModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="restore-form">
            <div class="form-group">
              <label class="form-label">Yedek Dosyası Seçin:</label>
              <input 
                type="file" 
                @change="handleBackupFile" 
                accept=".zip,.tar.gz"
                class="file-input"
              >
            </div>
            <div class="form-actions">
              <button class="btn-cancel" @click="closeRestoreModal">İptal</button>
              <button class="btn-restore-confirm" @click="confirmRestore">Geri Yükle</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

interface Backup {
  id: string
  name: string
  date: string
  size: string
}

// State
const activeTab = ref('general')
const showRestoreModal = ref(false)

// Settings tabs
const settingsTabs = [
  { id: 'general', name: 'Genel', icon: '⚙' },
  { id: 'v2ray', name: 'V2Ray', icon: '🚀' },
  { id: 'security', name: 'Güvenlik', icon: '🔒' },
  { id: 'monitoring', name: 'İzleme', icon: '📊' },
  { id: 'backup', name: 'Yedekleme', icon: '💾' }
]

// System info
const systemInfo = reactive({
  model: 'Keenetic Extra DSL KN2112',
  firmware: 'NDMS 3.7.8',
  k2rayVersion: '1.2.3',
  uptime: '7 gün 15 saat'
})

// Settings
const generalSettings = reactive({
  webPort: 8080,
  apiPort: 8081,
  httpsEnabled: true
})

const v2raySettings = reactive({
  corePath: '/opt/bin/v2ray',
  configDir: '/opt/etc/v2ray',
  logLevel: 'warning',
  autoStart: true,
  vmessPort: 1080,
  vlessPort: 1081,
  trojanPort: 1082
})

const securitySettings = reactive({
  loginAttemptLimit: 5,
  lockoutDuration: 15,
  sessionTimeout: 8,
  ipRestrictionEnabled: false
})

const monitoringSettings = reactive({
  logLevel: 'info',
  logRetentionDays: 30,
  maxLogSize: 100,
  cpuMonitoring: true,
  memoryMonitoring: true,
  networkMonitoring: true,
  monitoringInterval: 30
})

const backupSettings = reactive({
  autoBackupEnabled: true,
  backupFrequency: 'weekly',
  maxBackups: 5
})

// Certificate info
const certificateInfo = reactive({
  valid: true,
  expiryDate: '2025-10-07',
  subject: 'CN=keenetic.local'
})

// Backup list
const backupList = ref<Backup[]>([])

// Methods
const createBackup = async () => {
  console.log('Creating backup...')
  // API call to create backup
  const newBackup: Backup = {
    id: `backup-${Date.now()}`,
    name: `Backup-${new Date().toISOString().split('T')[0]}`,
    date: new Date().toISOString(),
    size: '2.5 MB'
  }
  backupList.value.unshift(newBackup)
}

const closeRestoreModal = () => {
  showRestoreModal.value = false
}

const handleBackupFile = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    console.log('Backup file selected:', target.files[0].name)
  }
}

const confirmRestore = () => {
  console.log('Restoring backup...')
  // API call to restore backup
  closeRestoreModal()
}

const restoreBackup = (backupId: string) => {
  if (confirm('Bu yedekten geri yüklemek istediğinizden emin misiniz?')) {
    console.log(`Restoring backup: ${backupId}`)
    // API call to restore specific backup
  }
}

const deleteBackup = (backupId: string) => {
  if (confirm('Bu yedeği silmek istediğinizden emin misiniz?')) {
    const index = backupList.value.findIndex(b => b.id === backupId)
    if (index !== -1) {
      backupList.value.splice(index, 1)
    }
  }
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString('tr-TR')
}

const saveSettings = async () => {
  console.log('Saving settings...')
  console.log('General:', generalSettings)
  console.log('V2Ray:', v2raySettings)
  console.log('Security:', securitySettings)
  console.log('Monitoring:', monitoringSettings)
  console.log('Backup:', backupSettings)
  // API call to save all settings
}

const loadBackupList = () => {
  // Mock backup data
  const mockBackups: Backup[] = [
    {
      id: 'backup-001',
      name: 'Backup-2024-10-07',
      date: '2024-10-07T10:30:00Z',
      size: '2.5 MB'
    },
    {
      id: 'backup-002',
      name: 'Backup-2024-10-01',
      date: '2024-10-01T10:30:00Z',
      size: '2.3 MB'
    },
    {
      id: 'backup-003',
      name: 'Backup-2024-09-24',
      date: '2024-09-24T10:30:00Z',
      size: '2.1 MB'
    }
  ]
  
  backupList.value = mockBackups
}

onMounted(() => {
  loadBackupList()
})
</script>

<style scoped>
.settings-page {
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

.btn-backup, .btn-restore {
  @apply flex items-center space-x-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-200 rounded-lg transition-colors duration-200;
}

.backup-icon, .restore-icon {
  @apply text-lg;
}

.settings-nav {
  @apply flex flex-wrap gap-2 bg-slate-800/50 border border-slate-700 rounded-lg p-2;
}

.nav-tab {
  @apply flex items-center space-x-2 px-4 py-2 text-slate-400 hover:text-slate-200 hover:bg-slate-700/50 rounded-lg transition-colors duration-200;
}

.nav-tab.active {
  @apply bg-cyan-600 text-white;
}

.tab-icon {
  @apply text-lg;
}

.tab-text {
  @apply font-medium;
}

.settings-content {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg p-6;
}

.settings-section {
  @apply space-y-6;
}

.section-title {
  @apply text-xl font-bold text-cyan-400 mb-4;
}

.settings-grid {
  @apply grid grid-cols-1 lg:grid-cols-2 gap-6;
}

.setting-card {
  @apply bg-slate-700/30 border border-slate-600 rounded-lg overflow-hidden;
}

.setting-header {
  @apply p-4 bg-slate-600/30 border-b border-slate-600;
}

.setting-title {
  @apply text-lg font-semibold text-slate-200;
}

.setting-content {
  @apply p-4 space-y-4;
}

.info-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-3;
}

.info-item {
  @apply flex justify-between;
}

.info-label {
  @apply text-sm text-slate-400;
}

.info-value {
  @apply text-sm text-slate-200 font-medium;
}

.setting-group {
  @apply flex flex-col space-y-2;
}

.setting-label {
  @apply text-sm font-semibold text-slate-300;
}

.setting-input, .setting-select {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
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

.certificate-info {
  @apply space-y-3;
}

.cert-status {
  @apply flex items-center space-x-2;
}

.cert-status.valid .status-indicator {
  @apply w-3 h-3 bg-green-400 rounded-full;
  box-shadow: 0 0 10px rgba(34, 197, 94, 0.5);
}

.cert-status:not(.valid) .status-indicator {
  @apply w-3 h-3 bg-red-400 rounded-full;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);
}

.status-text {
  @apply text-sm font-medium;
}

.cert-details {
  @apply space-y-2;
}

.cert-item {
  @apply flex justify-between;
}

.cert-label {
  @apply text-sm text-slate-400;
}

.cert-value {
  @apply text-sm text-slate-200 font-mono;
}

.btn-upload-cert {
  @apply w-full py-2 text-sm text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.backup-list {
  @apply space-y-3 max-h-64 overflow-y-auto;
}

.backup-item {
  @apply flex items-center justify-between p-3 bg-slate-800/50 border border-slate-600 rounded;
}

.backup-info {
  @apply flex-1 space-y-1;
}

.backup-name {
  @apply font-medium text-slate-200;
}

.backup-date {
  @apply text-sm text-slate-400;
}

.backup-size {
  @apply text-xs text-slate-500;
}

.backup-actions {
  @apply flex space-x-2;
}

.btn-restore-backup, .btn-delete-backup {
  @apply px-3 py-1 text-xs rounded transition-colors duration-200;
}

.btn-restore-backup {
  @apply text-cyan-400 border border-cyan-400 hover:bg-cyan-400 hover:text-slate-900;
}

.btn-delete-backup {
  @apply text-red-400 border border-red-400 hover:bg-red-400 hover:text-white;
}

.settings-footer {
  @apply flex justify-end;
}

.btn-save-settings {
  @apply flex items-center space-x-2 px-6 py-3 bg-cyan-600 hover:bg-cyan-500 text-white rounded-lg font-semibold transition-colors duration-200;
}

.save-icon {
  @apply text-lg;
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
  @apply p-4;
}

.restore-form {
  @apply space-y-4;
}

.form-group {
  @apply flex flex-col space-y-2;
}

.form-label {
  @apply text-sm font-semibold text-slate-300;
}

.file-input {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.form-actions {
  @apply flex space-x-2;
}

.btn-cancel {
  @apply flex-1 py-2 text-slate-400 border border-slate-600 rounded hover:bg-slate-700;
}

.btn-restore-confirm {
  @apply flex-1 py-2 bg-cyan-600 hover:bg-cyan-500 text-white rounded;
}

/* Responsive */
@media (max-width: 768px) {
  .settings-page {
    @apply p-4 space-y-4;
  }
  
  .page-header {
    @apply flex-col items-start space-y-2;
  }
  
  .header-actions {
    @apply w-full justify-end;
  }
  
  .settings-nav {
    @apply flex-col;
  }
  
  .nav-tab {
    @apply w-full justify-start;
  }
  
  .settings-grid {
    @apply grid-cols-1;
  }
  
  .info-grid {
    @apply grid-cols-1;
  }
}
</style>