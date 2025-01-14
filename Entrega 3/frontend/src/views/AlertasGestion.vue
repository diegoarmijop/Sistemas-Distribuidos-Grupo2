<template>
    <SideBar page-title="Gestión de Alertas">
        <v-container fluid>
            <v-row>
                <!-- Filtros y estadísticas -->
                <v-col cols="12">
                    <v-card class="mb-4">
                        <v-card-text class="d-flex align-center">
                            <v-select
    v-model="filtroEstado"
    :items="['Todas', 'Resuelta', 'Seria', 'Crítico', 'Extremo Peligro']"
    label="Estado"
    class="mr-4 estado-select"
    style="max-width: 200px"
    variant="outlined"
    density="comfortable"
    bg-color="white"
></v-select>
                            <v-spacer></v-spacer>
                            <v-btn color="success" prepend-icon="mdi-microsoft-excel" @click="downloadExcel"
                                :loading="downloadingExcel">
                                Exportar Excel
                            </v-btn>
                        </v-card-text>
                    </v-card>
                </v-col>

                <!-- Lista de alertas -->
                <v-col cols="12">
                    <v-data-table :headers="headers" :items="alertasFiltradas" :loading="loading" class="elevation-1">
                        <!-- Estado -->
                        <template #[`item.estado`]="{ item }">
                            <v-chip :color="getStatusColor(item.estado)" text-color="white">
                                {{ item.estado }}
                            </v-chip>
                        </template>

                        <!-- Tipo de alerta -->
                        <template #[`item.tipo_alerta`]="{ item }">
                            <div class="d-flex flex-wrap gap-1">
                                <v-chip v-for="(tipo, index) in splitTipoAlerta(item.tipo_alerta)" :key="index"
                                    :color="getAlertTypeColor(tipo)" text-color="white" size="small"
                                    class="ma-1 tipo-chip" variant="flat">
                                    {{ tipo }}
                                </v-chip>
                            </div>
                        </template>

                        <!-- Acciones -->
                        <template #[`item.actions`]="{ item }">
                            <v-btn icon variant="text" color="primary" @click="mostrarDetalles(item)">
                                <v-icon>mdi-eye</v-icon>
                            </v-btn>
                            <v-btn v-if="item.estado !== 'Resuelta'" icon variant="text" color="success"
                                @click="mostrarDialogoResolucion(item)">
                                <v-icon>mdi-check-circle</v-icon>
                            </v-btn>
                        </template>
                    </v-data-table>
                </v-col>
            </v-row>

            <!-- Diálogo de detalles -->
            <v-dialog v-model="dialogoDetalles" max-width="700">
                <v-card v-if="alertaSeleccionada">
                    <v-card-title class="d-flex align-center">
                        Detalles de la Alerta
                        <v-spacer></v-spacer>
                        <v-btn icon @click="dialogoDetalles = false">
                            <v-icon>mdi-close</v-icon>
                        </v-btn>
                    </v-card-title>

                    <v-card-text>
                        <v-row>
                            <v-col cols="6">
                                <strong>Estado:</strong>
                                <v-chip :color="getStatusColor(alertaSeleccionada.estado)" text-color="white"
                                    class="ml-2">
                                    {{ alertaSeleccionada.estado }}
                                </v-chip>
                            </v-col>
                            <v-col cols="6">
                                <strong>Fecha:</strong>
                                {{ formatDate(alertaSeleccionada.fecha_hora) }}
                            </v-col>
                            <v-col cols="12">
                                <strong>Tipos de Alerta:</strong>
                                <div class="d-flex flex-wrap gap-1 mt-2">
                                    <v-chip v-for="(tipo, index) in splitTipoAlerta(alertaSeleccionada.tipo_alerta)"
                                        :key="index" :color="getAlertTypeColor(tipo)" text-color="white" size="small"
                                        class="ma-1 tipo-chip" variant="flat">
                                        {{ tipo }}
                                    </v-chip>
                                </div>
                            </v-col>
                            <v-col cols="12">
                                <strong>Descripción:</strong>
                                <p class="mt-2">{{ alertaSeleccionada.descripcion }}</p>
                            </v-col>

                            <!-- Información de resolución si está resuelta -->
                            <template v-if="alertaSeleccionada.estado === 'Resuelta'">
                                <v-col cols="12">
                                    <v-divider class="my-4"></v-divider>
                                    <h3 class="text-h6 mb-3">Información de Resolución</h3>
                                </v-col>
                                <v-col cols="12">
                                    <strong>Solución Aplicada:</strong>
                                    <p class="mt-2">{{ alertaSeleccionada.solucion_aplicada }}</p>
                                </v-col>
                                <v-col cols="12">
                                    <strong>Comentarios:</strong>
                                    <p class="mt-2">{{ alertaSeleccionada.comentarios }}</p>
                                </v-col>
                                <v-col cols="12">
                                    <strong>Plan de Acción:</strong>
                                    <p class="mt-2">{{ alertaSeleccionada.plan_accion }}</p>
                                </v-col>
                                <v-col cols="6">
                                    <strong>Fecha de Resolución:</strong>
                                    <p class="mt-2">{{ formatDate(alertaSeleccionada.fecha_resolucion) }}</p>
                                </v-col>
                            </template>
                        </v-row>
                    </v-card-text>
                </v-card>
            </v-dialog>

            <!-- Diálogo de resolución -->
            <v-dialog v-model="dialogoResolucion" max-width="900">
                <v-card v-if="alertaSeleccionada">
                    <v-card-title class="d-flex align-center bg-grey-lighten-4 py-3"> <!-- Añadido py-3 -->
                        <span class="text-h6">Resolver Alerta</span>
                        <v-spacer></v-spacer>
                        <v-btn icon @click="dialogoResolucion = false">
                            <v-icon>mdi-close</v-icon>
                        </v-btn>
                    </v-card-title>

                    <v-divider></v-divider>

                    <v-card-text class="pa-4">
                        <!-- Información de la alerta actual -->
                        <div class="mb-4">
                            <v-chip :color="getStatusColor(alertaSeleccionada.estado)" text-color="white" class="mr-2">
                                {{ alertaSeleccionada.estado }}
                            </v-chip>
                            <span class="text-grey">{{ formatDate(alertaSeleccionada.fecha_hora) }}</span>
                        </div>

                        <v-row>
                            <!-- Panel de sugerencias -->
                            <v-col cols="12" md="5">
                                <v-card variant="outlined">
                                    <v-card-title class="text-subtitle-1 bg-grey-lighten-3">
                                        <v-icon start class="mr-2">mdi-lightbulb</v-icon>
                                        Sugerencias de Resolución
                                    </v-card-title>

                                    <v-card-text class="pa-6">
                                        <!-- Sugerencias para problemas múltiples -->
                                        <template v-if="sugerencias.combinadas.length">
                                            <div class="text-subtitle-2 font-weight-bold mb-3">
                                                Soluciones Prioritarias:
                                            </div>
                                            <v-list density="compact" class="bg-yellow-lighten-5 rounded-lg mb-4">
                                                <v-list-item v-for="(sugerencia, index) in sugerencias.combinadas"
                                                    :key="'combined-' + index" :value="sugerencia"
                                                    @click="aplicarSugerencia(sugerencia)" class="cursor-pointer">
                                                    <template v-slot:prepend>
                                                        <v-icon color="warning" size="small">
                                                            mdi-alert-circle
                                                        </v-icon>
                                                    </template>
                                                    <v-list-item-title>{{ sugerencia }}</v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </template>

                                        <!-- Sugerencias individuales -->
                                        <template v-if="sugerencias.individuales.length">
                                            <div class="text-subtitle-2 font-weight-bold mb-3">
                                                Soluciones Adicionales:
                                            </div>
                                            <v-list density="compact" class="bg-grey-lighten-5 rounded-lg">
                                                <v-list-item v-for="(sugerencia, index) in sugerencias.individuales"
                                                    :key="'individual-' + index" :value="sugerencia"
                                                    @click="aplicarSugerencia(sugerencia)" class="cursor-pointer">
                                                    <template v-slot:prepend>
                                                        <v-icon color="info" size="small">
                                                            mdi-lightbulb-outline
                                                        </v-icon>
                                                    </template>
                                                    <v-list-item-title>{{ sugerencia }}</v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </template>

                                        <!-- Mensaje cuando no hay sugerencias -->
                                        <div v-if="!sugerencias.combinadas.length && !sugerencias.individuales.length"
                                            class="text-center pa-4 text-grey">
                                            No hay sugerencias disponibles para este tipo de alerta
                                        </div>
                                    </v-card-text>
                                </v-card>

                                <!-- Soluciones históricas -->
                                <v-card variant="outlined" class="mt-4">
                                    <v-expansion-panels variant="accordion">
                                        <v-expansion-panel>
                                            <v-expansion-panel-title class="text-subtitle-2">
                                                <v-icon start class="mr-2">mdi-history</v-icon>
                                                Soluciones Anteriores
                                            </v-expansion-panel-title>
                                            <v-expansion-panel-text>
                                                <v-list density="compact" v-if="sugerenciasHistoricas.length">
                                                    <v-list-item v-for="(solucion, index) in sugerenciasHistoricas"
                                                        :key="index" @click="aplicarSolucionHistorica(solucion)"
                                                        class="cursor-pointer">
                                                        <v-list-item-title>
                                                            {{ solucion.solucion_aplicada }}
                                                        </v-list-item-title>
                                                        <v-list-item-subtitle>
                                                            {{ formatDate(solucion.fecha_resolucion) }}
                                                        </v-list-item-subtitle>
                                                    </v-list-item>
                                                </v-list>
                                                <div v-else class="text-center pa-2 text-grey">
                                                    No hay soluciones anteriores
                                                </div>
                                            </v-expansion-panel-text>
                                        </v-expansion-panel>
                                    </v-expansion-panels>
                                </v-card>
                            </v-col>
                            <!-- Formulario de resolución -->
                            <v-col cols="12" md="7">
                                <v-form ref="formResolucion" v-model="formularioValido">
                                    <div class="form-field">
                                        <div class="field-label">Solución a Aplicar</div>
                                        <v-textarea v-model="resolucion.solucion_aplicada"
                                            :rules="[v => !!v || 'Este campo es requerido']" rows="4"
                                            class="mb-4 solution-field" hide-details="auto" variant="outlined"
                                            density="comfortable"></v-textarea>
                                    </div>

                                    <div class="form-field">
                                        <div class="field-label">Comentarios Adicionales</div>
                                        <v-textarea v-model="resolucion.comentarios" rows="3"
                                            class="mb-4 comments-field" hide-details="auto" variant="outlined"
                                            density="comfortable"></v-textarea>
                                    </div>

                                    <div class="form-field">
                                        <div class="field-label">Plan de Acción y Seguimiento</div>
                                        <v-textarea v-model="resolucion.plan_accion"
                                            :rules="[v => !!v || 'Este campo es requerido']" rows="4"
                                            class="action-field" hide-details="auto" variant="outlined"
                                            density="comfortable"></v-textarea>
                                    </div>

                                    <v-card-actions class="mt-4 pa-0">
                                        <v-spacer></v-spacer>
                                        <v-btn color="grey-darken-1" variant="text" @click="dialogoResolucion = false">
                                            Cancelar
                                        </v-btn>
                                        <v-btn color="success" :loading="resolving" :disabled="!formularioValido"
                                            @click="resolverAlerta">
                                            Resolver Alerta
                                        </v-btn>
                                    </v-card-actions>
                                </v-form>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </v-dialog>
        </v-container>
    </SideBar>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import axios from 'axios'
import * as XLSX from 'xlsx'

export default {
    name: 'AlertasGestion',
    components: {
        SideBar
    },

    data: () => ({
        loading: false,
        alertas: [],
        headers: [
            { title: 'ID', key: 'id', width: '5%' },
            { title: 'Estado', key: 'estado', width: '10%' },
            { title: 'Tipo', key: 'tipo_alerta', width: '25%' },
            { title: 'Descripción', key: 'descripcion', width: '35%' },
            { title: 'Fecha', key: 'fecha_hora', width: '15%' },
            { title: 'Acciones', key: 'actions', width: '10%', sortable: false }
        ],
        filtroEstado: 'Todas',
        downloadingExcel: false,
        dialogoDetalles: false,
        dialogoResolucion: false,
        alertaSeleccionada: null,
        sugerencias: {
            combinadas: [],
            individuales: []
        },
        sugerenciasHistoricas: [],
        resolucion: {
            solucion_aplicada: '',
            comentarios: '',
            plan_accion: '',
            resuelta_por: null
        },
        resolving: false,
        formularioValido: false
    }),

    computed: {
        alertasFiltradas() {
        // Primero filtrar para excluir las alertas activas
        let alertasFiltradas = this.alertas.filter(alerta => alerta.estado !== 'Activa');
        
        // Luego aplicar el filtro de estado si no es 'Todas'
        if (this.filtroEstado !== 'Todas') {
            alertasFiltradas = alertasFiltradas.filter(alerta => alerta.estado === this.filtroEstado);
        }
        
        // Ordenar por fecha (más recientes primero)
        return alertasFiltradas.sort((a, b) => new Date(b.fecha_hora) - new Date(a.fecha_hora));
    }
},

    mounted() {
        this.cargarAlertas()
        const userData = JSON.parse(localStorage.getItem('userData') || '{}')
        this.resolucion.resuelta_por = userData.id
    },

    methods: {
        async cargarAlertas() {
    this.loading = true;
    try {
        const response = await axios.get('http://localhost:8080/api/alertas/');
        // Ordenar las alertas antes de asignarlas
        this.alertas = response.data.sort((a, b) => new Date(b.fecha_hora) - new Date(a.fecha_hora));
    } catch (error) {
        console.error('Error al cargar alertas:', error);
        this.$root.$emit('showMessage', {
            text: 'Error al cargar alertas',
            color: 'error'
        });
    } finally {
        this.loading = false;
    }
},

        async mostrarDetalles(alerta) {
            this.alertaSeleccionada = alerta
            this.dialogoDetalles = true
        },

        async mostrarDialogoResolucion(alerta) {
            this.alertaSeleccionada = alerta
            this.resolucion = {
                solucion_aplicada: '',
                comentarios: '',
                plan_accion: '',
                resuelta_por: this.resolucion.resuelta_por
            }

            // Cargar sugerencias y soluciones históricas
            await Promise.all([
                this.cargarSugerencias(alerta.tipo_alerta),
                this.cargarSolucionesHistoricas(alerta.tipo_alerta)
            ])

            this.dialogoResolucion = true
        },

        async cargarSugerencias(tipoAlerta) {
            try {
                const response = await axios.get(`http://localhost:8080/api/alertas/sugerencias?tipo=${tipoAlerta}`)
                // Organizamos las sugerencias por categoría
                if (response.data.tipos_detectados && response.data.tipos_detectados.length > 1) {
                    this.sugerencias = {
                        combinadas: response.data.sugerencias.slice(0, 4), // Primeras 4 sugerencias son para la combinación
                        individuales: response.data.sugerencias.slice(4) // El resto son individuales
                    }
                } else {
                    this.sugerencias = {
                        combinadas: [],
                        individuales: response.data.sugerencias
                    }
                }
            } catch (error) {
                console.error('Error al cargar sugerencias:', error)
                this.sugerencias = { combinadas: [], individuales: [] }
            }
        },

        async cargarSolucionesHistoricas(tipoAlerta) {
            try {
                // Obtener alertas resueltas del mismo tipo
                const response = await axios.get('http://localhost:8080/api/alertas/')
                this.sugerenciasHistoricas = response.data
                    .filter(alerta =>
                        alerta.estado === 'Resuelta' &&
                        alerta.tipo_alerta === tipoAlerta &&
                        alerta.solucion_aplicada
                    )
                    .slice(0, 5) // Mostrar solo las últimas 5 soluciones
            } catch (error) {
                console.error('Error al cargar soluciones históricas:', error)
                this.sugerenciasHistoricas = []
            }
        },

        aplicarSugerencia(sugerencia) {
            this.resolucion.solucion_aplicada = sugerencia
            this.resolucion.plan_accion = `Implementar y monitorear: ${sugerencia}\n\nPasos de seguimiento:\n1. Implementación inmediata de la solución\n2. Monitoreo cada 2 horas durante las primeras 24 horas\n3. Evaluación de efectividad al día siguiente\n4. Ajustes según sea necesario`

            // Mostrar feedback visual
            this.$root.$emit('showMessage', {
                text: 'Sugerencia aplicada al formulario',
                color: 'info'
            })
        },

        aplicarSolucionHistorica(solucion) {
            this.resolucion.solucion_aplicada = solucion.solucion_aplicada
            this.resolucion.plan_accion = solucion.plan_accion
            this.resolucion.comentarios = `Solución basada en caso anterior del ${this.formatDate(solucion.fecha_resolucion)}`
            // Mostrar feedback visual
            this.$root.$emit('showMessage', {
                text: 'Solución histórica aplicada',
                color: 'info'
            })
        },

        async resolverAlerta() {
            if (!this.$refs.formResolucion.validate()) return

            this.resolving = true
            try {
                await axios.post(
                    `http://localhost:8080/api/alertas/${this.alertaSeleccionada.id}/resolver`,
                    this.resolucion
                )

                this.dialogoResolucion = false
                await this.cargarAlertas()

                this.$root.$emit('showMessage', {
                    text: 'Alerta resuelta exitosamente',
                    color: 'success'
                })
            } catch (error) {
                console.error('Error al resolver la alerta:', error)
                this.$root.$emit('showMessage', {
                    text: 'Error al resolver la alerta',
                    color: 'error'
                })
            } finally {
                this.resolving = false
            }
        },

        getStatusColor(estado) {
            const colors = {
                'Seria': 'warning',
                'Crítico': 'error',
                'Extremo Peligro': 'deep-purple-darken-4',
                'Resuelta': 'success'
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
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            })
        },

        async downloadExcel() {
            this.downloadingExcel = true
            try {
                const excelData = this.alertas.map(alerta => ({
                    ID: alerta.id,
                    Estado: alerta.estado,
                    'Tipo de Alerta': alerta.tipo_alerta,
                    Descripción: alerta.descripcion,
                    'Fecha y Hora': this.formatDate(alerta.fecha_hora),
                    'Solución Aplicada': alerta.solucion_aplicada || '-',
                    'Fecha Resolución': alerta.fecha_resolucion ? this.formatDate(alerta.fecha_resolucion) : '-',
                    'Comentarios': alerta.comentarios || '-',
                    'Plan de Acción': alerta.plan_accion || '-'
                }))

                const ws = XLSX.utils.json_to_sheet(excelData)
                const wb = XLSX.utils.book_new()

                ws['!cols'] = [
                    { wch: 5 },  // ID
                    { wch: 15 }, // Estado
                    { wch: 30 }, // Tipo de Alerta
                    { wch: 50 }, // Descripción
                    { wch: 20 }, // Fecha y Hora
                    { wch: 40 }, // Solución Aplicada
                    { wch: 20 }, // Fecha Resolución
                    { wch: 40 }, // Comentarios
                    { wch: 40 }  // Plan de Acción
                ]

                XLSX.utils.book_append_sheet(wb, ws, 'Alertas')

                const fecha = new Date().toLocaleDateString('es-ES').replace(/\//g, '-')
                const nombreArchivo = `Alertas_${fecha}.xlsx`

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

/* Estilo para chips de tipo cuadrados */
:deep(.tipo-chip) {
    border-radius: 0 !important;
}

/* Estilos para el select de estado */
:deep(.estado-select) {
    .v-field__input {
        padding-top: 0 !important;
        min-height: 40px !important;
    }

    .v-field__field {
        min-height: 40px !important;
    }

    .v-field__outline {
        --v-field-border-width: 1px !important;
    }

    .v-field--variant-outlined .v-field__outline {
        border-color: rgba(0, 0, 0, 0.25) !important;
    }

    .v-select__selection-text {
        padding: 0;
        line-height: 40px;
    }
}

/* Estilos para el menú del select */
:deep(.v-overlay__content) {
    .v-list {
        padding: 4px;
        border-radius: 8px;
    }

    .v-list-item {
        min-height: 36px;
        border-radius: 4px;
    }

    .v-list-item-title {
        font-size: 0.95rem;
    }
}

/* Estilos para el panel de sugerencias */
.cursor-pointer {
    cursor: pointer;
}

.cursor-pointer:hover {
    background-color: rgba(var(--v-theme-primary), 0.05);
}

:deep(.v-expansion-panel-title) {
    font-size: 0.875rem;
}

/* Estilos para los campos del formulario */
.form-field {
    margin-bottom: 24px;
}

.form-field:last-child {
    margin-bottom: 0;
}

.field-label {
    font-size: 0.95rem;
    font-weight: 500;
    color: rgba(0, 0, 0, 0.7);
    margin-bottom: 8px;
}

/* Ajustes para los textareas */
:deep(.v-textarea .v-field__input) {
    font-size: 0.95rem !important;
    padding: 12px !important;
    min-height: 120px;
}

:deep(.v-field--variant-outlined) {
    --v-field-border-width: 1px;
    border-radius: 8px;
}

:deep(.v-textarea .v-field__outline) {
    --v-field-border-width: 1px !important;
    border-color: rgba(0, 0, 0, 0.25) !important;
}

:deep(.v-field--focused .v-field__outline) {
    border-color: rgb(var(--v-theme-primary)) !important;
}

/* Quitar labels internos */
:deep(.v-field__field) {
    min-height: unset !important;
}

:deep(.v-field__input) {
    min-height: unset !important;
    padding: 8px 12px !important;
}

/* Eliminar transiciones */
:deep(.v-field__outline),
:deep(.v-field__field),
:deep(.v-field__input) {
    transition: none !important;
}

/* Animaciones */
.v-card {
    transition: all 0.3s ease-in-out;
}

.v-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 25px 0 rgba(0, 0, 0, 0.1);
}

/* Estilo para las sugerencias seleccionadas */
:deep(.v-list-item--active) {
    background-color: rgba(var(--v-theme-primary), 0.1) !important;
}

/* Estilo para el panel de soluciones históricas */
:deep(.v-expansion-panel) {
    background-color: transparent !important;
}

:deep(.v-expansion-panel-text__wrapper) {
    padding: 0;
}

/* Ajustes para el diálogo de resolución */
:deep(.v-dialog > .v-card > .v-card-title) {
    min-height: 64px;
}

:deep(.v-dialog > .v-card > .v-card-text) {
    padding: 24px !important;
    margin-top: 0;
}

/* Estilos para los paneles de sugerencias */
:deep(.bg-yellow-lighten-5) {
    background-color: #fff8e1 !important;
}

:deep(.bg-grey-lighten-5) {
    background-color: #f5f5f5 !important;
}

.rounded-lg {
    border-radius: 8px;
}

.suggestion-divider {
    margin: 12px 0;
    opacity: 0.12;
}

/* Estilos para el encabezado del dialogo */
:deep(.bg-grey-lighten-4) {
    background-color: #f5f5f5 !important;
}

/* Mejoras visuales para los diálogos */
:deep(.v-dialog) {
    border-radius: 12px;
    overflow: hidden;
}

:deep(.v-card-title.bg-grey-lighten-4) {
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

/* Ajustes para los listados de sugerencias */
:deep(.v-list-item) {
    margin-bottom: 4px;
    border-radius: 6px;
}

:deep(.v-list-item:last-child) {
    margin-bottom: 0;
}

/* Estilos para mensajes de error */
:deep(.v-messages) {
    min-height: 0;
    font-size: 0.85rem;
    padding-top: 4px;
}

:deep(.error--text) {
    color: rgb(var(--v-theme-error)) !important;
}

:deep(.v-field--error .v-field__outline) {
    border-color: rgb(var(--v-theme-error)) !important;
    border-width: 2px !important;
}

/* Ajustes para dispositivos móviles */
@media (max-width: 600px) {
    :deep(.v-dialog) {
        margin: 0;
    }
    
    .v-card-title {
        font-size: 1.1rem;
    }
    
    .text-subtitle-2 {
        font-size: 0.9rem;
    }

    :deep(.v-dialog > .v-card > .v-card-text) {
        padding: 16px !important;
    }

    :deep(.v-dialog > .v-card > .v-card-title) {
        padding: 16px;
    }

    :deep(.v-textarea .v-field__input) {
        font-size: 0.9rem;
    }

    .field-label {
        font-size: 0.9rem;
    }
}
</style>