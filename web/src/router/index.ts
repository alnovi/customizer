import Vue from 'vue'
import VueRouter, {RouteConfig} from 'vue-router'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    component: () => import('../layout/Main/Main.vue'),
    children: [
      {
        name: 'default.static',
        path: 'default/statics',
        component: () => import('../views/default/static/Static.vue'),
      },
      {
        name: 'default.static.create',
        path: 'default/statics/create',
        component: () => import('../views/default/static/_Create.vue'),
      },
      {
        name: 'default.static.edit',
        path: 'default/statics/:name',
        component: () => import('../views/default/static/_Edit.vue'),
        props: true,
      },
      {
        name: 'clients',
        path: 'clients',
        component: () => import('../views/client/Client.vue'),
        props: true,
      },
      {
        name: 'clients.create',
        path: 'clients/create',
        component: () => import('../views/client/_Create.vue'),
        props: true,
      },
      {
        name: 'clients.statics',
        path: 'clients/:id/statics',
        component: () => import('../views/client/static/Static.vue'),
        props: true,
      },
      {
        name: 'clients.statics.create',
        path: 'clients/:id/statics/create',
        component: () => import('../views/client/static/_Create.vue'),
        props: true,
      },
      {
        name: 'clients.statics.edit',
        path: 'clients/:id/statics/:static_id',
        component: () => import('../views/client/static/_Edit.vue'),
        props: true,
      }
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
