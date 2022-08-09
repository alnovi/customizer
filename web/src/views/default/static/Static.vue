<template>
  <div class="mt-4">
    <div class="d-flex flex-row justify-content-between align-items-center">
      <h2 class="d-flex">Список статик</h2>
      <div class="d-flex">
        <router-link
          :to="{name:'default.static.create'}"
          v-slot="{ href, isActive }">
          <b-button :href="href">Добавить новую статику</b-button>
        </router-link>
      </div>
    </div>
    <table class="table mt-2">
      <thead>
        <tr>
          <th>Название</th>
          <th class="text-right">Элементов</th>
        </tr>
      </thead>
      <tbody>
      <tr v-for="item of configs" :key="item._id">
        <td>
          <router-link :to="{name: 'default.static.edit', params: {name: item._id}}">
            {{item._id}}
          </router-link>
        </td>

        <td class="text-right">{{ item.fields}}</td>
      </tr>
      </tbody>
    </table>
  </div>

</template>

<script>
export default {
  name: "DefaultStatic",
  data() {
    return {
      configs: [],
    }
  },
  methods: {
    load() {
      this.$api.getDefaultConfigs('static')
        .then(res => {
          return res.json()
        })
        .then(res => {
          this.configs = res
        })
        .catch(err => {
          console.log('Cann not load configs', err)
        })
    }
  },
  mounted() {
    this.load();
  }
}
</script>

<style scoped lang="scss">
.items {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  padding: 0 15px;

  .item {
    width: 100%;
    display: flex;
    flex-direction: row;
    height: 80px;
    font-size: 22px;
    align-items: center;
    border: 1px solid #333;
    margin: 10px 0;
    padding: 5px 10px;
    border-radius: 3px;
    justify-content: space-between;
    background-color: #a5a5a5;

    p {
      margin: 0;
      padding: 0;
    }

    .name {
      font-weight: bold;
      margin: 0 15px;
      color: #333333;
    }

    .count {
      color: #fcfcfc;
      font-weight: bold;
      margin: 0 15px;
    }
  }
}
</style>
