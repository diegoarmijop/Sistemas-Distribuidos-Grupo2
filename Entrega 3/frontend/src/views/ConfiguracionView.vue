<!-- src/views/ConfiguracionView.vue -->
<template>
    <SideBar page-title="Configuración del Sistema">
      <v-card class="mb-4">
        <v-card-title class="d-flex align-center">
          Configuración de Umbrales
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            prepend-icon="mdi-content-save"
            @click="saveConfig"
            :loading="saving"
          >
            Guardar Cambios
          </v-btn>
        </v-card-title>
  
        <v-card-text>
          <v-form ref="form">
            <v-row>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="config.umbral_temp"
                  label="Umbral de Temperatura (°C)"
                  type="number"
                  step="0.1"
                  :rules="[v => !!v || 'Este campo es requerido']"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="config.umbral_humedad"
                  label="Umbral de Humedad (%)"
                  type="number"
                  step="0.1"
                  :rules="[v => !!v || 'Este campo es requerido']"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="config.umbral_luminosidad"
                  label="Umbral de Luminosidad (lux)"
                  type="number"
                  step="0.1"
                  :rules="[v => !!v || 'Este campo es requerido']"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-form>
  
          <v-divider class="my-4"></v-divider>
  
          <v-row>
            <v-col cols="12">
              <h3 class="text-h6 mb-2">Historial de Cambios</h3>
              <v-data-table
                :headers="historyHeaders"
                :items="configHistory"
                :loading="loading"
              >
                <template #[`item.fecha`]="{ item }">
                  {{ formatDate(item.fecha) }}
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </SideBar>
  </template>
  
  <script>
  import SideBar from '@/components/SideBar.vue'
  
  export default {
    name: 'ConfiguracionView',
    components: {
      SideBar
    },
    
    data() {
      return {
        saving: false,
        loading: false,
        config: {
          umbral_temp: 30.0,
          umbral_humedad: 75.0,
          umbral_luminosidad: 1000.0
        },
        historyHeaders: [
          { title: 'Fecha', key: 'fecha', align: 'start' },
          { title: 'Usuario', key: 'usuario' },
          { title: 'Campo Modificado', key: 'campo' },
          { title: 'Valor Anterior', key: 'valor_anterior' },
          { title: 'Valor Nuevo', key: 'valor_nuevo' }
        ],
        configHistory: [
          {
            fecha: '2024-01-15T10:30:00',
            usuario: 'Admin',
            campo: 'Umbral de Temperatura',
            valor_anterior: '28.0',
            valor_nuevo: '30.0'
          },
          {
            fecha: '2024-01-14T15:45:00',
            usuario: 'Supervisor',
            campo: 'Umbral de Humedad',
            valor_anterior: '70.0',
            valor_nuevo: '75.0'
          }
        ]
      }
    },
  
    methods: {
      async saveConfig() {
        this.saving = true
        // Simular guardado
        await new Promise(resolve => setTimeout(resolve, 1000))
        this.saving = false
        // Mostrar mensaje de éxito
        this.$root.$emit('showMessage', {
          text: 'Configuración guardada exitosamente',
          color: 'success'
        })
      },
  
      formatDate(date) {
        return new Date(date).toLocaleString('es-ES')
      }
    }
  }
  </script>