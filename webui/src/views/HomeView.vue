<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                username: null,
                userid: localStorage.getItem("userid"),
                username: localStorage.getItem("username"),
                searcheduser: null,
                users: []
            }
        },
        methods: {
            async searchuser(searcheduser) {
                this.errormsg = null;
                this.users=[];
                try {
                    let response = await this.$axios.get("/users", {params: {username: this.searcheduser }});
                    response.data.userlist.forEach(user => {
                        this.users.push(user);
                    });
                } catch (e) {
                    this.errormsg = e.toString();
                }
                this.loading = false;
            }
        }
    }
</script>


<template>
    <div class="navbar-dark">
        <h1>{{username}}</h1>
        <div class="searchbox">
            <input class="searchbox-user" v-model="searcheduser" placeholder="username" style="color: whitesmoke; background: #332a2a;" @keyup.enter="searchuser(sercheduser)">
            <div v-if="users.length>0" class="searched-dropdown">
                <ul>
                    <li v-for="user in users":key="user.username">
                        {{user.username}}
                    </li>
                </ul>
            </div>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
body {
        background-color: rgb(33, 28, 28);
    }
    .navbar-dark{
        background-color: rgb(38, 38, 44);
        position: relative;
        display: flex;
        align-items: center;
        height: 60px;
        width: 100%;
    }
    .navbar-dark h1 {
    position: absolute; /* Posiziona l'elemento in modo indipendente dagli altri */
    margin: 10pt;          /* Elimina margini predefiniti */
    color: whitesmoke;       /* Colore del testo */
    }  
    .searchbox{
        position: absolute; /* Posiziona l'elemento in modo indipendente dagli altri */
        top: 50%;           /* Centra verticalmente */
        left: 50%;          /* Centra orizzontalmente */
        transform: translate(-50%, -50%); /* Sposta il punto di riferimento al centro dell'elemento */
        margin: 0;          /* Elimina margini predefiniti */
        color: white;       /* Colore del testo */
        height: 27px;        /* deve essere uguale all'altezza del searchbox */
    }
    .searchbox-user{
        width: 100%;
        margin: 0;          /* Elimina margini predefiniti */
    }
    .searched-dropdown {
        background-color: #332a2a;
        width: 100%;
        max-height: 200px;
        overflow-y: auto;
    }
    .searched-dropdown ul {
    list-style: none;
    padding: 0;
    margin: 0;
    color: whitesmoke;
    }

    .searched-dropdown li {
    padding: 10px;
    cursor: pointer;
    }

    .searched-dropdown li:hover {
    background-color: #695d5d;
    }
</style>