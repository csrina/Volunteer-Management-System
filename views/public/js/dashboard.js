function load() {    
    let req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
	let data = JSON.parse(req.response);
	input(data);
    });
    req.open("GET", "http://localhost:8080/api/v1/dashboard");
    req.send();
}
function input(data) {
    let done = document.getElementById("hoursDone");
    let booked = document.getElementById("hoursBooked");
    let table = document.getElementById("events");

    for (let i = 0; i < data.length; i++) {
	let row = document.createElement("tr");
	row.innerHTML = `
	    <td>${data[i].BlockStart} </td>
	    <td> ${data[i].BlockEnd}</td>`
	table.appendChild(row);
    }
}
console.log("hi")
load();
