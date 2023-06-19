import Cookies from "js-cookie";

const TokenKey = "X-Token";

const authutil = {
  getToken() {
    return Cookies.get(TokenKey);
  },
  setToken(token) {
    return Cookies.set(TokenKey, token);
  },
  removeToken() {
    return Cookies.remove(TokenKey);
  }
};
export default authutil;
