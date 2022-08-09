<template>
  <div class="mt-4">
    <div class="d-flex flex-row justify-content-between align-items-center">
      <h2 class="d-flex">Список клиентов</h2>
      <div class="d-flex">
        <router-link :to="{ name: 'clients.create' }" v-slot="{ href, isActive }">
          <b-button :href="href">Добавить нового клиента</b-button>
        </router-link>
      </div>
    </div>
    <table class="table mt-2">
      <thead>
        <tr>
          <th>Название</th>
          <th>Id</th>
          <th>Secret</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item of clients" :key="item._id">
          <td class="w-100">
            <router-link :to="{ name: 'clients.statics', params: { id: item.id, name: item.name } }">
              {{ item.name }}
            </router-link>
          </td>
          <td>{{ item.id }}</td>
          <td>{{ item.secret }}</td>
          <td class="text-right">
            <div style="cursor: pointer;" v-on:click="handleClientDelete(item.id, item.name)">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash-fill" viewBox="0 0 16 16">
                <path
                  d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1H2.5zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5zM8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5zm3 .5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 1 0z"
                />
              </svg>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'Client',
  data() {
    return {
      clients: [],
    }
  },
  methods: {
    load() {
      this.$api
        .clientList(1)
        .then(res => {
          return res.json()
        })
        .then(res => {
          this.clients = res
        })
        .catch(err => {
          console.log('can not load clients', err)
        })
    },
    handleClientDelete(clientId, clientName) {
      let isConfirmed = window.confirm(`Вы действительно хотите удалить клиента ${clientName}?`)
      if (isConfirmed) {
        this.$api
          .deleteClient(clientId)
          .then(result => {
            console.log('Successfully deleted client')
            // Ghetto refetch
            this.load()
          })
          .catch(error => console.log('Failed to delete a client'))
      }
    },
  },
  created() {
    this.load()
  },
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
