<template>
  <div class="bc">
    <div class="chain">
      <div class="search">
        <el-input class="search-input" :placeholder="searchTips" v-model="search">
          <el-select class="search-type" v-model="searchType" slot="prepend" placeholder="Select">
            <el-option label="Block Num" value="1"></el-option>
            <el-option label="Block Hash" value="2"></el-option>
            <el-option label="Tx Hash" value="3"></el-option>
          </el-select>
          <el-button slot="append" icon="el-icon-search" @click="fetchBlock"></el-button>
        </el-input>
      </div>
      <div class="info">
        <div class="info-item"> 
          <span>Network ID:</span>
          <span>{{chain.networkId}}</span>
        </div>
        <div class="info-item"> 
          <span>Block Num:</span>
          <span>{{chain.blockNum}}</span>
        </div>
        <div class="info-item"> 
          <span>Peer Count:</span>
          <span>{{chain.peerCount}}</span>
        </div>
      </div>
    </div>
    <div class="blocks">
      <el-table :data="blocks" border>
        <el-table-column label="Number">
          <template slot-scope="{row}">
            <span class="t-num" @click="jumpBlock(row.number)">{{row.number}}</span>
          </template>
        </el-table-column>
        <el-table-column label="Hash">
          <template slot-scope="{row}">
            <span class="t-hide">{{row.hash}}</span>
          </template>
        </el-table-column>
        <el-table-column label="Age">
          <template slot-scope="{row}">
            <span>{{ row.timestamp | parseAge }} seconds ago</span>
          </template>
        </el-table-column>
        <el-table-column prop="miner" label="Miner">
          <template slot-scope="{row}">
            <span class="t-miner">{{row.miner}}</span>
          </template>
        </el-table-column>
        <el-table-column label="Txs" width="80">
          <template slot-scope="{row}">
            <span>{{row.txs.length}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="uncles" label="Uncles" width="80" />
      </el-table>
    </div>
  </div>
</template>

<script>
import web3 from "web3"
export default {
  data() {
    return {
      chain: {},
      blocks: [],
      search: "",
      searchType: "1"
    }
  },
  filters: {
    parseBigInt(v) {
      return web3.utils.hexToNumber(v)
    },
    parseAge(ts) {
      const tsNow  = Math.floor(new Date() / 1000)
      const diff = tsNow - ts
      return diff
    }
  },
  computed: {
    searchTips() {
      switch (this.searchType) {
        case "1":
          return "Please input block number"
        case "2":
          return "Please input block hash"
        case "3":
          return "Please input transation hash"
        default:
          return ""
      }
    }
  },
  async created() {
    this.fetchChainInfo()
      .then(() => {
        return this.fetchLastestBlock()
      })
  },
  methods: {
    fetchChainInfo() {
      return this.$http.get("/chain/info")
        .then(res => {
          this.chain = res.data
        })
        .catch(() => {})
    },
    fetchLastestBlock() {
      const to = this.chain.blockNum + 1
      var from = to - 10
      if (from < 0) {
        from = 0
      }
      const postData = {
        from, to
      }
      this.$http.post("/chain/block/range", postData)
        .then(res => {
          this.blocks = res.data.reverse()
        })
        .catch( () => {})
    },
    fetchBlock() {
      switch (this.searchType) {
        case "1":
          this.$http.post("/chain/block/number", {blockNum: Number(this.search)})
            .then(res => {
              this.blocks = [res.data]
            })
            .catch(() => {})
          break
        case "2":
          this.$http.post("/chain/block/hash", {hash: this.search})
            .then(res => {
              this.blocks = [res.data]
            })
            .catch(() => {})
          break
        case "3":
          this.$http.post("/chain/block/txhash", {hash: this.search})
            .then(res => {
              this.blocks = [res.data]
            })
            .catch(() => {})
          break
        default:
          this.$message.error("haha")
      }
    },
    jumpBlock(num) {
      this.$router.push(`/block/${num}`)
    }
  },
}
</script>

<style lang="postcss" scoped>
.bc {
  width: 80vw;
  margin: auto;
  padding-bottom: 50px;
}

.chain {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.blocks {
  margin-top: 30px;
}

.search-input {
  width: 380px;
}

.search-type {
  width: 110px;
}

.info {
  display: flex;
}

.info-item {
  font-size: 12px;
  margin-right: 20px;
}

.info-item span:nth-child(2){
  font-size: 14px;
  font-weight: bold;
  margin-left: 5px;
}

.t-miner {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.t-num {
  color: var(--main-color);
  cursor: pointer;
}

.t-num:hover {
  text-decoration: underline;
}

.t-hide {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>