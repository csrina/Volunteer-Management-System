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
    let needed = 0;
    if (data.children == 1) {
	needed = 2.5;
    } else {
	needed = 5;
    }
    let done = document.getElementById("hoursDone");
    let booked = document.getElementById("hoursBooked");
    let table = document.getElementById("events");

    if (data.hoursDone/needed > 0.99) {
	done.style.color = "green"
    } else if (data.hoursDone/needed > 0.66) {
	done.style.color = "yellow"
    } else if (data.hoursDone/needed > 0.33) {
	done.style.color = "orange"
    } else {
	done.style.color = "red"
    }
    
    done.innerHTML = data.hoursDone;
    booked.innerHTML = data.hoursBooked;
    for (let i = 0; i < data.eventlist.length; i++) {
	let item = document.createElement("li");
	item.innerHTML = data.eventlist[i];
	table.appendChild(item)
    }
}
load();
