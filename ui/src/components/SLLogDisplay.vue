<template>
<div>
  <pre class="log-container">
<template v-for="item in logs.slice().reverse()">
{{ item.Appname }}| {{ item.Msg}}</template>
  </pre>
</div>
</template>

<script>
export default {
  name: 'SLLogDisplay',
  computed: {
    logs: {
      get() {
        return this.$store.state.log;
      },
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .log-container {
    padding: 16px 16px;
    margin: 0px;
    font-size: 12pt;
  }
</style>
