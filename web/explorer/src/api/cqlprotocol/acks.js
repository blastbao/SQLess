import api from './config'
import requests from './requests'

export default {
  getAckByHash (db, hash) {
    return api.get(`/v1/ack/${db}/${hash}`)
  },

  async getSQLQueryOfAck (db, ackHash) {
    let sqlQuery = []
    try {
      const ackResp = await this.getAckByHash(db, ackHash)
      const ack = ackResp.data.data.ack
      // Now, ignore "read" type
      if (ack.request.type === 'read') {
        return sqlQuery
      }

      const requestResp = await requests.getRequestByHash(db, ack.request.hash)
      const request = requestResp.data.data.request
      sqlQuery = request.queries.map(x => {
        return {
          ackHash: ackHash,
          requestHash: request.hash,
          nodeHash: request.node,
          type: request.type,
          sql: x,
          timestamp: request.timestamp
        }
      })
    } catch (error) {
      return [
        {
          ackHash: ackHash,
          requestHash: '',
          nodeHash: '',
          type: '',
          sql: '',
          timestamp: 0,
          error: error.response.data
        }
      ]
    }
    return sqlQuery
  },

  async gatherSQLQueriesOfAcks (db, ackHashes) {
    let queryList = []
    if (!ackHashes) {
      return queryList
    }

    let promises = ackHashes.map(x => this.getSQLQueryOfAck(db, x))
    let result = await Promise.all(promises)
    queryList = result.reduce((acc, x) => {
      acc = acc.concat(x)
      return acc
    }, [])
    return queryList
  }
}
