import api from './config'

export default {
  getBlockProxy (db, key) {
    if (key.length !== 64) {
      return this.getBlockByCount(db, key)
    }
    return this.getBlockByHash(db, key)
  },

  getBlockByCount (db, count) {
    return api.get(`/v1/count/${db}/${count}`)
  },

  getBlockByHash (db, hash) {
    return api.get(`/v1/block/${db}/${hash}`)
  },

  getBlockProxyV3 (db, key, params) {
    if (key.length !== 64) {
      return this.getBlockByCountV3(db, key, params)
    }
    return this.getBlockByHashV3(db, key, params)
  },

  getBlockByCountV3 (db, count, params) {
    return api.get(`/v3/count/${db}/${count}`, { params: params })
  },

  getBlockByHashV3 (db, hash, params) {
    return api.get(`/v3/block/${db}/${hash}`, { params: params })
  },

  async getBlockByCountIgnoreError (db, count) {
    try {
      const resp = await this.getBlockByCountV3(db, count, { size: 1 })
      let block = resp.data.data.block
      // if (block.count >= 90 && block.count < 94) {
      //   block.error = { status: 'mock error' }
      // }
      if (block.pagination) {
        block.totalQueries = block.pagination.total || 0
      } else {
        block.totalQueries = block.queries.length
      }
      return block
    } catch (error) {
      const errorResp = error.response.data
      console.error(`get block ${count} failed: ${JSON.stringify(errorResp)}`)
      return {
        count: count,
        hash: '',
        height: 0,
        producer: '',
        queries: [],
        timestamp: 0,
        version: 0,
        error: errorResp,
        totalQueries: 0
      }
    }
  },

  async getBlockList (db, startCount, endCount) {
    let blockList = []
    if (!db || startCount === endCount) {
      return blockList
    }

    let promises = []

    if (startCount < endCount) {
      for (let i = startCount; i < endCount; i++) {
        promises.push(this.getBlockByCountIgnoreError(db, i))
      }
    } else {
      for (let i = startCount - 1; i >= endCount; i--) {
        promises.push(this.getBlockByCountIgnoreError(db, i))
      }
    }

    blockList = await Promise.all(promises)
    return blockList
  },

  async getMaxCount (db) {
    let result = await api.get(`/v2/head/${db}`)
    return result.data.data.block.count
  }
}
