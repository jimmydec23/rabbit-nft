<template>
	<div class="header">
		<span class="title">Rabbit NFT</span>
    <div class="nav">
      <span 
        v-for="(item, i) in tabs" 
        :key="i" 
        :class="active(item.url) ? 'active' : ''"
        @click="jump(item.url)"
      >
        {{item.name}}
      </span>
    </div>
    <el-dropdown class="header-dropdown"  trigger="click">
      <span class="el-dropdown-link">
        <div class="user">
          <img class="user-image" src="@/assets/logo.png" alt="">
          <span v-if="userInfo.basic" class="user-account">{{userInfo.basic.nickname}}</span>
        </div>
      </span>
      <el-dropdown-menu slot="dropdown">
        <template v-if="userInfo.basic">
          <el-dropdown-item class="drop-tab" @click.native="jump('/txlog')">Tx Log</el-dropdown-item>
          <el-dropdown-item class="drop-tab" @click.native="logout">Logout</el-dropdown-item>
        </template>
        <template v-else>
          <el-dropdown-item class="drop-tab" @click.native="jump('/login')">Login</el-dropdown-item>
        </template>
      </el-dropdown-menu>
    </el-dropdown>
	</div>
</template>

<script>
import cookieUtil from "@/utils/cookie"
export default {
  data() {
    return {
      userInfo: {},
      tabs: [
        {name: "Market", url: "/market"},
        {name: "Blockchain", url: "/blockchain"},
        {name: "Create", url: "/create"},
        {name: "Profile", url: "/profile"},
      ]
    }
  },
  watch: {
    "$store.state.userLoaded": function(){
      this.readUserInfo()
    }
  },
  created() {
    this.readUserInfo() 
  },
  methods: {
    readUserInfo(){
      const userInfo = this.$session.get("userInfo")
      if (userInfo) {
        this.userInfo = userInfo
      }else{
        this.userInfo = {}
      } 
    },
    active(url) {
      const match = this.$route.path === url
      return match
    },
    jump(url) {
      this.$router.push(url).catch(() => {})
    },
    logout(){
      cookieUtil.removeToken()
      this.$session.remove("userInfo")
      this.$store.commit("userChange")
      this.$router.push("/login")
    }
  },
}
</script>

<style lang="postcss" scoped>
.header {
	height: 50px;
	border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.title {
  font-weight: bold;
}

.nav {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-right: 5%;
}

.nav span {
  display: inline-block;
  margin: 0 20px;
  cursor: pointer;
  position: relative;
}

.nav span:hover::after {
  content: "";
  height: 3px;
  width: 100%;
  left: 0;
  bottom: -13px;
  background: var(--main-color);
  position: absolute;
}

.user-image {
  width: 25px;
  height: 25px;
}

.user {
  display: flex;
  align-items: flex-end;
  cursor: pointer;
}

.user-account {
  margin-left: 5px;
  font-weight: bold;
}

.active {
  color: var(--main-color);
}

.header-dropdown {
  margin-right: 20px;
}

.drop-tab{
  width: 60px;
  text-align: center;
}
</style>