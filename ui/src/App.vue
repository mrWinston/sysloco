<template>

  <div class="content">
    <!--
    <div class="cell top-left">
      <img src="./assets/logo-wide-opt.svg" alt="" class="logo">
    </div>
    <div class="cell top-right"></div>
    -->
    <sl-header-bar class="top"></sl-header-bar>
    <sl-filter-section class="cell bottom-left"></sl-filter-section>

    <div id="logcontainer" class="cell bottom-right">
      <sl-log-display></sl-log-display>
    </div>
  </div>


</template>

<script>
import SLLogDisplay from './components/SLLogDisplay.vue';
import SLHeaderBar from './components/SLHeaderBar.vue';
import SLFilterSection from './components/SLFilterSection.vue';

export default {
  name: 'App',
  components: {
    'sl-log-display': SLLogDisplay,
    'sl-header-bar': SLHeaderBar,
    'sl-filter-section': SLFilterSection,
  },
  created() {
    this.$store.dispatch('getLog');
    if (this.$store.state.refresh) {
      this.$store.dispatch('runRefreshLoop');
    }
  },
  computed: {
    logs: {
      get() {
        return this.$store.state.log;
      },
    },
  },
  watch: {
    logs: function () {
      if (this.$store.state.autoscroll) {
        var container = this.$el.querySelector("#logcontainer");
        container.scrollTop = container.scrollHeight;
      }
    },
  },
};
</script>
  <style>

  @media (min-width: 1600px) {
    .content {
      margin: 0px 50px 0px 50px;
    }
  }

  html,body {
    font-family: 'Overpass Mono', monospace;
    font-size: 14pt;
    margin: 0px;
    position: relative;
    height: 100vh;
    width: 100vw;
  }

  .content {
    display: grid;
    grid-template-rows: 64px 1fr;
    grid-template-columns: 25% 75%;
    height: 100%;
    grid-gap: 16px;
  }
  .top {
    margin-top: 10px;
    grid-column: 1 / 3;
    grid-row: 1 / 2;
  }
  .top-left {
    margin-top: 10px;
    grid-column: 1 / 2;
    grid-row: 1 / 2;
  }
  .top-right {
    margin-top: 10px;
    grid-column: 2 / 3;
    grid-row: 1 / 2;
  }
  .bottom-left {
    margin-bottom: 10px;
    grid-column: 1 / 2;
    grid-row: 2 / 3;
  }
  .bottom-right {
    margin-top: 16px;
    margin-bottom: 10px;
    grid-column: 2 / 3;
    grid-row: 2 / 3;
    background-color: #ECECEC;
    overflow: auto;
  }

</style>
