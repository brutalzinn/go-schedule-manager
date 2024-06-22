const now = new Date();
const dayOfWeek = now.toLocaleString('en-us', { weekday: 'long' }).toLowerCase();
const socket = new WebSocket(`ws://localhost:8000/ws/${dayOfWeek}`);
const history = document.getElementById("history")
const next = document.getElementById("next")

socket.onmessage = function (event) {
    const message = JSON.parse(event.data);
    console.log(message)
    if (message.action == "initialEvents") {
        history.replaceChildren()
        next.replaceChildren()
        for (let i = 0; i < message.data.next.length; i++) {
            let shedule = message.data.next[i]
            addNext(shedule.data.html)
        }
        for (let i = 0; i < message.data.history.length; i++) {
            let shedule = message.data.history[i]
            addHistory(shedule.data.html)
        }
    }
    if (message.action == "audio") {
        const audio = new Audio(message.data.audio);
        audio.play();
    }
};


function addHistory(data) {
    const li = document.createElement("li");
    li.innerHTML = data;
    history.appendChild(li);
}

function addNext(data) {
    const li = document.createElement("li");
    li.innerHTML = data;
    next.appendChild(li);
}


socket.onerror = function (error) {
    console.error("WebSocket error: ", error);
};