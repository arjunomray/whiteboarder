let drawing = false;
let currentX = 0;
let currentY = 0;
let currentColor = getRandomColor(); // Generate a random color for each user

// Initialize WebSocket connection
const socket = new WebSocket("ws://localhost:8080/ws");

const canvas = document.getElementById("drawingCanvas");
const ctx = canvas.getContext("2d");
const clientCountDiv = document.getElementById("clientCount");

// Draw on canvas when receiving drawing data from the WebSocket server
socket.onmessage = function (event) {
    const data = JSON.parse(event.data);
    if (data.type === "clientCount") {
        clientCountDiv.innerText = `Connected Users: ${data.count}`;
        return;
    } else {
        drawLine(data.startX, data.startY, data.endX, data.endY, data.color);
    }
};

// Start drawing when mouse is pressed or touch starts
canvas.addEventListener("mousedown", (e) => {
    drawing = true;
    currentX = e.offsetX;
    currentY = e.offsetY;
});
canvas.addEventListener("touchstart", (e) => {
    drawing = true;
    const touch = e.touches[0];
    const rect = canvas.getBoundingClientRect();
    currentX = touch.clientX - rect.left;
    currentY = touch.clientY - rect.top;
});

// Stop drawing when mouse is released or touch ends
canvas.addEventListener("mouseup", () => {
    drawing = false;
});
canvas.addEventListener("touchend", () => {
    drawing = false;
});

// Draw on canvas when moving mouse or touch moves
canvas.addEventListener("mousemove", (e) => {
    if (!drawing) return;
    const x = e.offsetX;
    const y = e.offsetY;

    drawLine(currentX, currentY, x, y, currentColor);
    // Send drawing data to the server via WebSocket
    socket.send(JSON.stringify({ startX: currentX, startY: currentY, endX: x, endY: y, color: currentColor }));

    currentX = x;
    currentY = y;
});
canvas.addEventListener("touchmove", (e) => {
    if (!drawing) return;
    const touch = e.touches[0];
    const rect = canvas.getBoundingClientRect();
    const x = touch.clientX - rect.left;
    const y = touch.clientY - rect.top;

    drawLine(currentX, currentY, x, y, currentColor);
    // Send drawing data to the server via WebSocket
    socket.send(JSON.stringify({ startX: currentX, startY: currentY, endX: x, endY: y, color: currentColor }));

    currentX = x;
    currentY = y;
    e.preventDefault(); // Prevent scrolling when drawing
});

// Helper function to draw a line on the canvas
function drawLine(startX, startY, endX, endY, color) {
    ctx.strokeStyle = color;
    ctx.beginPath();
    ctx.moveTo(startX, startY);
    ctx.lineTo(endX, endY);
    ctx.stroke();
}

// Helper function to generate a random color
function getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}