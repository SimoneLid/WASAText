<script>
    export default {
        data: function() {
            return {
                errormsg: null,
                userid: localStorage.getItem("userid"),
                username: localStorage.getItem("username"),
                userphoto: localStorage.getItem("userphoto"),
                
                // for searching users
                searcheduser: null,
                users: [],

                // for changing username and photo
                changeusernameshown: false,
                changedinfo: false,
                newusername : null,
                newuserphoto: null,

                // for chats preview
                chats: []
            }
        },
        methods: {
            handleClickOutside(event) {
                // Check if the click is outside the search box to search users
                if (this.$refs.boxsearchuser && !this.$refs.boxsearchuser.contains(event.target)) {
                    this.users = [];
                    this.searcheduser = null;
                }
            },
            async searchUser(searcheduser) {
                this.errormsg = null;
                this.users=[];
                try {
                    let response = await this.$axios.get("/users", {params: {username: this.searcheduser }});
                    response.data.userlist.forEach(user => {
                        if (user.username != this.username){
                            this.users.push(user);
                        }
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;;
                }
            },
            changePhotoFileSelect(){
                const file = this.$refs.changeimageInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.newuserphoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            changePhotoButton(){
                this.$refs.changeimageInput.click();
            },
            async resetChangeUsernamePrompt(){
                this.newusername = null;
                this.newuserphoto = null;
                this.changeusernameshown = false;
            },
            async changeUsernamePhoto(){
                this.errormsg = null;
                if (this.newusername){
                    try {
                        let response = await this.$axios.put("/users/"+this.userid+"/name", {username: this.newusername.trim()},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        localStorage.setItem('username', this.newusername);
                        this.username = this.newusername;
                        this.changedinfo = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }
                if (this.newuserphoto){
                    try {
                        let response = await this.$axios.put("/users/"+this.userid+"/photo",{photo: this.newuserphoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        localStorage.setItem('userphoto', this.newuserphoto);
                        this.userphoto = this.newuserphoto;
                        this.changedinfo = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }
                this.newusername = null;
                this.newuserphoto = null;
                if (this.changedinfo){
                    this.changeusernameshown = false;
                    this.errormsg = null;
                }
            },
            async buildChatPreview(){
                this.errormsg = null;
                this.chats=[];
                try {
                    let response = await this.$axios.get("/chats",{headers:{"Authorization": `Bearer ${this.userid}`}});
                    response.data.forEach(chat => {
                        if (chat.lastmessage.text.length>20){
                            chat.lastmessage.text = chat.lastmessage.text.slice(0,20)+"...";
                        }
                        if (chat.lastmessage.photo.length>0 && chat.lastmessage.text.length===0){
                            console.log("Testo 0");
                            chat.lastmessage.text="Photo";
                        }
                        if (chat.lastmessage.username==this.username){
                            chat.lastmessage.username="You";
                        }
                        chat.lastmessage.timestamp = chat.lastmessage.timestamp.slice(11,16);
                        this.chats.push(chat);
                    });
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;;
                }
            }
        },
        mounted(){
            document.addEventListener('click', this.handleClickOutside);
            this.buildChatPreview();
        }
    }
</script>


<template>
    <div class="navbar-dark">
        <!-- User photo, name and button to change those -->
        <div class="user-info">
            <img class="img-circular":src="userphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
            <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{username}}</h3>
            <img @click="changeusernameshown = true" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;"/>
        </div>

        <!-- Searchbox to search users -->
        <div class="searchbox" ref="boxsearchuser">
            <input class="searchbox-user" v-model="searcheduser" placeholder="username" @keyup.enter="searchUser(sercheduser)">
            <div v-if="users.length>0" class="searched-dropdown">
                <ul>
                    <li v-for="user in users":key="user.username">
                        {{user.username}}
                    </li>
                </ul>
            </div>
        </div>
    </div>


    <div class="main-screen" style="display: flex;">
        <!-- sidebar with chats -->
        <div class="sidebar-chats" style="position: absolute; left: 0%">
            <div v-if="chats.length>0" class="chats-dropdown">
                <ul>
                    <li v-for="chat in chats":key="chat.userid">
                        <div class="chatpreview">
                            <div class="chatpreviewname">
                                <img class="img-circular" :src="chat.groupphoto" style="width: 32px; height: 32px;">
                                <h3 style="margin-left: 10px; margin-bottom: 0;">{{chat.groupname}}</h3>
                                <div class="timepreview">{{chat.lastmessage.timestamp}}</div>
                            </div>
                            <div class="messagepreview">
                                <b>{{chat.lastmessage.username}}: </b>
                                <img v-if="chat.lastmessage.photo.length>0" src="/assets/photo-icon.svg" style="height: 24px; width: 24px; margin-right: 10px; margin-left: 5px;">
                                &nbsp;{{chat.lastmessage.text}}
                                <img class="checkmark" v-if="chat.lastmessage.isallread && chat.lastmessage.userid==this.userid" src="/assets/double-check.svg" style="height: 24px; width: 24px; color: blue;">
                                <img class="checkmark" v-else-if="chat.lastmessage.isallreceived && chat.lastmessage.userid==this.userid" src="/assets/double-check.svg" style="height: 24px; width: 24px;">
                                <img class="checkmark" v-else-if="chat.lastmessage.userid==this.userid" src="/assets/single-check.svg" style="height: 24px; width: 24px;">
                            </div>
                        </div>
                    </li>
                </ul>
            </div>
        </div>


        <div class="box-container">
            <!-- Box to change username and photo -->
            <div v-if="changeusernameshown" class="blurred-box">
                Enter a new username:
                <input class="new-username" v-model="newusername" placeholder="new username">
                <input type="file" accept="image/*" ref="changeimageInput" style="display: none;" @change="changePhotoFileSelect"/>
                <div v-if="newuserphoto" style="display: flex; flex-direction: column; align-items: center;">
                    Preview profile pic:
                    <img class="img-circular" :src="newuserphoto" style="width: 64px; height: 64px; background-color: #695d5d;"/>
                </div>
                <button @click="changePhotoButton">Select Photo</button>
                <button class="change-button" @click="changeUsernamePhoto">Confirm</button>
                <button class="cancel-button" @click="resetChangeUsernamePrompt">Cancel</button>
                <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
            </div>
        </div>
    </div>
</template>

<style>
body {
        background-color: rgb(33, 28, 28);
    }

    /* Circular images */
    .img-circular{
        border-radius: 50%; 
        object-fit: cover;
    }
    

    /* Navbar */
    .navbar-dark{
        background-color: rgb(38, 38, 44);
        position: relative;
        display: flex;
        align-items: center;
        height: 10vh;
        width: 100%;
    }

    /* Info of user */
    .user-info{
        margin-left: 10px;
        color: whitesmoke;       /* Colore del testo */
        background-color: #695d5d;
        height: 36px;
        display: flex;
        justify-content: center;
        align-items: center;
        border-radius: 20px;
    }


    /* Searchbox */
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
        color: whitesmoke;
        background: #332a2a;
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


    /* Box change username and photo */
    .box-container {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 100%;
        height: 90vh;
    }
    .blurred-box {
        display: flex;
        width: 400px;
        height: 400px;
        background-color: rgba(0, 0, 255, 0.5);
        justify-content: center;
        align-items: center;
        flex-direction: column;
        color: whitesmoke;
        position: relative;
        z-index: 1000;
    }
    .new-username{
        color: whitesmoke;
        background: #332a2a;
    }
    .change-button{
        position: absolute;
        top: 90%;
        left: 80%;
        transform: translate(-50%, -50%);
    }
    .cancel-button{
        position: absolute;
        top: 90%;
        left: 20%;
        transform: translate(-50%, -50%);
    }

    /* Sidebar */
    .sidebar-chats{
        background-color: rgb(38, 38, 44);
        height: 90vh;
        width: 20vw;
    }
    .chats-dropdown {
        background-color: #332a2a;
        width: 100%;
    }
    .chats-dropdown ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }

    .chats-dropdown li {
        padding: 10px;
        cursor: pointer;
        height: 80px;
        display: flex;
        align-items: center;
    }

    .chats-dropdown li:hover {
        background-color: #695d5d;
    }

    .chatpreview{
        width: 100%;
    }
    .chatpreviewname{
        display: flex;
        flex-direction: row;
        align-items: center;
        position: relative;
        width: 100%;
    }
    .messagepreview{
        display: flex;
        flex-direction: row;
        align-items: center;
        position: relative;
        width: 100%;
    }
    .checkmark{
        position: absolute;
        left: 95%;
        top: 50%;
        transform: translate(-50%,-50%);
    }
    .timepreview{
        position: absolute;
        left: 95%;
        top: 50%;
        transform: translate(-50%,-50%);
    }
</style>