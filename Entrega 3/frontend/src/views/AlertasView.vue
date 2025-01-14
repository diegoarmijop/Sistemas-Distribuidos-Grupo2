<template>
    <SideBar page-title="Gestión de Alertas">
        <v-card>
            <v-card-title>
                Alertas del Sistema
                <v-spacer></v-spacer>
                <v-text-field 
                    v-model="search" 
                    append-icon="mdi-magnify" 
                    label="Buscar" 
                    single-line
                    hide-details>
                </v-text-field>
            </v-card-title>

            <v-card-text>
                <v-data-table 
                    :headers="headers" 
                    :items="alertas" 
                    :search="search" 
                    :loading="loading"
                    class="elevation-1">
                    <template #[`item.estado`]="{ item }">
                        <v-chip :color="getStatusColor(item.estado)" text-color="white">
                            {{ item.estado }}
                        </v-chip>
                    </template>

                    <template #[`item.tipo_alerta`]="{ item }">
                        <div class="d-flex flex-wrap gap-1">
                            <v-chip
                                v-for="(tipo, index) in splitTipoAlerta(item.tipo_alerta)"
                                :key="index"
                                :color="getAlertTypeColor(tipo)"
                                text-color="white"
                                size="small"
                                class="ma-1"
                            >
                                {{ tipo }}
                            </v-chip>
                        </div>
                    </template>

                    <template #[`item.descripcion`]="{ item }">
                        <v-tooltip bottom>
                            <template v-slot:activator="{ props }">
                                <div v-bind="props" class="text-truncate" style="max-width: 300px;">
                                    {{ item.descripcion }}
                                </div>
                            </template>
                            <span>{{ item.descripcion }}</span>
                        </v-tooltip>
                    </template>

                    <template #[`item.fecha_hora`]="{ item }">
                        {{ formatDate(item.fecha_hora) }}
                    </template>
                </v-data-table>
            </v-card-text>

            <v-card-actions class="justify-end">
                <v-btn
                    class="success-btn"
                    prepend-icon="mdi-microsoft-excel"
                    @click="downloadExcel"
                    :loading="downloadingExcel"
                >
                    Exportar Excel
                </v-btn>
            </v-card-actions>
        </v-card>
    </SideBar>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import axios from 'axios'
import * as XLSX from 'xlsx'

export default {
    name: 'AlertasView',
    components: {
        SideBar
    },

    data() {
        return {
            search: '',
            loading: false,
            downloadingExcel: false,
            alertas: [],
            headers: [
                { title: 'ID', key: 'id', align: 'start', width: '5%' },
                { title: 'Estado', key: 'estado', width: '10%' },
                { title: 'Tipo', key: 'tipo_alerta', width: '30%' },
                { title: 'Descripción', key: 'descripcion', width: '35%' },
                { title: 'Fecha', key: 'fecha_hora', width: '20%' }
            ]
        }
    },

    mounted() {
        this.loadAlertas()
        this.intervalId = setInterval(this.loadAlertas, 30000)
    },

    beforeUnmount() {
        if (this.intervalId) {
            clearInterval(this.intervalId)
        }
    },

    methods: {
        async loadAlertas() {
            this.loading = true
            try {
                const response = await axios.get('http://localhost:8080/api/alertas/')
                this.alertas = response.data.reverse()
            } catch (error) {
                console.error('Error al cargar alertas:', error)
            } finally {
                this.loading = false
            }
        },

        async downloadExcel() {
            this.downloadingExcel = true
            try {
                const response = await axios.get('http://localhost:8080/api/alertas/')
                const alertas = response.data

                // Preparar los datos para Excel
                const excelData = alertas.map(alerta => ({
                    ID: alerta.id,
                    Estado: alerta.estado,
                    'Tipo de Alerta': alerta.tipo_alerta,
                    Descripción: alerta.descripcion,
                    'Fecha y Hora': new Date(alerta.fecha_hora).toLocaleString('es-ES'),
                }))

                // Crear la hoja y el libro
                const ws = XLSX.utils.json_to_sheet(excelData)
                const wb = XLSX.utils.book_new()

                // Agregar estilos a los encabezados
                const range = XLSX.utils.decode_range(ws['!ref'])
                for (let C = range.s.c; C <= range.e.c; ++C) {
                    const address = XLSX.utils.encode_cell({ r: 0, c: C })
                    if (!ws[address]) continue
                    ws[address].s = {
                        font: { bold: true, color: { rgb: "FFFFFF" } },
                        fill: { fgColor: { rgb: "4472C4" } },
                        alignment: { horizontal: "center" }
                    }
                }

                // Ajustar anchos de columna
                ws['!cols'] = [
                    { wch: 5 },  // ID
                    { wch: 15 }, // Estado
                    { wch: 30 }, // Tipo de Alerta
                    { wch: 50 }, // Descripción
                    { wch: 20 }, // Fecha y Hora
                ]

                XLSX.utils.book_append_sheet(wb, ws, 'Alertas')

                // Generar nombre del archivo con la fecha actual
                const fecha = new Date().toLocaleDateString('es-ES').replace(/\//g, '-')
                const nombreArchivo = `Alertas_${fecha}.xlsx`

                // Descargar el archivo
                XLSX.writeFile(wb, nombreArchivo)

                this.$root.$emit('showMessage', {
                    text: 'Excel descargado exitosamente',
                    color: 'success'
                })
            } catch (error) {
                console.error('Error al descargar Excel:', error)
                this.$root.$emit('showMessage', {
                    text: 'Error al descargar el Excel',
                    color: 'error'
                })
            } finally {
                this.downloadingExcel = false
            }
        },

        getStatusColor(estado) {
            const colors = {
                'Seria': 'warning',
                'Crítico': 'error',
                'Extremo Peligro': 'deep-purple darken-4'
            }
            return colors[estado] || 'grey'
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
            return typeColors[tipo.trim()] || 'grey'
        },

        splitTipoAlerta(tipo) {
            return tipo ? tipo.split('/').map(t => t.trim()) : []
        },

        formatDate(date) {
            return new Date(date).toLocaleString('es-ES', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            })
        }
    }
}
</script>

<style scoped>
.gap-1 {
    gap: 4px;
}

/* Estilos para estados específicos */
:deep(.v-chip.v-theme--light.deep-purple) {
    background-color: #311B92 !important;
}

:deep(.v-data-table) {
    background-color: white;
    border-radius: 8px;
}

/* Estilo para el botón de Excel */
.success-btn {
    text-transform: none;
    background-color: #4CAF50 !important; /* Color verde */
    color: white !important;
    border-radius: 8px !important;
    padding: 10px 20px !important;
    font-weight: bold !important;
    transition: background-color 0.3s ease !important;
}

.success-btn:hover {
    background-color: #45A049 !important; /* Color verde más oscuro al pasar el ratón */
}
</style>