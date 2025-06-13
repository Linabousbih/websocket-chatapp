// Assumes `socket` is available globally
const input = document.getElementById("msg");
const sendBtn = document.getElementById("sendBtn");
const messagesDiv = document.getElementById("messages");

// Wait for socket to open before sending
socket.addEventListener("open", () => {
  console.log("✅ WebSocket connected");
});

socket.addEventListener("message", (event) => {
  const msg = document.createElement("div");
  msg.textContent = event.data;
  messagesDiv.appendChild(msg);
});

sendBtn.addEventListener("click", () => {
  sendMessage();
});

input.addEventListener("keyup", (event) => {
  if (event.key === "Enter") {
    sendMessage();
  }
});

function sendMessage() {
  if (socket.readyState !== WebSocket.OPEN) {
    console.error("❌ Socket is not open");
    return;
  }

  const text = input.value.trim();
  if (text === "") return;

  socket.send(text);
  input.value = "";
}
