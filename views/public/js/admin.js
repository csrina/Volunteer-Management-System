document.addEventListener("DOMContentLoaded", function() {
    if (window.location.href.split("/").pop() == "dashboard") {
        loadDash();
    }
})

//sets active category in top bars
function setActiveCategory() {
    let cat = window.location.href.split("/").pop();
    document.querySelector(`#${cat}Btn`).setAttribute('class','active');
}

function familyData() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
	console.log(xhttp.response);
	let httpData = JSON.parse(xhttp.response);
	
	var ctx = document.getElementById('skills').getContext('2d');
	var barData = {
	    labels: ['FirstWeek', 'SecondWeek', 'ThirdWeek',
		     'FourthWeek', 'FifthWeek', 'SixthWeek'],
	    //labels: ['Week'],
	    datasets: []
	};
	
	window.myBar = new Chart(ctx, {
	    type: 'bar',
	    data: barData,
	    options: {
		responsive: true,
		legend: {
		    position: 'right',
		},
		title: {
		    display: true,
		    text: 'Chart.js Horizontal Bar Chart'
		}
	    }
	});


	var colourList = ["#00FFFF", "#A52A2A", "#7FFF00", "#FF7F50",
			  "#006400", "#8B008B", "#FFD700", "#808080"]
	var total = 0;
	for (let i=0; i<httpData.length;i++) {
	    let name = httpData[i].familyName;
	    let hours = httpData[i].weeks;
	    barData.datasets.push({
		label: name,
		backgroundColor: colourList[total%8],
		borderWidth: 1,
		data: hours});
	    total ++;
	}
	window.myBar.update();
    });
    xhttp.open("GET", "http://localhost:8080/api/v1/charts");
    xhttp.send();
}
    
function loadDash() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let familyInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_familyList").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(familyInfo);

    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/dashboard`);

    xhttp.send();
}

function userList() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_listUsers").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);

        let userBtns = document.querySelectorAll("[id*='edit_']");
        for (let i = 0; i < userBtns.length; i++) {
            userBtns[i].addEventListener('click', loadEditUser);
        }
        let passBtns = document.querySelectorAll("[id*='pass_']");
        for (let i = 0; i < passBtns.length; i++) {
            passBtns[i].addEventListener('click', loadEditPassword);
        }
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/users`);
    xhttp.send();
}

function loadEditUser(e) {
    let userID = e.srcElement.id.split("_")[1];
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_editUser").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);
        document.querySelector("#cancel").addEventListener('click', userList);
        document.querySelector("#submit").addEventListener('click', submitUserEdit);
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/users?u=${userID}`);
    xhttp.send();
}

function loadEditPassword(e) {
    let userID = e.srcElement.id.split("_")[1];
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_password").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);
        document.querySelector("#cancel").addEventListener('click', userList);
        document.querySelector("#submit").addEventListener('click', submitPassword);
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/users?u=${userID}`);
    xhttp.send();
}


function familyList() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_listFamilies").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);
        let btns = document.querySelectorAll("[id*='edit_']");
        for(let i = 0; i < btns.length; i++) {
            btns[i].addEventListener("click", loadEditFamily);
        }
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/families`);
    xhttp.send();
}

function loadEditFamily(e) {
    let familyID = e.srcElement.id.split("_")[1];
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        console.log(xhttp.response);
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_editFamily").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);
        document.querySelector("#cancel").addEventListener('click', familyList);
        document.querySelector("#submit").addEventListener('click', submitFamilyEdit);
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/families?f=${familyID}`);
    xhttp.send();
}

function submitFamilyEdit() {
    let familyID = document.querySelector("#famId").value;
    let familyName = document.querySelector("#famName").value;
    let children = document.querySelector("#children").value;

    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        if (xhttp.status > 300) {
            alert('ERROR: Could not update family');
            return;
        }
        alert('SUCESS: Family updated');
        familyList();
    });
    xhttp.open("PUT", "http://localhost:8080/api/v1/admin/families");
    xhttp.send(JSON.stringify(family));
}

function newFamily() {
    let tmpl = document.querySelector("#tmpl_newFamily").innerHTML;
    document.querySelector("#displayData").innerHTML = tmpl;
    document.querySelector("#cancel").addEventListener('click', familyList);
    document.querySelector("#submit").addEventListener('click', submitNewFamily);
    document.querySelector("#addFacil").addEventListener('click', addParent);
    document.querySelector("#remFacil").addEventListener('click', () => {
        let p = document.querySelector("#parents").lastChild;
        document.querySelector("#parents").removeChild(p);
    });
}

function lonelyFacilitators() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let parents = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_newParent").innerHTML;
        func = doT.template(tmpl);
        document.querySelector("#parents").insertAdjacentHTML('beforeend', func(parents));
    });
    xhttp.open("GET", "http://localhost:8080/api/v1/admin/facilitators");
    xhttp.send();
}

function submitNewFamily() {
    let surname = document.querySelector("#famName").value;
    let numChild = parseInt(document.querySelector("#children").value);

    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        if (xhttp.status > 300) {
            alert('ERROR: Could not create family');
            return;
        }
        alert('SUCESS: Family created');
        familyList();
    });
    xhttp.open("POST", "http://localhost:8080/api/v1/admin/families");
    console.log(JSON.stringify(newFamily));
    xhttp.send(JSON.stringify(newFamily));
}

function familyInfo() {

}

function submitUserEdit() {  
        let fields = document.querySelectorAll("input");
        for(let i = 0; i < fields.length; i++) {
            if (fields[i].value == "") {
                alert('Please fill out all sections');
                return;
            } 
        }
        let uId = parseInt(document.querySelector("#IDNum").innerHTML);
        let newFName = document.querySelector("#fname").value;
        let newLName = document.querySelector("#lName").value;
        let newEmail = document.querySelector("#email").value;
        let newPhone = document.querySelector("#phoneNum").value; 
        let newUName = `${newLName.toLowerCase()}_${newFName.toLowerCase()}`;
        let xhttp = new XMLHttpRequest();
        xhttp.addEventListener("loadend", () => {
            if (xhttp.status > 300) {
                alert('ERROR: Could not update user.');
                return;
            }
            if (xhttp.status == 200) {
                alert('SUCCESS: User updated.');
                userList();
            }
        });
        xhttp.open("PUT", "http://localhost:8080/api/v1/admin/users");
        xhttp.send(JSON.stringify({userid:uId, username:newUName,
                    firstname: newFName, lastname:newLName,
                    email:newEmail, phoneNumber:newPhone}));
}

function newUser() {
    let tmpl = document.querySelector("#tmpl_addUser").innerHTML;
    document.querySelector("#displayData").innerHTML = tmpl;
    document.querySelector("#cancel").addEventListener('click', userList);
    document.querySelector("#submit").addEventListener('click', () => {
        
        let fields = document.querySelectorAll("input");
        for(let i = 0; i < fields.length; i++) {
            if (fields[i].value == "" && fields[i].id != "bonusNote") {
                alert('Please fill out all sections');
                return;
            } 
        }
        if (document.querySelector("#pass1").value.length < 8) {
            alert('Password must be longer than 8 characters');
            return;
        }
        if (document.querySelector("#pass1").value != document.querySelector("#pass2").value) {
            alert('Passwords do not match')
            return;
        }
        let newRole = parseInt(document.querySelector("#role").value);
        let newFName = document.querySelector("#fname").value;
        let newLName = document.querySelector("#lName").value;
        let newEmail = document.querySelector("#email").value;
        let newPhone = document.querySelector("#phoneNum").value;
        let newUName = `${newLName}${newFName}`.toLowerCase();
        let newPass = document.querySelector("#pass1").value;
        let bHours = parseInt(document.querySelector("#bonusHours").value);
        let bNote = document.querySelector("#bonusNote").value;
        let newPassData = [];
        for (let i = 0; i < newPass.length; i++) {
            newPassData.push(newPass.charCodeAt(i));
        }

        let xhttp = new XMLHttpRequest();
        xhttp.addEventListener("loadend", () => {
            if (xhttp.status > 300) {
                alert('ERROR: Could not create user.');
                return;
            }
            if (xhttp.status == 201) {
                alert('SUCCESS: User added to system');
                userList();
            }
        });
        xhttp.open("POST", "http://localhost:8080/api/v1/admin/users");
        xhttp.send(JSON.stringify({userrole:newRole, username:newUName,
                    password:newPassData, firstname: newFName, lastname:newLName,
                    email:newEmail, phoneNumber:newPhone, bonusHours:bHours, bonusNote:bNote}));
    });
}

function userInfo() {

}

function listClasses() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let classes = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_listClasses").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(classes);
        
        let btns = document.querySelectorAll("[id*='edit_']");
        for(let i = 0; i < btns.length; i++) {
            btns[i].addEventListener("click", loadEditClass);
        }
     });
    xhttp.open("GET", "http://localhost:8080/api/v1/admin/classes");
    xhttp.send();
}

function loadEditClass() {

}

function addClassRoom() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        let facilitators = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_createClass").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(facilitators);
        document.querySelector("#cancel").addEventListener('click', listClasses);
        document.querySelector("#submit").addEventListener('click', submitNewClass);
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/teachers`);
    xhttp.send();
}

function submitNewClass() {
    let newName = document.querySelector("#cName").value;
    let newTeacher = parseInt(document.querySelector("#cTeacher").value);
    let newRoomNum = document.querySelector("#cNum").value;

    if (newName == "") {
        alert('Class name cannot be empty');
        return;
    }
    if (newRoomNum == "") {
        alert('Room number cannot be empty');
        return;
    }

    let newClass = {"roomNAme":newName, "teacherId":newTeacher, "roomNum":newRoomNum};
    console.log(newClass);
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        if(xhttp.status > 300) {
            alert('ERROR: Could not create class');
            return;
        }
        alert('SUCCESS: Class created.');
        listClasses();
    });
    xhttp.open("POST", "http://localhost:8080/api/v1/admin/classes");
    xhttp.send(JSON.stringify(newClass));
}
