<template>
  <v-layout row wrap justify-center>
    <v-flex md8 sm12>
      <sp-error-card :error="error"></sp-error-card>
      <v-card v-if="ack">
        <v-card-title primary-title>
          <div class="headline">ACK - {{ substr(ack.hash, 12) }}</div>
          <v-container grid-list-md>
            <v-layout row wrap class="mono">
              <v-flex md2 sm2 caption>Time</v-flex>
              <v-flex md10 sm10>
                {{ humanReadableTime(ack.timestamp) }} ({{ ack.timestamp }})
              </v-flex>

              <v-flex md2 sm2 caption>Hash</v-flex>
              <v-flex md10 sm10>
                {{ ack.hash }}
              </v-flex>

              <v-flex md2 sm2 caption>Node</v-flex>
              <v-flex md10 sm10>{{ ack.node }}</v-flex>

              <v-flex md2 sm2 caption>Request</v-flex>
              <v-flex md10 sm10></v-flex>
              <v-flex md2 sm2 caption pl-3>-Time</v-flex>
              <v-flex md10 sm10>
                {{ humanReadableTime(ack.request.timestamp) }} ({{ ack.request.timestamp }})
              </v-flex>
              <v-flex md2 sm2 caption pl-3>-Hash</v-flex>
              <v-flex md10 sm10>
                <router-link :to="{name:'Request', params:{db:$route.params.db,hash:$route.params.hash}}">{{ ack.request.hash }}</router-link>
              </v-flex>
              <v-flex md2 sm2 caption pl-3>-Node
              </v-flex>
              <v-flex md10 sm10>{{ ack.request.node }}</v-flex>
              <v-flex md2 sm2 caption pl-3>-Type</v-flex>
              <v-flex md10 sm10>{{ ack.request.type }}</v-flex>
              <v-flex md2 sm2 caption pl-3>-Count</v-flex>
              <v-flex md10 sm10>{{ ack.request.count }}</v-flex>

              <v-flex md2 sm2 caption>Response</v-flex>
              <v-flex md10 sm10></v-flex>
              <v-flex md2 sm2 caption pl-3>-Time</v-flex>
              <v-flex md10 sm10>
                {{ humanReadableTime(ack.response.timestamp) }} ({{ ack.response.timestamp }})
              </v-flex>
              <v-flex md2 sm2 caption pl-3>-Hash</v-flex>
              <v-flex md10 sm10>
                {{ ack.response.hash }}
              </v-flex>
              <v-flex md2 sm2 caption pl-3>-Node</v-flex>
              <v-flex md10 sm10>
                {{ ack.response.node }}
              </v-flex>
              <v-flex md2 sm2 caption pl-3>-LogPosition</v-flex>
              <v-flex md10 sm10>
                {{ ack.response.log_position }}
              </v-flex>
            </v-layout>
          </v-container>
        </v-card-title>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
import { acks } from '@/api/cqlprotocol'
import SPErrorCard from '@/components/SPErrorCard'
import toolkit from '@/components/Utils/toolkit'

export default {
  mixins: [toolkit],

  components: {
    'sp-error-card': SPErrorCard
  },

  mounted () {
    this.currentDatabase = this.$route.params.db
    this.refreshAck()
  },

  data () {
    return {
      // ack: {
      //   hash: '',
      //   request: {
      //     hash: ''
      //   },
      //   response: {
      //     hash: ''
      //   }
      // },
      ack: null,
      error: null
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
    async refreshAck () {
      try {
        let resp = await acks.getAckByHash(
          this.$route.params.db,
          this.$route.params.hash
        )
        this.ack = resp.data.data.ack
      } catch (error) {
        this.error = error.response.data
      }
    }
  }
}
</script>
