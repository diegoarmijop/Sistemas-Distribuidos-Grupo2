// src/components/PlagueDialog.vue
<template>
  <v-dialog v-model="dialog" max-width="600px">
    <v-card>
      <v-card-title>
        {{ isEvent ? 'Evento de Plaga' : 'Tipo de Plaga' }}
      </v-card-title>

      <v-card-text>
        <v-container>
          <v-row>
            <!-- Campos para Tipo de Plaga -->
            <template v-if="!isEvent">
              <v-col cols="12">
                <v-text-field
                  v-model="formData.nombre_comun"
                  label="Nombre Común"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  v-model="formData.nombre_cientifico"
                  label="Nombre Científico"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="formData.descripcion"
                  label="Descripción"
                ></v-textarea>
              </v-col>
            </template>

            <!-- Campos para Evento de Plaga -->
            <template v-else>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.tipo_plaga_id"
                  :items="options.plagueTypes"
                  item-title="nombre_comun"
                  item-value="tipo_plaga_id"
                  label="Tipo de Plaga"
                  required
                ></v-select>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.campo_id"
                  :items="options.campos"
                  item-title="nombre"
                  item-value="id"
                  label="Campo"
                  required
                ></v-select>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.ubicacion"
                  label="Ubicación Específica"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.nivel_severidad"
                  :items="['Baja', 'Media', 'Alta', 'Crítica']"
                  label="Nivel de Severidad"
                  required
                ></v-select>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="formData.acciones_tomadas"
                  label="Acciones Tomadas"
                  rows="3"
                ></v-textarea>
              </v-col>
            </template>
          </v-row>
        </v-container>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="error" @click="close">Cancelar</v-btn>
        <v-btn color="success" @click="save">Guardar</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'PlagueDialog',
  props: {
    modelValue: Boolean,
    type: {
      type: String,
      default: 'type', // 'type' para tipo de plaga, 'event' para evento
      validator: value => ['type', 'event'].includes(value)
    },
    item: {
      type: Object,
      default: () => ({})
    },
    options: {
      type: Object,
      default: () => ({
        plagueTypes: [],
        campos: []
      })
    }
  },
  
  emits: ['update:modelValue', 'save'],

  data() {
    return {
      formData: {}
    }
  },

  computed: {
    dialog: {
      get() {
        return this.modelValue
      },
      set(val) {
        this.$emit('update:modelValue', val)
      }
    },
    isEvent() {
      return this.type === 'event'
    },
    defaultFormData() {
      return this.isEvent ? {
        id: null,
        fecha_deteccion: new Date().toISOString(),
        ubicacion: '',
        nivel_severidad: '',
        acciones_tomadas: '',
        tipo_plaga_id: null,
        campo_id: null
      } : {
        tipo_plaga_id: null,
        nombre_comun: '',
        nombre_cientifico: '',
        descripcion: ''
      }
    }
  },

  watch: {
    item: {
      handler(val) {
        this.formData = { ...this.defaultFormData, ...val }
      },
      immediate: true
    }
  },

  methods: {
    close() {
      this.dialog = false
      this.$emit('update:modelValue', false)
      this.formData = { ...this.defaultFormData }
    },
    save() {
      this.$emit('save', { ...this.formData })
      this.close()
    }
  }
}
</script>