<template>
    <SideBar page-title="Gestión de Alertas">
        <v-card>
            <v-card-title>
                Alertas del Sistema
                <v-spacer></v-spacer>
                <v-text-field v-model="search" append-icon="mdi-magnify" label="Buscar" single-line
                    hide-details></v-text-field>
            </v-card-title>

            <v-card-text>
                <v-data-table :headers="headers" :items="alertas" :search="search" :loading="loading">
                    <template #[`item.estado`]="{ item }">
                        <v-chip :color="getStatusColor(item.estado)" text-color="white">
                            {{ item.estado }}
                        </v-chip>
                    </template>

                    <template #[`item.tipo_alerta`]="{ item }">
                        <v-chip :color="getAlertTypeColor(item.tipo_alerta)" text-color="white">
                            {{ item.tipo_alerta }}
                        </v-chip>
                    </template>

                    <template #[`item.fecha_hora`]="{ item }">
                        {{ formatDate(item.fecha_hora) }}
                    </template>

                    <template #[`item.usuario`]="{ item }">
                        {{ item.usuario ? item.usuario.nombre : `ID: ${item.usuario_id}` }}
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
                { title: 'ID', key: 'id', align: 'start' },
                { title: 'Estado', key: 'estado' },
                { title: 'Tipo', key: 'tipo_alerta' },
                { title: 'Descripción', key: 'descripcion' },
                { title: 'Fecha', key: 'fecha_hora' },
                { title: 'Usuario', key: 'usuario' }
            ]
        }
    },

    mounted() {
        this.loadAlertas()
    },

    methods: {
        async loadAlertas() {
            this.loading = true
            try {
                const response = await axios.get('http://localhost:8080/api/alertas/')
                // Revertir el orden del array
                this.alertas = response.data.reverse()
            } catch (error) {
                console.error('Error al cargar alertas:', error)
            } finally {
                this.loading = false
            }
        },

        getStatusColor(estado) {
            const colors = {
                'Activa': 'warning',
                'Pendiente': 'info',
                'Resuelta': 'success',
                'Crítico': 'error'
            }
            return colors[estado] || 'grey'
        },

        getAlertTypeColor(tipo) {
            const colors = {
                'Urgente': 'red',
                'Advertencia': 'orange',
                'Información': 'blue',
                'Parámetros Alterados': 'purple'
            }
            return colors[tipo] || 'grey'
        },

        formatDate(date) {
            return new Date(date).toLocaleString('es-ES')
        }
    }
}
</script>