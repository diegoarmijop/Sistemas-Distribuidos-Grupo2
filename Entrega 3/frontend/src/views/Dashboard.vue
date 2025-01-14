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
            <v-card :class="['alert-card', `status-${alert.estado.toLowerCase().replace(' ', '-')}`]">
              <v-card-title :class="getStatusColorClass(alert.estado)">
                <v-icon :color="getStatusIconColor(alert.estado)" class="mr-2">
                  {{ getStatusIcon(alert.estado) }}
                </v-icon>
                {{ alert.estado }}
              </v-card-title>
              
              <v-card-text class="pt-4">
                <!-- Tipos de Alerta -->
                <div class="mb-3">
                  <div class="d-flex flex-wrap gap-2">
                    <v-chip
                      v-for="(tipo, index) in splitTipoAlerta(alert.tipo_alerta)"
                      :key="index"
                      :color="getAlertTypeColor(tipo)"
                      size="small"
                      label
                      class="ma-1"
                    >
                      <v-icon start size="small">{{ getTypeIcon(tipo) }}</v-icon>
                      {{ tipo }}
                    </v-chip>
                  </div>
                </div>

                <!-- Descripción -->
                <div class="text-body-1 mb-3">{{ alert.descripcion }}</div>

                <!-- Fecha -->
                <div class="text-caption d-flex align-center">
                  <v-icon size="small" class="mr-1">mdi-clock-outline</v-icon>
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
        const response = await axios.get('http://localhost:8080/api/alertas/')
        this.latestAlerts = response.data
      } catch (error) {
        console.error('Error al obtener alertas:', error)
        this.$root.$emit('showMessage', {
          text: 'Error al obtener alertas',
          color: 'error'
        })
      } finally {
        this.loading = false
      }
    },

    splitTipoAlerta(tipo) {
      return tipo ? tipo.split('/').map(t => t.trim()) : []
    },

    getStatusColorClass(estado) {
      const colors = {
        'Seria': 'bg-warning-lighten-4',
        'Crítico': 'bg-error-lighten-4',
        'Extremo Peligro': 'bg-deep-purple-lighten-4',
        'Normal': 'bg-success-lighten-4'
      }
      return colors[estado] || 'bg-grey-lighten-4'
    },

    getStatusIconColor(estado) {
      const colors = {
        'Seria': 'warning-darken-2',
        'Crítico': 'error',
        'Extremo Peligro': 'deep-purple-darken-4',
        'Normal': 'success'
      }
      return colors[estado] || 'grey'
    },

    getStatusIcon(estado) {
      const icons = {
        'Seria': 'mdi-alert',
        'Crítico': 'mdi-alert-octagon',
        'Extremo Peligro': 'mdi-alert-decagram',
        'Normal': 'mdi-check-circle'
      }
      return icons[estado] || 'mdi-help-circle'
    },

    getAlertTypeColor(tipo) {
      const typeColors = {
        'Humedad baja-alta': 'blue-darken-2',
        'Temperatura alta': 'deep-orange',
        'Temperatura baja': 'light-blue',
        'Nivel alto de insectos': 'red-darken-4',
        'Nivel de luz extremadamente alto': 'amber-darken-4',
        'Nivel de luz alto': 'amber'
      }
      return typeColors[tipo] || 'grey'
    },

    getTypeIcon(tipo) {
      const icons = {
        'Humedad baja-alta': 'mdi-water',
        'Temperatura alta': 'mdi-thermometer-high',
        'Temperatura baja': 'mdi-thermometer-low',
        'Nivel alto de insectos': 'mdi-bug',
        'Nivel de luz extremadamente alto': 'mdi-white-balance-sunny',
        'Nivel de luz alto': 'mdi-brightness-6'
      }
      return icons[tipo] || 'mdi-alert'
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
.alert-card {
  transition: all 0.3s ease-in-out;
  border-radius: 8px;
  overflow: hidden;
}

.alert-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 25px 0 rgba(0, 0, 0, 0.1);
}

.gap-2 {
  gap: 8px;
}

/* Estados específicos */
.status-seria {
  border-left: 4px solid #FB8C00;
}

.status-crítico {
  border-left: 4px solid #D32F2F;
}

.status-extremo-peligro {
  border-left: 4px solid #311B92;
}

.status-normal {
  border-left: 4px solid #2E7D32;
}

/* Añadir animación para nuevas alertas */
@keyframes highlight {
  0% {
    background-color: rgba(255, 193, 7, 0.2);
  }
  100% {
    background-color: transparent;
  }
}

.new-alert {
  animation: highlight 2s ease-out;
}
</style>