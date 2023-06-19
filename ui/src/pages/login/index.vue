<template>
  <div class="login">
    <div class="left">
      <h3 class="left-title">Rabbit NFT</h3>
      <p>Fun, easy use, user friendly nft platform.</p>
    </div>
    <div class="right">
      <el-form class="form" ref="form" :model="form" :rules="formRules"  label-width="80px">
        <div class="switch">
          <el-button class="switch-btn" type="text" @click="swichMode('Login')" :class="loginMode?'btn-active':''">Login</el-button>
          <el-button class="switch-btn" type="text" @click="swichMode('Register')" :class="!loginMode?'btn-active':''">Register</el-button>
        </div>
        <el-form-item label="Account:" prop="account">
          <el-input v-model="form.account"></el-input>
        </el-form-item>
        <el-form-item label="Password:" prop="password">
          <el-input v-model="form.password" show-password></el-input>
        </el-form-item>
        <el-form-item v-if="!loginMode" label="Nickname:" prop="nickname">
          <el-input v-model="form.nickname"></el-input>
        </el-form-item>
        <el-form-item class="btn-container">
          <el-button v-if="loginMode" @click="onLogin" :loading="loading">Login</el-button>
          <el-button v-else @click="onRegister" :loading="rLoading">Register</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import cookieUtil from '@/utils/cookie'
export default {
  data() {
    return {
      action: "Login",
      form: {
        account: "",
        nickname: "",
        password: ""
      },
      formRules: {
        account: [{required: true}],
        password: [{required: true}],
        nickname: [{required: true}]
      },
      loading: false,
      rLoading: false
    }
  },
  computed: {
    loginMode() {
      return this.action === 'Login'
    }
  },
  methods: {
    swichMode(action) {
      this.action = action
    },
    requestLogin() {
      const postData = {
        account: this.form.account,
        password: this.form.password
      }
      this.loading = true
      this.$http.post("/user/login", postData)
        .then( res => {
          this.$message({
            message: "login success",
            type: "success"
          })
          const token = res.data
          cookieUtil.setToken(token)
        })
        .then( () => {
          this.$router.push("/profile")
        })
        .catch(() => {})
        .finally( () => {
          this.loading = false
        })
    },
    requestRegister() {
      const postData = {
        account: this.form.account,
        nickname: this.form.nickname,
        password: this.form.password
      }
      this.rLoading = true
      this.$http.post("/user/register", postData)
        .then( () => {
          this.$message({
            message: "register success",
            type: "success"
          })
        })
        .catch(() => {})
        .finally( () => {
          this.rLoading = false
        })
    },
    

    onLogin() {
      this.$refs["form"].validate(valide => {
        if (!valide) {
          this.$message({
            message: "Form input not valid",
            type: "warning"
          })
        }else{
          this.requestLogin()
        }
      })
    },
    onRegister() {
      this.$refs["form"].validate(valide => {
        if (!valide) {
          this.$message({
            message: "Form input not valid",
            type: "warning"
          })
        }else{
          this.requestRegister()
        }
      })
    }
  }
}
</script>

<style scoped>
.login {
  position: absolute;
  display: flex;
  width: 100%;
  left: 0;
  top: 50px;
  bottom: 0;
}

.left { 
  background-color: rgb(250, 161, 179);
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #fff;
}

.left-title {
  font-size: 50px;
  margin: 0;
}

.right {
  flex: 1;
}

.form {
  flex: 1;
  margin: auto;
  margin-top: 20vh;
  width: 500px;
  padding: 30px 100px;
  box-sizing: border-box;
}

.form >>> .el-input__inner {
  border-top: none;
  border-left: none;
  border-right: none;
  border-radius: 0;
}

.switch {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 30px;
}

.switch-btn {
  color: #606266;
}

.btn-container {
  display: flex;
  justify-content: flex-end;
}

.btn-active {
  color: var(--main-color);
  border-bottom: 1px solid var(--main-color);
}

</style>