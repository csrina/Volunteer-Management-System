document.addEventListener("DOMContentLoaded", function() {
    setActiveCategory();
})

//sets active category in top bars
function setActiveCategory() {
    let cat = window.location.href.split("/").pop();
    document.querySelector(`#${cat}Btn`).setAttribute('class','active');
}

function userList() {
    let xhttp = new XMLHttpRequest();
    xhttp.addEventListener("loadend", () => {
        console.log(xhttp.response);
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
        console.log(xhttp.response);
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
        console.log(xhttp.response);
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
        console.log(xhttp.response);
        let userInfo = JSON.parse(xhttp.response);
        let tmpl = document.querySelector("#tmpl_listFamilies").innerHTML;
        let func = doT.template(tmpl);
        document.querySelector("#displayData").innerHTML = func(userInfo);
    });
    xhttp.open("GET", `http://localhost:8080/api/v1/admin/families`);
    xhttp.send();
}
