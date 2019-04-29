import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    log: [],
    refresh: true,
    refreshInterval: 2000,
    refreshIntervalId: -1,
    autoscroll: true,
    numLines: 500,
    filter: {
      app: '',
      msg: '',
    },
  },
  mutations: {
    updateAutoscroll(state, val) {
      state.autoscroll = val;
    },
    updateRefresh(state, val) {
      state.refresh = val;
    },
    updateRefreshInterval(state, val) {
      state.refreshInterval = val;
    },
    updateRefreshIntervalId(state, val) {
      state.refreshIntervalId = val;
    },
    updateNumLines(state, val) {
      state.numLines = val;
    },
    updateMsg(state, val) {
      state.filter.msg = val;
    },
    updateApp(state, val) {
      state.filter.app = val;
    },
    updateLog(state, val) {
      state.log = val;
    },
  },
  actions: {
    async getLog(context) {
      console.log('updating log');

      const data = new FormData();
      data.set('num', context.state.numLines);
      data.set('msg', context.state.filter.msg);
      data.set('app', context.state.filter.app);

      //      const data = {
      //        num: context.state.numLines,
      //      };
      console.log(data);
      let url = '/receiver/logs';
      //      if (context.state.filter.app !== '' || context.state.filter.msg !== '') {
      //        url = '/receiver/filter';
      //      }
      console.log(url);
      const config = {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      };
      console.log(context.state.numLines);
      console.log(context.state.filter.msg);
      console.log(context.state.filter.app);
      const response = await axios.post(url, data, config);
      console.log(response.data);
      context.commit('updateLog', response.data.Result);
    },
    async runRefreshLoop(context) {
      if (context.state.refreshIntervalId !== -1) {
        clearInterval(context.state.refreshIntervalId);
      }

      const intervalId = setInterval(() => {
        context.dispatch('getLog');
      }, context.state.refreshInterval);
      context.commit('updateRefreshIntervalId', intervalId);
    },
    async stopRefreshLoop(context) {
      if (context.state.refreshIntervalId !== -1) {
        clearInterval(context.state.refreshIntervalId);
      }
    },
  },
});
