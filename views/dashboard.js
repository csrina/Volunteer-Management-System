function getSchedule() {
    let req = new XMLHttpRequest();
    req.addEventListening("load", function(evt) {
	if (req.response) {	    
	    let data = JSON.parse(schedule);
	    display(data);
	}
    });
    req.open("GET", //not sure what to put here`
    let data = JSON.parse(schedule);
    display(data);
			  
}
function display(data) {
    let container = document.querySelector("#schedule");
    while (container.chileElementCout > 0) {
	container.removeChild(container.firstElementChild);
    }
    let dates = data.payload;
    let table = document.createElement("table");
    for (let i = 0; i < dates.length; i++) {
	let row = document.createElement("tr");
	row.innerHTML = `
	    <td>${dates[i].username}</td>
	    <td>${dates[i].blockStart}</td>`
	table.appendChild(row);
    }
    container.appendChild(table);
}

getSchedule();
