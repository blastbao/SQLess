import moment from 'moment'

export default {
  methods: {
    substr (s, n) {
      return s.substring(0, n)
    },

    humanReadableTime (t) {
      return moment(t).format('YYYY-MM-DD HH:mm:ss')
    }
  }
}
