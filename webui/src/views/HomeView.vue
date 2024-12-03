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
                chats: [],

                // for main chat
                mainchat: null,
                chatshown: false,
                messagetext: null,
                messagephoto: null,
                n_messageshown: 0
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

            // Button to change the photo of user handler
            changePhotoFileSelect(){
                const file = this.$refs.changePhotoInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.newuserphoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            changePhotoButton(){
                this.$refs.changePhotoInput.click();
            },


            async resetChangeUsernamePrompt(){
                this.newusername = null;
                this.newuserphoto = null;
                this.changeusernameshown = false;
                this.errormsg = null;
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
            },
            async buildMainChat(chatid){
                this.errormsg = null;
                this.mainchat = null;
                try {
                    let response = await this.$axios.get("/chats/"+chatid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.mainchat=response.data;
                    this.chatshown = true;
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async closeMainChat(){
                this.errormsg = null;
                this.mainchat = null;
                this.messagephoto = null;
                this.messagetext = null;
                this.n_messageshown = 0;
                this.chatshown = false;
            },

            // button to send a photo handler
            sendPhotoFileSelect(){
                const file = this.$refs.sendPhotoInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.messagephoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            sendPhotoButton(){
                this.$refs.sendPhotoInput.click();
            },

            async sendMessage(){
                try {
                    let response = await this.$axios.post("/chats/"+this.mainchat.chatid+"/messages",{text: this.messagetext,photo: this.messagephoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.messagetext = null;
                    this.messagephoto = null;
                    this.buildMainChat(this.mainchat.chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async sendMessageorCreateChat(){
                if(this.mainchat.chatid==-1){
                    try {
                        let response = await this.$axios.post("/newchat",{usernamelist:[this.username,this.mainchat.groupname],firstmessage:{text:this.messagetext}},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.messagetext = null;
                        this.messagephoto = null;
                        this.buildMainChat(response.data.chatid);
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }else{
                    this.sendMessage();
                }
                this.buildChatPreview();
            },

            async openChatFromUser(user){
                await this.buildChatPreview();
                var chatid = -1;
                for(let i=0;i<this.chats.length;i++){
                    if(this.chats[i].groupname==user.username && !this.chats[i].isgroup){
                        chatid = this.chats[i].chatid;
                        break;
                    }
                }
                if(chatid==-1){
                    this.mainchat={
                        chatid:-1,
                        groupname:user.username,
                        groupphoto:user.photo,
                        isgroup: false,
                        messagelist:[]
                    }
                    this.chatshown = true;
                }else{
                    this.buildMainChat(chatid);
                }
                this.searcheduser = null;
                this.users = [];
            }
        },
        mounted(){
            document.addEventListener('click', this.handleClickOutside);
            this.buildChatPreview();
        },
        /* Updater to check if the messagelist is shown and make the lit start from bottom */
        updated(){
            const div = document.querySelector('#messagelist');
            if (div && this.n_messageshown!=this.mainchat.messagelist.length) {
                div.scrollTop=div.scrollHeight;
                this.n_messageshown=this.mainchat.messagelist.length;
            }
        }
    }
</script>


<template>
    <div class="all-screen">
        <div class="navbar-dark">
            <!-- User photo, name and button to change those -->
            <div class="user-info">
                <img class="img-circular":src="userphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
                <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{username}}</h3>
                <img @click="changeusernameshown = true" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;" v-if="!changeusernameshown"/>
            </div>

            <!-- Searchbox to search users -->
            <div class="searchbox" ref="boxsearchuser">
                <img src="/assets/search.svg" style="width: 32px; height: 32px; margin-top: 2px; margin-left: 2px;">
                <div class="searchbox-userlist">
                    <input class="searchbox-user" v-model="searcheduser" placeholder="Search user" @keyup.enter="searchUser(sercheduser)">
                    <div v-if="users.length>0" class="searched-dropdown">
                        <ul>
                            <li v-for="user in users":key="user.username" @click="openChatFromUser(user)">
                                {{user.username}}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>


        <div class="main-screen">
            <!-- Sidebar with chats -->
            <div class="sidebar-chats">
                <div v-if="chats.length>0" class="chats-dropdown">
                    <ul>
                        <li v-for="chat in chats":key="chat.userid" @click="buildMainChat(chat.chatid)">
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
                                    <img class="checkmark" v-if="chat.lastmessage.isallread && chat.lastmessage.userid==this.userid" src="/assets/double-check-blue.svg" style="height: 24px; width: 24px;">
                                    <img class="checkmark" v-else-if="chat.lastmessage.isallreceived && chat.lastmessage.userid==this.userid" src="/assets/double-check.svg" style="height: 24px; width: 24px;">
                                    <img class="checkmark" v-else-if="chat.lastmessage.userid==this.userid" src="/assets/single-check.svg" style="height: 24px; width: 24px;">
                                </div>
                            </div>
                        </li>
                    </ul>
                </div>
            </div>

            <!-- Main chat screen -->
            <div class="main-chat" v-if="chatshown && mainchat">
                <!-- Topbar in mainchat -->
                <div class="topbar-chat">
                    <img class="backarrow" src="/assets/back-arrow.svg" style="width: 32px; height: 32px; cursor: pointer;" @click="closeMainChat"/>
                    <div class="user-info">
                        <img class="img-circular":src="mainchat.groupphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
                        <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{mainchat.groupname}}</h3>
                        <img v-if="mainchat.isgroup" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;"/>
                    </div>
                </div>
                <!-- Message screen in mainchat -->
                <div class="message-screen">
                    <div class="messagelist" id="messagelist">
                        <ul>
                            <li v-for="message in mainchat.messagelist":key="message.messageid">
                                <span v-if="message.userid==this.userid" style="display:flex; flex-direction: row-reverse; width: calc(100vw - 360px); height: 100%; ">
                                    <div class="messagebox-you">
                                        <div v-if="message.isforwarded" class="forwarded-info" style="display: flex; justify-content: right;">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px;">
                                            Forwarded {{message.userid}}{{this.userid}}
                                        </div>
                                        <div class="messagebox-username" style="text-align: right;">
                                            <b><h3 style="margin-bottom: 0;">You</h3></b>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                    </div>
                                </span>
                                <span v-else style="display:flex; width: calc(100vw - 360px); height: 100%;">
                                    <div class="messagebox-other">
                                        <div v-if="message.isforwarded" class="forwarded-info">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px;">
                                            Forwarded
                                        </div>
                                        <div class="messagebox-username">
                                            <b><h3 style="margin-bottom: 0;">{{message.username}}</h3></b>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                    </div>
                                </span>
                            </li>
                        </ul>
                    </div>
                </div>
                <!-- Bottom bar in mainchat -->
                <div class="bottombar-chat">
                    <input type="file" accept="image/*" ref="sendPhotoInput" style="display: none;" @change="sendPhotoFileSelect"/>
                    <img  v-if="!messagephoto" src="/assets/photo-icon.svg" @click="sendPhotoButton" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;">
                    <img v-else src="/assets/cross.svg" @click="messagephoto=null" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;">
                    <div class="message-text">
                        <input class="message-textbox" v-model="messagetext" placeholder="Write a message">
                    </div>
                    <img v-if="messagetext || messagephoto" src="/assets/send.svg" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;" @click="sendMessageorCreateChat">
                    <div v-if="messagephoto" class="messagephoto-preview">
                        <img :src="messagephoto" style="width: 250px; height: 250px;">
                    </div>
                </div>
            </div>
        </div>



        <!-- Box to change username and photo -->
        <div class="box-container" v-if="changeusernameshown">
            <div class="blurred-box">
                Enter a new username:
                <input class="new-username" v-model="newusername" placeholder="new username">
                <input type="file" accept="image/*" ref="changePhotoInput" style="display: none;" @change="changePhotoFileSelect"/>
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
    
    /* All page */
    .all-screen{
        width: 100vw;
        height: 100vh;
    }



    /* Navbar */
    .navbar-dark{
        background-color: rgb(38, 38, 44);
        position: relative;
        display: flex;
        align-items: center;
        height: 60px;
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
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        margin: 0;
        height: 36px;
        border-radius: 20px;
        background-color: #695d5d;
        display: flex;
        flex-direction: row;
        z-index: 100;
    }
    .searchbox-userlist{
        position: relative;
        margin-right: 34px;
        margin-top: 6.5px;
        width: 200px;
    }
    .searchbox-user{
        color: whitesmoke;
        width: 100%;
        border: none;
        background-color: #695d5d;
        margin-bottom: 6.5px;
    }
    .searchbox-user:focus{
        outline: none;
    }
    .searched-dropdown {
        position: relative;
        background-color: #332a2a;
        width: 200px;
        max-height: 120px;
        overflow-y: auto;
    }
    .searched-dropdown ul {
        list-style: none;
        padding: 0;
        margin: 0;
        color: whitesmoke;
    }

    .searched-dropdown li {
        height: 40px;
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
        height: calc(100vh - 60px);
        position: absolute;
        top: 60px;
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


    /* Main screen */
    .main-screen{
        display: flex;
        flex-direction: row;
        height: calc(100vh - 60px);
        width: 100%;
    }


    /* Sidebar */
    .sidebar-chats{
        background-color: rgb(38, 38, 44);
        height: 100%;
        width: 350px;
    }
    .chats-dropdown {
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


    /* Main chat */
    .main-chat{
        height: 100%;
        width: calc(100% - 350px);
    }

    /* Topbar */
    .topbar-chat{
        display: flex;
        align-items: center;
        background-color: rgb(39, 35, 35);
        height: 60px;
        width: 100%;
        position: relative;
    }
    .backarrow{
        margin-left: 10px;
    }

    /* Message screen */
    .message-screen{
        height: calc(100% - 120px);
        width: 100%;
        overflow-y: scroll;
    }
    .messagelist{
        height: 100%;
        width: 100%;
        overflow-y: scroll;
        display: flex;
        flex-direction: column;
    }
    .messagelist ul{
        list-style: none;
        padding: 0;
        color: whitesmoke;
    }
    .messagelist li {
        padding: 5px;
        position: relative;
    }
    .messagebox-you{
        width: max-content;
        background-color: green;
        border-radius: 20px;
        border-top-right-radius: 0;
        margin-right: 20px;
        max-width: 600px;
    }
    .messagebox-other{
        width: max-content;
        background-color: green;
        border-radius: 20px;
        border-top-left-radius: 0;
        margin-left: 20px;
        max-width: 600px;
    }
    .forwarded-info{
        margin-left: 15px;
        margin-right: 15px;
    }
    .messagebox-text{
        margin-left: 15px;
        margin-right: 15px;
        margin-bottom: 10px;
        word-break: break-word; /* Permette di spezzare parole lunghe */
    }
    .messagebox-username{
        margin-left: 15px;
        margin-right: 15px;
    }


    /* Bottom bar */
    .bottombar-chat{
        display: flex;
        align-items: center;
        background-color: rgb(39, 35, 35);
        height: 60px;
        width: 100%;
        position: relative;
    }
    .message-text{
        height: 36px;
        border-radius: 20px; 
        background-color: #695d5d;
        display: flex;
        align-items: center;
        width: 80%;
    }
    .message-textbox{
        margin-left: 10px;
        margin-right: 10px;
        color: whitesmoke;
        background: #695d5d;
        border: none;
        width: 100%;
    }
    .message-textbox:focus{
        outline: none;
    }
    .messagephoto-preview{
        position: absolute;
        top: 0%;
        height: 300px;
        width: 300px;
        transform: translateY(-100%);
        background-color: rgb(39, 35, 35);
        border-top-right-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
</style>