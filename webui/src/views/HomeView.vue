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

                // variable to make appear the blurred boxes
                boxshown: 0,
                /*  1 = change profile
                    2 = change group
                    3 = add users to group
                    4 = crate group
                */

                // for changing username and photo
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
                n_messageshown: 0,

                // for comments
                commentshown: 0,
                commentemoji: null,

                // for forwarding message
                messageToforward: 0,

                // for changing group name and photo
                changedgroupinfo: false,
                newgroupname : null,
                newgroupphoto: null,
            }
        },
        methods: {
            handleClickOutside(event) {
                // Check if the click is outside the search box to search users
                if (this.$refs.boxsearchuser && !this.$refs.boxsearchuser.contains(event.target)) {
                    this.users = [];
                    this.searcheduser = null;
                }
                if (this.messageToforward!=0 && event.target.id != "forwardbutton" && !this.$refs.chatlist.contains(event.target)){
                    this.messageToforward = 0;
                    console.log(this.messageToforward);
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
                this.boxshown = 0;
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
                        this.errormsg = e.response.data;
                    }
                }
                if (this.newuserphoto){
                    try {
                        let response = await this.$axios.put("/users/"+this.userid+"/photo",{photo: this.newuserphoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        localStorage.setItem('userphoto', this.newuserphoto);
                        this.userphoto = this.newuserphoto;
                        this.changedinfo = true;
                    } catch (e) {
                        this.errormsg = e.response.data;
                    }
                }
                this.newusername = null;
                this.newuserphoto = null;
                if (this.changedinfo){
                    this.boxshown = 0;
                    this.errormsg = null;
                }
            },
            async buildChatPreview(){
                this.errormsg = null;
                this.chats=[];
                try {
                    let response = await this.$axios.get("/chats",{headers:{"Authorization": `Bearer ${this.userid}`}});
                    response.data.forEach(chat => {
                        if (chat.groupname.length>16){
                            chat.groupname = chat.groupname.slice(0,16)+"...";
                        }
                        if (chat.lastmessage.text.length>18){
                            chat.lastmessage.text = chat.lastmessage.text.slice(0,18)+"...";
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
                    //this.errormsg = e.response.status + ": " + e.response.data;
                    this.errormsg = toString(e);
                }
            },
            async buildMainChat(chatid){
                this.errormsg = null;
                this.mainchat = null;
                if(this.messageToforward != 0){
                    this.forwardMessage(chatid);
                }else{
                    try {
                        let response = await this.$axios.get("/chats/"+chatid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.mainchat=response.data;
                        if (this.mainchat.groupname.length>16){
                            this.mainchat.groupname = this.mainchat.groupname.slice(0,16)+"...";
                        }
                        this.mainchat.messagelist.forEach( message => {
                            message.timestamp = message.timestamp.slice(11,16);
                        });
                        this.chatshown = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                    this.buildChatPreview();
                }
            },
            async closeMainChat(){
                this.errormsg = null;
                this.mainchat = null;
                this.messagephoto = null;
                this.messagetext = null;
                this.n_messageshown = 0;
                this.commentshown = 0;
                this.commentemoji = null;
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
                if(this.messagetext==null && this.messagephoto==null){
                    return
                }
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
            },
            async forwardMessage(chatid){
                try {
                    let response = await this.$axios.post("/chats/"+chatid+"/forwardedmessages",{messageid: this.messageToforward},{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.messageToforward = 0;
                    this.buildMainChat(chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
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
            },
            async deleteMessage(message){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                    if(this.mainchat.messagelist.length>1){
                        this.buildMainChat(message.chatid);
                    }else{
                        this.mainchat = null;
                        this.buildChatPreview();
                    }
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async showComments(message){
                if(this.commentshown!=message.messageid){
                    this.commentshown=message.messageid;
                }else{
                    this.commentshown=0;
                }
            },
            async commentMessage(message){
                const emojiRegex = /[\u{1F600}-\u{1F64F}|\u{1F300}-\u{1F5FF}|\u{1F680}-\u{1F6FF}|\u{1F700}-\u{1F77F}|\u{1F780}-\u{1F7FF}|\u{1F800}-\u{1F8FF}|\u{1F900}-\u{1F9FF}|\u{1FA00}-\u{1FA6F}|\u{2600}-\u{26FF}|\u{2700}-\u{27BF}|\u{FE00}-\u{FE0F}]/gu;
                if(emojiRegex.test(this.commentemoji)){
                    try{
                        let response = await this.$axios.put("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid+"/comments",{emoji: this.commentemoji},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.buildMainChat(message.chatid);
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                    
                }
                this.commentemoji = null;
            },
            async deleteComment(message){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/messages/"+message.messageid+"/comments",{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.buildMainChat(message.chatid);
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },
            async startForwardingMessage(messageid){
                this.messageToforward = messageid;
                console.log(this.messageToforward);
            },
            async leaveGroup(){
                try{
                    let response = await this.$axios.delete("/chats/"+this.mainchat.chatid+"/users/"+this.userid,{headers:{"Authorization": `Bearer ${this.userid}`}});
                    this.mainchat = null;
                    this.buildChatPreview();
                } catch (e) {
                    this.errormsg = e.response.status + ": " + e.response.data;
                }
            },

            // Button to change the photo of group handler
            changeGroupPhotoFileSelect(){
                const file = this.$refs.changeGroupPhotoInput.files[0];
                if (file) {
                    const reader = new FileReader();
                    reader.onload = (e) => {
                        this.newgroupphoto = e.target.result;
                    };
                    reader.readAsDataURL(file);
                }
            },
            changeGroupPhotoButton(){
                this.$refs.changeGroupPhotoInput.click();
            },


            async resetChangeGroupPrompt(){
                this.newgroupname = null;
                this.newgroupphoto = null;
                this.boxshown = 0;
                this.errormsg = null;
            },
            async changeGroupNamePhoto(){
                this.errormsg = null;
                if (this.newgroupname){
                    try {
                        let response = await this.$axios.put("/chats/"+this.mainchat.chatid+"/name", {groupname: this.newgroupname.trim()},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.changedgroupinfo = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }
                if (this.newgroupphoto){
                    try {
                        let response = await this.$axios.put("/chats/"+this.mainchat.chatid+"/photo",{photo: this.newgroupphoto},{headers:{"Authorization": `Bearer ${this.userid}`}});
                        this.changedgroupinfo = true;
                    } catch (e) {
                        this.errormsg = e.response.status + ": " + e.response.data;
                    }
                }
                this.newgroupname = null;
                this.newgroupphoto = null;
                if (this.changedgroupinfo){
                    this.boxshown = 0;
                    this.errormsg = null;
                    this.buildMainChat(this.mainchat.chatid);
                }
            },


            // function to refresh the views
            async refresh(){
                this.messageToforward = 0;
                this.n_messageshown = 0;
                if(this.mainchat){
                    this.buildMainChat(this.mainchat.chatid);
                }else{
                    this.buildChatPreview();
                }
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
                <img class="img-circular" :src="userphoto" style="width: 32px; height: 32px; margin-left: 2px;"/>
                <h3 style="margin-left: 10px; margin-bottom: 0; margin-right: 10px;">{{username}}</h3>
                <img @click="boxshown = 1" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;" v-if="boxshown != 1"/>
            </div>

            <!-- Searchbox to search users -->
            <div class="searchbox" ref="boxsearchuser">
                <img src="/assets/search.svg" style="width: 32px; height: 32px; margin-top: 2px; margin-left: 2px;">
                <div class="searchbox-userlist">
                    <input class="searchbox-user" v-model="searcheduser" placeholder="Search user" @keyup.enter="searchUser(sercheduser)">
                    <div v-if="users.length>0" class="searched-dropdown">
                        <ul>
                            <li v-for="user in users" :key="user.username" @click="openChatFromUser(user)">
                                {{user.username}}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>


        <div class="main-screen">
            <!-- Sidebar with chats -->
            <div class="sidebar-chats" ref="sidebar">
                <div class="sidebar-buttons">
                    <img src="/assets/add-square.svg" style="width: 32px; height: 32px; margin-left: 10px; margin-right: 10px; cursor: pointer;">
                    <div v-if="messageToforward!=0" style="color: whitesmoke;">
                        Forward to...<br>
                        (click ouside sidebar to cancel)
                    </div>
                    <img src="/assets/refresh.svg" style="width: 32px; height: 32px; margin-left: 10px; margin-right: 10px; cursor: pointer;" @click="refresh">
                </div>
                <div v-if="chats.length>0" class="chats-dropdown" ref="chatlist">
                    <ul>
                        <li v-for="chat in chats" :key="chat.userid" @click="buildMainChat(chat.chatid)">
                            <div class="chatpreview">
                                <div class="chatpreviewname">
                                    <img class="img-circular" :src="chat.groupphoto" style="width: 32px; height: 32px;">
                                    <h4 style="margin-left: 10px; margin-bottom: 0;">{{chat.groupname}}</h4>
                                    <div class="timepreview">{{chat.lastmessage.timestamp}}</div>
                                </div>
                                <div class="messagepreview">
                                    <b>{{chat.lastmessage.username}}: </b>
                                    <img v-if="chat.lastmessage.photo.length>0" src="/assets/photo-icon.svg" style="height: 24px; width: 24px; margin-left: 5px;">
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
                        <img v-if="mainchat.isgroup && boxshown != 2" src="/assets/pencil.svg" style="width: 16px; height: 16px; cursor: pointer; margin-right: 10px;" @click="boxshown = 2"/>
                    </div>
                    <button v-if="mainchat.isgroup" class="leavebutton" @click="leaveGroup">
                        <img src="/assets/leave.svg" style="height: 32px; width: 32px;">
                        Leave
                    </button>
                </div>
                <!-- Message screen in mainchat -->
                <div class="message-screen">
                    <div class="messagelist" id="messagelist">
                        <ul>
                            <li v-for="message in mainchat.messagelist" :key="message.messageid">
                                <span v-if="message.userid==this.userid" style="display:flex; flex-direction: row-reverse; width: calc(100vw - 360px); height: 100%; ">
                                    <div class="messagebox-you">
                                        <div v-if="message.isforwarded" class="forwarded-info" style="display: flex; justify-content: right;">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px; margin-right: 5px;">
                                            Forwarded
                                        </div>
                                        <div class="messagebox-username" style="text-align: right;">
                                            <b><h3 style="margin-bottom: 0;">You</h3></b>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                        <div class="messagebox-time">
                                            <img class="messagebox-checkmark" v-if="message.isallread" src="/assets/double-check-blue.svg" style="height: 24px; width: 24px;">
                                            <img class="messagebox-checkmark" v-else-if="message.isallreceived" src="/assets/double-check.svg" style="height: 24px; width: 24px;">
                                            <img class="messagebox-checkmark" v-else src="/assets/single-check.svg" style="height: 24px; width: 24px;">
                                            {{message.timestamp}}
                                        </div>
                                        <div class="messagebox-buttons">
                                            <img src="/assets/forward.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="startForwardingMessage(message.messageid)" id="forwardbutton">
                                            <img src="/assets/trashcan.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="deleteMessage(message)">
                                            <div>
                                                <img src="/assets/comment.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="showComments(message)">
                                                {{message.commentlist.length}}
                                                </div>
                                        </div>
                                        <div v-if="commentshown==message.messageid" class="messagebox-comment">
                                            <input class="commenttext" v-model="commentemoji" maxlength="2" placeholder="Emoji" @keyup.enter="commentMessage(message)">
                                            <div class="commentlist">
                                                <ul>
                                                    <li v-for="comment in message.commentlist" :key="comment.userid">
                                                        {{comment.username}}: {{comment.emoji}}
                                                        <img v-if="comment.userid==this.userid" src="/assets/trashcan.svg" style="height: 16px; width: 16px; cursor: pointer;" @click="deleteComment(message)">
                                                    </li>
                                                </ul>
                                            </div>
                                        </div>
                                    </div>
                                </span>
                                <span v-else style="display:flex; width: calc(100vw - 360px); height: 100%;">
                                    <div class="messagebox-other">
                                        <div v-if="message.isforwarded" class="forwarded-info">
                                            <img src="/assets/forward.svg" style="width: 24px; height: 24px; margin-right: 5px;">
                                            Forwarded
                                        </div>
                                        <div class="messagebox-username">
                                            <b><h3 style="margin-bottom: 0;">{{message.username}}</h3></b>
                                        </div>
                                        <img v-if="message.photo" :src="message.photo" style="max-width: 200px; max-height: 200px; margin: 10px;">
                                        <div class="messagebox-text">
                                            {{message.text}}
                                        </div>
                                        <div class="messagebox-time">
                                            {{message.timestamp}}
                                        </div>
                                        <div class="messagebox-buttons">
                                            <div>
                                            <img src="/assets/comment.svg" style="height: 24px; width: 24px; cursor: pointer;" @click="showComments(message)">
                                            {{message.commentlist.length}}
                                            </div>
                                            <img src="/assets/forward.svg" style="height: 24px; width: 24px;">
                                        </div>
                                        <div v-if="commentshown==message.messageid" class="messagebox-comment">
                                            <input class="commenttext" v-model="commentemoji" maxlength="2" placeholder="Emoji" @keyup.enter="commentMessage(message)">
                                            <div class="commentlist">
                                                <ul>
                                                    <li v-for="comment in message.commentlist" :key="comment.userid">
                                                        {{comment.username}}: {{comment.emoji}}
                                                        <img v-if="comment.userid==this.userid" src="/assets/trashcan.svg" style="height: 16px; width: 16px; cursor: pointer;" @click="deleteComment(message)">
                                                    </li>
                                                </ul>
                                            </div>
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
                        <input class="message-textbox" v-model="messagetext" placeholder="Write a message" @keyup.enter="sendMessageorCreateChat">
                    </div>
                    <img v-if="messagetext || messagephoto" src="/assets/send.svg" style="width: 32px; height: 32px; cursor: pointer; margin-left: 10px; margin-right: 10px;" @click="sendMessageorCreateChat">
                    <div v-if="messagephoto" class="messagephoto-preview">
                        <img :src="messagephoto" style="width: 250px; height: 250px;">
                    </div>
                </div>
            </div>
        </div>
    </div>


    <!-- Box to change username and photo -->
    <div class="box-container" v-if="boxshown == 1">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Profile Info</h1>
            <div class="new-username-box">
                Enter a new username:
                <div class="new-username-container">
                    <input class="new-username" v-model="newusername" placeholder="New username">
                </div>
            </div>
            <input type="file" accept="image/*" ref="changePhotoInput" style="display: none;" @change="changePhotoFileSelect"/>
            <button class="selectphoto-button" @click="changePhotoButton">Select Photo</button>
            <div v-if="newuserphoto" style="display: flex; flex-direction: column; align-items: center;">
                Preview profile pic:
                <img class="img-circular" :src="newuserphoto" style="width: 64px; height: 64px; background-color: #695d5d;"/>
            </div>
            <button class="confirm-button" @click="changeUsernamePhoto">Confirm</button>
            <button class="cancel-button" @click="resetChangeUsernamePrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        </div>
    </div>

    <!-- Box to change group name and photo -->
    <div class="box-container" v-if="boxshown == 2">
        <div class="blurred-box">
            <h1 style="margin-top: 20px;">Group Info</h1>
            <div class="new-username-box">
                Enter a new group name:
                <div class="new-username-container">
                    <input class="new-username" v-model="newgroupname" placeholder="New group name">
                </div>
            </div>
            <input type="file" accept="image/*" ref="changeGroupPhotoInput" style="display: none;" @change="changeGroupPhotoFileSelect"/>
            <button class="selectphoto-button" @click="changeGroupPhotoButton">Select Photo</button>
            <div v-if="newgroupphoto" style="display: flex; flex-direction: column; align-items: center;">
                Preview group photo:
                <img class="img-circular" :src="newgroupphoto" style="width: 64px; height: 64px; background-color: #695d5d;"/>
            </div>
            <button class="confirm-button" @click="changeGroupNamePhoto">Confirm</button>
            <button class="cancel-button" @click="resetChangeGroupPrompt">Cancel</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
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
        height: 100%;
        position: absolute;
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
        position: relative;
        color: whitesmoke;
    }
    .new-username-box{
        margin: 20px;
        width: 80%;
    }
    .new-username-container{
        width: 100%;
        padding: 5px;
        background-color: #171717;
        border-radius: 20px;
    }
    .new-username{
        background: none;
        border: none;
        outline: none;
        width: 100%;
        color: whitesmoke;
    }
    .selectphoto-button {
        padding: 5px;
        padding-left: 20px;
        padding-right: 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #252525;
        color: white;
        margin-bottom: 10px;
    }
    .selectphoto-button:hover {
        background-color: black;
    }
    .confirm-button{
        position: absolute;
        top: 90%;
        left: 80%;
        transform: translate(-50%, -50%);
        padding: 5px;
        padding-left: 20px;
        padding-right: 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #252525;
        color: white;
    }
    .confirm-button:hover{
        background-color: black;
    }
    .cancel-button{
        position: absolute;
        top: 90%;
        left: 20%;
        transform: translate(-50%, -50%);
        padding: 5px;
        padding-left: 20px;
        padding-right: 20px;
        border-radius: 20px;
        border: none;
        outline: none;
        background-color: #252525;
        color: white;
    }
    .cancel-button:hover{
        background-color: rgb(182, 52, 52);
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
    .sidebar-buttons{
        width: 100%;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    .chats-dropdown {
        width: 100%;
        height: calc(100% - 36px); /* minus height of buttons */
        overflow-y: scroll;
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
    .leavebutton{
        position: absolute;
        left: 91%;
        background-color: transparent;
        color: whitesmoke;
        width: 120px;
        border-radius: 20px;
        border-color: rgb(182, 52, 52);
    }
    .leavebutton:hover{
        background-color: rgb(182, 52, 52);
    }

    /* Message screen */
    .message-screen{
        height: calc(100% - 120px);
        width: 100%;
        overflow-y: auto;
        overflow-x: hidden;
    }
    .messagelist{
        height: 100%;
        width: 100%;
        overflow-y: auto;
        overflow-x: hidden;
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

    /* Box containing each message */
    .messagebox-you{
        width: max-content;
        background-color: green;
        border-radius: 20px;
        border-top-right-radius: 0;
        margin-right: 20px;
        max-width: 600px;
        display: flex;
        flex-direction: column;
        align-items: end;
    }
    .messagebox-other{
        width: max-content;
        background-color: green;
        border-radius: 20px;
        border-top-left-radius: 0;
        margin-left: 20px;
        max-width: 600px;
        display: flex;
        flex-direction: column;
        align-items: start;
    }
    .forwarded-info{
        margin-left: 15px;
        margin-right: 15px;
    }
    .messagebox-text{
        margin-left: 15px;
        margin-right: 15px;
        word-break: break-word; /* Permette di spezzare parole lunghe */
        font-size: 0.875rem;
        font-family: sans-serif;
    }
    .messagebox-username{
        margin-left: 15px;
        margin-right: 15px;
    }
    .messagebox-time{
        display: flex;
        justify-content: right;
        align-items: center;
        width: 100%;
        padding-right: 15px;
    }
    .messagebox-checkmark{
        margin-right: 5px;
    }
    .messagebox-buttons{
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 120px;
        margin-bottom: 5px;
        margin-left: 15px;
        margin-right: 15px;
    }

    /* Box containing comments */
    .messagebox-comment{
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        border-top: 2px solid black
    }
    .commenttext{
        width: 100px;
        margin-top: 10px;
        margin-bottom: 5px;
    }
    .commentlist{
        width: 100%;
        overflow-y: scroll;
        display: flex;
        flex-direction: column;
        max-height: 150px;
    }
    .commentlist ul{
        list-style: none;
        padding: 0;
        color: whitesmoke;
    }
    .commentlist li {
        padding: 5px;
        padding-left: 10px;
        padding-right: 10px;
        height: 30px;
        display: flex;
        justify-content: space-between;
        align-items: center;
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
        height: 23px;
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