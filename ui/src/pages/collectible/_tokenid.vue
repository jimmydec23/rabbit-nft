<template>
  <div class="col page">
    <div class="info">
      <img class="info-img" :src="collectible.url" alt="" />
    </div>
    <div class="market">
      <h3>Basic Info</h3>
      <div class="market-desc">
        <span>Name:</span>
        <span>{{collectible.name}}</span>
      </div>
      <div class="market-desc">
        <span>ID:</span>
        <span>#{{collectible.tokenId}}</span>
      </div>
      <div class="market-desc">
        <span>Owner:</span>
        <span v-if="isMine"> Me </span>
        <span v-else>
          {{ collectible.owner }}
        </span>
      </div>
      <div class="market-desc">
        <span>Hash:</span>
        <span>{{collectible.hash}}</span>
      </div>
      <div class="market-desc">
        <span>Description:</span>
        <el-input v-model="collectible.description" readonly :rows="3" type="textarea"></el-input>
      </div>

      <h3>Market Info</h3>
      <div class="market-desc">
        <span>On Sale:</span>
        <span>{{collectible.onsale ? 'Yes' : 'No'}}</span>
        <i v-if="collectible.onsale" class="el-icon-success onsale"></i>
      </div>
      <div class="market-desc" v-if="collectible.onsale">
        <span>Price:</span>
        <span class="market-price">{{collectible.price | toEth}}</span>
        <span>eth</span>
      </div>
      <div v-if="!isMine && collectible.onsale" class="purchase">
        <h3>Purchase this collection</h3>
        <el-button class="btn-purchase" type="primary" @click="purchaseCollectible">Purchase</el-button>
      </div>
      <div v-if="isMine && !collectible.onsale" class="sell">
        <h3>Sell this collection</h3>
        <el-form class="sell-form" :model="sellForm" label-width="50px" label-position="left">
          <el-form-item label="Price">
            <el-input-number v-model="sellForm.price" :precision="2" :min="0.01" :step="0.01">
              <template slot="append">eth</template>
            </el-input-number>
          </el-form-item>
          <el-form-item>
            <el-button 
              class="btn-sell"
              type="primary" 
              @click="sellCollectible" 
              :disabled="!sellForm.price" >
              Sell
            </el-button>
          </el-form-item>
        </el-form>
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
      tokenId: "",
      collectible: {},
      sellForm: {
      price: "",
      },
    };
  },
  filters: {
    toEth(val) {
      return web3.utils.fromWei(val)
    }
  },
  computed: {
    isMine() {
      if (!this.userInfo.basic) {
        return false;
      }
      return this.userInfo.basic.address === this.collectible.owner;
    },
  },
  async created() {
    const tokenId = this.$route.params.tokenid;
    this.tokenId = Number(tokenId);
    Promise.all([this.fetchUserInfo(), this.fetchCollectible(this.tokenId)]);
  },
  methods: {
    fetchUserInfo() {
      return this.$http
        .get("/user/info")
        .then((res) => {
        this.userInfo = res.data;
        })
        .catch(() => {});
    },
    fetchCollectible(tokenId) {
      this.$http
        .post("/collectible/get", { tokenId: tokenId })
        .then((res) => {
        this.collectible = res.data;
        })
        .catch(() => {});
    },
    sellCollectible() {
      const priceWei = web3.utils.toWei(String(this.sellForm.price))
      const postData = {
        tokenId: this.tokenId,
        price: String(priceWei),
      };
      this.$http
        .post("/collectible/sell", postData)
        .then(() => {
        return this.fetchCollectible(this.tokenId);
        })
        .catch(() => {});
    },
    purchaseCollectible() {
      const postData = {
        tokenId: this.tokenId,
      };
      this.$http
        .post("/collectible/purchase", postData)
        .then(() => {
        this.fetchCollectible(this.tokenId);
        })
        .catch(() => {});
    },
  },
};
</script>

<style lang="postcss" scoped>
.col {
  display: flex;
  width: 1000px;
  margin: auto;
  margin-top: 10vh;
  align-items: center;
}

.info {
  flex: 2;
}

.market{
  flex: 3;
  box-sizing: border-box;
  padding: 0 20px;
}

.market-desc {
  padding: 0 20px;
  box-sizing: border-box;
  font-size: 14px;
}

.market-desc span:first-child {
  width: 100px; 
  display: inline-block;
}

.info-img {
  width: 400px;
  height: 400px;
  box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
  object-fit: contain;
}

.sell-form {
  width: 300px;
  margin-left: 20px;
}

.btn-sell {
  width: 130px;
}

.btn-purchase {
  margin-left: 20px;
}

.market-title {
  font-size: 30px;
  margin: 0 0;
}

.market-price {
  color: red;
  font-weight: bold;
  font-size: 18px;
  margin-right: 5px;
}

.onsale {
  color: #67C23A;
  margin-left: 5px;
}

</style>