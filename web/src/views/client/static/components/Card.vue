<template>
  <div class="card-body">
    <form @submit="store">
      <div class="input-group my-2">
        <div class="input-group-prepend">
          <label class="input-group-text" for="typeGroup">Тип</label>
        </div>
        <select class="custom-select" id="typeGroup" disabled>
          <option v-for="(val, key) in types" :value="key" :selected="key === type">{{ val }}</option>
        </select>
      </div>

      <div class="input-group my-2">
        <div class="input-group-prepend">
          <label class="input-group-text">ИД</label>
        </div>
        <input type="text" class="form-control" placeholder="ИД элемента" :value="key" disabled />
      </div>

      <div class="input-group mb-2" v-if="['text', 'link'].includes(type)">
        <div class="input-group-prepend">
          <label class="input-group-text">Значение</label>
        </div>
        <input type="text" class="form-control" placeholder="Значение (Текст)" v-model="value" />
      </div>

      <div class="form-group" v-if="['color'].includes(type)">
        <div class="input-group mb-2">
          <div class="input-group-prepend">
            <label class="input-group-text">Значение</label>
          </div>
          <div class="input-group-prepend">
            <div class="input-group-text input-color" :style="{ 'background-color': value }"></div>
          </div>
          <input type="text" class="form-control" placeholder="Значение (Цвет)" v-model="value" />
        </div>
      </div>

      <div class="form-group" v-if="['image'].includes(type)">
        <div class="input-group mb-2">
          <div class="input-group-prepend">
            <label class="input-group-text">Значение</label>
          </div>
          <div class="px-2 input-group-text" v-on:click="preview(value)">
            <img style="width: 50px; height: 50px; object-fit: contain;" :src="value" class="mx-2" />
          </div>
          <div class="d-flex align-items-center ml-2">
            <input type="file" class="form-control-file" @change="uploadFile($event)" />
          </div>
        </div>
      </div>

      <div class="form-group" v-if="['svg'].includes(type)">
        <div class="input-group mb-2">
          <div class="input-group-prepend">
            <label class="input-group-text">Значение</label>
          </div>
          <div class="px-2 input-group-text" v-on:click="preview(svgToBase64(value))">
            <img style="width: 50px; height: 50px; object-fit: contain;" :src="svgToBase64(value)" alt="" />
          </div>
          <textarea style="flex:1;" v-model="value" class="form-control" placeholder="Значене элемента"></textarea>
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

      <div v-on:click="handleDelete" class="delete-wrapper">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash-fill" viewBox="0 0 16 16">
          <path
            d="M2.5 1a1 1 0 0 0-1 1v1a1 1 0 0 0 1 1H3v9a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2V4h.5a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1H10a1 1 0 0 0-1-1H7a1 1 0 0 0-1 1H2.5zm3 4a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 .5-.5zM8 5a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-1 0v-7A.5.5 0 0 1 8 5zm3 .5v7a.5.5 0 0 1-1 0v-7a.5.5 0 0 1 1 0z"
          />
        </svg>
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
  name: 'Card',
  props: {
    clientId: String,
    collectionName: String,
    keyItem: String,
    item: Object,
  },
  data() {
    return {
      types: TYPES,
      hash: this.makeHash(this.$props.keyItem, this.$props.item.value, this.$props.item.description),
      oldHash: this.makeHash(this.$props.keyItem, this.$props.item.value, this.$props.item.description),
      key: this.$props.keyItem,
      type: this.$props.item.type,
      value: this.$props.item.value,
      description: this.$props.item.description,
    }
  },
  methods: {
    preview(value) {
      let w = window.open('about:blank')
      let image = new Image()
      image.src = value
      setTimeout(function() {
        w.document.write(image.outerHTML)
      }, 0)
    },
    handleDelete() {
      const isConfirmed = window.confirm(`Вы действительно хотите удалить элемент ${this.key}?`)
      if (isConfirmed) {
        this.$api
          .deleteStaticElement(this.clientId, this.collectionName, this.key)
          .then(result => {
            this.$emit('update', {
              updateType: 'delete',
              key: this.key,
              value: this.value,
              description: this.description,
            })
          })
          .catch(error => {
            console.log(`Could not delete item with key ${this.key}`)
          })
      }
    },
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
    store(event) {
      event.preventDefault()

      let data = {
        [this.key]: {
          value: this.value,
          description: this.description,
        },
      }

      this.$api
        .storeClientConfig(this.clientId, 'static', this.collectionName, '', data)
        .then(res => {
          let hash = this.makeHash(this.key, this.value, this.description)

          this.hash = hash
          this.oldHash = hash

          this.$emit('update', {
            updateType: 'update',
            key: this.key,
            value: this.value,
            type: this.type,
            description: this.description,
          })
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
      this.hash = this.makeHash(this.key, this.value, newVal)
    },
  },
  created() {},
}
</script>

<style scoped lang="scss">
.card {
  margin: 16px 0;
  .input-group-prepend {
    label {
      min-width: 100px;
    }
  }

  .input-color {
    width: 38px;
  }
  .delete-wrapper {
    cursor: pointer;
    position: absolute;
    top: 6px;
    right: 2px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}
</style>
