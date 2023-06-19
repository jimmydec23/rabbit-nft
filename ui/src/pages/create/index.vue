<template>
  <div class="create">
    <div class="content">
      <h2>Create New Item</h2>
      <p>Image, Video, Audio, or 3D Model</p>
      <span class="content-ext">File types supported: JPG, PNG, GIF, SVG, WEBM. Max size: 30 MB</span>
      <div class="image-box"></div>
      <el-upload 
        class="avatar-uploader"
        :headers="headers"
        :action="uploadurl"
        :limit="1"
        :file-list="fileList"
        :on-change="onFileChange"
        :on-exceed="onExceed"
        :on-success="onSuccess"
        :before-upload="handleBeforeUpload"
        :accept="acceptFile"
      >
        <img v-if="imageUrl" :src="imageUrl" class="avatar">
        <i v-else class="el-icon-plus avatar-uploader-icon"></i>
      </el-upload>
      <p>Name</p>
      <el-input v-model="form.name" placeholder="Item name"></el-input>
      <p>Description</p>
      <el-input type="textarea" v-model="form.desc" placeholder="Provide a detailed description of your item."></el-input>
      <p>Supply</p>
      <el-input v-model="form.supply" readonly></el-input>
      <el-button class="create-btn" type="primary" @click="onMint" :disabled="hash===''">Create</el-button>
    </div>
  </div>
</template>

<script>
import cookieUtil from "@/utils/cookie.js"
export default {
  data() {
    return {
      uploadurl: process.env.VUE_APP_BASE_API + '/file/upload',
      fileList: [],
      hash: "",
      imageUrl: "",
      fundLoading: false,
      form: {
        name: "",
        desc: "",
        supply: 1
      },
      accept: [".jpg", ".jpeg", ".png", ".gif", ".svg", ".webm", ".ico"]
    }
  },
  computed: {
    headers() {
      const h = {"X-Token": cookieUtil.getToken()}
      return h
    },
    acceptFile() {
      return this.accept.join(",")
    }
  },
  methods: {
    onMint() {
      const postData = {
        hash: this.hash,
        name: this.form.name,
        description: this.form.desc
      }
      this.$http.post("/collectible/mint", postData)
        .then(res => {
          this.$message({
            message: `Mint success, tx: ${res.data}`,
            type: "success"
          })
        })
        .catch(() => {})
    },
    onFileChange(file, fileList) {
      this.fileList = fileList
    },
    onExceed() {
      this.$message.warning("Only one file allowed")
    },
    onSuccess(response){
      if (response.code !== 0) {
        this.$alert(response.msg, "Result", {type: 'error'})
        this.fileList = []
        return
      }
      this.hash = response.data.hash
      this.imageUrl = response.data.url
    },
    handleBeforeUpload(file) {
      const exceedSizeLimit = file.size / 1024 / 1024 < 30
      if (!exceedSizeLimit) {
        this.$alert("File size exceed 30m limit.", "Warn", {type: 'warning'})
        return false
      }
    }
  },
}
</script>

<style lang="postcss" scoped>
.content{
  width: 600px;
  margin: auto auto 30px auto;
}

.content-ext {
  font-size: 12px;
}

.avatar-uploader >>> .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader >>> .el-upload:hover {
  border-color: #409EFF;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
.create-btn{
  margin-top: 20px;
}
</style>