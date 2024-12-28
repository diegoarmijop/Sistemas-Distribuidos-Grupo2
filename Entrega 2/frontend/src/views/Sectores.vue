<template>
  <SideBar page-title="Monitoreo de Campos y Plagas">
    <!-- Tabs para alternar entre vistas -->
    <v-tabs v-model="activeTab" class="mb-4">
      <v-tab value="campos">Campos</v-tab>
      <v-tab value="plagas">Tipos de Plaga</v-tab>
      <v-tab value="eventos">Eventos de Plaga</v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <!-- Vista de Campos -->
      <v-window-item value="campos">
        <v-row class="d-flex justify-center">
          <v-col v-for="campo in campos" :key="campo.id" cols="12" sm="6" lg="4">
            <v-card class="mx-auto" elevation="4" hover>
              <v-card-item>
                <div class="d-flex align-center">
                  <v-icon color="primary" size="36" class="me-4">mdi-farm</v-icon>
                  <div>
                    <v-card-title class="headline">{{ campo.nombre }}</v-card-title>
                    <v-card-subtitle class="text-body-1">
                      <div><strong>Superficie:</strong> {{ campo.superficie }} ha</div>
                      <div><strong>Tipo de cultivo:</strong> {{ campo.tipo_cultivo }}</div>
                      <div><strong>Ubicación:</strong> {{ campo.ubicacion }}</div>
                    </v-card-subtitle>
                  </div>
                </div>
                <v-divider class="my-3"></v-divider>
                <v-card-actions class="justify-space-between">
                  <v-btn text color="primary" @click="verDetalles(campo)">
                    Ver detalles
                  </v-btn>
                  <v-chip :color="campo.sensor_id ? 'success' : 'error'" text-color="white">
                    {{ campo.sensor_id ? 'Sensor activo' : 'Sin sensor' }}
                  </v-chip>
                </v-card-actions>
              </v-card-item>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>

      <!-- Vista de Tipos de Plaga -->
      <v-window-item value="plagas">
        <v-card>
          <v-data-table :headers="plagueHeaders" :items="plagueTypes" :loading="loading" hover>
            <template #[`item.actions`]="{ item }">
              <v-icon size="small" class="me-2" @click="editItem(item.raw)">
                mdi-pencil
              </v-icon>
              <v-icon size="small" color="error" @click="deleteItem(item.raw)">
                mdi-delete
              </v-icon>
            </template>
          </v-data-table>
        </v-card>
      </v-window-item>

      <!-- Vista de Eventos de Plaga -->
      <v-window-item value="eventos">
        <v-card>
          <v-data-table :headers="eventHeaders" :items="plagueEvents" :loading="loading" hover>
            <template #[`item.nivel_severidad`]="{ item }">
              <v-chip :color="getSeverityColor(item.nivel_severidad)">
                {{ item.nivel_severidad }}
              </v-chip>
            </template>
            <template #[`item.fecha_deteccion`]="{ item }">
              {{ formatDate(item.fecha_deteccion) }}
            </template>
            <template #[`item.tipo_plaga`]="{ item }">
              {{ item.tipo_plaga.nombre_comun }}
            </template>
            <template #[`item.campo`]="{ item }">
              {{ item.campo.nombre }}
            </template>
            <template #[`item.actions`]="{ item }">
              <v-icon size="small" class="me-2" @click="editItem(item)">
                mdi-pencil
              </v-icon>
              <v-icon size="small" color="error" @click="deleteItem(item)">
                mdi-delete
              </v-icon>
            </template>
          </v-data-table>
        </v-card>
      </v-window-item>
    </v-window>
  </SideBar>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import axios from 'axios'

export default {
  name: 'MonitoreoView',
  components: {
    SideBar
  },

  data() {
    return {
      activeTab: 'campos',
      loading: false,
      campos: [],
      plagueTypes: [],
      plagueEvents: [],

      // Headers para las tablas
      plagueHeaders: [
        { title: 'ID', key: 'tipo_plaga_id', align: 'start' },
        { title: 'Nombre Común', key: 'nombre_comun' },
        { title: 'Nombre Científico', key: 'nombre_cientifico' },
        { title: 'Descripción', key: 'descripcion' },
        { title: 'Acciones', key: 'actions', sortable: false }
      ],

      eventHeaders: [
      { title: 'Fecha', key: 'fecha_deteccion', align: 'start' },
      { title: 'Campo', key: 'campo' },
      { title: 'Tipo de Plaga', key: 'tipo_plaga' },
      { title: 'Ubicación', key: 'ubicacion' },
      { title: 'Severidad', key: 'nivel_severidad' },
      { title: 'Acciones Tomadas', key: 'acciones_tomadas' },
      { title: 'Acciones', key: 'actions', sortable: false }
    ]
    }
  },

  mounted() {
    this.loadData()
  },

  methods: {
    async loadData() {
      try {
        this.loading = true
        const [camposRes, plagasRes, eventosRes] = await Promise.all([
          axios.get('http://localhost:8080/api/campos/'),
          axios.get('http://localhost:8080/api/tipoPlaga/'),
          axios.get('http://localhost:8080/api/eventoPlagas/')
        ])

        this.campos = camposRes.data
        this.plagueTypes = plagasRes.data
        this.plagueEvents = eventosRes.data

      } catch (error) {
        console.error('Error cargando datos:', error)
      } finally {
        this.loading = false
      }
    },

    getSeverityColor(severity) {
      const colors = {
        'Alto': 'error',
        'Medio': 'warning',
        'Bajo': 'success'
      }
      return colors[severity] || 'grey'
    },

    formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
},

    verDetalles(campo) {
      console.log('Ver detalles del campo:', campo)
      // Implementar lógica para ver detalles
    },

    editItem(item) {
      console.log('Editar:', item)
      // Implementar lógica de edición
    },

    deleteItem(item) {
      console.log('Eliminar:', item)
      // Implementar lógica de eliminación
    }
  }
}
</script>

<style scoped>
.v-data-table {
  margin-top: 16px;
}

.v-card {
  margin-bottom: 16px;
}
</style>