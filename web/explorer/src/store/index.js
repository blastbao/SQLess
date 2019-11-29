// https://vuex.vuejs.org/guide/structure.html
// Here we assemble modules and export the store.

import Vue from 'vue'
import Vuex from 'vuex'

import databases from './modules/databases'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    databases
  }
})
