// database history

const state = {
  databaseList: [],
  currentDatabase: null
}

const mutations = {
  LOAD_DATABASE_LIST (state) {
    let addrs = []
    try {
      addrs = JSON.parse(localStorage.getItem('addrs'))
      if (!Array.isArray(addrs)) {
        addrs = []
      }
    } catch (error) {}
    state.databaseList = addrs
    console.debug('load database list:', JSON.stringify(state.databaseList))
  },

  SET_CURRENT_DATABASE (state, db) {
    if (!db || db.length !== 64) {
      state.currentDatabase = ''
      return
    }
    state.currentDatabase = db
    console.debug('store.databases.SET_CURRENT_DATABASE:', db)
  },

  ADD_NEW_DATABASE (state, db) {
    if (!db || db.length !== 64) {
      return
    }
    if (!state.databaseList.includes(db)) {
      let copyList = state.databaseList.slice()
      copyList.push(db)
      state.databaseList = copyList
      console.debug('store.databases.ADD_NEW_DATABASE:', db)
    }
  },

  REMOVE_DATABASE (state, db) {
    let found = state.databaseList.indexOf(db)
    if (found > -1) {
      state.databaseList.splice(found, 1)
      console.debug('store.databases.REMOVE_DATABASE:', db)
    }
  }
}

const actions = {
  setCurrentDatabase ({ commit }, payload) {
    commit('SET_CURRENT_DATABASE', payload)
    commit('ADD_NEW_DATABASE', payload)
  },

  loadDatabaseList ({ commit }) {
    commit('LOAD_DATABASE_LIST')
  },

  addNewDatabase ({ commit }, payload) {
    commit('ADD_NEW_DATABASE', payload)
  },

  removeDatabase ({ commit }, payload) {
    commit('REMOVE_DATABASE', payload)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
