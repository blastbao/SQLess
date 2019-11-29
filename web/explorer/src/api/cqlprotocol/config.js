import axios from 'axios'

const api = axios.create({
  baseURL: process.env.API_CQL,
  headers: {
    'X-User-Agent': 'aletheia-webui'
  },
  withCredentials: true
})

console.log('api.cqlprotocol: ', process.env.API_CQL)

export default api
