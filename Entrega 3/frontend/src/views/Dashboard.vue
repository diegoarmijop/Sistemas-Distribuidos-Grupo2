<!-- src/views/Dashboard.vue -->
<template>
  <SideBar page-title="Panel de Alertas">
    <v-card>
      <v-card-title class="d-flex align-center">
        Alertas del Sistema
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          prepend-icon="mdi-refresh"
          :loading="loading"
          @click="refreshAlerts"
        >
          Actualizar
        </v-btn>
      </v-card-title>

      <!-- Loading Skeleton -->
      <v-card-text v-if="loading && !latestAlerts.length">
        <v-skeleton-loader
          v-for="n in 3"
          :key="n"
          class="mb-4"
          type="card"
        ></v-skeleton-loader>
      </v-card-text>

      <!-- Sin alertas -->
      <v-card-text v-else-if="latestAlerts.length === 0">
        <v-alert
          type="info"
          variant="tonal"
          class="text-center"
        >
          <v-icon size="64" color="info">mdi-check-circle</v-icon>
          <div class="text-h6 mt-4">No hay alertas activas</div>
          <div class="text-body-1">El sistema está funcionando correctamente</div>
        </v-alert>
      </v-card-text>

      <!-- Lista de alertas -->
      <v-card-text v-else>
        <v-row>
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
      </v-card-text>
    </v-card>
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
  }),
  computed: {
    sortedAlerts() {
      return [...this.latestAlerts].sort((a, b) => 
        new Date(b.fecha_hora) - new Date(a.fecha_hora)
      )
    }
  },
  methods: {
    async fetchAlerts() {
      if (this.loading) return
      
      this.loading = true
      try {
        const token = localStorage.getItem('token') || sessionStorage.getItem('token')
        if (!token) {
          this.$router.push('/login')
          return
        }

        const response = await axios.get('http://localhost:8080/api/alertas/', {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          }
        })
        
        this.latestAlerts = response.data

      } catch (error) {
        console.error('Error al obtener alertas:', error)
        let message = 'Error al obtener alertas'
        
        if (error.response?.status === 401) {
          this.$router.push('/login')
          message = 'Sesión expirada'
        } else if (error.response?.data?.error) {
          message = error.response.data.error
        }
        
        this.$root.$emit('showMessage', {
          text: message,
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    },

    refreshAlerts() {
      this.stopPolling()
      this.fetchAlerts()
      this.startPolling()
    },

    startPolling() {
      if (!this.pollingInterval) {
        this.pollingInterval = setInterval(() => {
          this.fetchAlerts()
        }, 5000)
      }
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