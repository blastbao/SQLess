<template>
  <v-select v-model="currentDatabase" placeholder="Database Address" :items="databaseList" prepend-icon="mdi-database" dense hide-details :menu-props="{minWidth:'375px'}">
    <v-list-tile slot="prepend-item">
      <v-text-field v-model="newDatabase" placeholder="Use new database address" prepend-icon="mdi-database" @change="useNewDatabaseAddress"></v-text-field>
    </v-list-tile>
    <v-divider slot="prepend-item" class="mt-3"></v-divider>

    <template slot="item" slot-scope="{ item, tile }">
      {{ item.text }}
      <v-spacer></v-spacer>
      <v-btn small flat icon color="red darken-2" @click="removeDatabaseAddr(currentDatabase, item.value)">
        <v-icon small>mdi-delete</v-icon>
      </v-btn>
    </template>
  </v-select>
</template>

<script>
export default {
  data () {
    return {
      newDatabase: ''
    }
  },

  created () {
    this.$store.dispatch('databases/loadDatabaseList')
  },

  watch: {
    currentDatabase: function (to, from) {
      console.debug('watch(currentDatabase): ', from, ' --> ', to)
      if (from === null) {
        return
      }
      this.changeDatabase(to)
    }
  },

  computed: {
    currentDatabase: {
      get () {
        return this.$store.state.databases.currentDatabase
      },
      set (newValue) {
        this.$store.dispatch('databases/setCurrentDatabase', newValue)
      }
    },

    databaseList () {
      return this.$store.state.databases.databaseList.map(item => {
        return { value: item, text: `${item.substring(0, 16)} ... ${item.substring(56)}` }
      })
    }
  },

  methods: {
    changeDatabase (todb) {
      if (todb) {
        this.$router.push({name: 'Database', params: {db: todb}})
      } else {
        this.$router.push({name: 'Home'})
      }
    },

    useNewDatabaseAddress () {
      // this.$store.dispatch('databases/addNewDatabase', this.newDatabase)
      this.currentDatabase = this.newDatabase
      this.newDatabase = ''
    },

    removeDatabaseAddr (atAddr, removeAddr) {
      console.debug(`removeDatabaseAddr: ${removeAddr}, from ${atAddr}`)
      const toAddr = atAddr === removeAddr ? '' : atAddr
      setTimeout(() => {
        this.changeDatabase(toAddr)
        this.$store.dispatch('databases/removeDatabase', removeAddr)
      }, 300)
    }
  }
}
</script>
