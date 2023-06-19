import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    userLoaded: false
  },
  mutations: {
    userChange(state) {
      state.userLoaded = !state.userLoaded
    }
  },
  actions: {

  }
})
