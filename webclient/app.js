new Vue({
    el: '#app',

    data: {
        ws_chat: null, 
        pesanBaru: '',
        isiChat: '',
        username: null, 
        joined: false 
    },

    created: function() {
        var self = this;
        this.ws_chat = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws_chat.addEventListener('message', function(e) {
            var pesan = JSON.parse(e.data);
            self.isiChat += '<div class="chip" style="background-color: #fbfcfc">'
                    + '<img src= " ' + 'https://ui-avatars.com/api/?name='+pesan.username + '">' 
                    + '<b>'+pesan.username+'</b>'
                + '<br/>'+ '</div>'
				+'<span style="background-color: #b3f7d9">'
                + pesan.message+'</span>'+ '<br/>'; 

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight;
        });
    },

    methods: {
        send: function () {
            if (this.pesanBaru != '') {
                this.ws_chat.send(
                    JSON.stringify({
                        username: this.username,
                        message: $('<p>').html(this.pesanBaru).text() 
                    }
                ));
                this.pesanBaru = ''; 
            }
        },

        join: function () {
            
            if (!this.username) {
                Materialize.toast('Masukkan Username anda', 2000);
                return
            }
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        
    }
});