<template>
  <v-layout rows wrap justify-center>
    <v-flex md8 xs12 v-if="error">
      <sp-error-card :error="error"></sp-error-card>
    </v-flex>
    <v-card v-else>
      <v-card-title primary-title>
        <div class="headline">
          Block #{{ block.count }}
        </div>
      </v-card-title>
      <v-progress-linear :active="pageLoading" :indeterminate="true" color="primary" height="3"></v-progress-linear>
      <v-card-text>
        <v-container grid-list-md pt-0>
          <v-layout row wrap class="mono" justify-center>
            <!-- block info panel -->
            <v-flex md8 xs12>
              <v-layout row wrap>
                <v-flex md2 xs2 caption>Height</v-flex>
                <v-flex md10 xs10>
                  <router-link :to="{name: 'Block', params: {db: currentDatabase, hash: block.count-1}}" class="block-nav-link">&lt;&lt;</router-link>
                  {{ block.count }}
                  <router-link :to="{name: 'Block', params: {db: currentDatabase, hash: block.count+1}}" class="block-nav-link">&gt;&gt;</router-link>

                  <router-link :to="{name: 'BlockList', params: {db: currentDatabase}}">List</router-link>
                </v-flex>

                <v-flex md2 xs2 caption>Time</v-flex>
                <v-flex md10 xs10>
                  {{ humanReadableTime(block.timestamp) }} ({{ block.timestamp }})
                </v-flex>

                <v-flex md2 xs2 caption>Hash</v-flex>
                <v-flex md10 xs10>
                  {{ block.hash }}
                </v-flex>

                <v-flex md2 xs2 caption>Version</v-flex>
                <v-flex md10 xs10>{{ block.version }}</v-flex>

                <v-flex md2 xs2 caption>Producer</v-flex>
                <v-flex md10 xs10>
                  {{ block.producer }}
                </v-flex>
              </v-layout>
            </v-flex>

            <!-- divider -->
            <v-flex md12 xs12 mt-3>
              <v-divider></v-divider>
            </v-flex>

            <!-- SQL list -->
            <v-flex md12 xs12 title my-3>
              SQL
            </v-flex>

            <v-flex md12 xs12 v-if="sqlList.length === 0" grey--text>
              None
            </v-flex>
            <v-flex md12 xs12 v-else>
              <v-data-table :headers="sqlListHeaders" :items="sqlList" item-key="requestHash" :total-items="total" :pagination.sync="pagination" :loading="loading" :rows-per-page-items="[10,20,30,50]">
                <template slot="items" slot-scope="props">
                  <tr>
                    <td>{{ humanReadableTime(props.item.timestamp) }}</td>
                    <td>
                      <router-link :to="{name: 'Request', params: {db: currentDatabase, hash: props.item.requestHash}}">
                        {{ substr(props.item.requestHash, 16) }}
                      </router-link>
                    </td>
                    <td :class="`sql-type-${props.item.type}`">{{ props.item.type }}</td>
                    <td>
                      <pre>{{ props.item.sql.pattern }}</pre>
                    </td>
                    <td>{{ props.item.sql.queries }}</td>
                    <td>
                      <v-btn flat icon color="primary" :to="{name: 'Request', params: {db: currentDatabase, hash: props.item.requestHash}}">
                        <v-icon small>mdi-eye</v-icon>
                      </v-btn>
                    </td>
                  </tr>
                </template>
              </v-data-table>
            </v-flex>
          </v-layout>
        </v-container>
      </v-card-text>
    </v-card>
  </v-layout>
</template>

<script>
import { blocks } from '@/api/cqlprotocol'
import SPErrorCard from '@/components/SPErrorCard'
import toolkit from '@/components/Utils/toolkit'

export default {
  mixins: [toolkit],
  components: {
    'sp-error-card': SPErrorCard
  },

  mounted () {
    this.initialize()
  },

  data () {
    return {
      block: {
        count: 0,
        height: 0,
        genesis_hash: '',
        hash: '',
        producer: '',
        queries: [],
        timestamp: 0,
        version: 0
      },
      sqlList: [],
      sqlListHeaders: [
        { text: 'Time', value: 'time', sortable: false },
        { text: 'Hash', value: 'hash', sortable: false },
        { text: 'Type', value: 'type', sortable: false },
        { text: 'SQL Preview', value: 'sql', sortable: false },
        { text: 'Queries', value: 'queries', sortable: false },
        { text: 'Details', value: '', sortable: false }
      ],
      error: null,
      pagination: {
        page: 1,
        rowsPerPage: 10,
        descending: false,
        sortBy: null,
        totalItem: 0
      },
      total: 0,
      loading: false,
      pageLoading: true
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
    initialize () {
      this.error = null
      this.currentDatabase = this.$route.params.db
      this.refreshBlock()
    },

    async refreshBlock () {
      const { page, rowsPerPage } = this.pagination
      this.loading = true
      let params = {
        page: page,
        size: rowsPerPage
      }
      try {
        // NB: the ":hash" variable also can be an integer (block count)
        let resp = await blocks.getBlockProxyV3(
          this.currentDatabase,
          this.$route.params.hash,
          params
        )
        this.block = resp.data.data.block
        this.sqlList = this.block.queries.map(x => {
          return {
            requestHash: x.request.hash,
            type: x.request.type,
            timestamp: x.request.timestamp,
            sql: {
              pattern: (() => {
                let overflow = false
                let patternExcerpt = ''
                for (let i = 0; i < x.request.queries.length; i++) {
                  if (patternExcerpt !== '') {
                    patternExcerpt += '\n'
                  }
                  patternExcerpt += x.request.queries[i].pattern
                  if (patternExcerpt.length > 200) {
                    overflow = true
                    break
                  }
                }
                return patternExcerpt.substring(0, 200) + (overflow ? '…' : '')
              })(),
              queries: x.request.queries.length
              // args: x.request.queries.reduce((acc, x) => acc + x.args.length, 0)
            }
          }
        })
        if (this.block.pagination) {
          this.total = this.block.pagination.total
        } else {
          this.total = this.sqlList.length
        }
        console.debug('Gather SQL list:', this.sqlList)
      } catch (error) {
        console.error(error)
        // this.error = error.response.error
      }
      this.loading = false
      this.pageLoading = false
    }
  },

  watch: {
    $route (to, from) {
      console.debug(`watch($route): ${from.path} -> ${to.path}`)
      this.initialize()
    },

    pagination: {
      handler () {
        this.refreshBlock()
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.block-nav-link {
  text-decoration: none;
  font-weight: bold;
}
</style>
