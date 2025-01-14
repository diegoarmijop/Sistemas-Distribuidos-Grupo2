<template>
  <v-app>
      <!-- Barra lateral -->
      <v-navigation-drawer
          v-model="drawer"
          app
          color="grey-lighten-4"
          class="elevation-3"
      >
          <v-list>
              <v-list-item
                  :title="getUserName"
                  :subtitle="getUserEmail"
              >
                  <template v-slot:prepend>
                      <v-icon size="40" color="grey">mdi-account-circle</v-icon>
                  </template>
              </v-list-item>
          </v-list>

          <v-divider></v-divider>

          <v-list density="compact" nav>
              <v-list-item
                  v-for="item in menuItems"
                  :key="item.path"
                  :value="item"
                  :prepend-icon="item.icon"
                  :title="item.title"
                  :to="item.path"
                  :active="currentPath === item.path"
                  active-color="primary"
              ></v-list-item>
          </v-list>
      </v-navigation-drawer>

      <!-- Barra superior -->
      <v-app-bar app color="white" elevation="1">
          <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
          <v-toolbar-title class="text-h6 font-weight-bold">
              {{ pageTitle }}
          </v-toolbar-title>
          <v-spacer></v-spacer>

          <!-- Menú de notificaciones -->
          <v-menu
              v-model="showNotifications"
              :close-on-content-click="false"
              location="bottom end"
              max-width="300"
          >
              <template v-slot:activator="{ props }">
                  <v-btn icon class="mr-2" v-bind="props">
                      <v-badge
                          :content="unreadNotifications.length"
                          :model-value="unreadNotifications.length > 0"
                          color="error"
                      >
                          <v-icon>mdi-bell</v-icon>
                      </v-badge>
                  </v-btn>
              </template>

              <v-card>
                  <v-card-title class="d-flex align-center py-2">
                      Notificaciones
                      <v-spacer></v-spacer>
                      <v-btn
                          v-if="unreadNotifications.length"
                          variant="text"
                          density="compact"
                          @click="markAllAsRead"
                      >
                          Marcar todo como leído
                      </v-btn>
                  </v-card-title>
                  <v-divider></v-divider>
                  
                  <v-list v-if="notifications.length">
                    <template v-for="notification in notifications" :key="notification.id">                          <v-list-item
                              :class="{ 'unread-notification': !notification.read }"
                              @click="handleNotificationClick(notification)"
                          >
                              <template v-slot:prepend>
                                  <v-icon :color="getNotificationColor(notification.estado)">
                                      {{ getNotificationIcon(notification.estado) }}
                                  </v-icon>
                              </template>
                              
                              <v-list-item-title class="text-subtitle-2">
                                  {{ notification.estado }}
                              </v-list-item-title>
                              <v-list-item-subtitle class="text-caption">
                                  {{ notification.descripcion }}
                              </v-list-item-subtitle>
                              <v-list-item-subtitle class="text-caption text-grey">
                                  {{ formatDate(notification.fecha_hora) }}
                              </v-list-item-subtitle>
                          </v-list-item>
                          <v-divider></v-divider>
                      </template>
                  </v-list>
                  
                  <v-card-text v-else class="text-center py-5">
                      <v-icon size="40" color="grey">mdi-bell-off</v-icon>
                      <div class="mt-2">No hay notificaciones</div>
                  </v-card-text>
              </v-card>
          </v-menu>

          <!-- Menú de usuario -->
          <v-menu min-width="200px" rounded>
              <template v-slot:activator="{ props }">
                  <v-btn icon v-bind="props">
                      <v-icon size="32">mdi-account-circle</v-icon>
                  </v-btn>
              </template>
              <v-card>
                  <v-card-text>
                      <div class="d-flex align-center mb-3">
                          <v-icon size="40" color="grey" class="mr-3">mdi-account-circle</v-icon>
                          <div>
                              <div class="text-subtitle-1">{{ getUserName }}</div>
                              <div class="text-caption">{{ getUserEmail }}</div>
                          </div>
                      </div>
                      <v-divider class="mb-3"></v-divider>
                      <v-list density="compact" nav>
                          <v-list-item
                              prepend-icon="mdi-account-cog"
                              title="Configuración"
                              to="/settings"
                          ></v-list-item>
                          <v-list-item
                              prepend-icon="mdi-logout"
                              title="Cerrar sesión"
                              @click="logout"
                          ></v-list-item>
                      </v-list>
                  </v-card-text>
              </v-card>
          </v-menu>
      </v-app-bar>

      <!-- Contenido principal -->
      <v-main class="grey-lighten-3">
          <v-container fluid>
              <slot></slot>
          </v-container>
      </v-main>

      <!-- Snackbar para notificaciones -->
      <v-snackbar
          v-model="snackbar.show"
          :color="snackbar.color"
          :timeout="3000"
      >
          {{ snackbar.text }}
          <template v-slot:actions>
              <v-btn
                  color="white"
                  variant="text"
                  @click="snackbar.show = false"
              >
                  Cerrar
              </v-btn>
          </template>
      </v-snackbar>
  </v-app>
</template>

<script>
import axios from 'axios'

export default {
  name: 'SideBar',
  props: {
      pageTitle: {
          type: String,
          default: 'Dashboard'
      }
  },
  data: () => ({
      drawer: true,
      userData: null,
      currentPath: '',
      showNotifications: false,
      notifications: [],
      notificationPolling: null,
      snackbar: {
          show: false,
          text: '',
          color: 'success'
      },
      menuItems: [
          { title: 'Dashboard', icon: 'mdi-view-dashboard', path: '/dashboard' },
          { title: 'Gestión de Alertas', icon: 'mdi-alert', path: '/gestion-alertas' },
          { title: 'Sectores', icon: 'mdi-seed', path: '/sectors' },
          { title: 'Reportes', icon: 'mdi-chart-bar', path: '/reportes' },
          { title: 'Usuarios', icon: 'mdi-account-group', path: '/users' },
          { title: 'Configuración', icon: 'mdi-cog', path: '/configuracion' }
      ]
  }),

  computed: {
      getUserName() {
          return this.userData?.name || 'Usuario'
      },
      getUserEmail() {
          return this.userData?.email || 'usuario@email.com'
      },
      unreadNotifications() {
          return this.notifications.filter(n => !n.read)
      }
  },

  mounted() {
      document.addEventListener('visibilitychange', this.handleVisibilityChange)
  },

  async created() {
      this.currentPath = this.$route.path
      await this.loadUserData()
      this.startNotificationPolling()
  },

  beforeUnmount() {
      this.stopNotificationPolling()
      document.removeEventListener('visibilitychange', this.handleVisibilityChange)
  },

  watch: {
      '$route.path'(newPath) {
          this.currentPath = newPath
      }
  },

  methods: {
      handleVisibilityChange() {
          if (document.visibilityState === 'visible') {
              this.fetchNotifications()
          }
      },

      async loadUserData() {
          const userStr = localStorage.getItem('userData') || sessionStorage.getItem('userData')
          if (userStr) {
              try {
                  this.userData = JSON.parse(userStr)
              } catch (e) {
                  console.error('Error parsing user data:', e)
              }
          }
      },

      startNotificationPolling() {
          this.fetchNotifications()
          this.notificationPolling = setInterval(() => {
              this.fetchNotifications()
          }, 5000)
      },

      stopNotificationPolling() {
          if (this.notificationPolling) {
              clearInterval(this.notificationPolling)
          }
      },

      async fetchNotifications() {
        try {
        const response = await axios.get('http://localhost:8080/api/alertas/')
        let newAlerts = response.data
        
        // Ordenar las alertas más recientes primero
        newAlerts.sort((a, b) => new Date(b.fecha_hora) - new Date(a.fecha_hora))
        
        // Filtrar alertas que no sean 'Activa'
        const filteredAlerts = newAlerts.filter(alert => alert.estado !== 'Activa')
        
        // Encontrar nuevas notificaciones
        const existingIds = new Set(this.notifications.map(n => n.id))
        const newNotifications = filteredAlerts.filter(alert => !existingIds.has(alert.id))
        
        if (newNotifications.length > 0) {
            // Preparar las nuevas notificaciones
            const notificationsToAdd = newNotifications.map(n => ({
                ...n,
                read: false
            }))
            
            // Agregar las nuevas notificaciones al inicio
            this.notifications = [...notificationsToAdd, ...this.notifications]
            
            // Reordenar todas las notificaciones por fecha
            this.notifications.sort((a, b) => new Date(b.fecha_hora) - new Date(a.fecha_hora))
            
            // Limitar a un número máximo de notificaciones
            if (this.notifications.length > 50) {
                this.notifications = this.notifications.slice(0, 50)
            }
                  
                  // Mostrar notificación del sistema
                  if (Notification.permission === 'granted') {
                      new Notification('Nueva Alerta', {
                          body: `Tienes ${newNotifications.length} nueva(s) alerta(s)`
                      })
                  }

                  // Emitir evento para otros componentes
                  this.$root.$emit('newNotifications', notificationsToAdd)
              }

              // Actualizar el estado de las alertas existentes manteniendo el orden
        const updatedNotifications = this.notifications.map(notification => {
            const updatedAlert = filteredAlerts.find(a => a.id === notification.id)
            return updatedAlert ? { ...updatedAlert, read: notification.read } : notification
        })
        
        // Volver a ordenar después de la actualización
        this.notifications = updatedNotifications.sort((a, b) => new Date(b.fecha_hora) - new Date(a.fecha_hora))

    } catch (error) {
        console.error('Error fetching notifications:', error)
    }
      },

      getNotificationColor(estado) {
          const colors = {
              'Seria': 'warning',
              'Crítico': 'error',
              'Extremo Peligro': 'deep-purple-darken-4'
          }
          return colors[estado] || 'grey'
      },

      getNotificationIcon(estado) {
          const icons = {
              'Seria': 'mdi-alert',
              'Crítico': 'mdi-alert-octagon',
              'Extremo Peligro': 'mdi-alert-decagram'
          }
          return icons[estado] || 'mdi-bell'
      },

      handleNotificationClick(notification) {
          if (!notification.read) {
              notification.read = true
          }
          if (this.$route.path !== '/alertas') {
              this.$router.push('/alertas')
          }
          this.showNotifications = false
      },

      markAllAsRead() {
          this.notifications.forEach(n => n.read = true)
      },

      formatDate(date) {
          return new Date(date).toLocaleString('es-ES', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit'
          })
      },

      logout() {
          localStorage.removeItem('userData')
          sessionStorage.removeItem('userData')
          this.stopNotificationPolling()
          this.$router.push('/login')
      }
  }
}
</script>

<style scoped>
.unread-notification {
  background-color: rgba(var(--v-theme-primary), 0.05);
}

.v-list-item:hover {
  background-color: rgba(var(--v-theme-primary), 0.1);
}

.notification-badge {
  position: absolute;
  top: -8px;
  right: -8px;
}

.v-navigation-drawer {
  transition: 0.2s ease-in-out;
}
</style>