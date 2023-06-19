<template>
  <div class="log">
    <h3>Transation Logs</h3>
    <el-table :data="log" border>
      <el-table-column label="Tx">
        <template slot-scope="{row}">
          <div style="display:flex; align-items:center">
            <i class="t-copy el-icon-document-copy" @click="copyTx(row.tx)"/>
            <span class="t-hide">{{row.tx}}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Type">
        <template slot-scope="{row}">
          <span>{{ row.type | typeName}}</span>
        </template>
      </el-table-column>
      <el-table-column label="From">
        <template slot-scope="{row}">
          <span class="t-hide">{{row.from}}</span>
        </template>
      </el-table-column>
      <el-table-column label="To">
        <template slot-scope="{row}">
          <span class="t-hide">{{row.to}}</span>
        </template>
      </el-table-column>
      <el-table-column prop="created" label="Created">
        <template slot-scope="{row}">
          <span>{{ row.created | timeFormat }}</span>
        </template>
      </el-table-column>
    </el-table>
    <div class="log-pagi">
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
export default {
  data() {
    return {
      log: [],
      total: 0,
      currentPage: 1
    }
  },
  filters: {
    typeName(t) {
      switch(t){
        case 0:
          return "Mint"
        case 1:
          return "Transfer"
        case 2:
          return "Sell"
        case 3:
          return "Purchase"
        default:
          return "Undefine"
      }
    },
    timeFormat(ts) {
      const d = new Date(ts * 1000)
      return d.toISOString().
        replace(/T/, ' ').      
        replace(/\..+/, '')
    }
  },
  async created() {
    this.fetchLog()
  },
  methods: {
    fetchLog(){
      this.$http.post("/txlog/list", {page:this.currentPage, limit: 10})
        .then(res => {
          this.log = res.data
          this.total = res.total
        })   
        .catch(() => {})
    },
    handleCurrentChange() {
      return this.fetchLog()
    },
    copyTx(hash) {
      navigator.clipboard.writeText(hash)
      this.$message.info("Tx hash copied.")
    }
  },
}
</script>

<style lang="postcss" scoped>
.log {
  width: 1000px;
  margin: auto;
}

.log-pagi {
  margin: 20px 0;
  display: flex;
  justify-content: center;
}

.t-hide {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  position: relative;
}

.t-copy {
  margin-right: 5px;
  cursor: pointer;
}
</style>