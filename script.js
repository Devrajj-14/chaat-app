let ws;

function connect() {
    ws = new WebSocket('ws://127.0.0.1:9090/ws');

    ws.onopen = () => {
        console.log('Connected to the WebSocket server');
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        displayMessage(message.body);
    };

    ws.onerror = (event) => {
        console.error('WebSocket error observed:', event);
    };

    ws.onclose = () => {
        console.log('WebSocket connection closed. Reconnecting...');
        setTimeout(connect, 1000);
    };
}

function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value.trim();

    if (message) {
        ws.send(JSON.stringify({ type: 1, body: message }));
        messageInput.value = '';
    }
}

function displayMessage(message) {
    const messages = document.getElementById('messages');
    const messageElement = document.createElement('li');
    messageElement.className = 'list-group-item';
    messageElement.textContent = message;
    messages.appendChild(messageElement);
    messages.scrollTop = messages.scrollHeight;
}

document.addEventListener('DOMContentLoaded', () => {
    connect();
});
