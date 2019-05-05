<template>
  <div>
    <div class="form-element">
      <div class="label">
        Filter by App
      </div>
      <input type="text" :value="filter.app" @change="updateApp">
    </div>
    <div class="form-element">
      <div class="label">
        Filter by Message
      </div>
      <input type="text" :value="filter.msg" @change="updateMsg">
    </div>
    <div class="form-element">
      <input type="checkbox" v-model="autoscroll" id="autoscroll-check">
      <label for="autoscroll-check" class="label">
        Auto-Scroll
      </label>
    </div>
    <div class="form-element">
      <input type="checkbox" v-model="refresh" @change="triggerRefresh" id="refresh-check">
      <label for="refresh-check" class="label">
        Auto-Refresh
      </label>
    </div>
    <div class="form-element">
      <label for="input-refresh-interval">Refresh Interval (in ms)</label>
      <br/>
      <input type="number" id="input-refresh-interval" :disabled="!refresh" v-model="refreshInterval">
    </div>

    <div class="form-element">
      <label for="input-lines">Number of Lines</label>
      <br/>
      <input v-model.lazy="numLines" type="number" value="500">
    </div>
  </div>
</template>
<script>

export default {
  name: 'SLFilterSection',
  methods: {
    updateMsg(event) {
      this.$store.commit('updateMsg', event.target.value);
    },
    updateApp(event) {
      this.$store.commit('updateApp', event.target.value);
    },
    triggerRefresh() {
      if (this.refresh) {
        this.$store.dispatch('runRefreshLoop');
      } else {
        this.$store.dispatch('stopRefreshLoop');
      }
    },
  },
  computed: {
    autoscroll: {
      get() {
        return this.$store.state.autoscroll;
      },
      set(value) {
        this.$store.commit('updateAutoscroll', value);
      },
    },
    refresh: {
      get() {
        return this.$store.state.refresh;
      },
      set(value) {
        this.$store.commit('updateRefresh', value);
      },
    },
    refreshInterval: {
      get() {
        return this.$store.state.refreshInterval;
      },
      set(value) {
        this.$store.commit('updateRefreshInterval', value);
      },
    },
    numLines: {
      get() {
        return this.$store.state.numLines;
      },
      set(value) {
        this.$store.commit('updateNumLines', value);
      },
    },
    filter: {
      get() {
        return this.$store.state.filter;
      },
    },
  },
};

</script>
<style scoped>
  .form-element {
    padding: 16px 0px;
  }

  input[type="text"] {
    font-family: 'Overpass Mono', monospace;
    font-size: 14pt;
    padding: 5px;
    margin: 15px 0px;
  }
  input[type="number"] {
    font-family: 'Overpass Mono', monospace;
    font-size: 14pt;
    padding: 5px;
    margin: 15px 0px;
  }
  input[type="checkbox"] {
    display: none;
  }
  input[type="checkbox"] + label::before {
    width: 16px;
    height: 16px;
    background-color: #FFF;
    display: block;
    border-style: solid;
    border-color: #000000;
    border-width: 1px;
    content: "";
    float: left;
    margin-right: 5px;
  }
  input[type="checkbox"]:checked+label::before {
    box-shadow: inset 0px 0px 0px 3px #fff;
    background-color: #000;
  }
</style>
