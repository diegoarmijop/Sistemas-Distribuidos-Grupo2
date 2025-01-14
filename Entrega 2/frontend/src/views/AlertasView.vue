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
        </v-card>
    </SideBar>
</template>

<script>
import SideBar from '@/components/SideBar.vue'
import axios from 'axios'

export default {
    name: 'AlertasView',
    components: {
        SideBar
    },

    data() {
        return {
            search: '',
            loading: false,
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
        // Actualizar cada 30 segundos
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
</style>