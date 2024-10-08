;(() => {
    const htcon = htmx.createWebSocket('wss://' + document.location.host + '/api/empl/ws')

    window.onload = function () {

    }
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

            // htcon.addEventListener('htmx:wsBeforeMessage', (evt) => {
            //     //TODO Handle messages that just been received
            // })
            // htcon.addEventListener('htmx:wsAfterMessage', (evt) => {
            //     const messagesDiv = document.getElementById('messages');
            //     const newMessage = document.createElement('p');
            //     newMessage.textContent = `Received: ${evt.data}`;
            //     messagesDiv.appendChild(newMessage);
            // })
            htcon.onmessage = function(evt) {

                try {
                    // Parse the JSON string from the WebSocket message
                    const data = JSON.parse(evt.data);

                    // Create an event object
                    const event = new Event(data.type, data.payload);

                    // Pass the event to your routing logic
                    routeEvents(event);
                } catch (error) {
                    console.error("Error parsing WebSocket message:", error);
                }
            }



        }


    }
    dial()

    class Event {
        constructor(type, payload) {
            this.type = type;
            this.payload = payload;
        }
    }

    function routeEvents(event){
        if (event.type === undefined) {
            alert('no type field in the event')
        }
        switch (event.type) {
            case "message":
                console.log("Received: message");
                break;
            case "new_order":
                console.log("Received: new order");

                const notifCountElement = document.querySelector('.notif .count');

                // Check if the element exists
                if (!notifCountElement) {
                    console.error('Notification count element not found.');
                    return;
                }

                // Get the current count from the element's text
                let currentCount = parseInt(notifCountElement.textContent) || 0;
                console.log('Current count is:', currentCount);

                // Increment the count
                const newCount = currentCount + 1;

                // Update the notification count element with the new count
                notifCountElement.textContent = newCount;

                console.log('Updated count is:', newCount);
                break;
            case "assign_address":
                console.log("Received: updated assigned address");
                break;
            default:
                alert(`unsupported type ${event.type}`);
                break;
        }
    }

    window.sendEvent = function (eventName, payload) {
        if (htcon.readyState === htcon.OPEN) {
            const event = new Event(eventName,payload);
            htcon.send(JSON.stringify(event));
        } else {
            console.warn("WebSocket is not open. Attempting to reconnect...");
            //TODO:reconnectWebSocket();

            // Optionally, warn the user about the disconnection
            alert("WebSocket connection lost. Please wait while we try to reconnect.");
        }

    }

    window.isWebSocketOpen = function () {
        return htcon.readyState === htcon.OPEN;
    }

    function updateNotificationCount() {
        const notifCountElement = document.querySelector('.notif .count');

        if (!notifCountElement) {
            console.error('Notification count element not found');
            return;
        }

        // Get the current count, parse it to an integer, and increment by 1
        let currentCount = parseInt(notifCountElement.textContent) || 0;
        notifCountElement.textContent = currentCount + 1;

        console.log('Updated notification count:', currentCount + 1); // Debug log
    }

    document.querySelector('.notif').addEventListener('click', function() {
        const notifCountElement = document.querySelector('.notif .count');
        notifCountElement.textContent = '0'; // Reset the count to 0
    });

})()

