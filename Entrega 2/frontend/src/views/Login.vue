<!-- src/views/Dashboard.vue -->
<template>
  <SideBar page-title="Panel de Alertas">
    <!-- Encabezado con contador de alertas -->
    <v-card class="mb-6" color="grey-darken-3">
      <v-card-title class="text-h5 white--text d-flex align-center">
        Sistema de Monitoreo de Alertas
        <v-chip class="ml-4" color="error" size="large">
          {{ latestAlerts.length }} Alertas Activas
        </v-chip>
        <v-spacer></v-spacer>
        <v-btn 
          color="white" 
          variant="outlined" 
          @click="refreshAlerts"
          :loading="loading"
        >
          Actualizar
          <v-icon right class="ml-2">mdi-refresh</v-icon>
        </v-btn>
      </v-card-title>
    </v-card>

    <!-- Lista de Alertas -->
    <v-row>
      <v-col cols="12">
        <!-- Loading Skeleton -->
        <template v-if="loading && !latestAlerts.length">
          <v-skeleton-loader
            v-for="n in 3"
            :key="n"
            class="mb-4"
            type="card"
          ></v-skeleton-loader>
        </template>

        <!-- Sin alertas -->
        <v-card v-else-if="latestAlerts.length === 0" class="text-center pa-6">
          <v-icon size="64" color="success">mdi-check-circle</v-icon>
          <div class="text-h6 mt-4">No hay alertas activas</div>
          <div class="text-body-1">El sistema está funcionando correctamente</div>
        </v-card>

        <!-- Lista de alertas -->
        <v-row v-else>
          <v-col v-for="alert in sortedAlerts" :key="alert.id" cols="12" md="6">
            <v-card :color="getAlertColor(alert.estado)" class="mx-2">
              <v-card-title class="text-white">
                <v-icon left color="white" class="mr-2">{{ getAlertIcon(alert.tipo_alerta) }}</v-icon>
                {{ alert.tipo_alerta }}
              </v-card-title>
              <v-card-text class="white bg-grey-lighten-4">
                <div class="text-body-1 font-weight-medium mb-2">{{ alert.descripcion }}</div>
                <v-chip :color="getAlertColor(alert.estado)" label class="mb-2">
                  {{ alert.estado }}
                </v-chip>
                <div class="text-caption mt-2">
                  <v-icon small class="mr-1">mdi-clock-outline</v-icon>
                  {{ formatDate(alert.fecha_hora) }}
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <!-- Snackbar para notificaciones -->
    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      :timeout="4000"
      location="top"
      style="top: 20px"
    >
      <div class="d-flex align-center">
        <v-icon :color="snackbarColor" class="mr-2">
          {{ snackbarIcon }}
        </v-icon>
        {{ snackbarText }}
      </div>
      <template v-slot:actions>
        <v-btn
          variant="text"
          @click="snackbar = false"
        >
          Cerrar
        </v-btn>
      </template>
    </v-snackbar>
  </SideBar>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import axios from 'axios'

export default {
  name: 'DashboardView',
  components: {
    SideBar,
  },
  data: () => ({
    latestAlerts: [],
    pollingInterval: null,
    loading: false,
    snackbar: false,
    snackbarText: '',
    snackbarColor: 'success',
    snackbarIcon: 'mdi-check-circle',
  }),
  computed: {
    sortedAlerts() {
      return [...this.latestAlerts].sort((a, b) => 
        new Date(b.fecha_hora) - new Date(a.fecha_hora)
      )
    }
  },
  methods: {
    showNotification(message, type = 'error') {
      this.snackbarColor = type
      this.snackbarIcon = type === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
      this.snackbarText = message
      this.snackbar = true
    },

    async fetchAlerts() {
      if (this.loading) return
      
      this.loading = true
      try {
        // Verificar si hay token
        const token = localStorage.getItem('token') || sessionStorage.getItem('token')
        if (!token) {
          this.$router.push('/login')
          return
        }

        const response = await axios.get('http://localhost:8080/api/alertas', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        this.latestAlerts = response.data
        // Solo mostrar notificación si es una actualización manual
        if (!this.pollingInterval) {
          this.showNotification('Alertas actualizadas correctamente', 'success')
        }

      } catch (error) {
        let errorMessage = 'Error al obtener las alertas'
        
        if (error.response?.status === 401) {
          this.$router.push('/login')
          errorMessage = 'Sesión expirada'
        } else if (error.response?.status === 404) {
          errorMessage = 'No se encontraron alertas'
        } else if (error.response?.data?.error) {
          errorMessage = error.response.data.error
        }
        
        this.showNotification(errorMessage, 'error')
      } finally {
        this.loading = false
      }
    },

    refreshAlerts() {
      // Detener el polling actual
      this.stopPolling()
      // Hacer la petición manual
      this.fetchAlerts()
      // Reiniciar el polling
      this.startPolling()
    },

    startPolling() {
      this.pollingInterval = setInterval(() => {
        this.fetchAlerts()
      }, 5000)
    },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval)
        this.pollingInterval = null
      }
    },

    formatDate(date) {
      return new Date(date).toLocaleString('es-ES', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    },

    getAlertColor(estado) {
      switch (estado.toLowerCase()) {
        case 'crítico':
          return 'error'
        case 'alto':
          return 'orange-darken-2'
        case 'medio':
          return 'warning'
        case 'bajo':
          return 'info'
        default:
          return 'grey'
      }
    },

    getAlertIcon(tipo) {
      switch (tipo.toLowerCase()) {
        case 'parámetros alterados':
          return 'mdi-alert-circle'
        case 'presencia de plagas':
          return 'mdi-bug'
        case 'error del sistema':
          return 'mdi-alert-octagon'
        default:
          return 'mdi-alert'
      }
    }
  },
  mounted() {
    this.fetchAlerts()
    this.startPolling()
  },
  beforeUnmount() {
    this.stopPolling()
  }
}
</script>

<style scoped>
.v-card {
  transition: all 0.3s ease-in-out;
}

.v-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 25px 0 rgba(0, 0, 0, 0.1);
}

/* Estilo para las tarjetas de alerta */
.v-card.error {
  border-left: 5px solid var(--v-error-base);
}
.v-card.warning {
  border-left: 5px solid var(--v-warning-base);
}
.v-card.info {
  border-left: 5px solid var(--v-info-base);
}
</style>