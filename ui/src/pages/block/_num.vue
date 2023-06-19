<template>
  <div class="block">
    <el-tabs class="card" type="card" @tab-click="handleClick">
      <el-tab-pane label="Overview">
        <div class="o-item">
          <span>Hash:</span>
          <span>{{ block.hash }}</span>
        </div>
        <div class="o-item">
          <span>Miner:</span>
          <span>{{ block.miner }}</span>
        </div>
        <div class="o-item">
          <span>Difficulty:</span>
          <span>{{ block.difficulty }}</span>
        </div>
        <div class="o-item">
          <span>Gas Usage:</span>
          <span>{{ block.gasUsed }}</span>
        </div>
        <div class="o-item">
          <span>Gas Limit:</span>
          <span>{{ block.gasLimit }}</span>
        </div>
        <div class="o-item">
          <span>Nonce:</span>
          <span>{{ block.nonce }}</span>
        </div>
        <div class="o-item">
          <span>Time:</span>
          <span>{{ block.timestamp | timeFormat}}</span>
        </div>
      </el-tab-pane>
      <el-tab-pane>
        <span slot="label">
          Transations
          <el-badge class="badge" :value="txCount" :max="99" />
        </span>
        <el-table :data="block.txs" border>
          <el-table-column prop="hash" label="Hash">
            <template slot-scope="{row}">
              <span class="t-hide">{{row.hash}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="from" label="From">
            <template slot-scope="{row}">
              <span class="t-hide">{{row.from}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="to" label="To">
            <template slot-scope="{row}">
              <span class="t-hide">{{row.to}}</span>
            </template>
          </el-table-column>
          <el-table-column label="Cost">
            <template slot-scope="{row}">
              <span>{{row.cost | toEth}} eth</span>
            </template>
          </el-table-column>
          <el-table-column prop="gasUsage" label="Gas Usage" />
        </el-table>
      </el-tab-pane>
  </el-tabs>
  </div>
</template>

<script>
import web3 from "web3"
export default {
  data() {
    return {
      blockNum: 0,
      block: {}
    }
  },
  filters: {
    timeFormat(ts) {
      if (!ts) {
        return ""
      }
      const d = new Date(ts * 1000)
      return d.toISOString().
        replace(/T/, ' ').      
        replace(/\..+/, '')
    },
    toEth(val) {
      return web3.utils.fromWei(val)
    }
  },
  computed: {
    txCount() {
      if (this.block.txs) {
        return this.block.txs.length
      }else{
        return 0
      }
    }
  },
  async created() {
    this.blockNum = Number(this.$route.params.num)
    this.fetchBlock()
  },
  methods: {
    fetchBlock() {
      this.$http.post("/chain/block/number", {blockNum: this.blockNum})
        .then(res => {
          this.block = res.data
        })
        .catch(() => {})
    },
    handleClick() {}
  },
}
</script>

<style lang="postcss" scoped>
.card {
  margin: auto;
  width: 80vw;
  margin-top: 10vh;
}

.o-item {
  height: 50px;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-sizing: border-box;
}

.o-item span:first-child{
  width: 120px;
  display: inline-block;
}

.t-hide {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge {
  position: absolute;
  top: -6px;
  right: 0;
}
</style>