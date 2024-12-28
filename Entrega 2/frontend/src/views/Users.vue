<template>
    <SideBar page-title="Gestión de Usuarios">
      <v-card>
        <v-card-title class="d-flex align-center">
          Usuarios
          <v-spacer></v-spacer>
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Buscar usuario"
            single-line
            hide-details
            class="mx-4"
          ></v-text-field>
          <v-btn
            color="primary"
            prepend-icon="mdi-plus"
            @click="openDialog('create')"
          >
            Nuevo Usuario
          </v-btn>
        </v-card-title>
  
        <!-- Tabla de usuarios -->
        <v-data-table
          :headers="headers"
          :items="users"
          :loading="loading"
          :search="search"
        >
          <!-- Columna de acciones -->
          <template #[`item.actions`]="{ item }">
            <v-btn
              icon="mdi-pencil"
              size="small"
              color="primary"
              variant="text"
              @click="openDialog('edit', item.raw)"
            ></v-btn>
            <v-btn
              icon="mdi-delete"
              size="small"
              color="error"
              variant="text"
              @click="confirmDelete(item.raw)"
            ></v-btn>
          </template>
        </v-data-table>
      </v-card>
  
      <!-- Diálogo para crear/editar usuario -->
      <v-dialog v-model="dialog" max-width="500px">
        <v-card>
          <v-card-title>
            <span class="text-h5">{{ formTitle }}</span>
          </v-card-title>
  
          <v-card-text>
            <v-container>
              <v-row>
                <v-col cols="12">
                  <v-text-field
                    v-model="editedItem.nombre"
                    label="Nombre"
                    :rules="[v => !!v || 'El nombre es requerido']"
                  ></v-text-field>
                </v-col>
                <v-col cols="12">
                  <v-text-field
                    v-model="editedItem.email"
                    label="Email"
                    :rules="[v => !!v || 'El email es requerido', v => /.+@.+\..+/.test(v) || 'El email debe ser válido']"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" v-if="dialogMode === 'create'">
                  <v-text-field
                    v-model="editedItem.password"
                    label="Contraseña"
                    :type="showPassword ? 'text' : 'password'"
                    :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                    @click:append="showPassword = !showPassword"
                    :rules="[v => !!v || 'La contraseña es requerida']"
                  ></v-text-field>
                </v-col>
              </v-row>
            </v-container>
          </v-card-text>
  
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="grey-darken-1"
              variant="text"
              @click="closeDialog"
            >
              Cancelar
            </v-btn>
            <v-btn
              color="primary"
              variant="text"
              @click="saveUser"
              :loading="saving"
            >
              Guardar
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
  
      <!-- Diálogo de confirmación de eliminación -->
      <v-dialog v-model="deleteDialog" max-width="400px">
        <v-card>
          <v-card-title class="text-h5">
            ¿Eliminar usuario?
          </v-card-title>
          <v-card-text>
            ¿Está seguro de que desea eliminar este usuario? Esta acción no se puede deshacer.
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              color="grey-darken-1"
              variant="text"
              @click="deleteDialog = false"
            >
              Cancelar
            </v-btn>
            <v-btn
              color="error"
              variant="text"
              @click="deleteUser"
              :loading="deleting"
            >
              Eliminar
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </SideBar>
  </template>
  
  <script>
  import SideBar from '@/components/SideBar.vue'
  import axios from 'axios'
  
  export default {
    name: 'UsersView',
    components: {
      SideBar
    },
    data() {
      return {
        search: '',
        loading: false,
        dialog: false,
        deleteDialog: false,
        saving: false,
        deleting: false,
        showPassword: false,
        dialogMode: 'create', // 'create' o 'edit'
        headers: [
          { title: 'ID', align: 'start', key: 'id' },
          { title: 'Nombre', align: 'start', key: 'nombre' },
          { title: 'Email', key: 'email' },
          { title: 'Acciones', key: 'actions', sortable: false }
        ],
        users: [],
        editedItem: {
          id: null, // Asegúrate de tener el id en el objeto
          nombre: '',
          email: '',
          password: ''
        },
        defaultItem: {
          id: null,
          nombre: '',
          email: '',
          password: ''
        },
        userToDelete: null
      }
    },
    computed: {
      formTitle() {
        return this.dialogMode === 'create' ? 'Nuevo Usuario' : 'Editar Usuario'
      }
    },
    mounted() {
      this.loadUsers()
    },
    methods: {
      openDialog(mode, user = null) {
        this.dialogMode = mode
        if (mode === 'edit' && user) {
          this.editedItem = { ...user } // Asignar los datos del usuario al formulario de edición
        } else {
          this.editedItem = { ...this.defaultItem } // Limpiar el formulario si es un nuevo usuario
        }
        this.dialog = true // Mostrar el diálogo
      },
      closeDialog() {
        this.dialog = false // Cerrar el diálogo
        this.editedItem = { ...this.defaultItem } // Limpiar el formulario
      },
      async loadUsers() {
        this.loading = true
        try {
          console.log('Cargando usuarios...')
          const response = await axios.get('http://localhost:8080/api/users/')
          console.log('Respuesta:', response.data)
          this.users = response.data
        } catch (error) {
          console.error('Error al cargar usuarios:', error)
          this.$root.$emit('showMessage', {
            text: 'Error al cargar usuarios',
            color: 'error'
          })
        } finally {
          this.loading = false
        }
      },
      async saveUser() {
        if (!this.editedItem.nombre || !this.editedItem.email) return
  
        this.saving = true
        try {
          if (this.dialogMode === 'create') {
            await axios.post('http://localhost:8080/api/users/', this.editedItem)
          } else {
            // Asegúrate de que el ID se esté pasando correctamente en la solicitud PUT
            if (!this.editedItem.id) {
              this.$root.$emit('showMessage', {
                text: 'El ID del usuario es necesario para editar',
                color: 'error'
              })
              return
            }
            await axios.put(`http://localhost:8080/api/users/${this.editedItem.id}`, this.editedItem)
          }
  
          this.closeDialog()
          await this.loadUsers() // Recargamos la lista después de guardar
          this.$root.$emit('showMessage', {
            text: `Usuario ${this.dialogMode === 'create' ? 'creado' : 'actualizado'} correctamente`,
            color: 'success'
          })
        } catch (error) {
          console.error('Error al guardar usuario:', error)
          this.$root.$emit('showMessage', {
            text: `Error al ${this.dialogMode === 'create' ? 'crear' : 'actualizar'} usuario`,
            color: 'error'
          })
        } finally {
          this.saving = false
        }
      },
      async deleteUser() {
        if (!this.userToDelete) return
  
        this.deleting = true
        try {
          // Asegúrate de que el ID se pase correctamente al backend
          await axios.delete(`http://localhost:8080/api/users/${this.userToDelete.id}`)
          this.deleteDialog = false
          await this.loadUsers() // Recargamos la lista después de eliminar
          this.$root.$emit('showMessage', {
            text: 'Usuario eliminado correctamente',
            color: 'success'
          })
        } catch (error) {
          console.error('Error al eliminar usuario:', error)
          this.$root.$emit('showMessage', {
            text: 'Error al eliminar usuario',
            color: 'error'
          })
        } finally {
          this.deleting = false
          this.userToDelete = null
        }
      },
      confirmDelete(user) {
        this.userToDelete = user
        this.deleteDialog = true
      }
    }
  }
  </script>
  