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
    <div class="container">
        <div class="box">
            <h1>WasaText</h1>
            <input class="user-password" v-model="username" placeholder="username" @keyup.enter="login">
            <button class="btn-hover" @click="login">Login</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>
</template>

<style>
    body {
        background-color: rgb(33, 28, 28);
    }

    .container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
    }

    .box {
        display: flex;
        width: 400px;
        height: 400px;
        background-color: blue;
        justify-content: center;
        align-items: center;
        flex-direction: column;
    }
</style>