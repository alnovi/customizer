<template>
  <div class="card card-body" style="margin: 20px 0">
    <form @submit="send">
      <div class="form-group">
        <label for="client-id">ИД клиента</label>
        <input ref="client-id"
               pattern="[a-zA-Z0-9-]+"
               minlength="10"
               maxlength="10"
               type="text"
               class="form-control"
               id="client-id"
               v-model="id"
               required>
      </div>
      <div class="form-group">
        <label for="client-name">Название клиента</label>
        <input ref="client-name"
               type="text"
               maxlength="60"
               class="form-control"
               id="client-name"
               v-model="name"
               required>
      </div>

      <div class="form-group">
        <label for="client-token">Токен клиента</label>
        <input ref="client-token"
               pattern="[a-zA-Z0-9-]+"
               minlength="36"
               maxlength="36"
               type="text"
               class="form-control"
               id="client-token"
               v-model="token"
               required>
      </div>

      <button style="margin: 20px 0" type="submit" class="btn btn-outline-primary">Сохранить</button>
    </form>
  </div>
</template>

<script>
export default {
  name: "Client_Create",
  data() {
    return {
      id: '',
      name: '',
      token: '',
    }
  },
  methods: {
    strRand(length) {
      const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
      const charactersLength = characters.length;
      let result = '';

      for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
      }

      return result;
    },
    send(event) {
      event.preventDefault()

      this.$api.storeClient(this.id, this.name, this.token)
        .then(res => {
          this.$router.push({name: 'clients'})
        })
        .catch(err => {
          console.log('Can not save', err)
        })
    }
  },
  created() {
    this.id = this.strRand(10)
    this.token = this.strRand(36)
  }
}
</script>

<style scoped>

</style>
