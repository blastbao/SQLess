<template>
  <v-layout row wrap justify-center fill-count>
    <v-flex md8 sm8>
      <v-layout column>
        <v-card>
          <v-card-title primary-title>
            <v-layout column>
              <v-flex headline primary--text text--darken-1 text-sm-center text-md-center>
              </v-flex>
              <v-text-field v-model.trim="searchContent" hide-details label="Hashes / Block Count / Kayak Offset ..." single-line append-icon="mdi-magnify" color="primary darken-1" @input="search"></v-text-field>
            </v-layout>
          </v-card-title>
        </v-card>

        <v-card class="mt-1" v-if="!currentDatabase">
          <v-card-title primary-title>
            <v-flex caption text-sm-center red--text text--darken-1 font-weight-bold>NO DATABASE SPECIFIED</v-flex>
          </v-card-title>
        </v-card>

        <v-card class="mt-1" v-if="currentDatabase">
          <v-card-title primary-title>
            <!-- the search results -->
            <v-layout column v-if="searchContent !== ''">
              <v-flex caption>SEARCH RESULT</v-flex>
              <v-list dense class="mono">
                <v-list-tile avatar v-for="item in results" :key="item.target" :to="item.href">
                  <v-list-tile-content>
                    <v-list-tile-title>
                      <v-layout row wrap>
                        <v-flex md3 sm3 caption :class="'result-label-' + item.target">
                          {{ item.target }}
                        </v-flex>
                        <v-flex md9 sm9>
                          {{ item.result }}
                        </v-flex>
                      </v-layout>
                    </v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
            </v-layout>

            <!-- the latest blocks -->
            <v-layout column v-else>
              <v-flex caption>
                <v-layout row wrap>
                  <v-flex md6 xs6>
                    LATEST BLOCKS
                  </v-flex>
                  <v-flex md6 xs6 text-md-right text-xs-right>
                    <router-link class="font-weight-bold" :to="{name: 'BlockList', params:{db: currentDatabase}}">MORE</router-link>
                  </v-flex>
                </v-layout>
              </v-flex>
              <v-list dense class="mono">
                <v-list-tile avatar v-for="item in latestBlocks" :key="item.count" :to="{name:'Block', params:{db:currentDatabase, hash:item.count}}" :class="!!item.error ? 'red lighten-4' : ''">
                  <v-list-tile-title>
                    <v-layout row wrap>
                      <v-flex md2 sm2>#{{ item.count }}</v-flex>
                      <v-flex md6 sm6>
                        {{ humanReadableTime(item.timestamp) }}
                      </v-flex>
                      <v-flex md4 sm4 v-if="!!!item.error">
                        {{ substr(item.hash, 16) }}
                        <span class="has-sql primary--text" v-if="item.queries.length">
                          (SQL)
                        </span>
                      </v-flex>
                      <v-flex md4 sm4 v-else class="error--text text--darken-2">
                        {{ item.error.status }}
                      </v-flex>
                    </v-layout>
                  </v-list-tile-title>
                </v-list-tile>
              </v-list>
            </v-layout>
          </v-card-title>
        </v-card>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
// import { createNamespacedHelpers } from 'vuex'
// const { mapState } = createNamespacedHelpers('databases')
import { blocks, acks, requests } from '@/api/cqlprotocol'
import toolkit from '@/components/Utils/toolkit'

const NUM_SHOW_RECENT_BLOCKS = 10

export default {
  mixins: [toolkit],

  components: {
  },

  mounted () {
    this.currentDatabase = this.$route.params.db
    this.reload()
  },

  data () {
    return {
      searchContent: '',
      results: [],
      latestBlocks: []
    }
  },

  computed: {
    // ...mapState({
    //   currentDatabase: state => state.currentDatabase
    // })
    currentDatabase: {
      get () {
        return this.$store.state.databases.currentDatabase
      },
      set (newValue) {
        this.$store.dispatch('databases/setCurrentDatabase', newValue)
      }
    }
  },

  watch: {
    currentDatabase: function (to, from) {
      console.debug('watch(currentDatabase): ', from, ' --> ', to)
      this.reload()
    }
  },

  methods: {
    reload () {
      this.refreshLatestBlocks()
    },

    search () {
      console.log('searching...')
      if (!this.$store.state.databases.currentDatabase) {
        console.error('no database specified')
        return
        // TODO: give some error information
      }

      this.results = []
      if (this.searchContent === '') {
        return
      }

      console.info(`search: ${this.searchContent}`)
      // Define result object and show them in a list
      blocks
        .getBlockProxyV3(this.currentDatabase, this.searchContent, {size: 1})
        .then(resp => {
          let block = resp.data.data.block
          this.results.push({
            target: 'BLOCK',
            result: `${this.substr(block.hash, 32)} #${block.count}`,
            error: block.error,
            href: {
              name: 'Block',
              params: {
                db: this.currentDatabase,
                hash: block.count
              }
            }
          })
        })

      acks.getAckByHash(this.currentDatabase, this.searchContent).then(resp => {
        let ack = resp.data.data.ack
        this.results.push({
          target: 'ACK',
          result: this.substr(ack.hash, 32),
          href: {
            name: 'Ack',
            params: {
              db: this.currentDatabase,
              hash: ack.hash
            }
          }
        })
      })

      requests
        .getRequestProxy(this.currentDatabase, this.searchContent)
        .then(resp => {
          let request = resp.data.data.request
          this.results.push({
            target: 'REQUEST',
            result: this.substr(request.hash, 32),
            href: {
              name: 'Request',
              params: {
                db: this.currentDatabase,
                hash: request.hash
              }
            }
          })
        })
    },

    async refreshLatestBlocks () {
      this.latestBlocks = []
      let maxCount = await blocks.getMaxCount(this.currentDatabase)
      let startCount = maxCount + 1
      let endCount = startCount - NUM_SHOW_RECENT_BLOCKS
      endCount = endCount < 0 ? 0 : endCount
      let result = await blocks.getBlockList(
        this.currentDatabase,
        startCount,
        endCount
      )
      this.latestBlocks = result
    }
  }
}
</script>

<style lang="scss" scoped>
.result-label-BLOCK {
  color: #0f37d6;
  font-weight: bold;
}
.result-label-ACK {
  color: #ef3412;
  font-weight: bold;
}
.result-label-REQUEST {
  color: green;
  font-weight: bold;
}
.has-sql {
  font-size: 0.75em;
  font-weight: bold;
}
#logo {
  fill: #1272e4;
  height: 48px;
}
#logo:hover {
  fill: #0c60c3;
}
</style>
