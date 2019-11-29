<template>
  <v-layout row wrap justify-center>
    <v-flex md10 sm12>
      <v-data-table :headers="blockListHeaders" :items="blockList" :total-items="pagination.totalItems" :pagination.sync="pagination" :loading="loading" class="elevation-1" item-key="count" :rows-per-page-items="paginationSizeList" disable-page-reset>
        <template slot="items" slot-scope="props">
          <tr class="mono" :class="!!props.item.error ? 'red lighten-4' : ''">
            <td>{{ props.item.count }}</td>
            <td>{{ humanReadableTime(props.item.timestamp) }}</td>
            <td v-if="!!!props.item.error">
              {{ substr(props.item.hash, 32) }}
            </td>
            <td v-else class="error--text text--darken-2">
              {{ props.item.error.status }}
            </td>
            <td :class="props.item.totalQueries > 0 ? 'primary--text font-weight-bold' : 'grey--text'">{{ props.item.totalQueries > 0 ? props.item.totalQueries : 'none' }}</td>
            <td>
              <v-btn flat icon color="primary" :to="{name: 'Block', params: {db: currentDatabase, hash: props.item.count}}" :disabled="!!props.item.error">
                <v-icon small>mdi-eye</v-icon>
              </v-btn>
            </td>
          </tr>
        </template>
      </v-data-table>
    </v-flex>
  </v-layout>
</template>

<script>
import { blocks } from '@/api/cqlprotocol'
import toolkit from '@/components/Utils/toolkit'

export default {
  mixins: [toolkit],

  mounted () {
    this.currentDatabase = this.$route.params.db
    // this.loadURL()
  },

  data () {
    return {
      blockList: [],
      blockListHeaders: [
        { text: '# Count', value: 'count' },
        { text: 'Time', value: 'time', sortable: false },
        { text: 'Hash', value: 'hash', sortable: false },
        { text: 'SQL', value: 'sql', sortable: false },
        { text: 'Details', Value: 'view', sortable: false }
      ],
      paginationSizeList: [10, 20, 30, 50],
      pagination: {
        descending: true,
        page: 1,
        rowsPerPage: 10,
        sortBy: 'count',
        totalItems: 0
      },
      loading: false
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

  watch: {
    'pagination': {
      handler: function () {
        this.refreshBlockList()
      }
    }
  },

  methods: {
    loadURL () {
      this.pagination.descending = !this.$route.query.asc
      this.pagination.page = parseInt(this.$route.query.page) || 1
      let size = parseInt(this.$route.query.size) || 10
      this.pagination.rowsPerPage = this.paginationSizeList.includes(size) ? size : 10
    },

    async refreshBlockList () {
      const { sortBy, descending, page, rowsPerPage } = this.pagination
      this.loading = true
      this.pagination.totalItems = (await blocks.getMaxCount(this.currentDatabase)) + 1
      let startCount = (page - 1) * rowsPerPage
      let endCount = startCount + rowsPerPage
      if (endCount > this.pagination.totalItems) {
        endCount = this.pagination.totalItems
      }
      if (startCount < 0) {
        startCount = 0
      }

      if (sortBy === 'count' && descending) {
        startCount = this.pagination.totalItems - (page - 1) * rowsPerPage
        endCount = startCount - rowsPerPage
        endCount = endCount < 0 ? 0 : endCount
      }
      console.debug(
        `refreshBlockList, total=${
          this.pagination.totalItems
        }, page=${page}, rowsPerPage=${rowsPerPage}, startCount=${startCount}, endCount=${endCount}, sortBy=${sortBy}`
      )
      this.blockList = await blocks.getBlockList(
        this.currentDatabase,
        startCount,
        endCount
      )
      this.loading = false
    }
  }
}
</script>
