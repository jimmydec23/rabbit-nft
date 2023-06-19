<template>
  <div class="market page">
  <!-- my collectible -->
	<h2 class="collectible-title" v-if="this.collectible.length">Hot Collections</h2>
  <div v-else class="market-empty">
    <img src="@/assets/logo.png" alt="">
    <span>No collections on sale now.</span>
  </div>
	<div class="list">
		<div
			class="item"
			v-for="(item, i) in collectible"
			:key="i"
			@click="jump(item.tokenId)"
      @mouseover="overItem=item.tokenId"
		>
      <img :src="item.url" alt="" />
      <div class="item-content">
        <div class="item-content-part">
          <span>{{item.name}}</span>
          <span>{{ item.name}} #{{item.tokenId}}</span>
        </div>
        <div class="item-content-part item-content-part__right">
          <span>Price</span>
          <span>{{ item.price | toEth}} ETH</span>
        </div>
      </div>
      <span v-if="overItem===item.tokenId" class="item-buy">Buy now</span>
		</div>
	</div>
  <div class="market-pg">
    <el-pagination 
      v-if="total > 10"
      background layout="prev, pager, next" 
      @current-change="handleCurrentChange"
      :current-page.sync="currentPage"
      :total="total">
    </el-pagination>
  </div>
	</div>
</template>

<script>
import web3 from "web3"
export default {
	data() {
		return {
      total: 0,
      currentPage: 1,
			collectible: [],
      overItem: -1
		};
	},
  filters: {
    toEth(val) {
      return web3.utils.fromWei(val)
    }
  },
	async created() {
		this.fetchCollectible();
	},
	methods: {
		fetchCollectible() {
			this.$http
				.post("/market/list", { page: this.currentPage, limit: 9 })
				.then((res) => {
					this.collectible = res.data;
          this.total = res.total;
				})
				.catch(() => {});
		},
		jump(tokenId) {
			this.$router.push(`/collectible/${tokenId}`);
		},
    handleCurrentChange() {
      this.fetchCollectible()
    }
	},
};
</script>

<style lang="postcss" scoped>
.collectible-title {
  text-align: center;
}

.list {
	display: flex;
	flex-wrap: wrap;
  width: 1000px;
  margin: auto;
}

.item {
	width: 300px;
	height: 380px;
	display: flex;
	flex-direction: column;
	word-break: break-all;
	margin: 5px;
	box-sizing: border-box;
	box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
	cursor: pointer;
}

.item-content {
	width: 100%;
	display: flex;
  justify-content: space-between;
  padding: 0 10px;
  box-sizing: border-box;
}

.item-content-part {
	width: 100%;
	display: flex;
  flex-direction: column;
  font-size: 14px;
}

.item-content-part span:first-child{
  color: #707A83;
  margin-top: 5px;
}

.item-content-part span:nth-child(2) {
  font-weight: bold;
}

.item-content-part__right {
  text-align: right;
}

.item img {
	width: 300px;
	height: 300px;
  object-fit: contain;
}

.item-buy {
  padding: 0 10px;
  color: var(--main-color);
  font-size: 14px;
  font-weight: bold;
  margin-top: 5px;
}

.market-pg {
  margin: 50px 0 30px 0;
  display: flex;
  justify-content: center;
}

.market-empty {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin-top: 20vh;
}

.market-empty img {
  width: 200px;
  margin: 10px 0;
}

</style>

