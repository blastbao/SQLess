import api from './config'

export default {
  getRequestProxy (db, key) {
    if (key.length !== 64) {
      return this.getRequestByKayakOffset(db, key)
    }
    return this.getRequestByHash(db, key)
  },

  getRequestByKayakOffset (db, offset) {
    return api.get(`/v1/offset/${db}/${offset}`)
  },

  getRequestByHash (db, hash) {
    return api.get(`/v1/request/${db}/${hash}`)
  }
}
