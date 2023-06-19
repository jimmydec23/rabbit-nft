<template>
  <div class="profile page">
    <div v-if="userInfo.basic" class="basic">
      <div class="basic-left">
        <img src="@/assets/logo.png" alt="">
      </div>
      <div class="basic-right">
        <div>
          <span class="basic-name">{{userInfo.basic.nickname}}</span>
          <span class="basic-address">{{userInfo.basic.address}}</span>
        </div>
        <div>
          <span>Balance:</span>
          <span  class="basic-balance">{{userInfo.basic.balance | fromWei }}</span>
          <span>eth</span>
          <el-button 
            class="fund-btn"
            type="primary" 
            @click="onFundMe" 
            :loading="fundLoading" 
            size="mini">Fund Me</el-button>
        </div>
        <div>
          <span>Collection Amount: </span>
          <span class="basic-count">{{collectible.length}}</span>
        </div>
      </div>
    </div>
    <!-- my collectible -->
    <h3 v-if="collectible.length===0" class="collectible-title">
      You don't have any collection now.
    </h3>
    <h3 v-else class="collectible-title">My Collections</h3>
    <div class="list">
      <div class="item" v-for="(item, i) in collectible" :key="i" @click="jump(item.tokenId)">
        <img :src="item.url" alt="">
        <div class="item-name">
          <span>{{item.name}}</span>
          <span>#{{item.tokenId}}</span>
        </div>
        <span v-if="item.onsale" class="item-onsale">Onsale</span>
      </div>
    </div>
  </div>
</template>

<script>
import web3 from "web3"
export default {
  data() {
    return {
      userInfo: {},
      fundLoading: false,
      collectible: []
    }
  },
  filters: {
    fromWei(val) {
      return web3.utils.fromWei(web3.utils.toBN(val))
    }
  },
  async created() {
    Promise.all([
      this.fetchUserInfo(),
      this.fetchCollectibles(),
    ])
  },
  methods: {
    fetchUserInfo() {
      return this.$http
        .get("/user/info")
        .then(res => {
          this.userInfo = res.data
          this.$session.set("userInfo", res.data)
          this.$store.commit("userChange")
        })
        .catch(() => {});
    },
    fetchCollectibles() {
      this.$http.post("/collectible/list")
        .then(res => {
          this.collectible = res.data
        })
        .catch(() => {})
    },
    onFundMe() {
      this.fundLoading = true
      this.$http.post("/user/faucet")
        .then(() => {
          this.$message({
            message: "Congratuation! You have got some fund!",
            type: "success"
          })
          return this.fetchUserInfo()
        })
        .catch(() => {})
        .finally(() => {
          this.fundLoading = false
        })
    },
    jump(tokenId) {
      this.$router.push(`/collectible/${tokenId}`)
    }
  },
}
</script>

<style lang="postcss" scoped>
.profile {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.basic {
  display: flex;
  align-items: center;
  font-size: 14px;
  margin-top: 5vh;
}

.basic-left, .basic-right {
  padding: 0 20px;
  box-sizing: border-box;
}

.basic-left img {
  width: 80px;
  height: 80px;
}

.basic-right {
  display: flex;
  flex-direction: column;
}

.basic-name {
  font-weight: bold;
}

.basic-address {
  font-size: 12px;
  margin-left: 20px;
}

.basic-balance {
  font-weight: bold;
  margin-right: 5px;
  margin-left: 20px;
  font-size: 18px;
  color: #67C23A;
}

.basic-count {
  margin-left: 10px;
}

.fund-btn {
  margin-left: 20px;
  font-size: 10px;
  height: 20px;
  padding: 0 8px;
}

.fund, .nft {
  margin-top: 20px;
}

.collectible-title {
  margin-top: 50px;
}

.list{
  display: flex;
  flex-wrap: wrap;
  width: 1000px;
}

.item {
  width: 300px;
  height: 350px;
  display: flex;
  flex-direction: column;
  word-break: break-all;
  margin: 5px;
  box-sizing: border-box;
  box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
  cursor: pointer;
  position: relative;
}

.item img {
  width: 300px;
  height: 300px;
  object-fit: contain;
}

.item-onsale{
  background-color: #E6A23C;
  color: #fff;
  font-size: 12px;
  text-align: center;
  position: absolute;
  right: 0;
  top: 2px;
  padding: 2px 5px 2px 10px;
  z-index: 10;
  border-bottom-left-radius: 20px;
}

.item-name {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  text-align: center;
}
</style>