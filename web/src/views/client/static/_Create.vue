<template>
  <div style="margin: 20px 0">
    <form @submit="send">
      <div class="form-group">
        <label for="static-name">Название блока</label>
        <input ref="static-name" pattern="[a-zA-Z0-9-]+" type="text" class="form-control" id="static-name" v-model="name" required />
      </div>

      <div class="input-group my-2">
        <div class="input-group-prepend">
          <label class="input-group-text" for="default-static-configs">Родитель</label>
        </div>
        <select class="custom-select" id="default-static-configs" v-model="parent_static">
          <option v-for="(val, key) in default_statics" :value="val._id" :selected="parent_static === val._id">
            {{ val._id }}
          </option>
        </select>
      </div>

      <vue-json-editor @json-change="onJsonChange" @has-error="onJsonError" ref="jsonEdit" v-model="json" :modes="['view', 'code']" :mode="'code'" />

      <button :disabled="isErrorInJson" style="margin: 20px 0" type="submit" class="btn btn-outline-primary">Сохранить</button>
    </form>
  </div>
</template>

<script>
import vueJsonEditor from 'vue-json-editor'

export default {
  name: 'ClientStatic_Create',
  props: {
    id: String,
  },
  data() {
    return {
      default_statics: [],
      name: '',
      parent_static: '',
      json: {},
      isErrorInJson: false,
    }
  },
  components: {
    vueJsonEditor,
  },
  methods: {
    send(event) {
      event.preventDefault()
      this.$api
        .storeClientConfig(this.id, 'static', this.name, this.parent_static, this.json)
        .then(res => {
          this.$router.push({ name: 'clients.statics' })
        })
        .catch(err => {
          console.log('Can not save', err)
        })
    },
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
    loadDefaultStatics() {
      this.$api
        .getDefaultConfigs('static')
        .then(res => {
          return res.json()
        })
        .then(res => {
          this.default_statics = res

          if (this.default_statics.length > 0) {
            this.parent_static = this.default_statics[0]._id ?? ''
          }
        })
        .catch(err => {
          console.log('Cann not load configs', err)
        })
    },
  },
  created() {
    this.loadDefaultStatics()
  },
}
</script>

<style scoped></style>
