function clearErrors() {
    var ngroup = document.querySelector("#usrgroup");
    var pgroup = document.querySelector("#pwdgroup");
    ngroup.classList.remove('has-error');
    pgroup.classList.remove('has-error');
}

function missingCheck(name, pass) {
    if (name == "" && pass == "") {
        console.log("both missing");
        var theDiv = document.querySelector("#errorbox");
        theDiv.innerHTML= `<div class="alert alert-danger" role="alert">Username and password are missing</div>`;
        var ngroup = document.querySelector("#usrgroup");
        var pgroup = document.querySelector("#pwdgroup");
        ngroup.classList.add('has-error');
        pgroup.classList.add('has-error');
        return false;
    }
    if (name == "") {
        console.log("name missing");
        var theDiv = document.querySelector("#errorbox");
        theDiv.innerHTML= `<div class="alert alert-danger" role="alert">Username is missing</div>`
        var ngroup = document.querySelector("#usrgroup");
        ngroup.classList.add('has-error');
        return false;
    }
    
    if (pass == "") {
        console.log("password missing");
        var theDiv = document.querySelector("#errorbox");
        theDiv.innerHTML= `<div class="alert alert-danger" role="alert">Password is missing</div>`
        var pgroup = document.querySelector("#pwdgroup");
        pgroup.classList.add('has-error');
        return false;
    }
    return true;
}

function getApiCall(cur) {
    var u = "/api/v1/login/";
    if(cur == "/login/facilitator") {
        u = u.concat("facilitator/")
    } else if (cur == "/login/teacher") {
        u = u.concat("teacher/")
    } else if (cur == "/login/admin") {
        u = u.concat("admin/")
    }
    return u;
}

function login(){
    clearErrors();
    
    var name= document.querySelector("#usr").value;
    var pass = document.querySelector("#pwd").value;
    var filled = missingCheck(name, pass)
    if (filled == false) {
        return
    }
    var u = getApiCall(window.location.pathname);

    var xmlhttp = new XMLHttpRequest();   // new HttpRequest instance 
    xmlhttp.onreadystatechange= function() {
        if (this.readyState!==4) return; // not ready yet
        if (this.status!==202) { // HTTP 200 OK
            console.log(this.status);
            console.log("Login attempt failed");
            var theDiv = document.querySelector("#errorbox");
            theDiv.innerHTML= `<div class="alert alert-danger" role="alert">Login attempt failed, invalid username or password</div>`
        } else {
            console.log(this.status);
            console.log("Login attempt pass");
            // var theDiv = document.querySelector("#errorbox");
            // theDiv.innerHTML= `<div class="alert alert-success" role="alert">Login attempt was a success</div>`
            window.location="/dashboard";
        }
    };
    xmlhttp.open("POST", u, true);
    xmlhttp.setRequestHeader("Content-Type", "application/json");
    
    var data = [];
    for (var i = 0; i < pass.length; i++){  
        data.push(pass.charCodeAt(i));
    }
    
    xmlhttp.send(JSON.stringify({username:name, password:data}));
    
}

function logout() {
    window.location="/logout";
    console.log("test")
}