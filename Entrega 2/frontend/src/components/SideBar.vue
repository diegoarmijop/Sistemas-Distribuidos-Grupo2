<template>
    <v-app>
      <!-- Barra lateral -->
      <v-navigation-drawer
        v-model="drawer"
        app
        color="grey-lighten-4"
        class="elevation-3"
      >
        <v-list>
          <v-list-item
            prepend-avatar="https://randomuser.me/api/portraits/men/78.jpg"
            :title="userData.name || 'Usuario'"
            :subtitle="userData.email || 'usuario@usach.cl'"
          ></v-list-item>
        </v-list>
  
        <v-divider></v-divider>
  
        <v-list density="compact" nav>
          <v-list-item
            v-for="item in menuItems"
            :key="item.path"
            :value="item"
            :prepend-icon="item.icon"
            :title="item.title"
            :to="item.path"
            :active="currentPath === item.path"
            active-color="primary"
          ></v-list-item>
        </v-list>
      </v-navigation-drawer>
  
      <!-- Barra superior -->
      <v-app-bar app color="white" elevation="1">
        <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
        <v-toolbar-title class="text-h6 font-weight-bold">
          {{ pageTitle }}
        </v-toolbar-title>
        <v-spacer></v-spacer>
  
        <!-- Botones de la barra superior -->
        <v-btn icon class="mr-2">
          <v-badge dot color="error">
            <v-icon>mdi-bell</v-icon>
          </v-badge>
        </v-btn>
  
        <v-menu min-width="200px" rounded>
          <template v-slot:activator="{ props }">
            <v-btn
              icon
              v-bind="props"
            >
              <v-avatar color="grey-lighten-1" size="32">
                <v-icon dark>mdi-account-circle</v-icon>
              </v-avatar>
            </v-btn>
          </template>
          <v-card>
            <v-card-text>
              <div class="d-flex align-center mb-3">
                <v-avatar class="mr-3">
                  <v-icon>mdi-account-circle</v-icon>
                </v-avatar>
                <div>
                  <div class="text-subtitle-1">{{ userData.name || 'Usuario' }}</div>
                  <div class="text-caption">{{ userData.email || 'usuario@email.com' }}</div>
                </div>
              </div>
              <v-divider class="mb-3"></v-divider>
              <v-list density="compact" nav>
                <v-list-item
                  prepend-icon="mdi-account-cog"
                  title="Configuración"
                  to="/settings"
                ></v-list-item>
                <v-list-item
                  prepend-icon="mdi-logout"
                  title="Cerrar sesión"
                  @click="logout"
                ></v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-menu>
      </v-app-bar>
  
      <!-- Contenido principal -->
      <v-main class="grey-lighten-3">
        <v-container fluid>
          <slot></slot>
        </v-container>
      </v-main>
  
      <!-- Snackbar para notificaciones -->
      <v-snackbar
        v-model="snackbar.show"
        :color="snackbar.color"
        :timeout="3000"
      >
        {{ snackbar.text }}
        <template v-slot:actions>
          <v-btn
            color="white"
            variant="text"
            @click="snackbar.show = false"
          >
            Cerrar
          </v-btn>
        </template>
      </v-snackbar>
    </v-app>
  </template>
  
  <script>
  export default {
    name: 'AppLayout',
    props: {
      pageTitle: {
        type: String,
        default: 'Dashboard'
      }
    },
    data: () => ({
      drawer: true,
      userData: {},
      currentPath: '',
      snackbar: {
        show: false,
        text: '',
        color: 'success'
      },
      menuItems: [
        { title: 'Dashboard', icon: 'mdi-view-dashboard', path: '/dashboard' },
        { title: 'Usuarios', icon: 'mdi-account-group', path: '/users' },
        { title: 'Sectores', icon: 'mdi-seed', path: '/sectors' },
        { title: 'Reportes', icon: 'mdi-chart-bar', path: '/reportes' },
        { title: 'Configuración', icon: 'mdi-cog', path: '/configuracion' },
      ]
    }),
  
    created() {
      this.currentPath = this.$route.path
      this.loadUserData()
    },
  
    watch: {
      '$route.path'(newPath) {
        this.currentPath = newPath
      }
    },
  
    methods: {
      loadUserData() {
        const userStr = localStorage.getItem('user') || sessionStorage.getItem('user')
        if (userStr) {
          try {
            this.userData = JSON.parse(userStr)
          } catch (e) {
            console.error('Error parsing user data:', e)
          }
        }
      },
  
      logout() {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        sessionStorage.removeItem('token')
        sessionStorage.removeItem('user')
        this.$router.push('/login')
      },
  
      showNotification(text, color = 'success') {
        this.snackbar = {
          show: true,
          text,
          color
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .v-navigation-drawer {
    transition: 0.2s ease-in-out;
  }
  
  .v-list-item {
    transition: all 0.3s;
  }
  
  .v-list-item:hover {
    background-color: rgba(var(--v-theme-primary), 0.1);
  }
  </style>