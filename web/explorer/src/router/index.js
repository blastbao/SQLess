import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Block from '@/components/Blocks/Block'
import BlockList from '@/components/Blocks/BlockList'
import Ack from '@/components/Acks/Ack'
import Request from '@/components/Requests/Request'

Vue.use(Router)

const rv = {
  template: `<router-view></router-view>`
}

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/dbs/:db',
      component: rv,
      children: [
        {
          path: '',
          name: 'Database',
          component: Home
        },
        {
          path: 'blocks',
          component: rv,
          children: [
            {
              path: '',
              name: 'BlockList',
              component: BlockList
            },
            {
              path: ':hash',
              name: 'Block',
              component: Block
            }
          ]
        },
        {
          path: 'acks',
          component: rv,
          children: [
            {
              path: ':hash',
              name: 'Ack',
              component: Ack
            }
          ]
        },
        {
          path: 'requests',
          component: rv,
          children: [
            {
              path: ':hash',
              name: 'Request',
              component: Request
            }
          ]
        }
      ]
    },
    {
      path: '*',
      redirect: '/'
    }
  ]
})
