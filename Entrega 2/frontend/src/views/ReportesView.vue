<!-- src/views/ReportesView.vue -->
<template>
    <SideBar page-title="Reportes y Estadísticas">
      <!-- Filtros -->
      <v-card class="mb-4">
        <v-card-text>
          <v-row>
            <v-col cols="12" md="3">
              <v-select
                v-model="selectedCampo"
                :items="campos"
                label="Campo"
                item-title="nombre"
                item-value="id"
              ></v-select>
            </v-col>
            <v-col cols="12" md="3">
              <v-select
                v-model="selectedMetric"
                :items="metrics"
                label="Métrica"
              ></v-select>
            </v-col>
            <v-col cols="12" md="3">
              <v-menu
                ref="menuFechaInicio"
                v-model="menuFechaInicio"
                :close-on-content-click="false"
                transition="scale-transition"
              >
                <template v-slot:activator="{ props }">
                  <v-text-field
                    v-model="fechaInicio"
                    label="Fecha Inicio"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-bind="props"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="fechaInicio"
                  @click="menuFechaInicio = false"
                ></v-date-picker>
              </v-menu>
            </v-col>
            <v-col cols="12" md="3">
              <v-menu
                ref="menuFechaFin"
                v-model="menuFechaFin"
                :close-on-content-click="false"
                transition="scale-transition"
              >
                <template v-slot:activator="{ props }">
                  <v-text-field
                    v-model="fechaFin"
                    label="Fecha Fin"
                    prepend-icon="mdi-calendar"
                    readonly
                    v-bind="props"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-model="fechaFin"
                  @click="menuFechaFin = false"
                ></v-date-picker>
              </v-menu>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
  
      <!-- Gráficos -->
      <v-row>
        <v-col cols="12" md="8">
          <v-card>
            <v-card-title>Evolución Temporal</v-card-title>
            <v-card-text>
              <v-chart class="chart" :option="chartOption" autoresize />
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" md="4">
          <v-card>
            <v-card-title>Resumen Estadístico</v-card-title>
            <v-card-text>
              <v-list>
                <v-list-item>
                  <v-list-item-title>Promedio</v-list-item-title>
                  <v-list-item-subtitle>{{ stats.promedio }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Máximo</v-list-item-title>
                  <v-list-item-subtitle>{{ stats.maximo }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Mínimo</v-list-item-title>
                  <v-list-item-subtitle>{{ stats.minimo }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Desviación Estándar</v-list-item-title>
                  <v-list-item-subtitle>{{ stats.desviacion }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
  
      <!-- Tabla de Datos -->
      <v-card class="mt-4">
        <v-card-title class="d-flex align-center">
          Datos Detallados
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            prepend-icon="mdi-download"
            @click="downloadReport"
          >
            Exportar
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-data-table
            :headers="headers"
            :items="mediciones"
            :loading="loading"
          >
            <template #[`item.fecha`]="{ item }">
              {{ formatDate(item.fecha) }}
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </SideBar>
  </template>
  
  <script>
  import SideBar from '@/components/SideBar.vue'
  import { use } from 'echarts/core'
  import { CanvasRenderer } from 'echarts/renderers'
  import { LineChart } from 'echarts/charts'
  import {
    GridComponent,
    TooltipComponent,
    LegendComponent,
  } from 'echarts/components'
  import VChart from 'vue-echarts'
  
  use([
    CanvasRenderer,
    LineChart,
    GridComponent,
    TooltipComponent,
    LegendComponent,
  ])
  
  export default {
    name: 'ReportesView',
    components: {
      SideBar,
      VChart
    },
    
    data() {
      return {
        loading: false,
        selectedCampo: null,
        selectedMetric: null,
        fechaInicio: '2024-01-01',
        fechaFin: '2024-01-31',
        menuFechaInicio: false,
        menuFechaFin: false,
  
        campos: [
          { id: 1, nombre: 'Campo Norte' },
          { id: 2, nombre: 'Campo Sur' }
        ],
        metrics: [
          'Temperatura',
          'Humedad',
          'Luminosidad'
        ],
  
        stats: {
          promedio: '25.4°C',
          maximo: '32.1°C',
          minimo: '18.3°C',
          desviacion: '2.8°C'
        },
  
        headers: [
          { title: 'Fecha', key: 'fecha' },
          { title: 'Valor', key: 'valor' },
          { title: 'Umbral', key: 'umbral' },
          { title: 'Estado', key: 'estado' }
        ],
  
        mediciones: [
          {
            fecha: '2024-01-15T10:00:00',
            valor: '26.5°C',
            umbral: '30.0°C',
            estado: 'Normal'
          },
          {
            fecha: '2024-01-15T11:00:00',
            valor: '31.2°C',
            umbral: '30.0°C',
            estado: 'Alerta'
          }
        ],
  
        chartOption: {
          tooltip: {
            trigger: 'axis'
          },
          xAxis: {
            type: 'category',
            data: ['10:00', '11:00', '12:00', '13:00', '14:00']
          },
          yAxis: {
            type: 'value'
          },
          series: [{
            data: [26.5, 31.2, 28.7, 27.3, 29.1],
            type: 'line',
            smooth: true
          }]
        }
      }
    },
  
    methods: {
      formatDate(date) {
        return new Date(date).toLocaleString('es-ES')
      },
  
      downloadReport() {
        // Simular descarga
        this.$root.$emit('showMessage', {
          text: 'Reporte descargado exitosamente',
          color: 'success'
        })
      }
    }
  }
  </script>
  
  <style scoped>
  .chart {
    height: 400px;
  }
  </style>