// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import store from './store'
import router from './router'

import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import '@mdi/font/css/materialdesignicons.min.css'

Vue.use(Vuetify, {
  theme: {
    primary: '#1272e4'
  }
})

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  components: { App },
  template: '<App/>'
})

window.onbeforeunload = function () {
  localStorage.setItem('lastAddr', store.state.databases.currentDatabase)
  localStorage.setItem(
    'addrs',
    JSON.stringify(store.state.databases.databaseList)
  )
}
