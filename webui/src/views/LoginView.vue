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
                    this.errormsg = e.response.data;
                }
            }
        }
    }
</script>

<template>
        <div class="blurred-box">
            <h1 style="position: absolute; top:10%;">WasaText</h1>
            <div class="input-field">
                <input class="username-input" v-model="username" placeholder="Username" @keyup.enter="login">
            </div>
            <ErrorMsg v-if="errormsg" :msg="errormsg" style="position: absolute; top: 60%;"></ErrorMsg>
            <button class="login-button" @click="login">Login</button>
        </div>
</template>

<style>
    body {
        background-color: rgb(33, 28, 28);
    }
    .blurred-box {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 400px;
        height: 400px;
        background-color: rgba(80, 59, 59, 0.5);
        border-radius: 25px;
        z-index: 1000;
    }
    .input-field{
        position: absolute;
        top: 45%;
        border-radius: 20px;
        width: 70%;
        padding: 5px;
        background-color: #171717;
    }
    .username-input {
        background: none;
        border: none;
        outline: none;
        width: 100%;
        color: whitesmoke;
    }
    .login-button {
        position: absolute;
        top: 85%;
        padding: 5px;
        padding-left: 20px;
        padding-right: 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #252525;
        color: white;
    }

    .login-button:hover {
        background-color: black;
    }
</style>