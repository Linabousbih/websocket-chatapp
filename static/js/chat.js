const params = new URLSearchParams(window.location.search);
const room = params.get("room");

if (!room) {
  alert("No room specified. Redirecting to homepage...");
  window.location.href = "/";
}

const socket = new WebSocket(`ws://${location.host}/room?room=${room}`);

socket.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data); // parse JSON string into object
    const msg = document.createElement("div");
    msg.textContent = `${data.name}: ${data.message}`;  // show nicely formatted message
    document.getElementById("messages").appendChild(msg);

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
