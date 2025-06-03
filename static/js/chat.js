const params = new URLSearchParams(window.location.search);
const room = params.get("room");

if (!room) {
  alert("No room specified. Redirecting to homepage...");
  window.location.href = "/";
}

const socket = new WebSocket(`ws://${location.host}/room?room=${room}`);

socket.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);

    // Create the container div
    const msgContainer = document.createElement("div");
    msgContainer.classList.add("message-container");

    // Create the username div
    const usernameDiv = document.createElement("div");
    usernameDiv.classList.add("username");
    usernameDiv.textContent = data.name;

    // Create the message div
    const messageDiv = document.createElement("div");
    messageDiv.classList.add("message");
    messageDiv.textContent = data.message;

    // Append username and message in correct order
    msgContainer.appendChild(usernameDiv);
    msgContainer.appendChild(messageDiv);

    // Append the whole message container to the messages div
    document.getElementById("messages").appendChild(msgContainer);

    // Auto-scroll
    const messagesDiv = document.getElementById("messages");
    messagesDiv.scrollTop = messagesDiv.scrollHeight;

  } catch (err) {
    console.error("Invalid JSON received:", event.data);
  }
};

function sendMessage() {
  const input = document.getElementById("msg");
  if (input.value.trim() !== "") {
    socket.send(input.value);
    input.value = "";
  }
}

document.getElementById("sendBtn").addEventListener("click", sendMessage);

document.getElementById("msg").addEventListener("keyup", function (event) {
  if (event.key === "Enter") {
    sendMessage();
  }
});
