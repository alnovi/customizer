import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import {BootstrapVue, IconsPlugin} from 'bootstrap-vue'

import './assets/scss/app.scss'
import {Api} from "@/services/api";

// Install BootstrapVue
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)

Vue.config.productionTip = false

fetch('/config.json')
  .then((res) => res.json())
  .then((config) => {
    config.host = location.origin
    Vue.prototype.$config = config
    Vue.prototype.$api = new Api({
      host: config.host
    })

    new Vue({
      router,
      store,
      render: h => h(App)
    }).$mount('#app')
  })
