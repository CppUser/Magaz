;(() => {

    function dial(otp) {
        if(window["WebSocket"]){
            console.log("Supports WebSocket Client");
            // const conn = new WebSocket('wss://' + document.location.host + '/api/empl/ws?otp=' + otp);
            //
            // conn.addEventListener('open', (evt) => {})
            // conn.addEventListener('close', (evt) => {})
            //
            // conn.addEventListener('message', (evt) => {
            //
            // })


            ////////////////////////////////////////////////////////////////////////////////////
            //      Some htmx events
            //https://v1.htmx.org/extensions/web-sockets/
            ////////////////////////////////////////////////////////////////////////////////////
            const htcon = htmx.createWebSocket('wss://' + document.location.host + '/api/empl/ws')

            //htcon.wsConnecting = function (evt) {}
            htcon.addEventListener('htmx:wsConnecting', (evt) => {
                console.log('Connecting!')

            })
            //htcon.wsOpen = function (evt) {}
            htcon.addEventListener('htmx:wsOpen', (evt) => {
                console.log('WebSocket Connected!')
            })

            htcon.addEventListener('htmx:wsClose', (evt) => {
                console.log('Closed!')
            })

            htcon.addEventListener('htmx:wsError', (evt) => {

            })

            htcon.addEventListener('htmx:wsBeforeMessage', (evt) => {
                const messagesDiv = document.getElementById('messages');
                const newMessage = document.createElement('p');
                newMessage.textContent = `Received: ${evt.data}`;
                messagesDiv.appendChild(newMessage);
            })

            document.getElementById('messageForm').addEventListener('submit', function (event) {
                event.preventDefault();
                const message = document.getElementById('message').value;

                if (htcon.readyState === htcon.OPEN) {
                    htcon.send(message);
                    console.log('Message sent: ', message);
                } else {
                    console.error('WebSocket connection is not open');
                }

                document.getElementById('message').value = ''; // Clear the input field
            })

        }


    }
    dial()

})()

