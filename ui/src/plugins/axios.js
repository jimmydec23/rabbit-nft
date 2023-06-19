import Vue from 'vue'
import axios from 'axios'
import cookieUtil from '@/utils/cookie'
import { MessageBox, Message } from "element-ui";
import router from "@/router";

var showing = false

const instance = axios.create( {
    baseURL: process.env.VUE_APP_BASE_API
})

// set token
instance.interceptors.request.use(config => {
  config.headers["X-Token"] = cookieUtil.getToken()
  return config
})

// handle error
instance.interceptors.response.use(resp => {
  const r = resp.data
  const cfg = resp.config
  if (r.code === 0) {
    return r
  }
  if (r.code == 1) {
    if (!cfg.hideErrorTips) {
      Message({
        message: `Server response error: ${r.msg}`,
        type: "error"
      })
    }
    return Promise.reject(new Error(r.msg || "error"))
  }
  if (r.code == 2) {
    if (!showing) {
      showing = true
      MessageBox.alert(
        "Account offline, please login again",
        "Account Offline",
        {type: "warning"}
      ).finally(() => {
        router.push({ path: "/login" }).catch(() => {});
        showing = false
      })
    }
    return Promise.reject(new Error(r.msg || "error"))
  }
}, error => {
  Message({
    message: `Server response error: ${error}`,
    type: "error"
  })
  return error
})

// set vue plugin
const plugin = {};
plugin.install = function(Vue) {
  Vue.prototype.$http = instance;
};

Vue.use(plugin);

export default instance;