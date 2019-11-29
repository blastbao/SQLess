<template>
  <v-layout row wrap justify-center>
    <v-flex md8 xs12 v-if="error">
      <sp-error-card :error="error"></sp-error-card>
    </v-flex>
    <v-flex v-if="request"><v-card>
      <v-card-title primary-title>
        <div class="headline">Request - {{ request.hash.substring(0, 12) }}</div>
      </v-card-title>
      <v-card-text>
        <v-container grid-list-md pt-0>
          <v-layout row wrap class="mono">
            <v-flex md2 xs2 caption>Time</v-flex>
            <v-flex md10 xs10>
              {{ humanReadableTime(request.timestamp) }} ({{ request.timestamp }})
            </v-flex>

            <v-flex md2 xs2 caption>Hash</v-flex>
            <v-flex md10 xs10>
              {{ request.hash }}
            </v-flex>

            <v-flex md2 xs2 caption>Node</v-flex>
            <v-flex md10 xs10>{{ request.node }}</v-flex>

            <v-flex md2 xs2 caption>Type</v-flex>
            <v-flex md10 xs10 :class="`sql-type-${request.type}`">{{ request.type }}</v-flex>

            <v-flex md2 xs2 caption>Count</v-flex>
            <v-flex md10 xs10>{{ request.count }}</v-flex>

            <v-flex md2 xs2 caption>Queries</v-flex>
            <v-flex md10 xs10>
              <v-data-table
                :headers="queryListHeaders"
                :items="request.queries"
                hide-actions
              >
                <template slot="items" slot-scope="props">
                  <tr>
                    <td>{{ props.index }}</td>
                    <td class="primary--text font-weight-bold">
                      <pre>{{ props.item.pattern }}</pre>
                    </td>
                    <td>
                      <pre>{{ props.item.args }}</pre>
                    </td>
                  </tr>
                </template>
              </v-data-table>
            </v-flex>
          </v-layout>
        </v-container>
      </v-card-text>
    </v-card></v-flex>
  </v-layout>
</template>

<script>
import { requests } from '@/api/cqlprotocol'
import SPErrorCard from '@/components/SPErrorCard'
import toolkit from '@/components/Utils/toolkit'

export default {
  mixins: [toolkit],

  components: {
    'sp-error-card': SPErrorCard
  },

  mounted () {
    this.currentDatabase = this.$route.params.db
    this.refreshRequest()
  },

  data () {
    return {
      request: null,
      error: null,
      queryListHeaders: [
        { text: 'Index', value: '', sortable: false },
        { text: 'SQL', value: 'pattern', sortable: false },
        { text: 'Args', value: 'args', sortable: false }
      ]
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
    }
  },

  methods: {
    async refreshRequest () {
      try {
        let resp = await requests.getRequestProxy(
          this.$route.params.db,
          this.$route.params.hash
        )
        this.request = resp.data.data.request
      } catch (error) {
        this.error = error.response.data
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
