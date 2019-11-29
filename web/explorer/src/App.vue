<template>
  <div id="app">
    <v-app>
      <v-toolbar color="primary" app fixed prominent dark class="pt-1">
        <v-toolbar-title>
          <router-link :to="{name:'Database', params:{db: currentDatabase}}" class="site-title">CQL Explorer</router-link>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items>
          <database-selector></database-selector>
        </v-toolbar-items>
      </v-toolbar>
      <v-content>
        <v-container fluid>
          <router-view></router-view>
        </v-container>
      </v-content>
    </v-app>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex'
import DatabaseSelector from '@/components/DatabaseSelector'
const { mapState } = createNamespacedHelpers('databases')

export default {
  name: 'App',

  components: {
    'database-selector': DatabaseSelector
  },

  computed: {
    ...mapState({
      currentDatabase: state => state.currentDatabase
    })
  }
}
</script>

<style>
.mono {
  font-family: Menlo, Consolas, monospace;
}
.site-title {
  color: white;
  text-decoration: none;
}
.sql-type-read {
  color: #388e3c;
}
.sql-type-write {
  color: #d32f2f;
}
</style>
