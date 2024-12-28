<!-- src/views/Dashboard.vue -->
<template>
    <SideBar page-title="Dashboard">
      <!-- Tarjetas de resumen -->
      <v-row>
        <v-col
          v-for="card in summaryCards"
          :key="card.title"
          cols="12"
          sm="6"
          lg="3"
        >
          <v-card class="mx-auto" elevation="2" hover>
            <v-card-item>
              <template v-slot:prepend>
                <v-icon
                  :color="card.color"
                  size="large"
                  :icon="card.icon"
                ></v-icon>
              </template>
              <v-card-title>{{ card.title }}</v-card-title>
              <div class="text-primary text-h4 mt-2">
                {{ card.value }}
              </div>
              <v-card-subtitle class="mt-2">
                {{ card.subtitle }}
                <v-icon
                  :color="card.trend === 'up' ? 'success' : 'error'"
                  size="small"
                >
                  {{ card.trend === 'up' ? 'mdi-trending-up' : 'mdi-trending-down' }}
                </v-icon>
              </v-card-subtitle>
            </v-card-item>
          </v-card>
        </v-col>
      </v-row>
  
      <!-- Gráfico y tabla -->
      <v-row class="mt-6">
        <v-col cols="12" md="8">
          <v-card>
            <v-card-title class="text-h6">
              Actividad Reciente
              <v-spacer></v-spacer>
              <v-btn
                variant="text"
                color="primary"
                size="small"
              >
                Ver más
              </v-btn>
            </v-card-title>
            <v-card-text>
              <v-chart class="chart" :option="chartOption" autoresize />
            </v-card-text>
          </v-card>
        </v-col>
  
        <v-col cols="12" md="4">
          <v-card>
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
                <template v-slot:prepend>
                  <v-avatar :color="activity.color" size="36">
                    <v-icon dark>{{ activity.icon }}</v-icon>
                  </v-avatar>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </SideBar>
  </template>
  
  <script>
  import SideBar from '@/components/SideBar.vue'
  import { use } from "echarts/core"
  import { CanvasRenderer } from "echarts/renderers"
  import { LineChart } from "echarts/charts"
  import { GridComponent, TooltipComponent, LegendComponent } from "echarts/components"
  import VChart from "vue-echarts"
  
  use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])
  
  export default {
  name: 'DashboardView',
  components: {
    SideBar,
    VChart
  },
  data: () => ({
    summaryCards: [
      {
        title: 'Total Campos',
        value: '24',
        subtitle: '2 campos nuevos este mes',
        icon: 'mdi-terrain',  // Icono de terreno
        color: 'success',
        trend: 'up'
      },
      {
        title: 'Alertas de Plagas',
        value: '8',
        subtitle: '3 alertas menos esta semana',
        icon: 'mdi-alert-circle', // Círculo con signo de exclamación
        // O podrías usar alguna de estas alternativas:
        // icon: 'mdi-alert', // Triángulo con signo de exclamación
        // icon: 'mdi-alert-octagon', // Octágono con signo de exclamación
        // icon: 'mdi-alert-octagram', // Estrella con signo de exclamación
        // icon: 'mdi-alert-decagram', // Decágono con signo de exclamación
        // icon: 'mdi-warning', // Señal de advertencia
        color: 'error',
        trend: 'down'
        },
      {
        title: 'Cosechas Activas',
        value: '15',
        subtitle: '5% más que el mes anterior',
        icon: 'mdi-sprout',  // Icono de planta
        color: 'primary',
        trend: 'up'
      },
      {
        title: 'Control de Plagas',
        value: '92%',
        subtitle: 'Efectividad del tratamiento',
        icon: 'mdi-shield-check',  // Icono de escudo
        color: 'info',
        trend: 'up'
      }
    ],
    recentActivities: [
      {
        title: 'Alerta de plaga detectada en Campo Norte',
        time: 'Hace 2 horas',
        icon: 'mdi-alert',
        color: 'error'
      },
      {
        title: 'Fumigación completada en Campo Sur',
        time: 'Hace 5 horas',
        icon: 'mdi-check-circle',
        color: 'success'
      },
      {
        title: 'Nueva cosecha iniciada en Campo Este',
        time: 'Hace 1 día',
        icon: 'mdi-sprout',
        color: 'primary'
      },
      {
        title: 'Inspección programada para mañana',
        time: 'En 12 horas',
        icon: 'mdi-calendar-check',
        color: 'info'
      }
    ],
    chartOption: {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['Incidencias de Plagas', 'Efectividad del Control']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul']
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '{value}%'
        }
      },
      series: [
        {
          name: 'Incidencias de Plagas',
          type: 'line',
          smooth: true,
          color: '#ff4757',
          data: [25, 32, 28, 15, 12, 8, 5]
        },
        {
          name: 'Efectividad del Control',
          type: 'line',
          smooth: true,
          color: '#2ed573',
          data: [75, 78, 82, 89, 92, 95, 98]
        }
      ]
    }
  })
}
</script>
  
  <style scoped>
  .chart {
    height: 400px;
  }
  
  .v-card {
    transition: all 0.3s ease-in-out;
  }
  
  .v-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 4px 25px 0 rgba(0, 0, 0, 0.1);
  }
  </style>