<template>
  <b-container>
    <div>
      <b-navbar toggleable="lg" type="dark" variant="dark">
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <router-link :to="{ name: 'default.static' }" v-slot="{ href, isActive }">
              <b-nav-item :active="isActive" :href="href">Статики</b-nav-item>
            </router-link>
            <router-link :to="{ name: 'clients' }" v-slot="{ href, isActive }">
              <b-nav-item :active="isActive" :href="href">Клиенты</b-nav-item>
            </router-link>
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div>

    <router-view></router-view>
  </b-container>
</template>

<script>
export default {
  name: 'LayoutMain',
}
</script>

<style scoped lang="scss"></style>
