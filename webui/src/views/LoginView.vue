<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                username: null
            }
        },
        methods: {
            async login() {
                this.errormsg = null;
                try {
                    let response = await this.$axios.post("/session", {
                        username: this.username.trim()
                    });
                    let userinfo = response.data;
                    localStorage.setItem('username', this.username);
                    localStorage.setItem('userid', userinfo.userid);
                    localStorage.setItem('userphoto',userinfo.photo);
                    this.$router.push({
                        path: "/chats"
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            }
        }
    }
</script>

<template>
        <div class="blurred-box">
            <h1 style="position: absolute; top:10%;">WasaText</h1>
            <input class="user-password" v-model="username" placeholder="username" @keyup.enter="login">
            <button class="btn-hover" @click="login">Login</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
</template>

<style>
    body {
        background-color: rgb(33, 28, 28);
    }
    .blurred-box {
        display: flex;
        width: 400px;
        height: 400px;
        background-color: rgba(80, 59, 59, 0.5);
        border-radius: 20px;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        color: whitesmoke;
        position: relative;
        z-index: 1000;
    }
</style>