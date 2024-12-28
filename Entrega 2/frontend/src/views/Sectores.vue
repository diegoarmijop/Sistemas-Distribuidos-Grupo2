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
          <!-- Tarjetas de resumen de campos -->
          <v-row class="d-flex justify-center">
            <v-col
              v-for="campo in campos"
              :key="campo.id"
              cols="12"
              sm="6"
              lg="4"
            >
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
                    <v-btn text color="primary" @click="showDetails(campo)">
                      Ver detalles
                    </v-btn>
                    <v-chip
                      :color="campo.sensor && campo.sensor.id ? 'success' : 'error'"
                      text-color="white"
                    >
                      {{ campo.sensor && campo.sensor.id ? 'Sensor activo' : 'Sin sensor' }}
                    </v-chip>
                  </v-card-actions>
                </v-card-item>
              </v-card>
            </v-col>
          </v-row>
        </v-window-item>
  
        <!-- Vista de Tipos de Plaga -->
        <v-window-item value="plagas">
          <v-card class="mb-4">
            <v-card-title class="d-flex align-center">
              <span>Tipos de Plaga</span>
              <v-spacer></v-spacer>
              <v-btn color="primary" prepend-icon="mdi-plus" @click="openPlagueDialog">
                Nuevo Tipo
              </v-btn>
            </v-card-title>
  
            <v-card-text>
              <v-data-table
                :headers="plagueHeaders"
                :items="plagueTypes"
                :loading="loading"
                hover
              >
                <template #[`item.actions`]="{ item }">
                  <v-icon size="small" class="me-2" @click="editPlagueType(item.raw)">
                    mdi-pencil
                  </v-icon>
                  <v-icon size="small" color="error" @click="deletePlagueType(item.raw)">
                    mdi-delete
                  </v-icon>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-window-item>
  
        <!-- Vista de Eventos de Plaga -->
        <v-window-item value="eventos">
          <v-card class="mb-4">
            <v-card-title class="d-flex align-center">
              <span>Eventos de Plaga</span>
              <v-spacer></v-spacer>
              <v-btn color="primary" prepend-icon="mdi-plus" @click="openEventDialog">
                Nuevo Evento
              </v-btn>
            </v-card-title>
  
            <v-card-text>
              <v-data-table
                :headers="eventHeaders"
                :items="plagueEvents"
                :loading="loading"
                hover
              >
                <template #[`item.nivel_severidad`]="{ item }">
                  <v-chip :color="getSeverityColor(item.raw.nivel_severidad)">
                    {{ item.raw.nivel_severidad }}
                  </v-chip>
                </template>
                <template #[`item.fecha_deteccion`]="{ item }">
                  {{ formatDate(item.raw.fecha_deteccion) }}
                </template>
                <template #[`item.actions`]="{ item }">
                  <v-icon size="small" class="me-2" @click="editEvent(item.raw)">
                    mdi-pencil
                  </v-icon>
                  <v-icon size="small" color="error" @click="deleteEvent(item.raw)">
                    mdi-delete
                  </v-icon>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-window-item>
      </v-window>
  
      <!-- Gráfico y Actividades Recientes -->
      <v-row class="mt-6">
        <v-col cols="12" md="8">
          <v-card class="elevation-2">
            <v-card-title class="text-h6">
              Eventos de Plaga Recientes
              <v-spacer></v-spacer>
              <v-btn variant="text" color="primary" size="small">
                Ver más
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-chart class="chart" :option="chartOption" autoresize />
            </v-card-text>
          </v-card>
        </v-col>
  
        <v-col cols="12" md="4">
          <v-card class="elevation-2">
            <v-card-title class="text-h6">
              Últimas Actividades
            </v-card-title>
            <v-list lines="two">
              <v-list-item
                v-for="(activity, i) in recentActivities"
                :key="i"
                :title="activity.title"
                :subtitle="activity.time"
              >
                <template #prepend>
                  <v-avatar :color="activity.color" size="36">
                    <v-icon dark>{{ activity.icon }}</v-icon>
                  </v-avatar>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
  
      <!-- Diálogo unificado -->
        <PlagueDialog
            v-model="dialogVisible"
            :type="dialogType"
            :item="editedItem"
            :options="dialogOptions"
            @save="saveItem"
        />
    </SideBar>
  </template>
  
  <script>
import SideBar from '@/components/SideBar.vue';
import PlagueDialog from '@/components/PlagueDialog.vue';
import axios from 'axios';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, BarChart } from 'echarts/charts';
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
} from 'echarts/components';
import VChart from 'vue-echarts';

// Registrar componentes de ECharts
use([
  CanvasRenderer,
  LineChart,
  BarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
]);

export default {
  name: 'SectoresView',
  components: {
    SideBar,
    PlagueDialog,
    VChart
  },
  
  mounted() {
    this.loadData();
  },

  methods: {
    // Método para cargar datos
    async loadData() {
      try {
        this.loading = true;
        const [typesResponse, eventsResponse, camposResponse] = await Promise.all([
          axios.get('http://localhost:8080/api/tipoPlaga'),
          axios.get('http://localhost:8080/api/eventoPlagas'),
          axios.get('http://localhost:8080/api/campos')
        ]);

        this.plagueTypes = typesResponse.data;
        this.plagueEvents = eventsResponse.data;
        this.campos = camposResponse.data;

        // Actualizar opciones del diálogo
        this.dialogOptions = {
          plagueTypes: this.plagueTypes,
          campos: this.campos
        };

        // Actualizar gráfico y actividades
        this.updateChartData();
        this.updateRecentActivities();
      } catch (error) {
        console.error('Error cargando datos:', error);
      } finally {
        this.loading = false;
      }
    },

    // Método para guardar items
    async saveItem(item) {
      try {
        if (this.dialogType === 'type') {
          if (item.tipo_plaga_id) {
            await axios.put(`http://localhost:8080/api/tipoPlaga/${item.tipo_plaga_id}`, item);
          } else {
            await axios.post('http://localhost:8080/api/tipoPlaga', item);
          }
        } else {
          if (item.id) {
            await axios.put(`http://localhost:8080/api/eventoPlagas/${item.id}`, item);
          } else {
            await axios.post('http://localhost:8080/api/eventoPlagas', item);
          }
        }
        await this.loadData(); // Recargar datos
      } catch (error) {
        console.error('Error al guardar:', error);
      }
    },

    async deletePlagueType(item) {
        if (confirm('¿Está seguro de eliminar este tipo de plaga?')) {
            try {
            await axios.delete(`http://localhost:8080/api/tipoPlaga/${item.tipo_plaga_id}`);
            await this.loadData();
            } catch (error) {
            console.error('Error al eliminar tipo de plaga:', error);
            }
        }
        },

        async deleteEvent(item) {
        if (confirm('¿Está seguro de eliminar este evento?')) {
            try {
            await axios.delete(`http://localhost:8080/api/eventoPlagas/${item.id}`);
            await this.loadData();
            } catch (error) {
            console.error('Error al eliminar evento:', error);
            }
        }
        },

    // Método para actualizar el gráfico
    updateChartData() {
      const eventsByMonth = {};
      const months = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'];
      
      this.plagueEvents.forEach(event => {
        const date = new Date(event.fecha_deteccion);
        const month = months[date.getMonth()];
        eventsByMonth[month] = (eventsByMonth[month] || 0) + 1;
      });

      this.chartOption = {
        title: {
          text: 'Eventos de Plaga por Mes',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        xAxis: {
          type: 'category',
          data: months,
          axisLabel: {
            rotate: 45
          }
        },
        yAxis: {
          type: 'value',
          name: 'Número de Eventos'
        },
        series: [
          {
            name: 'Eventos',
            type: 'bar',
            data: months.map(month => eventsByMonth[month] || 0),
            itemStyle: {
              color: '#1976D2'
            }
          }
        ]
      };
    },

    // Método para actualizar actividades recientes
    updateRecentActivities() {
      this.recentActivities = this.plagueEvents
        .sort((a, b) => new Date(b.fecha_deteccion) - new Date(a.fecha_deteccion))
        .slice(0, 5)
        .map(event => ({
          title: `${event.tipo_plaga.nombre_comun} en ${event.campo.nombre}`,
          time: this.getRelativeTime(new Date(event.fecha_deteccion)),
          icon: 'mdi-bug-alert',
          color: this.getSeverityColor(event.nivel_severidad)
        }));
    },

    // Método para obtener tiempo relativo
    getRelativeTime(date) {
      const now = new Date();
      const diffTime = Math.abs(now - date);
      const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
      
      if (diffHours < 24) {
        return `Hace ${diffHours} horas`;
      } else {
        const diffDays = Math.floor(diffHours / 24);
        return `Hace ${diffDays} días`;
      }
    },

    // ... resto de tus métodos existentes ...
  }
};
</script>

<style scoped>
.chart {
  min-height: 400px;
  width: 100%;
}

.v-data-table {
  margin-top: 16px;
}

.v-card {
  margin-bottom: 16px;
}
</style>