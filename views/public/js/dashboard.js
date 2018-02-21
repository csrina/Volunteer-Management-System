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
    let time = 0;
    for (let i = 0; i < data.length; i++) {
	let row = document.createElement("tr");
	let start = new Date(data[i].BlockStart);	
	let startMonth = start.getMonth();
	let startDate = start.getDate();
	let startHours = start.getHours();
	let startMinutes = start.getMinutes();
	let end = new Date(data[i].BlockEnd);
	let endMonth = end.getMonth();
	let endDate = end.getDate();
	let endHours = end.getHours();
	let endMinutes = end.getMinutes();
	row.innerHTML = `
	    <td style="border: solid">${startMonth}, ${startDate} at ${startHours}:${startMinutes}</td>
	    <td style="border: solid">${endMonth}, ${endDate}</td>`
	time += Math.abs(new Date(data[i].BlockEnd)- new Date(data[i].BlockStart))
	table.appendChild(row);
    }
}
load();
