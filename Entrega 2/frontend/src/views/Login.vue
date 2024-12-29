<template>
    <v-container fluid class="fill-height bg-grey-lighten-3">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="8" md="6" lg="4">
          <v-card class="elevation-12 rounded-lg pa-4" style="backdrop-filter: blur(10px);">
            <!-- Logo y título -->
            <div class="text-center mb-4">
              <v-icon size="64" color="primary">mdi-account-circle</v-icon>
              <h1 class="text-h4 font-weight-bold mt-2 primary--text">Bienvenido</h1>
              <p class="text-subtitle-1 text-medium-emphasis">Ingresa tus credenciales para continuar</p>
            </div>
  
            <v-card-text>
              <v-form ref="form" v-model="valid" @submit.prevent="login">
                <!-- Campo de email con animación -->
                <v-slide-y-transition>
                  <v-text-field
                    v-model="email"
                    :rules="emailRules"
                    label="Correo electrónico"
                    prepend-inner-icon="mdi-email"
                    variant="outlined"
                    :loading="loading"
                    class="mb-2"
                    hide-details="auto"
                    @keyup.enter="focusPassword"
                  ></v-text-field>
                </v-slide-y-transition>
  
                <!-- Campo de contraseña con animación -->
                <v-slide-y-transition>
                  <v-text-field
                    v-model="password"
                    ref="passwordField"
                    :rules="passwordRules"
                    label="Contraseña"
                    prepend-inner-icon="mdi-lock"
                    :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                    :type="showPassword ? 'text' : 'password'"
                    variant="outlined"
                    :loading="loading"
                    class="mb-2"
                    hide-details="auto"
                    @click:append-inner="showPassword = !showPassword"
                    @keyup.enter="login"
                  ></v-text-field>
                </v-slide-y-transition>
  
                <!-- Opciones adicionales -->
                <div class="d-flex justify-space-between align-center my-2">
                  <v-checkbox
                    v-model="rememberMe"
                    label="Recordarme"
                    hide-details
                    class="mt-0"
                  ></v-checkbox>
                  <v-btn
                    variant="text"
                    color="primary"
                    class="text-caption"
                    @click="forgotPassword"
                  >
                    ¿Olvidaste tu contraseña?
                  </v-btn>
                </div>
  
                <!-- Botón de login con efecto de loading -->
                <v-expand-transition>
                  <v-btn
                    block
                    color="primary"
                    size="large"
                    :loading="loading"
                    :disabled="!valid || loading"
                    height="50"
                    class="mt-6"
                    @click="login"
                  >
                    {{ loading ? 'Iniciando sesión...' : 'Iniciar Sesión' }}
                    <template v-slot:loader>
                      <v-progress-linear indeterminate color="white"></v-progress-linear>
                    </template>
                  </v-btn>
                </v-expand-transition>
              </v-form>
              
            </v-card-text>
          </v-card>
  
          <!-- Snackbar para notificaciones -->
          <v-snackbar
            v-model="snackbar"
            :color="snackbarColor"
            :timeout="4000"
            location="top"
            style="top: 20px"
          >
            <div class="d-flex align-center">
              <v-icon :color="snackbarColor" class="mr-2">
                {{ snackbarIcon }}
              </v-icon>
              {{ snackbarText }}
            </div>
            <template v-slot:actions>
              <v-btn
                variant="text"
                @click="snackbar = false"
              >
                Cerrar
              </v-btn>
            </template>
          </v-snackbar>
        </v-col>
      </v-row>
    </v-container>
  </template>
  
  <script>
  import axios from 'axios'
  
  export default {
    name: 'LoginView',
    data: () => ({
      valid: false,
      loading: false,
      email: '',
      password: '',
      showPassword: false,
      rememberMe: false,
      snackbar: false,
      snackbarText: '',
      snackbarColor: 'success',
      snackbarIcon: 'mdi-check-circle',
      emailRules: [
        v => !!v || 'El correo es requerido',
        v => /.+@.+\..+/.test(v) || 'Ingresa un correo electrónico válido'
      ],
      passwordRules: [
        v => !!v || 'La contraseña es requerida',
        v => v.length >= 6 || 'La contraseña debe tener al menos 6 caracteres'
      ]
    }),
  
    methods: {
      focusPassword() {
        this.$refs.passwordField.$el.querySelector('input').focus()
      },
  
      showNotification(message, type = 'error') {
        this.snackbarColor = type
        this.snackbarIcon = type === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
        this.snackbarText = message
        this.snackbar = true
      },
  
      async login() {
        if (!this.$refs.form.validate()) return
  
        this.loading = true
        try {
          const response = await axios.post('http://localhost:8080/api/login', {
            email: this.email,
            password: this.password
          })
  
          // Guardar el token
          if (this.rememberMe) {
            localStorage.setItem('token', response.data.token)
          } else {
            sessionStorage.setItem('token', response.data.token)
          }
          
          this.showNotification('¡Inicio de sesión exitoso!', 'success')
          
          // Pequeña pausa para mostrar la animación
          setTimeout(() => {
            this.$router.push('/dashboard')
          }, 1000)
  
        } catch (error) {
          let errorMessage = 'Error al iniciar sesión'
          
          if (error.response?.status === 401) {
            errorMessage = 'Credenciales incorrectas'
          } else if (error.response?.status === 404) {
            errorMessage = 'Usuario no encontrado'
          } else if (error.response?.data?.error) {
            errorMessage = error.response.data.error
          }
          
          this.showNotification(errorMessage, 'error')
        } finally {
          this.loading = false
        }
      },
  
      forgotPassword() {
        // Implementar lógica para recuperar contraseña
        this.$router.push('/forgot-password')
      }
  
    }
  }
  </script>
  
  <style scoped>
.v-container {
  background-image: url('../assets/campo_principal.png'); /* Ajusta la ruta según tu estructura */
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  position: relative;
}

/* Efecto blur para el fondo */
.v-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: inherit;
  filter: blur(8px); /* Ajusta el valor para más o menos blur */
  backdrop-filter: blur(8px);
  z-index: 0;
}

.v-card {
  border: 1px solid rgba(var(--v-border-color), 0.1);
  background: rgba(255, 255, 255, 0.8); /* Fondo semi-transparente */
  backdrop-filter: blur(16px);
  position: relative;
  z-index: 1;
}

.v-text-field :deep(.v-field__outline__start) {
  border-radius: 10px 0 0 10px;
}

.v-text-field :deep(.v-field__outline__end) {
  border-radius: 0 10px 10px 0;
}

.v-btn {
  text-transform: none;
  letter-spacing: 0.5px;
}

/* Asegura que el contenido esté por encima del blur */
.v-row {
  position: relative;
  z-index: 1;
}
</style>