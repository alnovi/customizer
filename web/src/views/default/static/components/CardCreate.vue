<template>
  <div class="card-body">
    <form @submit="store">
      <div class="input-group my-2">
        <div class="input-group-prepend">
          <label class="input-group-text" for="typeGroup">Тип</label>
        </div>
        <select class="custom-select" id="typeGroup" v-model="type">
          <option v-for="(val, key) in types" :value="key" :selected="key === type">{{ val }}</option>
        </select>
      </div>
      <div class="form-group">
        <input type="text" class="form-control" placeholder="ИД элемента" v-model="key" />
      </div>

      <div class="form-group" v-if="['text', 'link'].includes(type)">
        <input type="text" class="form-control" placeholder="Значение (Текст)" v-model="value" />
      </div>

      <div class="form-group" v-if="['color'].includes(type)">
        <div class="input-group mb-2">
          <div class="input-group-prepend">
            <div class="input-group-text input-color" :style="{ 'background-color': value }"></div>
          </div>
          <input type="text" class="form-control" placeholder="Значение (Цвет)" v-model="value" />
        </div>
      </div>

      <div class="form-group" v-if="['image'].includes(type)">
        <div class="input-group mb-2">
          <div class="d-flex" style="align-items: center">
            <img style="width: 50px; height: 50px; object-fit: contain;" v-if="value" :src="value" class="mx-2" />
            <input type="file" class="form-control-file" @change="uploadFile($event)" />
          </div>
        </div>
      </div>

      <div class="form-group" v-if="['svg'].includes(type)">
        <div class="input-group mb-2">
          <div class="d-flex" style="align-items: center; width: 100%">
            <img style="width: 50px; height: 50px; object-fit: contain;" v-if="value" :src="svgToBase64(value)" alt="" />
            <textarea v-model="value" class="form-control" placeholder="Значене элемента"></textarea>
          </div>
        </div>
      </div>

      <div class="form-group" v-if="['script', 'style'].includes(type)">
        <div class="input-group mb-2">
          <div class="d-flex" style="align-items: center; width: 100%">
            <textarea v-model="value" class="form-control" placeholder="Значене элемента"></textarea>
          </div>
        </div>
      </div>

      <div class="form-group">
        <textarea v-model="description" class="form-control" placeholder="Описание элемента"></textarea>
      </div>

      <button v-if="hash !== oldHash" type="submit" class="btn btn-primary">Сохранить</button>
    </form>
  </div>
</template>

<script>
import Hash from 'object-hash'

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
  name: 'CardCreate',
  props: {
    collectionName: String,
  },
  data() {
    return {
      types: TYPES,
      hash: this.makeHash('', '', ''),
      oldHash: this.makeHash('', '', ''),
      key: '',
      type: 'text',
      value: '',
      description: '',
    }
  },
  methods: {
    makeHash(key, value, description) {
      return Hash.MD5({
        key: key,
        value: value,
        description: description,
      })
    },
    utf8ToB64(str) {
      return btoa(unescape(encodeURIComponent(str)))
    },
    b64ToUtf8(str) {
      return decodeURIComponent(escape(atob(str)))
    },
    fileAsBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.readAsDataURL(file)
        reader.onload = () => resolve(reader.result)
        reader.onerror = error => reject(error)
      })
    },
    svgToBase64(data) {
      return `data:image/svg+xml;base64, ${this.utf8ToB64(data)}`
    },
    uploadFile(event) {
      this.fileAsBase64(event.target.files[0])
        .then(res => {
          this.value = res
        })
        .catch(err => {
          console.log('can not convert img to base 64', err)
        })
    },
    reset() {
      this.hash = this.makeHash('', '', '')
      this.oldHash = this.makeHash('', '', '')
      this.key = ''
      this.type = 'text'
      this.value = ''
      this.description = ''
    },
    store(event) {
      event.preventDefault()

      let key = [this.type, this.key].join('.')

      let data = {
        [key]: {
          value: this.value,
          description: this.description,
        },
      }

      this.$api
        .storeDefaultConfig('statics', this.collectionName, data)
        .then(res => {
          let hash = this.makeHash(this.key, this.value, this.description)

          this.hash = hash
          this.oldHash = hash

          this.$emit('update', {
            updateType: 'create',
            key: key,
            value: this.value,
            description: this.description,
            type: this.type,
          })
          this.reset()
        })
        .catch(err => {
          console.log('Can not save', err)
        })
    },
  },
  watch: {
    value: function(newVal, oldVal) {
      this.hash = this.makeHash(this.key, newVal, this.description)
    },
    description: function(newVal, oldVal) {
      this.hash = this.makeHash(this.key, newVal, this.description)
    },
    key: function(newVal, oldVal) {
      this.hash = this.makeHash(this.key, newVal, this.description)
    },
  },
  created() {},
}
</script>

<style scoped lang="scss">
.card {
  margin: 10px 0;

  .input-color {
    width: 38px;
  }
}
</style>
