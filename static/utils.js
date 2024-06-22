function refreshTime() {
    var dateString = new Date().toLocaleTimeString('en-US', { hour12: false });
    var formattedString = dateString.replace(", ", " - ");
    let timeDisplay = document.getElementById("time");
    timeDisplay.innerHTML = formattedString;
}

setInterval(refreshTime, 1000);