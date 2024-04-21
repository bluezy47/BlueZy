let webSocketConn = null;

// Chat List DOM elements
const chatListContainer = document.querySelector(".chat-selection-wrapper .chat-selection-sub-wrapper");
const chatListItems = chatListContainer.querySelectorAll("div.select-chat-user-wrapper div.chat-user");

// Middle chat container DOM elements
const middleChatContainer = document.querySelector(".selected-chat-wrapper");
const chatInput = middleChatContainer.querySelector(".chat-input-wrapper input.chat-input");
const chatSendBtn = middleChatContainer.querySelector(".chat-input-wrapper button.chat-send-btn");
const openedChatsContainer = middleChatContainer.querySelector(".selected-sub-chat-wrapper div.chat-list-wrapper");
const selectedChatHeader = middleChatContainer.querySelector(".selected-user-wrapper .selected-user-info");

//
const initWebSocketConnection = () => {
    webSocketConn = new WebSocket('ws://localhost:5050/ws/helloworld?email=user1@gmail.com'); // todo: add dynamic email fucker!

    webSocketConn.onopen = function(event) {
        console.log(':: WebSocket Connection was established! ::');
    };
    
    webSocketConn.onmessage = function(event) {
        let data = JSON.parse(event.data);
        console.log('Message received from server:', data);
        appendMessageToChat(data.message, "received");
    };
    
    webSocketConn.onerror = function(event) {
        console.error('WebSocket error:', event);
    };
    
    webSocketConn.onclose = function(event) {
        console.log(':: WebSocket Connection was closed! ::');
    };
};
//
// use to append the chat message to the chat list
const appendMessageToChat = (data, type="sent") => {
    let testMessage = data.message;
    let messageWrapper = document.createElement('div');
    let messageElement = document.createElement('p');
    messageWrapper.classList.add('single-chat-wrapper');
    if (type == "received")
        messageWrapper.classList.add("received");
    else
        messageWrapper.classList.add("sent");
    //
    messageElement.innerHTML = testMessage;
    messageWrapper.appendChild(messageElement);
    // Append message to chat
    openedChatsContainer.appendChild(messageWrapper);
}
//
// use to send the message to the server
const sendMessage = (data) => {
    // check websocket connection
    if (webSocketConn) {
        let sendData = {
            'message': data,
            'to': data.to
        }
        webSocketConn.send(JSON.stringify(sendData));
        return true;
    }
    console.log(':: WebSocket Connection is not established! ::');
    return false;
}
chatSendBtn.onclick = (e) => {
    let message = chatInput.value;
    if (message) {
        let data = {
            'message': message,
            'to': 'user1@gmail.com'
        }
        sendMessage(data);
        appendMessageToChat(data);
    }
    chatInput.value = '';
}

// Click event for the chat list items
const viewChat = (chatUser) => {
    let profilePicElm = selectedChatHeader.querySelector("img");
    let userNameElm = selectedChatHeader.querySelectorAll("div.user-info span")[0];
    let onlineStatusElm = selectedChatHeader.querySelectorAll("div.user-info span")[1];
    //
    let proPic = chatUser.querySelector("img").src;
    let userName = chatUser.getAttribute("full-name");
    let onlineStatus = chatUser.getAttribute("last-online");
    //
    profilePicElm.src = proPic;
    userNameElm.innerHTML = userName;
    onlineStatusElm.innerHTML = onlineStatus;

    // clear the Previous chat
    openedChatsContainer.innerHTML = '';
};
chatListItems.forEach((chatItem) => {
    chatItem.onclick = (e) => {
        viewChat(chatItem);
    }
});


//
// ----- INITILIZE -----
initWebSocketConnection();
