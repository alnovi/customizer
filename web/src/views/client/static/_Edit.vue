<template>
  <div>
    <h2>{{ `${name}/${static_id}` }}</h2>
    <h3 class="mt-4">Новый элемент</h3>
    <div class="card">
      <card-create :client-id="id" :collection-name="static_id" @update="updateItem"></card-create>
    </div>
    <h3 class="mt-4">Список элементов</h3>
    <ul class="nav nav-tabs mt-4">
      <li class="nav-item" v-for="tab in tabs" :key="tab">
        <a v-on:click="setActiveTab(tab)" class="nav-link" :class="{ active: activeTab === tab }">
          {{ types[tab] }}
        </a>
      </li>
    </ul>
    <div class="list-wrapper">
      <div v-for="(item, index) in filteredCollection" class="card" :key="index">
        <card :keyItem="index" :client-id="id" :item="item" :collection-name="static_id" @update="updateItem"></card>
      </div>
    </div>
  </div>
</template>

<script>
import Card from './components/Card'
import CardCreate from '@/views/client/static/components/CardCreate'

const TYPES = {
  text: 'Текст',
  image: 'Изображение',
  color: 'Цвет',
  script: 'Скрипт',
  style: 'Стиль',
  svg: 'SVG',
  link: 'Ссылка',
}

export default {
  name: 'ClientStatic_Edit',
  components: {
    CardCreate,
    Card,
  },
  props: {
    id: String,
    static_id: String,
    name: String,
  },
  data() {
    return {
      collection: {},
      tabs: [],
      types: TYPES,
      activeTab: null,
    }
  },
  computed: {
    filteredCollection: function() {
      // if (this.cache[this.activeTab]) {
      //   return this.cache[this.activeTab]
      // }
      const result = {}
      Object.entries(this.collection).forEach(([key, value]) => {
        let subKey = key.split('.')[0]
        if (subKey === this.activeTab) {
          result[key] = value
        }
      })
      // this.cache[this.activeTab] = result
      return result
    },
  },
  methods: {
    load() {
      this.$api
        .getClientConfigStaticById(this.id, this.static_id)
        .then(res => {
          return res.json()
        })
        .then(res => {
          const { collection, presentKeys } = this.prepared(res)
          this.collection = collection
          this.tabs = presentKeys
          this.activeTab = presentKeys[0]
        })
        .catch(err => {
          console.log('Can not load collection', err)
        })
    },
    prepared(collection) {
      const KNOWN_KEYS = ['text', 'image', 'color', 'script', 'style', 'svg', 'link']
      const presentKeys = []
      for (let key in collection) {
        let subKey = key.split('.')[0]
        collection[key].type = KNOWN_KEYS.includes(subKey) ? subKey : 'text'
        if (!presentKeys.includes(subKey)) {
          presentKeys.push(subKey)
        }
      }
      // Sort keys, so Text will be the default active tab
      const sortedKeys = presentKeys.sort((a, b) => b.localeCompare(a))

      return { collection, presentKeys: sortedKeys }
    },
    updateItem(event) {
      if (event.updateType === 'delete') {
        this.$delete(this.collection, event.key)
      }
      if (event.updateType === 'create' || event.updateType === 'update') {
        this.$set(this.collection, event.key, {
          value: event.value,
          description: event.description,
          type: event.type,
        })
      }
    },
    setActiveTab(tab) {
      this.activeTab = tab
    },
  },
  mounted() {
    this.load()
  },
  watch: {},
}
</script>

<style scoped lang="scss">
.list-wrapper > div:first-of-type {
  margin-top: 0;
  border-top: transparent;
  border-top-left-radius: 0;
  border-top-right-radius: 0;
}
.nav-item {
  cursor: pointer;
}
.card {
  margin: 16px 0;

  .input-color {
    width: 38px;
  }
}
</style>
