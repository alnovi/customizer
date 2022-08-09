<template>
  <div style="margin: 20px 0">
    <form @submit="send">
      <div class="form-group">
        <label for="default-static-name">Название блока</label>
        <input ref="default-static-name" pattern="[a-zA-Z0-9-]+" type="text" class="form-control" id="default-static-name" v-model="name" required />
      </div>

      <vue-json-editor @json-change="onJsonChange" @has-error="onJsonError" ref="jsonEdit" v-model="json" :modes="['view', 'code']" :mode="'code'" />

      <button :disabled="isErrorInJson" style="margin: 20px 0" type="submit" class="btn btn-outline-primary">Сохранить</button>
    </form>
  </div>
</template>

<script>
import vueJsonEditor from 'vue-json-editor'

export default {
  name: 'DefaultStatic_Create',
  data() {
    return {
      name: '',
      json: {},
      isErrorInJson: false,
    }
  },
  components: {
    vueJsonEditor,
  },
  methods: {
    onJsonChange() {
      // Reset error state on json change.
      // If there is still an error, onJsonError will handle it.
      this.isErrorInJson = false
    },
    onJsonError(errorText) {
      // errorText returns an Error instance with details of error
      // Cast to boolean which is user in validation.
      this.isErrorInJson = Boolean(errorText)
    },
    send(event) {
      event.preventDefault()

      this.$api
        .storeDefaultConfig('statics', this.name, this.json)
        .then(res => {
          this.$router.push({ name: 'default.static' })
        })
        .catch(err => {
          console.log('Can not save', err)
        })
    },
  },
}
</script>

<style scoped lang="scss"></style>
