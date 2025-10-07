<template>
  <div class="profile-page">
    <div class="page-header">
      <h1 class="page-title">Kullanıcı Profili</h1>
      <div class="profile-actions">
        <button class="btn-edit" @click="toggleEditMode">
          <span class="edit-icon">✏</span>
          {{ isEditMode ? 'İptal' : 'Düzenle' }}
        </button>
        <button v-if="isEditMode" class="btn-save" @click="saveProfile">
          <span class="save-icon">💾</span>
          Kaydet
        </button>
      </div>
    </div>

    <div class="profile-content">
      <!-- Profile Information -->
      <div class="profile-section">
        <div class="section-header">
          <h2 class="section-title">Kişisel Bilgiler</h2>
        </div>
        
        <div class="profile-card">
          <div class="profile-avatar">
            <div class="avatar-circle">
              <span class="avatar-text">{{ avatarText }}</span>
            </div>
            <button v-if="isEditMode" class="change-avatar-btn">
              Avatar Değiştir
            </button>
          </div>
          
          <div class="profile-info">
            <div class="info-grid">
              <div class="info-item">
                <label class="info-label">Kullanıcı Adı:</label>
                <input 
                  v-if="isEditMode"
                  type="text" 
                  v-model="profileData.username" 
                  class="info-input"
                >
                <span v-else class="info-value">{{ profileData.username }}</span>
              </div>
              
              <div class="info-item">
                <label class="info-label">E-posta:</label>
                <input 
                  v-if="isEditMode"
                  type="email" 
                  v-model="profileData.email" 
                  class="info-input"
                >
                <span v-else class="info-value">{{ profileData.email }}</span>
              </div>
              
              <div class="info-item">
                <label class="info-label">Ad Soyad:</label>
                <input 
                  v-if="isEditMode"
                  type="text" 
                  v-model="profileData.fullName" 
                  class="info-input"
                >
                <span v-else class="info-value">{{ profileData.fullName }}</span>
              </div>
              
              <div class="info-item">
                <label class="info-label">Rol:</label>
                <span class="info-value role-badge">{{ profileData.role }}</span>
              </div>
              
              <div class="info-item">
                <label class="info-label">Son Giriş:</label>
                <span class="info-value">{{ formatDate(profileData.lastLogin) }}</span>
              </div>
              
              <div class="info-item">
                <label class="info-label">Hesap Oluşturma:</label>
                <span class="info-value">{{ formatDate(profileData.createdAt) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security Settings -->
      <div class="profile-section">
        <div class="section-header">
          <h2 class="section-title">Güvenlik Ayarları</h2>
        </div>
        
        <div class="security-card">
          <!-- Password Change -->
          <div class="security-item">
            <div class="security-info">
              <h3 class="security-title">Şifre Değiştir</h3>
              <p class="security-desc">Hesabınızın güvenliği için düzenli olarak şifrenizi değiştirin</p>
            </div>
            <button class="btn-change-password" @click="showPasswordModal = true">
              Şifre Değiştir
            </button>
          </div>
          
          <!-- Two-Factor Authentication -->
          <div class="security-item">
            <div class="security-info">
              <h3 class="security-title">İki Faktörlü Doğrulama</h3>
              <p class="security-desc">Hesabınız için ek güvenlik katmanı ekleyin</p>
            </div>
            <div class="security-toggle">
              <label class="toggle-switch">
                <input 
                  type="checkbox" 
                  v-model="securitySettings.twoFactorEnabled" 
                  @change="toggle2FA"
                >
                <span class="slider"></span>
              </label>
              <span class="toggle-status">
                {{ securitySettings.twoFactorEnabled ? 'Aktif' : 'Pasif' }}
              </span>
            </div>
          </div>
          
          <!-- Session Management -->
          <div class="security-item">
            <div class="security-info">
              <h3 class="security-title">Aktif Oturumlar</h3>
              <p class="security-desc">{{ activeSessions.length }} aktif oturum bulunuyor</p>
            </div>
            <button class="btn-manage-sessions" @click="showSessionsModal = true">
              Oturumları Yönet
            </button>
          </div>
        </div>
      </div>

      <!-- Preferences -->
      <div class="profile-section">
        <div class="section-header">
          <h2 class="section-title">Tercihler</h2>
        </div>
        
        <div class="preferences-card">
          <div class="preference-item">
            <label class="preference-label">Dil:</label>
            <select v-model="preferences.language" class="preference-select">
              <option value="tr">Türkçe</option>
              <option value="en">English</option>
              <option value="de">Deutsch</option>
            </select>
          </div>
          
          <div class="preference-item">
            <label class="preference-label">Tema:</label>
            <select v-model="preferences.theme" class="preference-select">
              <option value="dark">Koyu</option>
              <option value="light">Açık</option>
              <option value="auto">Otomatik</option>
            </select>
          </div>
          
          <div class="preference-item">
            <label class="preference-label">Bildirimler:</label>
            <div class="notification-settings">
              <label class="notification-item">
                <input type="checkbox" v-model="preferences.notifications.email">
                <span class="notification-text">E-posta Bildirimleri</span>
              </label>
              <label class="notification-item">
                <input type="checkbox" v-model="preferences.notifications.push">
                <span class="notification-text">Push Bildirimleri</span>
              </label>
              <label class="notification-item">
                <input type="checkbox" v-model="preferences.notifications.security">
                <span class="notification-text">Güvenlik Uyarıları</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Password Change Modal -->
    <div v-if="showPasswordModal" class="modal-overlay" @click="closePasswordModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Şifre Değiştir</h3>
          <button class="modal-close" @click="closePasswordModal">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="changePassword" class="password-form">
            <div class="form-group">
              <label class="form-label">Mevcut Şifre:</label>
              <input 
                type="password" 
                v-model="passwordForm.current" 
                class="form-input"
                required
              >
            </div>
            <div class="form-group">
              <label class="form-label">Yeni Şifre:</label>
              <input 
                type="password" 
                v-model="passwordForm.new" 
                class="form-input"
                required
                minlength="8"
              >
            </div>
            <div class="form-group">
              <label class="form-label">Yeni Şifre Tekrar:</label>
              <input 
                type="password" 
                v-model="passwordForm.confirm" 
                class="form-input"
                required
              >
            </div>
            <div class="form-actions">
              <button type="button" class="btn-cancel" @click="closePasswordModal">
                İptal
              </button>
              <button type="submit" class="btn-confirm">
                Şifreyi Değiştir
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Sessions Modal -->
    <div v-if="showSessionsModal" class="modal-overlay" @click="closeSessionsModal">
      <div class="modal-content sessions-modal" @click.stop>
        <div class="modal-header">
          <h3>Aktif Oturumlar</h3>
          <button class="modal-close" @click="closeSessionsModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="sessions-list">
            <div 
              v-for="session in activeSessions" 
              :key="session.id"
              class="session-item"
              :class="{ 'current': session.current }"
            >
              <div class="session-info">
                <div class="session-device">
                  <span class="device-icon">{{ getDeviceIcon(session.device) }}</span>
                  <span class="device-name">{{ session.device }}</span>
                  <span v-if="session.current" class="current-badge">Mevcut</span>
                </div>
                <div class="session-details">
                  <div class="session-detail">
                    <span class="detail-label">IP:</span>
                    <span class="detail-value">{{ session.ip }}</span>
                  </div>
                  <div class="session-detail">
                    <span class="detail-label">Konum:</span>
                    <span class="detail-value">{{ session.location }}</span>
                  </div>
                  <div class="session-detail">
                    <span class="detail-label">Son Aktivite:</span>
                    <span class="detail-value">{{ formatDate(session.lastActivity) }}</span>
                  </div>
                </div>
              </div>
              <div class="session-actions">
                <button 
                  v-if="!session.current"
                  class="btn-terminate"
                  @click="terminateSession(session.id)"
                >
                  Sonlandır
                </button>
              </div>
            </div>
          </div>
          <div class="sessions-actions">
            <button class="btn-terminate-all" @click="terminateAllSessions">
              Diğer Oturumları Sonlandır
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'

interface Session {
  id: string
  device: string
  ip: string
  location: string
  lastActivity: string
  current: boolean
}

// State
const isEditMode = ref(false)
const showPasswordModal = ref(false)
const showSessionsModal = ref(false)

// Profile data
const profileData = reactive({
  username: 'admin',
  email: 'admin@keenetic.local',
  fullName: 'Sistem Yöneticisi',
  role: 'Yönetici',
  lastLogin: '2024-10-07T15:30:00Z',
  createdAt: '2024-01-15T10:00:00Z'
})

// Security settings
const securitySettings = reactive({
  twoFactorEnabled: false
})

// Preferences
const preferences = reactive({
  language: 'tr',
  theme: 'dark',
  notifications: {
    email: true,
    push: false,
    security: true
  }
})

// Password form
const passwordForm = reactive({
  current: '',
  new: '',
  confirm: ''
})

// Active sessions
const activeSessions = ref<Session[]>([])

// Computed
const avatarText = computed(() => {
  const names = profileData.fullName.split(' ')
  return names.map(name => name[0]).join('').toUpperCase()
})

// Methods
const toggleEditMode = () => {
  isEditMode.value = !isEditMode.value
}

const saveProfile = async () => {
  console.log('Saving profile:', profileData)
  // API call to save profile
  isEditMode.value = false
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleString('tr-TR')
}

const toggle2FA = () => {
  console.log('2FA toggled:', securitySettings.twoFactorEnabled)
  // API call to enable/disable 2FA
}

const closePasswordModal = () => {
  showPasswordModal.value = false
  // Reset form
  Object.assign(passwordForm, {
    current: '',
    new: '',
    confirm: ''
  })
}

const changePassword = async () => {
  if (passwordForm.new !== passwordForm.confirm) {
    alert('Yeni şifreler eşleşmiyor!')
    return
  }
  
  if (passwordForm.new.length < 8) {
    alert('Yeni şifre en az 8 karakter olmalıdır!')
    return
  }
  
  console.log('Changing password...')
  // API call to change password
  
  closePasswordModal()
}

const closeSessionsModal = () => {
  showSessionsModal.value = false
}

const getDeviceIcon = (device: string): string => {
  if (device.includes('Mobile')) return '📱'
  if (device.includes('Windows')) return '🖥'
  if (device.includes('Mac')) return '🖥'
  if (device.includes('Linux')) return '🐧'
  return '💻'
}

const terminateSession = async (sessionId: string) => {
  console.log(`Terminating session: ${sessionId}`)
  // API call to terminate session
  const index = activeSessions.value.findIndex(s => s.id === sessionId)
  if (index !== -1) {
    activeSessions.value.splice(index, 1)
  }
}

const terminateAllSessions = async () => {
  if (confirm('Diğer tüm oturumları sonlandırmak istediğinizden emin misiniz?')) {
    console.log('Terminating all other sessions...')
    // API call to terminate all sessions except current
    activeSessions.value = activeSessions.value.filter(s => s.current)
  }
}

const loadActiveSessions = () => {
  // Mock sessions data
  const mockSessions: Session[] = [
    {
      id: 'session-1',
      device: 'Windows 11 - Chrome',
      ip: '192.168.1.100',
      location: 'İstanbul, Türkiye',
      lastActivity: '2024-10-07T15:30:00Z',
      current: true
    },
    {
      id: 'session-2',
      device: 'Android Mobile - Chrome',
      ip: '192.168.1.105',
      location: 'İstanbul, Türkiye',
      lastActivity: '2024-10-07T14:45:00Z',
      current: false
    },
    {
      id: 'session-3',
      device: 'MacOS - Safari',
      ip: '192.168.1.110',
      location: 'İstanbul, Türkiye',
      lastActivity: '2024-10-07T12:20:00Z',
      current: false
    }
  ]
  
  activeSessions.value = mockSessions
}

onMounted(() => {
  loadActiveSessions()
})
</script>

<style scoped>
.profile-page {
  @apply p-6 space-y-6;
}

.page-header {
  @apply flex items-center justify-between;
}

.page-title {
  @apply text-2xl font-bold text-cyan-400;
}

.profile-actions {
  @apply flex space-x-2;
}

.btn-edit, .btn-save {
  @apply flex items-center space-x-2 px-4 py-2 bg-slate-700 hover:bg-slate-600 text-slate-200 rounded-lg transition-colors duration-200;
}

.btn-save {
  @apply bg-cyan-600 hover:bg-cyan-500 text-white;
}

.edit-icon, .save-icon {
  @apply text-lg;
}

.profile-content {
  @apply space-y-6;
}

.profile-section {
  @apply space-y-4;
}

.section-header {
  @apply border-b border-slate-700 pb-2;
}

.section-title {
  @apply text-lg font-semibold text-cyan-400;
}

.profile-card {
  @apply flex flex-col md:flex-row bg-slate-800/50 border border-slate-700 rounded-lg p-6 space-y-4 md:space-y-0 md:space-x-6;
}

.profile-avatar {
  @apply flex flex-col items-center space-y-3;
}

.avatar-circle {
  @apply w-24 h-24 bg-gradient-to-br from-cyan-500 to-blue-600 rounded-full flex items-center justify-center;
}

.avatar-text {
  @apply text-2xl font-bold text-white;
}

.change-avatar-btn {
  @apply text-sm text-cyan-400 hover:text-cyan-300;
}

.profile-info {
  @apply flex-1;
}

.info-grid {
  @apply grid grid-cols-1 md:grid-cols-2 gap-4;
}

.info-item {
  @apply flex flex-col space-y-1;
}

.info-label {
  @apply text-sm font-semibold text-slate-400;
}

.info-input {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.info-value {
  @apply text-sm text-slate-200;
}

.role-badge {
  @apply px-2 py-1 text-xs font-semibold bg-cyan-500/20 text-cyan-400 rounded-full;
}

.security-card {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg divide-y divide-slate-700;
}

.security-item {
  @apply flex items-center justify-between p-4;
}

.security-info {
  @apply flex-1;
}

.security-title {
  @apply text-base font-semibold text-slate-200;
}

.security-desc {
  @apply text-sm text-slate-400 mt-1;
}

.btn-change-password, .btn-manage-sessions {
  @apply px-4 py-2 text-sm text-cyan-400 border border-cyan-400 rounded hover:bg-cyan-400 hover:text-slate-900 transition-colors duration-200;
}

.security-toggle {
  @apply flex items-center space-x-3;
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

.toggle-status {
  @apply text-sm text-slate-300;
}

.preferences-card {
  @apply bg-slate-800/50 border border-slate-700 rounded-lg p-4 space-y-4;
}

.preference-item {
  @apply flex flex-col space-y-2;
}

.preference-label {
  @apply text-sm font-semibold text-slate-300;
}

.preference-select {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.notification-settings {
  @apply space-y-2;
}

.notification-item {
  @apply flex items-center space-x-2 cursor-pointer;
}

.notification-text {
  @apply text-sm text-slate-300;
}

/* Modal Styles */
.modal-overlay {
  @apply fixed inset-0 bg-black/50 flex items-center justify-center z-50;
}

.modal-content {
  @apply bg-slate-800 border border-slate-600 rounded-lg max-w-md w-full mx-4;
}

.sessions-modal {
  @apply max-w-2xl;
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

.password-form {
  @apply space-y-4;
}

.form-group {
  @apply flex flex-col space-y-1;
}

.form-label {
  @apply text-sm font-semibold text-slate-300;
}

.form-input {
  @apply px-3 py-2 bg-slate-700 border border-slate-600 text-slate-200 rounded focus:border-cyan-500 focus:outline-none;
}

.form-actions {
  @apply flex space-x-2 pt-4;
}

.btn-cancel {
  @apply flex-1 py-2 text-slate-400 border border-slate-600 rounded hover:bg-slate-700;
}

.btn-confirm {
  @apply flex-1 py-2 bg-cyan-600 hover:bg-cyan-500 text-white rounded;
}

.sessions-list {
  @apply space-y-3 max-h-64 overflow-y-auto;
}

.session-item {
  @apply flex items-center justify-between p-3 bg-slate-700/30 border border-slate-600 rounded;
}

.session-item.current {
  @apply border-cyan-500/50 bg-cyan-500/10;
}

.session-info {
  @apply flex-1 space-y-2;
}

.session-device {
  @apply flex items-center space-x-2;
}

.device-icon {
  @apply text-lg;
}

.device-name {
  @apply font-medium text-slate-200;
}

.current-badge {
  @apply px-2 py-0.5 text-xs font-semibold bg-cyan-500/20 text-cyan-400 rounded-full;
}

.session-details {
  @apply space-y-1;
}

.session-detail {
  @apply flex items-center space-x-2 text-sm;
}

.detail-label {
  @apply text-slate-400;
}

.detail-value {
  @apply text-slate-300;
}

.session-actions {
  @apply flex flex-col space-y-1;
}

.btn-terminate {
  @apply px-3 py-1 text-xs text-red-400 border border-red-400 rounded hover:bg-red-400 hover:text-white;
}

.sessions-actions {
  @apply pt-4 border-t border-slate-600 mt-4;
}

.btn-terminate-all {
  @apply w-full py-2 text-red-400 border border-red-400 rounded hover:bg-red-400 hover:text-white;
}

/* Responsive */
@media (max-width: 768px) {
  .profile-page {
    @apply p-4 space-y-4;
  }
  
  .page-header {
    @apply flex-col items-start space-y-2;
  }
  
  .profile-actions {
    @apply w-full justify-end;
  }
  
  .info-grid {
    @apply grid-cols-1;
  }
  
  .security-item {
    @apply flex-col items-start space-y-2;
  }
}
</style>